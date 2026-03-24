package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	integrationDomain "am-erp-go/internal/module/integration/domain"
)

type ProviderSummary struct {
	Code        string `json:"code"`
	Type        string `json:"type"`
	DisplayName string `json:"display_name"`
}

type StartAuthorizationInput struct {
	ProviderCode string
	AccountAlias string
	OperatorID   *uint64
}

type StartAuthorizationResult struct {
	AuthorizationID uint64 `json:"authorization_id"`
	ProviderCode    string `json:"provider_code"`
	AuthorizeURL    string `json:"authorize_url"`
	OAuthState      string `json:"oauth_state"`
	ExpireAt        string `json:"expire_at"`
}

type OAuthCallbackInput struct {
	ProviderCode      string
	OAuthState        string
	AuthorizationCode string
	SellerPartnerID   string
	OAuthError        string
	OAuthErrorDesc    string
}

type AuthorizationOptions struct {
	RefreshRetryTimes    int
	RefreshFailThreshold uint8
}

type AuthorizationRefreshSummary struct {
	Total   int `json:"total"`
	Success int `json:"success"`
	Failed  int `json:"failed"`
}

type AuthorizationUsecase struct {
	repo                 integrationDomain.AuthorizationRepository
	providers            map[string]integrationDomain.AuthorizationProvider
	nowFn                func() time.Time
	refreshRetryTimes    int
	refreshFailThreshold uint8
}

func NewAuthorizationUsecase(
	repo integrationDomain.AuthorizationRepository,
	providers []integrationDomain.AuthorizationProvider,
	nowFn func() time.Time,
	options ...AuthorizationOptions,
) *AuthorizationUsecase {
	if nowFn == nil {
		nowFn = time.Now
	}
	opt := AuthorizationOptions{
		RefreshRetryTimes:    3,
		RefreshFailThreshold: 3,
	}
	if len(options) > 0 {
		if options[0].RefreshRetryTimes > 0 {
			opt.RefreshRetryTimes = options[0].RefreshRetryTimes
		}
		if options[0].RefreshFailThreshold > 0 {
			opt.RefreshFailThreshold = options[0].RefreshFailThreshold
		}
	}

	reg := map[string]integrationDomain.AuthorizationProvider{}
	for _, provider := range providers {
		if provider == nil {
			continue
		}
		code := normalizeProviderCode(provider.Code())
		if code == "" {
			continue
		}
		reg[code] = provider
	}

	return &AuthorizationUsecase{
		repo:                 repo,
		providers:            reg,
		nowFn:                nowFn,
		refreshRetryTimes:    opt.RefreshRetryTimes,
		refreshFailThreshold: opt.RefreshFailThreshold,
	}
}

func (u *AuthorizationUsecase) ListProviders() []ProviderSummary {
	out := make([]ProviderSummary, 0, len(u.providers))
	for code, provider := range u.providers {
		display := code
		if provider.Type() != "" {
			display = fmt.Sprintf("%s (%s)", code, strings.ToUpper(provider.Type()))
		}
		out = append(out, ProviderSummary{
			Code:        code,
			Type:        strings.ToLower(provider.Type()),
			DisplayName: display,
		})
	}
	return out
}

