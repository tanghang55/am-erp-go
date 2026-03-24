package domain

import (
	"context"
	"time"
)

type AuthorizationStatus string

const (
	AuthorizationStatusPending    AuthorizationStatus = "PENDING"
	AuthorizationStatusAuthorized AuthorizationStatus = "AUTHORIZED"
	AuthorizationStatusFailed     AuthorizationStatus = "FAILED"
	AuthorizationStatusDisabled   AuthorizationStatus = "DISABLED"
)

type IntegrationAuthorization struct {
	ID                   uint64              `json:"id" gorm:"primaryKey;autoIncrement"`
	ProviderCode         string              `json:"provider_code" gorm:"column:provider_code;size:64;not null"`
	ProviderType         string              `json:"provider_type" gorm:"column:provider_type;size:32;not null"`
	AccountAlias         *string             `json:"account_alias" gorm:"column:account_alias;size:100"`
	SellerPartnerID      *string             `json:"seller_partner_id" gorm:"column:seller_partner_id;size:64"`
	Status               AuthorizationStatus `json:"status" gorm:"column:status;size:20;not null"`
	OAuthState           *string             `json:"oauth_state,omitempty" gorm:"column:oauth_state;size:128"`
	OAuthStateExpireAt   *time.Time          `json:"oauth_state_expire_at,omitempty" gorm:"column:oauth_state_expire_at"`
	RefreshToken         *string             `json:"-" gorm:"column:refresh_token;size:2048"`
	AccessToken          *string             `json:"-" gorm:"column:access_token;size:4096"`
	AccessTokenExpireAt  *time.Time          `json:"access_token_expire_at" gorm:"column:access_token_expire_at"`
	TokenScope           *string             `json:"token_scope" gorm:"column:token_scope;size:512"`
	LastAuthorizedAt     *time.Time          `json:"last_authorized_at" gorm:"column:last_authorized_at"`
	LastRefreshAt        *time.Time          `json:"last_refresh_at" gorm:"column:last_refresh_at"`
	RefreshFailCount     uint8               `json:"refresh_fail_count" gorm:"column:refresh_fail_count;not null;default:0"`
	LastRefreshAttemptAt *time.Time          `json:"last_refresh_attempt_at" gorm:"column:last_refresh_attempt_at"`
	LastRefreshFailedAt  *time.Time          `json:"last_refresh_failed_at" gorm:"column:last_refresh_failed_at"`
	LastErrorMessage     *string             `json:"last_error_message" gorm:"column:last_error_message;size:500"`
	CreatedBy            *uint64             `json:"created_by" gorm:"column:created_by"`
	UpdatedBy            *uint64             `json:"updated_by" gorm:"column:updated_by"`
	GmtCreate            time.Time           `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified          time.Time           `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (IntegrationAuthorization) TableName() string {
	return "integration_authorization"
}

type ListAuthorizationParams struct {
	Page         int
	PageSize     int
	ProviderCode string
	Status       string
}

type OAuthToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	Scope        string
}

type AuthorizationProvider interface {
	Code() string
	Type() string
	BuildAuthorizeURL(state string) (string, error)
	ExchangeAuthorizationCode(ctx context.Context, code string) (*OAuthToken, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*OAuthToken, error)
}

type AuthorizationRepository interface {
	CreateAuthorization(record *IntegrationAuthorization) error
	UpdateAuthorization(record *IntegrationAuthorization) error
	GetAuthorizationByID(id uint64) (*IntegrationAuthorization, error)
	GetAuthorizationByProviderAndState(providerCode string, oauthState string) (*IntegrationAuthorization, error)
	ListAuthorizations(params *ListAuthorizationParams) ([]IntegrationAuthorization, int64, error)
	ListNeedRefresh(before time.Time, limit int) ([]IntegrationAuthorization, error)
}
