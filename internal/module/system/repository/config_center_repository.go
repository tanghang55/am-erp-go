package repository

import (
	systemdomain "am-erp-go/internal/module/system/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type configCenterRepository struct {
	db *gorm.DB
}

func NewConfigCenterRepository(db *gorm.DB) systemdomain.ConfigCenterRepository {
	return &configCenterRepository{db: db}
}

func (r *configCenterRepository) ListDefinitions(moduleCode string) ([]*systemdomain.ConfigDefinition, error) {
	items := make([]*systemdomain.ConfigDefinition, 0)
	query := r.db.Model(&systemdomain.ConfigDefinition{}).Where("is_active = 1")
	if moduleCode != "" {
		query = query.Where("module_code = ?", moduleCode)
	}
	if err := query.Order("module_code ASC, group_code ASC, sort ASC, id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *configCenterRepository) ListValues(scopeType string, scopeRefID uint64, keys []string) ([]*systemdomain.ConfigValue, error) {
	items := make([]*systemdomain.ConfigValue, 0)
	query := r.db.Model(&systemdomain.ConfigValue{}).Where("scope_type = ? AND scope_ref_id = ?", scopeType, scopeRefID)
	if len(keys) > 0 {
		query = query.Where("config_key IN ?", keys)
	}
	if err := query.Order("config_key ASC, id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *configCenterRepository) UpsertValues(items []*systemdomain.ConfigValue) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "config_key"},
			{Name: "scope_type"},
			{Name: "scope_ref_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"config_value", "updated_by", "gmt_modified"}),
	}).Create(&items).Error
}
