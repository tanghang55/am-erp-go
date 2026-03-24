package repository

import (
	"errors"
	"strings"
	"time"

	integrationDomain "am-erp-go/internal/module/integration/domain"

	"gorm.io/gorm"
)

type authorizationRepository struct {
	db *gorm.DB
}

func NewAuthorizationRepository(db *gorm.DB) *authorizationRepository {
	return &authorizationRepository{db: db}
}

func (r *authorizationRepository) CreateAuthorization(record *integrationDomain.IntegrationAuthorization) error {
	if record == nil {
		return nil
	}
	return r.db.Create(record).Error
}

func (r *authorizationRepository) UpdateAuthorization(record *integrationDomain.IntegrationAuthorization) error {
	if record == nil {
		return nil
	}
	return r.db.Save(record).Error
}

func (r *authorizationRepository) GetAuthorizationByID(id uint64) (*integrationDomain.IntegrationAuthorization, error) {
	var record integrationDomain.IntegrationAuthorization
	if err := r.db.First(&record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (r *authorizationRepository) GetAuthorizationByProviderAndState(providerCode string, oauthState string) (*integrationDomain.IntegrationAuthorization, error) {
	var record integrationDomain.IntegrationAuthorization
	if err := r.db.
		Where("provider_code = ? AND oauth_state = ?", strings.TrimSpace(providerCode), strings.TrimSpace(oauthState)).
		Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (r *authorizationRepository) ListAuthorizations(params *integrationDomain.ListAuthorizationParams) ([]integrationDomain.IntegrationAuthorization, int64, error) {
	var list []integrationDomain.IntegrationAuthorization
	var total int64

	page := 1
	pageSize := 20
	if params != nil {
		if params.Page > 0 {
			page = params.Page
		}
		if params.PageSize > 0 {
			pageSize = params.PageSize
		}
	}

	query := r.db.Model(&integrationDomain.IntegrationAuthorization{})
	if params != nil {
		if code := strings.TrimSpace(params.ProviderCode); code != "" {
			query = query.Where("provider_code = ?", strings.ToUpper(code))
		}
		if status := strings.TrimSpace(params.Status); status != "" {
			query = query.Where("status = ?", strings.ToUpper(status))
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []integrationDomain.IntegrationAuthorization{}, 0, nil
	}

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *authorizationRepository) ListNeedRefresh(before time.Time, limit int) ([]integrationDomain.IntegrationAuthorization, error) {
	if limit <= 0 {
		limit = 200
	}
	var list []integrationDomain.IntegrationAuthorization
	err := r.db.
		Model(&integrationDomain.IntegrationAuthorization{}).
		Where("status = ? AND refresh_token IS NOT NULL AND refresh_token <> ''", integrationDomain.AuthorizationStatusAuthorized).
		Where("access_token_expire_at IS NULL OR access_token_expire_at <= ?", before).
		Order("COALESCE(access_token_expire_at, '1970-01-01') ASC, id ASC").
		Limit(limit).
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return []integrationDomain.IntegrationAuthorization{}, nil
	}
	return list, nil
}
