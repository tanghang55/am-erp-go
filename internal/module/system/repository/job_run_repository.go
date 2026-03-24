package repository

import (
	systemdomain "am-erp-go/internal/module/system/domain"
	"time"

	"gorm.io/gorm"
)

type jobRunRepository struct {
	db *gorm.DB
}

func NewJobRunRepository(db *gorm.DB) systemdomain.JobRunRepository {
	return &jobRunRepository{db: db}
}

func (r *jobRunRepository) Create(run *systemdomain.JobRun) error {
	return r.db.Create(run).Error
}

func (r *jobRunRepository) Update(run *systemdomain.JobRun) error {
	return r.db.Save(run).Error
}

func (r *jobRunRepository) DeleteOlderThan(cutoff time.Time) (int64, error) {
	result := r.db.Where("gmt_create < ?", cutoff).Delete(&systemdomain.JobRun{})
	return result.RowsAffected, result.Error
}
