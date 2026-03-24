package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"am-erp-go/internal/infrastructure/auth"
	"am-erp-go/internal/infrastructure/response"
	integrationDomain "am-erp-go/internal/module/integration/domain"
	integrationUsecase "am-erp-go/internal/module/integration/usecase"
	systemUsecase "am-erp-go/internal/module/system/usecase"

	"github.com/gin-gonic/gin"
)

type AuthorizationHandler struct {
	usecase     *integrationUsecase.AuthorizationUsecase
	auditLogger AuthorizationAuditLogger
}

func NewAuthorizationHandler(usecase *integrationUsecase.AuthorizationUsecase) *AuthorizationHandler {
	return &AuthorizationHandler{usecase: usecase}
}

type AuthorizationAuditLogger interface {
	RecordFromContext(c *gin.Context, payload systemUsecase.AuditLogPayload) error
}

func (h *AuthorizationHandler) BindAuditLogger(logger AuthorizationAuditLogger) {
	h.auditLogger = logger
}

func (h *AuthorizationHandler) ListProviders(c *gin.Context) {
	response.Success(c, h.usecase.ListProviders())
}

func (h *AuthorizationHandler) ListAuthorizations(c *gin.Context) {
	params := &integrationDomain.ListAuthorizationParams{
		Page:         parseIntOrDefault(c.Query("page"), 1),
		PageSize:     parseIntOrDefault(c.Query("page_size"), 20),
		ProviderCode: c.Query("provider_code"),
		Status:       c.Query("status"),
	}
	list, total, err := h.usecase.ListAuthorizations(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, list, total, params.Page, params.PageSize)
}

type startAuthorizationRequest struct {
	ProviderCode string `json:"provider_code" binding:"required"`
	AccountAlias string `json:"account_alias"`
}

func (h *AuthorizationHandler) StartAuthorization(c *gin.Context) {
	var req startAuthorizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	operatorID := parseOperatorIDFromContext(c)
	result, err := h.usecase.StartAuthorization(&integrationUsecase.StartAuthorizationInput{
		ProviderCode: req.ProviderCode,
		AccountAlias: req.AccountAlias,
		OperatorID:   operatorID,
	})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	h.recordAudit(c, "CREATE_PLATFORM_AUTHORIZATION", "IntegrationAuthorization", fmt.Sprintf("%d", result.AuthorizationID), nil, map[string]any{
		"id":            result.AuthorizationID,
		"provider_code": result.ProviderCode,
		"account_alias": strings.TrimSpace(req.AccountAlias),
		"status":        integrationDomain.AuthorizationStatusPending,
	})
	response.Success(c, result)
}

func (h *AuthorizationHandler) ManualRefresh(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	operatorID := parseOperatorIDFromContext(c)
	before, _ := h.usecase.GetAuthorization(id)
	record, err := h.usecase.ManualRefresh(c.Request.Context(), id, operatorID)
	if err != nil {
		after, _ := h.usecase.GetAuthorization(id)
		h.recordAuditIfChanged(c, "MANUAL_REFRESH_TOKEN", "IntegrationAuthorization", fmt.Sprintf("%d", id), before, after)
		response.BadRequest(c, err.Error())
		return
	}
	h.recordAuditIfChanged(c, "MANUAL_REFRESH_TOKEN", "IntegrationAuthorization", fmt.Sprintf("%d", id), before, record)
	response.Success(c, record)
}

func (h *AuthorizationHandler) HandleOAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	record, err := h.usecase.HandleCallback(c.Request.Context(), &integrationUsecase.OAuthCallbackInput{
		ProviderCode:      provider,
		OAuthState:        c.Query("state"),
		AuthorizationCode: c.Query("spapi_oauth_code"),
		SellerPartnerID:   c.Query("selling_partner_id"),
		OAuthError:        c.Query("error"),
		OAuthErrorDesc:    c.Query("error_description"),
	})
	if err != nil {
		writeOAuthResultHTML(c, false, err.Error())
		return
	}
	writeOAuthResultHTML(c, true, fmt.Sprintf("授权成功，提供方: %s", record.ProviderCode))
}

