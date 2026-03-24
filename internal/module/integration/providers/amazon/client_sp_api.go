package amazon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	integrationDomain "am-erp-go/internal/module/integration/domain"
)

const (
	defaultSPAPIEndpoint          = "https://sellingpartnerapi-na.amazon.com"
	defaultAmazonAuthorizeBaseURL = "https://sellercentral.amazon.com/apps/authorize/consent"
)

type SPAPIConfig struct {
	Endpoint          string
	AppID             string
	AuthorizeBaseURL  string
	RedirectURI       string
	MarketplaceIDs    []string
	LWAClientID       string
	LWAClientSecret   string
	LWARefreshToken   string
	RequestTimeoutSec int
}

type SPAPIClient struct {
	cfg        SPAPIConfig
	httpClient *http.Client

	tokenMu        sync.Mutex
	accessToken    string
	accessTokenExp time.Time
}

func NewSPAPIClient(cfg SPAPIConfig) *SPAPIClient {
	if strings.TrimSpace(cfg.Endpoint) == "" {
		cfg.Endpoint = defaultSPAPIEndpoint
	}
	if strings.TrimSpace(cfg.AuthorizeBaseURL) == "" {
		cfg.AuthorizeBaseURL = defaultAmazonAuthorizeBaseURL
	}
	if cfg.RequestTimeoutSec <= 0 {
		cfg.RequestTimeoutSec = 20
	}
	return &SPAPIClient{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.RequestTimeoutSec) * time.Second,
		},
	}
}

func (c *SPAPIClient) ListOrders(ctx context.Context, req integrationDomain.ExternalListOrdersRequest) ([]integrationDomain.ExternalOrder, error) {
	token, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	marketplaceIDs := req.MarketplaceIDs
	if len(marketplaceIDs) == 0 {
		marketplaceIDs = c.cfg.MarketplaceIDs
	}
	if len(marketplaceIDs) == 0 {
		return nil, fmt.Errorf("amazon marketplace ids are required")
	}
	if req.LastUpdatedAfter.IsZero() {
		return nil, fmt.Errorf("LastUpdatedAfter is required")
	}

	orders := make([]integrationDomain.ExternalOrder, 0)
	nextToken := ""
	for {
		pageOrders, pageNextToken, err := c.listOrdersPage(ctx, token, marketplaceIDs, req.LastUpdatedAfter.UTC(), nextToken)
		if err != nil {
			return nil, err
		}
		for _, order := range pageOrders {
			items, err := c.listOrderItems(ctx, token, order.OrderID)
			if err != nil {
				return nil, err
			}
			order.Items = items
			orders = append(orders, order)
		}
		if strings.TrimSpace(pageNextToken) == "" {
			break
		}
		nextToken = pageNextToken
	}

	return orders, nil
}

func (c *SPAPIClient) ListRefunds(ctx context.Context, req integrationDomain.ExternalListRefundsRequest) ([]integrationDomain.ExternalRefund, error) {
	token, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	if req.PostedAfter.IsZero() {
		return nil, fmt.Errorf("PostedAfter is required")
	}

	refunds := make([]integrationDomain.ExternalRefund, 0)
	nextToken := ""
	refundSeq := 0
	for {
		pageRefunds, pageNextToken, err := c.listRefundEventsPage(ctx, token, req, nextToken, &refundSeq)
		if err != nil {
			return nil, err
		}
		refunds = append(refunds, pageRefunds...)
		if strings.TrimSpace(pageNextToken) == "" {
			break
		}
		nextToken = pageNextToken
	}
	return refunds, nil
}

func (c *SPAPIClient) BuildAuthorizeURL(state string) (string, error) {
	baseURL := strings.TrimSpace(c.cfg.AuthorizeBaseURL)
	appID := strings.TrimSpace(c.cfg.AppID)
	redirectURI := strings.TrimSpace(c.cfg.RedirectURI)
	if baseURL == "" || appID == "" || redirectURI == "" {
		return "", fmt.Errorf("amazon oauth config incomplete: authorize_base_url/app_id/redirect_uri required")
	}
	if strings.TrimSpace(state) == "" {
		return "", fmt.Errorf("oauth state is required")
	}

	query := url.Values{}
	query.Set("application_id", appID)
	query.Set("state", state)
	query.Set("redirect_uri", redirectURI)
	query.Set("version", "beta")

	if strings.Contains(baseURL, "?") {
		return baseURL + "&" + query.Encode(), nil
	}
	return baseURL + "?" + query.Encode(), nil
}

