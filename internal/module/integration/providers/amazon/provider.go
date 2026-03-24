package amazon

import (
	"context"
	"strings"

	integrationDomain "am-erp-go/internal/module/integration/domain"
)

type OrdersClient interface {
	ListOrders(ctx context.Context, req integrationDomain.ExternalListOrdersRequest) ([]integrationDomain.ExternalOrder, error)
	ListRefunds(ctx context.Context, req integrationDomain.ExternalListRefundsRequest) ([]integrationDomain.ExternalRefund, error)
	BuildAuthorizeURL(state string) (string, error)
	ExchangeAuthorizationCode(ctx context.Context, code string) (*integrationDomain.OAuthToken, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*integrationDomain.OAuthToken, error)
}

type OrdersProvider struct {
	code   string
	client OrdersClient
}

func NewOrdersProvider(code string, client OrdersClient) *OrdersProvider {
	if strings.TrimSpace(code) == "" {
		code = "AMAZON"
	}
	return &OrdersProvider{
		code:   strings.ToUpper(strings.TrimSpace(code)),
		client: client,
	}
}

func (p *OrdersProvider) Code() string {
	return p.code
}

func (p *OrdersProvider) Type() string {
	return "amazon"
}

func (p *OrdersProvider) ListOrders(ctx context.Context, req integrationDomain.ExternalListOrdersRequest) ([]integrationDomain.ExternalOrder, error) {
	return p.client.ListOrders(ctx, req)
}

func (p *OrdersProvider) ListRefunds(ctx context.Context, req integrationDomain.ExternalListRefundsRequest) ([]integrationDomain.ExternalRefund, error) {
	return p.client.ListRefunds(ctx, req)
}

func (p *OrdersProvider) BuildAuthorizeURL(state string) (string, error) {
	return p.client.BuildAuthorizeURL(state)
}

func (p *OrdersProvider) ExchangeAuthorizationCode(ctx context.Context, code string) (*integrationDomain.OAuthToken, error) {
	return p.client.ExchangeAuthorizationCode(ctx, code)
}

func (p *OrdersProvider) RefreshAccessToken(ctx context.Context, refreshToken string) (*integrationDomain.OAuthToken, error) {
	return p.client.RefreshAccessToken(ctx, refreshToken)
}