func parseOperatorIDFromContext(c *gin.Context) *uint64 {
	if userID, ok := c.Get(auth.UserIDKey); ok {
		if id, castOK := userID.(uint64); castOK {
			return &id
		}
	}
	return nil
}

func writeOAuthResultHTML(c *gin.Context, success bool, msg string) {
	title := "授权失败"
	if success {
		title = "授权成功"
	}
	escaped := htmlEscape(msg)
	html := fmt.Sprintf(`<!doctype html><html><head><meta charset="utf-8"><title>%s</title></head>
<body style="font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;padding:24px;">
<h3 style="margin:0 0 12px 0;">%s</h3>
<p style="margin:0 0 16px 0;color:#333;">%s</p>
<script>
if (window.opener) {
  window.opener.postMessage({ type: 'integration-oauth-callback', success: %t, message: %q }, '*');
}
setTimeout(function(){ window.close(); }, 1000);
</script>
</body></html>`, title, title, escaped, success, msg)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func htmlEscape(input string) string {
	s := strings.ReplaceAll(input, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func (h *AuthorizationHandler) recordAudit(c *gin.Context, action, entityType, entityID string, before, after any) {
	if h == nil || h.auditLogger == nil || c == nil {
		return
	}
	_ = h.auditLogger.RecordFromContext(c, systemUsecase.AuditLogPayload{
		Module:     "Integration",
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Before:     before,
		After:      after,
	})
}

func (h *AuthorizationHandler) recordAuditIfChanged(c *gin.Context, action, entityType, entityID string, before, after any) {
	beforeDiff, afterDiff, changed := buildIntegrationAuditDiff(before, after)
	if !changed {
		return
	}
	h.recordAudit(c, action, entityType, entityID, beforeDiff, afterDiff)
}

func buildIntegrationAuditDiff(before, after any) (any, any, bool) {
	normalizedBefore := normalizeIntegrationAuditValue(before)
	normalizedAfter := normalizeIntegrationAuditValue(after)
	return diffIntegrationAuditValues(normalizedBefore, normalizedAfter)
}

func normalizeIntegrationAuditValue(value any) any {
	if value == nil {
		return nil
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return value
	}
	var normalized any
	if err := json.Unmarshal(raw, &normalized); err != nil {
		return value
	}
	return scrubIntegrationAuditValue(normalized)
}

func scrubIntegrationAuditValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		cleaned := make(map[string]any, len(typed))
		for key, item := range typed {
			switch key {
			case "oauth_state", "oauth_state_expire_at", "access_token", "refresh_token", "token_scope",
				"gmt_create", "gmt_modified", "created_at", "updated_at":
				continue
			}
			cleaned[key] = scrubIntegrationAuditValue(item)
		}
		return cleaned
	case []any:
		items := make([]any, 0, len(typed))
		for _, item := range typed {
			items = append(items, scrubIntegrationAuditValue(item))
		}
		return items
	default:
		return value
	}
}

func diffIntegrationAuditValues(before, after any) (any, any, bool) {
	if reflect.DeepEqual(before, after) {
		return nil, nil, false
	}
	beforeMap, beforeMapOK := before.(map[string]any)
	afterMap, afterMapOK := after.(map[string]any)
	if beforeMapOK && afterMapOK {
		beforeDiff := map[string]any{}
		afterDiff := map[string]any{}
		keys := make(map[string]struct{}, len(beforeMap)+len(afterMap))
		for key := range beforeMap {
			keys[key] = struct{}{}
		}
		for key := range afterMap {
			keys[key] = struct{}{}
		}
		for key := range keys {
			childBefore, childAfter, changed := diffIntegrationAuditValues(beforeMap[key], afterMap[key])
			if !changed {
				continue
			}
			beforeDiff[key] = childBefore
			afterDiff[key] = childAfter
		}
		if len(beforeDiff) == 0 && len(afterDiff) == 0 {
			return nil, nil, false
		}
		return beforeDiff, afterDiff, true
	}
	return before, after, true
}