func (c *SPAPIClient) ExchangeAuthorizationCode(ctx context.Context, code string) (*integrationDomain.OAuthToken, error) {
	return c.fetchToken(ctx, url.Values{
		"grant_type":    []string{"authorization_code"},
		"code":          []string{strings.TrimSpace(code)},
		"client_id":     []string{c.cfg.LWAClientID},
		"client_secret": []string{c.cfg.LWAClientSecret},
	})
}

func (c *SPAPIClient) RefreshAccessToken(ctx context.Context, refreshToken string) (*integrationDomain.OAuthToken, error) {
	return c.fetchToken(ctx, url.Values{
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{strings.TrimSpace(refreshToken)},
		"client_id":     []string{c.cfg.LWAClientID},
		"client_secret": []string{c.cfg.LWAClientSecret},
	})
}

func (c *SPAPIClient) listOrdersPage(
	ctx context.Context,
	token string,
	marketplaceIDs []string,
	lastUpdatedAfter time.Time,
	nextToken string,
) ([]integrationDomain.ExternalOrder, string, error) {
	query := url.Values{}
	if nextToken != "" {
		query.Set("NextToken", nextToken)
	} else {
		query.Set("MarketplaceIds", strings.Join(marketplaceIDs, ","))
		query.Set("LastUpdatedAfter", lastUpdatedAfter.Format(time.RFC3339))
		query.Set("MaxResultsPerPage", "100")
	}

	respBody, err := c.doSPAPIRequest(ctx, http.MethodGet, "/orders/v0/orders", query, nil, token)
	if err != nil {
		return nil, "", err
	}

	var parsed struct {
		Payload struct {
			Orders []struct {
				AmazonOrderID  string `json:"AmazonOrderId"`
				MarketplaceID  string `json:"MarketplaceId"`
				PurchaseDate   string `json:"PurchaseDate"`
				LastUpdateDate string `json:"LastUpdateDate"`
				OrderTotal     struct {
					CurrencyCode string `json:"CurrencyCode"`
				} `json:"OrderTotal"`
			} `json:"Orders"`
			NextToken string `json:"NextToken"`
		} `json:"payload"`
	}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, "", fmt.Errorf("decode orders response failed: %w", err)
	}

	out := make([]integrationDomain.ExternalOrder, 0, len(parsed.Payload.Orders))
	for _, item := range parsed.Payload.Orders {
		out = append(out, integrationDomain.ExternalOrder{
			OrderID:       item.AmazonOrderID,
			MarketplaceID: item.MarketplaceID,
			PurchaseAt:    parseTimeRFC3339(item.PurchaseDate),
			LastUpdatedAt: parseTimeRFC3339(item.LastUpdateDate),
			Currency:      item.OrderTotal.CurrencyCode,
		})
	}
	return out, parsed.Payload.NextToken, nil
}