func (u *AuthorizationUsecase) ListAuthorizations(params *integrationDomain.ListAuthorizationParams) ([]integrationDomain.IntegrationAuthorization, int64, error) {
	if params == nil {
		params = &integrationDomain.ListAuthorizationParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.ProviderCode != "" {
		params.ProviderCode = normalizeProviderCode(params.ProviderCode)
	}
	return u.repo.ListAuthorizations(params)
}

func (u *AuthorizationUsecase) GetAuthorization(id uint64) (*integrationDomain.IntegrationAuthorization, error) {
	if id == 0 {
		return nil, fmt.Errorf("id is required")
	}
	return u.repo.GetAuthorizationByID(id)
}

func (u *AuthorizationUsecase) StartAuthorization(input *StartAuthorizationInput) (*StartAuthorizationResult, error) {
	if input == nil {
		return nil, fmt.Errorf("input is required")
	}
	providerCode := normalizeProviderCode(input.ProviderCode)
	provider, err := u.getProvider(providerCode)
	if err != nil {
		return nil, err
	}

	state, err := generateOAuthState()
	if err != nil {
		return nil, err
	}
	authorizeURL, err := provider.BuildAuthorizeURL(state)
	if err != nil {
		return nil, err
	}

	now := u.nowFn().UTC()
	expireAt := now.Add(15 * time.Minute)
	accountAlias := strPtrOrNil(input.AccountAlias)

	record := &integrationDomain.IntegrationAuthorization{
		ProviderCode:       providerCode,
		ProviderType:       strings.ToLower(provider.Type()),
		AccountAlias:       accountAlias,
		Status:             integrationDomain.AuthorizationStatusPending,
		OAuthState:         &state,
		OAuthStateExpireAt: &expireAt,
		CreatedBy:          input.OperatorID,
		UpdatedBy:          input.OperatorID,
	}
	if err := u.repo.CreateAuthorization(record); err != nil {
		return nil, err
	}

	return &StartAuthorizationResult{
		AuthorizationID: record.ID,
		ProviderCode:    providerCode,
		AuthorizeURL:    authorizeURL,
		OAuthState:      state,
		ExpireAt:        expireAt.Format(time.RFC3339),
	}, nil
}

func (u *AuthorizationUsecase) HandleCallback(ctx context.Context, input *OAuthCallbackInput) (*integrationDomain.IntegrationAuthorization, error) {
	if input == nil {
		return nil, fmt.Errorf("callback input is required")
	}
	providerCode := normalizeProviderCode(input.ProviderCode)
	provider, err := u.getProvider(providerCode)
	if err != nil {
		return nil, err
	}

	oauthState := strings.TrimSpace(input.OAuthState)
	if oauthState == "" {
		return nil, fmt.Errorf("oauth state is required")
	}
	record, err := u.repo.GetAuthorizationByProviderAndState(providerCode, oauthState)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, fmt.Errorf("authorization session not found")
	}

	now := u.nowFn().UTC()
	if record.OAuthStateExpireAt != nil && now.After(record.OAuthStateExpireAt.UTC()) {
		msg := "oauth state expired"
		record.Status = integrationDomain.AuthorizationStatusFailed
		record.RefreshFailCount = 0
		record.LastRefreshAttemptAt = nil
		record.LastRefreshFailedAt = nil
		record.LastErrorMessage = &msg
		record.OAuthState = nil
		record.OAuthStateExpireAt = nil
		_ = u.repo.UpdateAuthorization(record)
		return nil, errors.New(msg)
	}

	if strings.TrimSpace(input.OAuthError) != "" {
		msg := strings.TrimSpace(input.OAuthError)
		if desc := strings.TrimSpace(input.OAuthErrorDesc); desc != "" {
			msg = msg + ": " + desc
		}
		record.Status = integrationDomain.AuthorizationStatusFailed
		record.RefreshFailCount = 0
		record.LastRefreshAttemptAt = nil
		record.LastRefreshFailedAt = nil
		record.LastErrorMessage = &msg
		record.OAuthState = nil
		record.OAuthStateExpireAt = nil
		if err := u.repo.UpdateAuthorization(record); err != nil {
			return nil, err
		}
		return nil, errors.New(msg)
	}

	authCode := strings.TrimSpace(input.AuthorizationCode)
	if authCode == "" {
		return nil, fmt.Errorf("authorization code is required")
	}

	token, err := provider.ExchangeAuthorizationCode(ctx, authCode)
	if err != nil {
		msg := err.Error()
		record.Status = integrationDomain.AuthorizationStatusFailed
		record.RefreshFailCount = 0
		record.LastRefreshAttemptAt = nil
		record.LastRefreshFailedAt = nil
		record.LastErrorMessage = &msg
		record.OAuthState = nil
		record.OAuthStateExpireAt = nil
		_ = u.repo.UpdateAuthorization(record)
		return nil, err
	}

	expiresAt := now.Add(time.Duration(maxI64(token.ExpiresIn, 1)) * time.Second)
	sellerPartnerID := strPtrOrNil(input.SellerPartnerID)
	refreshToken := token.RefreshToken
	if strings.TrimSpace(refreshToken) == "" && record.RefreshToken != nil {
		refreshToken = *record.RefreshToken
	}

	record.SellerPartnerID = sellerPartnerID
	record.AccessToken = strPtrOrNil(token.AccessToken)
	record.RefreshToken = strPtrOrNil(refreshToken)
	record.AccessTokenExpireAt = &expiresAt
	record.TokenScope = strPtrOrNil(token.Scope)
	record.Status = integrationDomain.AuthorizationStatusAuthorized
	record.LastAuthorizedAt = &now
	record.LastRefreshAt = &now
	record.RefreshFailCount = 0
	record.LastRefreshAttemptAt = &now
	record.LastRefreshFailedAt = nil
	record.LastErrorMessage = nil
	record.OAuthState = nil
	record.OAuthStateExpireAt = nil
	if err := u.repo.UpdateAuthorization(record); err != nil {
		return nil, err
	}

	return record, nil
}

