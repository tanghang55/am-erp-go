package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	integrationDomain "am-erp-go/internal/module/integration/domain"
)

func TestStartAuthorizationCreatesPendingRecord(t *testing.T) {
	now := time.Date(2026, 3, 4, 12, 0, 0, 0, time.UTC)
	repo := &fakeAuthorizationRepo{}
	provider := &fakeAuthorizationProvider{code: "AMAZON_US", typ: "amazon"}
	uc := NewAuthorizationUsecase(repo, []integrationDomain.AuthorizationProvider{provider}, func() time.Time { return now })

	result, err := uc.StartAuthorization(&StartAuthorizationInput{
		ProviderCode: "AMAZON_US",
		AccountAlias: "主店铺",
	})
	if err != nil {
		t.Fatalf("StartAuthorization error: %v", err)
	}
	if result.AuthorizationID == 0 {
		t.Fatalf("expected authorization id > 0")
	}
	if repo.lastCreated == nil || repo.lastCreated.Status != integrationDomain.AuthorizationStatusPending {
		t.Fatalf("expected pending record created")
	}
	if provider.lastState == "" {
		t.Fatalf("expected oauth state generated")
	}
}

func TestHandleCallbackAndRefreshSuccess(t *testing.T) {
	now := time.Date(2026, 3, 4, 12, 0, 0, 0, time.UTC)
	state := "abc123state"
	repo := &fakeAuthorizationRepo{
		byState: &integrationDomain.IntegrationAuthorization{
			ID:                 11,
			ProviderCode:       "AMAZON_US",
			ProviderType:       "amazon",
			Status:             integrationDomain.AuthorizationStatusPending,
			OAuthState:         &state,
			OAuthStateExpireAt: ptrTime(now.Add(10 * time.Minute)),
		},
	}
	provider := &fakeAuthorizationProvider{
		code: "AMAZON_US",
		typ:  "amazon",
		exchangeToken: &integrationDomain.OAuthToken{
			AccessToken:  "token-1",
			RefreshToken: "refresh-1",
			ExpiresIn:    3600,
			Scope:        "sellingpartnerapi::orders",
		},
		refreshTokenResp: &integrationDomain.OAuthToken{
			AccessToken:  "token-2",
			RefreshToken: "refresh-2",
			ExpiresIn:    3600,
		},
	}
	uc := NewAuthorizationUsecase(repo, []integrationDomain.AuthorizationProvider{provider}, func() time.Time { return now })

	record, err := uc.HandleCallback(context.Background(), &OAuthCallbackInput{
		ProviderCode:      "AMAZON_US",
		OAuthState:        state,
		AuthorizationCode: "code-123",
		SellerPartnerID:   "A1SELLER",
	})
	if err != nil {
		t.Fatalf("HandleCallback error: %v", err)
	}
	if record.Status != integrationDomain.AuthorizationStatusAuthorized {
		t.Fatalf("expected AUTHORIZED, got %s", record.Status)
	}
	if record.RefreshToken == nil || *record.RefreshToken != "refresh-1" {
		t.Fatalf("expected refresh token saved")
	}

	refreshed, err := uc.ManualRefresh(context.Background(), record.ID, nil)
	if err != nil {
		t.Fatalf("ManualRefresh error: %v", err)
	}
	if refreshed.AccessToken == nil || *refreshed.AccessToken != "token-2" {
		t.Fatalf("expected refreshed access token")
	}
	if refreshed.RefreshToken == nil || *refreshed.RefreshToken != "refresh-2" {
		t.Fatalf("expected refreshed refresh token")
	}
}

func TestManualRefreshRetriesAndMarksFailedAfterThreshold(t *testing.T) {
	now := time.Date(2026, 3, 20, 9, 0, 0, 0, time.UTC)
	refreshToken := "refresh-1"
	repo := &fakeAuthorizationRepo{
		byID: &integrationDomain.IntegrationAuthorization{
			ID:           101,
			ProviderCode: "AMAZON_US",
			ProviderType: "amazon",
			Status:       integrationDomain.AuthorizationStatusAuthorized,
			RefreshToken: &refreshToken,
		},
	}
	provider := &fakeAuthorizationProvider{
		code:       "AMAZON_US",
		typ:        "amazon",
		refreshErr: errors.New("temporary network error"),
	}
	uc := NewAuthorizationUsecase(repo, []integrationDomain.AuthorizationProvider{provider}, func() time.Time { return now })

	for i := 0; i < 3; i++ {
		if _, err := uc.ManualRefresh(context.Background(), 101, nil); err == nil {
			t.Fatalf("expected refresh error on attempt %d", i+1)
		}
	}

	if repo.byID == nil {
		t.Fatalf("expected repo record to be updated")
	}
	if repo.byID.RefreshFailCount != 3 {
		t.Fatalf("expected refresh_fail_count=3, got %d", repo.byID.RefreshFailCount)
	}
	if repo.byID.Status != integrationDomain.AuthorizationStatusFailed {
		t.Fatalf("expected status FAILED after threshold, got %s", repo.byID.Status)
	}
	if repo.byID.LastRefreshFailedAt == nil {
		t.Fatalf("expected last_refresh_failed_at to be set")
	}
	if repo.byID.LastErrorMessage == nil || *repo.byID.LastErrorMessage == "" {
		t.Fatalf("expected last_error_message to be set")
	}
	if provider.refreshCalls != 9 {
		t.Fatalf("expected 9 refresh calls (3 retries x 3 rounds), got %d", provider.refreshCalls)
	}
}

