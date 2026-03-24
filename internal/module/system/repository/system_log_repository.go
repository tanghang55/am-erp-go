package repository

import (
	systemdomain "am-erp-go/internal/module/system/domain"
	"time"

	"gorm.io/gorm"
)

type systemLogRepository struct {
	db *gorm.DB
}

func NewSystemLogRepository(db *gorm.DB) systemdomain.SystemLogRepository {
	return &systemLogRepository{db: db}
}

func (r *systemLogRepository) Create(log *systemdomain.SystemLog) error {
	return r.db.Create(log).Error
}

func (r *systemLogRepository) DeleteOlderThan(cutoff time.Time) (int64, error) {
	result := r.db.Where("gmt_create < ?", cutoff).Delete(&systemdomain.SystemLog{})
	return result.RowsAffected, result.Error
}