func (c *SPAPIClient) listRefundEventsPage(
	ctx context.Context,
	token string,
	req integrationDomain.ExternalListRefundsRequest,
	nextToken string,
	refundSeq *int,
) ([]integrationDomain.ExternalRefund, string, error) {
	query := url.Values{}
	if strings.TrimSpace(nextToken) != "" {
		query.Set("NextToken", strings.TrimSpace(nextToken))
	} else {
		query.Set("PostedAfter", req.PostedAfter.UTC().Format(time.RFC3339))
		query.Set("MaxResultsPerPage", "100")
	}

	respBody, err := c.doSPAPIRequest(ctx, http.MethodGet, "/finances/v0/financialEvents", query, nil, token)
	if err != nil {
		return nil, "", err
	}
	var parsed struct {
		Payload struct {
			RefundEventList []struct {
				AmazonOrderID              string `json:"AmazonOrderId"`
				PostedDate                 string `json:"PostedDate"`
				ShipmentItemAdjustmentList []struct {
					SellerSKU                string      `json:"SellerSKU"`
					OrderItemID              string      `json:"OrderItemId"`
					QuantityShipped          interface{} `json:"QuantityShipped"`
					ItemChargeAdjustmentList []struct {
						ChargeAmount struct {
							CurrencyCode   string `json:"CurrencyCode"`
							CurrencyAmount string `json:"CurrencyAmount"`
						} `json:"ChargeAmount"`
					} `json:"ItemChargeAdjustmentList"`
				} `json:"ShipmentItemAdjustmentList"`
			} `json:"RefundEventList"`
			NextToken string `json:"NextToken"`
		} `json:"payload"`
	}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, "", fmt.Errorf("decode refund events response failed: %w", err)
	}
	out := make([]integrationDomain.ExternalRefund, 0)
	marketplaceID := ""
	if len(req.MarketplaceIDs) > 0 {
		marketplaceID = strings.TrimSpace(req.MarketplaceIDs[0])
	}
	for _, event := range parsed.Payload.RefundEventList {
		postedAt := parseTimeRFC3339(event.PostedDate)
		if postedAt.IsZero() {
			postedAt = req.PostedAfter
		}
		for _, item := range event.ShipmentItemAdjustmentList {
			qty := parsePositiveUint(item.QuantityShipped)
			if qty == 0 {
				qty = 1
			}
			amount := 0.0
			currency := ""
			for _, charge := range item.ItemChargeAdjustmentList {
				v, parseErr := strconv.ParseFloat(strings.TrimSpace(charge.ChargeAmount.CurrencyAmount), 64)
				if parseErr != nil {
					continue
				}
				if v < 0 {
					v = -v
				}
				amount += v
				if currency == "" {
					currency = strings.TrimSpace(charge.ChargeAmount.CurrencyCode)
				}
			}
			if amount == 0 {
				continue
			}
			seq := 0
			if refundSeq != nil {
				*refundSeq = *refundSeq + 1
				seq = *refundSeq
			}
			refundID := fmt.Sprintf("%s|%s|%s|%s|%d", strings.TrimSpace(event.AmazonOrderID), postedAt.UTC().Format(time.RFC3339), strings.TrimSpace(item.OrderItemID), strings.TrimSpace(item.SellerSKU), seq)
			out = append(out, integrationDomain.ExternalRefund{
				RefundID:      refundID,
				OrderID:       strings.TrimSpace(event.AmazonOrderID),
				OrderItemID:   strings.TrimSpace(item.OrderItemID),
				SellerSKU:     strings.TrimSpace(item.SellerSKU),
				MarketplaceID: marketplaceID,
				Quantity:      qty,
				Amount:        amount,
				Currency:      strings.TrimSpace(currency),
				PostedAt:      postedAt.UTC(),
			})
		}
	}
	return out, parsed.Payload.NextToken, nil
}

func (c *SPAPIClient) listOrderItems(ctx context.Context, token string, orderID string) ([]integrationDomain.ExternalOrderItem, error) {
	if strings.TrimSpace(orderID) == "" {
		return nil, fmt.Errorf("orderID is required")
	}

	items := make([]integrationDomain.ExternalOrderItem, 0)
	nextToken := ""
	escapedOrderID := url.PathEscape(orderID)

	for {
		query := url.Values{}
		if nextToken != "" {
			query.Set("NextToken", nextToken)
		}

		path := fmt.Sprintf("/orders/v0/orders/%s/orderItems", escapedOrderID)
		respBody, err := c.doSPAPIRequest(ctx, http.MethodGet, path, query, nil, token)
		if err != nil {
			return nil, err
		}

		var parsed struct {
			Payload struct {
				OrderItems []struct {
					OrderItemID     string `json:"OrderItemId"`
					SellerSKU       string `json:"SellerSKU"`
					QuantityOrdered uint64 `json:"QuantityOrdered"`
					ItemPrice       struct {
						Amount string `json:"Amount"`
					} `json:"ItemPrice"`
				} `json:"OrderItems"`
				NextToken string `json:"NextToken"`
			} `json:"payload"`
		}
		if err := json.Unmarshal(respBody, &parsed); err != nil {
			return nil, fmt.Errorf("decode order items response failed: %w", err)
		}

		for _, item := range parsed.Payload.OrderItems {
			amount, err := strconv.ParseFloat(strings.TrimSpace(item.ItemPrice.Amount), 64)
			if err != nil {
				amount = 0
			}
			items = append(items, integrationDomain.ExternalOrderItem{
				OrderItemID: item.OrderItemID,
				SellerSKU:   item.SellerSKU,
				Quantity:    item.QuantityOrdered,
				Amount:      amount,
			})
		}

		if strings.TrimSpace(parsed.Payload.NextToken) == "" {
			break
		}
		nextToken = parsed.Payload.NextToken
	}

	return items, nil
}