func (u *AuthorizationUsecase) ManualRefresh(ctx context.Context, id uint64, operatorID *uint64) (*integrationDomain.IntegrationAuthorization, error) {
	record, err := u.repo.GetAuthorizationByID(id)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, fmt.Errorf("authorization not found")
	}
	if record.Status != integrationDomain.AuthorizationStatusAuthorized && record.Status != integrationDomain.AuthorizationStatusFailed {
		return nil, fmt.Errorf("authorization status %s cannot refresh", record.Status)
	}
	return u.refreshOne(ctx, record, operatorID)
}

func (u *AuthorizationUsecase) RefreshExpiringTokens(ctx context.Context) (*AuthorizationRefreshSummary, error) {
	summary := &AuthorizationRefreshSummary{}
	before := u.nowFn().UTC().Add(10 * time.Minute)
	records, err := u.repo.ListNeedRefresh(before, 300)
	if err != nil {
		return nil, err
	}
	summary.Total = len(records)
	for i := range records {
		record := records[i]
		if _, err := u.refreshOne(ctx, &record, nil); err != nil {
			summary.Failed++
			// 不中断批处理
			continue
		}
		summary.Success++
	}
	return summary, nil
}

func (u *AuthorizationUsecase) refreshOne(ctx context.Context, record *integrationDomain.IntegrationAuthorization, operatorID *uint64) (*integrationDomain.IntegrationAuthorization, error) {
	if record == nil {
		return nil, fmt.Errorf("authorization record is required")
	}
	provider, err := u.getProvider(record.ProviderCode)
	if err != nil {
		return nil, err
	}
	refreshToken := ""
	if record.RefreshToken != nil {
		refreshToken = strings.TrimSpace(*record.RefreshToken)
	}
	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token is empty")
	}

	attemptAt := u.nowFn().UTC()
	var token *integrationDomain.OAuthToken
	var refreshErr error
	attempts := int(maxI64(int64(u.refreshRetryTimes), 1))
	for i := 0; i < attempts; i++ {
		token, refreshErr = provider.RefreshAccessToken(ctx, refreshToken)
		if refreshErr == nil {
			break
		}
	}
	if refreshErr != nil {
		msg := refreshErr.Error()
		record.LastRefreshAttemptAt = &attemptAt
		record.LastRefreshFailedAt = &attemptAt
		record.RefreshFailCount = saturatingIncUint8(record.RefreshFailCount)
		if record.RefreshFailCount >= u.refreshFailThreshold {
			record.Status = integrationDomain.AuthorizationStatusFailed
		}
		record.LastErrorMessage = &msg
		record.UpdatedBy = operatorID
		_ = u.repo.UpdateAuthorization(record)
		return nil, refreshErr
	}

	now := u.nowFn().UTC()
	expiresAt := now.Add(time.Duration(maxI64(token.ExpiresIn, 1)) * time.Second)
	nextRefreshToken := refreshToken
	if strings.TrimSpace(token.RefreshToken) != "" {
		nextRefreshToken = strings.TrimSpace(token.RefreshToken)
	}

	record.AccessToken = strPtrOrNil(token.AccessToken)
	record.RefreshToken = strPtrOrNil(nextRefreshToken)
	record.AccessTokenExpireAt = &expiresAt
	record.TokenScope = strPtrOrNil(token.Scope)
	record.Status = integrationDomain.AuthorizationStatusAuthorized
	record.LastRefreshAt = &now
	record.LastRefreshAttemptAt = &now
	record.LastRefreshFailedAt = nil
	record.RefreshFailCount = 0
	record.LastErrorMessage = nil
	record.UpdatedBy = operatorID
	if err := u.repo.UpdateAuthorization(record); err != nil {
		return nil, err
	}
	return record, nil
}

func (u *AuthorizationUsecase) getProvider(providerCode string) (integrationDomain.AuthorizationProvider, error) {
	if u == nil || u.repo == nil {
		return nil, fmt.Errorf("authorization usecase not configured")
	}
	provider, ok := u.providers[normalizeProviderCode(providerCode)]
	if !ok || provider == nil {
		return nil, fmt.Errorf("provider not supported: %s", providerCode)
	}
	return provider, nil
}

func normalizeProviderCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func generateOAuthState() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func strPtrOrNil(v string) *string {
	s := strings.TrimSpace(v)
	if s == "" {
		return nil
	}
	return &s
}

func maxI64(v int64, fallback int64) int64 {
	if v > 0 {
		return v
	}
	return fallback
}

func saturatingIncUint8(v uint8) uint8 {
	if v >= 255 {
		return 255
	}
	return v + 1
}
