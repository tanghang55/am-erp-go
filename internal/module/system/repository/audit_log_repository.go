package repository

import (
	systemdomain "am-erp-go/internal/module/system/domain"

	"gorm.io/gorm"
)

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) systemdomain.AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) List(params systemdomain.AuditLogListParams) ([]*systemdomain.AuditLog, int64, error) {
	var logs []*systemdomain.AuditLog
	var total int64

	query := r.db.Model(&systemdomain.AuditLog{})

	if params.Module != "" {
		query = query.Where("module = ?", params.Module)
	}
	if params.Action != "" {
		query = query.Where("action = ?", params.Action)
	}
	if params.Username != "" {
		query = query.Where("username = ?", params.Username)
	}
	if params.EntityType != "" {
		query = query.Where("entity_type = ?", params.EntityType)
	}
	if params.EntityID != "" {
		query = query.Where("entity_id = ?", params.EntityID)
	}
	if params.Keyword != "" {
		like := "%" + params.Keyword + "%"
		query = query.Where("trace_id LIKE ? OR username LIKE ? OR entity_type LIKE ? OR entity_id LIKE ?", like, like, like, like)
	}
	if params.DateFrom != "" {
		query = query.Where("gmt_create >= ?", params.DateFrom+" 00:00:00")
	}
	if params.DateTo != "" {
		query = query.Where("gmt_create <= ?", params.DateTo+" 23:59:59")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := params.Page
	pageSize := params.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	if err := query.Order("gmt_create DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *auditLogRepository) Create(log *systemdomain.AuditLog) error {
	return r.db.Create(log).Error
}