func (c *SPAPIClient) doSPAPIRequest(
	ctx context.Context,
	method string,
	path string,
	query url.Values,
	body []byte,
	accessToken string,
) ([]byte, error) {
	req, err := c.newRequest(ctx, method, path, query, body, accessToken)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("sp-api http %d: %s", resp.StatusCode, trimBody(respBody))
	}

	var errResp struct {
		Errors []struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"errors"`
	}
	_ = json.Unmarshal(respBody, &errResp)
	if len(errResp.Errors) > 0 {
		return nil, fmt.Errorf("sp-api %s: %s", errResp.Errors[0].Code, errResp.Errors[0].Message)
	}

	return respBody, nil
}

func (c *SPAPIClient) getAccessToken(ctx context.Context) (string, error) {
	c.tokenMu.Lock()
	if c.accessToken != "" && time.Now().Before(c.accessTokenExp.Add(-60*time.Second)) {
		token := c.accessToken
		c.tokenMu.Unlock()
		return token, nil
	}
	c.tokenMu.Unlock()

	tokenResp, err := c.RefreshAccessToken(ctx, c.cfg.LWARefreshToken)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(tokenResp.AccessToken) == "" {
		return "", fmt.Errorf("empty access_token from refresh flow")
	}
	if tokenResp.ExpiresIn <= 0 {
		tokenResp.ExpiresIn = 3600
	}

	c.tokenMu.Lock()
	c.accessToken = tokenResp.AccessToken
	c.accessTokenExp = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	c.tokenMu.Unlock()
	return tokenResp.AccessToken, nil
}

func (c *SPAPIClient) newRequest(
	ctx context.Context,
	method string,
	path string,
	query url.Values,
	body []byte,
	accessToken string,
) (*http.Request, error) {
	endpoint := strings.TrimRight(c.cfg.Endpoint, "/")
	fullURL := endpoint + path
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = strings.NewReader(string(body))
	}
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-amz-access-token", accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func trimBody(body []byte) string {
	raw := strings.TrimSpace(string(body))
	if len(raw) <= 300 {
		return raw
	}
	return raw[:300]
}

func parseTimeRFC3339(raw string) time.Time {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return t
	}
	if t, err := time.Parse(time.RFC3339Nano, raw); err == nil {
		return t
	}
	return time.Time{}
}

func parsePositiveUint(v any) uint64 {
	switch typed := v.(type) {
	case float64:
		if typed <= 0 {
			return 0
		}
		return uint64(typed)
	case float32:
		if typed <= 0 {
			return 0
		}
		return uint64(typed)
	case int:
		if typed <= 0 {
			return 0
		}
		return uint64(typed)
	case int64:
		if typed <= 0 {
			return 0
		}
		return uint64(typed)
	case uint64:
		return typed
	case string:
		trimmed := strings.TrimSpace(typed)
		if trimmed == "" {
			return 0
		}
		if parsed, err := strconv.ParseFloat(trimmed, 64); err == nil && parsed > 0 {
			return uint64(parsed)
		}
	}
	return 0
}

func (c *SPAPIClient) fetchToken(ctx context.Context, form url.Values) (*integrationDomain.OAuthToken, error) {
	if strings.TrimSpace(form.Get("client_id")) == "" || strings.TrimSpace(form.Get("client_secret")) == "" {
		return nil, fmt.Errorf("amazon lwa client credentials required")
	}
	if strings.TrimSpace(form.Get("grant_type")) == "" {
		return nil, fmt.Errorf("grant_type is required")
	}
	if form.Get("grant_type") == "refresh_token" && strings.TrimSpace(form.Get("refresh_token")) == "" {
		return nil, fmt.Errorf("refresh_token is required")
	}
	if form.Get("grant_type") == "authorization_code" && strings.TrimSpace(form.Get("code")) == "" {
		return nil, fmt.Errorf("authorization code is required")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.amazon.com/auth/o2/token",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("lwa token http %d: %s", resp.StatusCode, trimBody(respBody))
	}

	var parsed struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
		Scope        string `json:"scope"`
	}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, err
	}
	if strings.TrimSpace(parsed.AccessToken) == "" {
		return nil, errors.New("empty access_token from lwa")
	}
	if parsed.ExpiresIn <= 0 {
		parsed.ExpiresIn = 3600
	}
	return &integrationDomain.OAuthToken{
		AccessToken:  strings.TrimSpace(parsed.AccessToken),
		RefreshToken: strings.TrimSpace(parsed.RefreshToken),
		ExpiresIn:    parsed.ExpiresIn,
		Scope:        strings.TrimSpace(parsed.Scope),
	}, nil
}