func TestManualRefreshSuccessResetsFailureAndRecoversStatus(t *testing.T) {
	now := time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC)
	refreshToken := "refresh-old"
	lastErr := "token expired"
	repo := &fakeAuthorizationRepo{
		byID: &integrationDomain.IntegrationAuthorization{
			ID:               201,
			ProviderCode:     "AMAZON_US",
			ProviderType:     "amazon",
			Status:           integrationDomain.AuthorizationStatusFailed,
			RefreshToken:     &refreshToken,
			RefreshFailCount: 2,
			LastErrorMessage: &lastErr,
		},
	}
	provider := &fakeAuthorizationProvider{
		code: "AMAZON_US",
		typ:  "amazon",
		refreshTokenResp: &integrationDomain.OAuthToken{
			AccessToken:  "token-new",
			RefreshToken: "refresh-new",
			ExpiresIn:    3600,
		},
	}
	uc := NewAuthorizationUsecase(repo, []integrationDomain.AuthorizationProvider{provider}, func() time.Time { return now })

	record, err := uc.ManualRefresh(context.Background(), 201, nil)
	if err != nil {
		t.Fatalf("ManualRefresh error: %v", err)
	}
	if record.Status != integrationDomain.AuthorizationStatusAuthorized {
		t.Fatalf("expected status AUTHORIZED, got %s", record.Status)
	}
	if record.RefreshFailCount != 0 {
		t.Fatalf("expected refresh_fail_count reset to 0, got %d", record.RefreshFailCount)
	}
	if record.LastErrorMessage != nil {
		t.Fatalf("expected last_error_message cleared, got %v", *record.LastErrorMessage)
	}
	if record.LastRefreshFailedAt != nil {
		t.Fatalf("expected last_refresh_failed_at cleared")
	}
}

type fakeAuthorizationRepo struct {
	lastCreated *integrationDomain.IntegrationAuthorization
	byID        *integrationDomain.IntegrationAuthorization
	byState     *integrationDomain.IntegrationAuthorization
}

func (f *fakeAuthorizationRepo) CreateAuthorization(record *integrationDomain.IntegrationAuthorization) error {
	record.ID = 1
	f.lastCreated = record
	f.byID = record
	f.byState = record
	return nil
}

func (f *fakeAuthorizationRepo) UpdateAuthorization(record *integrationDomain.IntegrationAuthorization) error {
	f.byID = record
	f.byState = record
	return nil
}

func (f *fakeAuthorizationRepo) GetAuthorizationByID(id uint64) (*integrationDomain.IntegrationAuthorization, error) {
	if f.byID != nil && f.byID.ID == id {
		return f.byID, nil
	}
	return nil, nil
}

func (f *fakeAuthorizationRepo) GetAuthorizationByProviderAndState(providerCode string, oauthState string) (*integrationDomain.IntegrationAuthorization, error) {
	if f.byState == nil || f.byState.OAuthState == nil {
		return nil, nil
	}
	if f.byState.ProviderCode == providerCode && *f.byState.OAuthState == oauthState {
		return f.byState, nil
	}
	return nil, nil
}

func (f *fakeAuthorizationRepo) ListAuthorizations(params *integrationDomain.ListAuthorizationParams) ([]integrationDomain.IntegrationAuthorization, int64, error) {
	_ = params
	return []integrationDomain.IntegrationAuthorization{}, 0, nil
}

func (f *fakeAuthorizationRepo) ListNeedRefresh(before time.Time, limit int) ([]integrationDomain.IntegrationAuthorization, error) {
	_ = before
	_ = limit
	return []integrationDomain.IntegrationAuthorization{}, nil
}

type fakeAuthorizationProvider struct {
	code             string
	typ              string
	lastState        string
	exchangeToken    *integrationDomain.OAuthToken
	refreshTokenResp *integrationDomain.OAuthToken
	refreshErr       error
	refreshCalls     int
}

func (f *fakeAuthorizationProvider) Code() string {
	return f.code
}

func (f *fakeAuthorizationProvider) Type() string {
	return f.typ
}

func (f *fakeAuthorizationProvider) BuildAuthorizeURL(state string) (string, error) {
	f.lastState = state
	return "https://example.com/oauth?state=" + state, nil
}

func (f *fakeAuthorizationProvider) ExchangeAuthorizationCode(ctx context.Context, code string) (*integrationDomain.OAuthToken, error) {
	_ = ctx
	if code == "" {
		return nil, fmt.Errorf("code required")
	}
	return f.exchangeToken, nil
}

func (f *fakeAuthorizationProvider) RefreshAccessToken(ctx context.Context, refreshToken string) (*integrationDomain.OAuthToken, error) {
	_ = ctx
	f.refreshCalls++
	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token required")
	}
	if f.refreshErr != nil {
		return nil, f.refreshErr
	}
	return f.refreshTokenResp, nil
}

func ptrTime(v time.Time) *time.Time {
	return &v
}
