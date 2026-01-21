package repository

import (
	systemdomain "am-erp-go/internal/module/system/domain"

	"gorm.io/gorm"
)

type fieldLabelRepository struct {
	db *gorm.DB
}

func NewFieldLabelRepository(db *gorm.DB) systemdomain.FieldLabelRepository {
	return &fieldLabelRepository{db: db}
}

func (r *fieldLabelRepository) GetAll() ([]*systemdomain.FieldLabel, error) {
	var labels []*systemdomain.FieldLabel
	err := r.db.Order("id ASC").Find(&labels).Error
	return labels, err
}

func (r *fieldLabelRepository) List(page, pageSize int, keyword string) ([]*systemdomain.FieldLabel, int64, error) {
	var labels []*systemdomain.FieldLabel
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	query := r.db.Model(&systemdomain.FieldLabel{})
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("label_key LIKE ?", like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&labels).Error; err != nil {
		return nil, 0, err
	}

	return labels, total, nil
}

func (r *fieldLabelRepository) GetByID(id uint64) (*systemdomain.FieldLabel, error) {
	var label systemdomain.FieldLabel
	if err := r.db.First(&label, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, systemdomain.ErrFieldLabelNotFound
		}
		return nil, err
	}
	return &label, nil
}

func (r *fieldLabelRepository) GetByKey(key string) (*systemdomain.FieldLabel, error) {
	var label systemdomain.FieldLabel
	if err := r.db.Where("label_key = ?", key).First(&label).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, systemdomain.ErrFieldLabelNotFound
		}
		return nil, err
	}
	return &label, nil
}

func (r *fieldLabelRepository) Create(label *systemdomain.FieldLabel) error {
	return r.db.Create(label).Error
}

func (r *fieldLabelRepository) Update(label *systemdomain.FieldLabel) error {
	return r.db.Save(label).Error
}

func (r *fieldLabelRepository) Delete(id uint64) error {
	return r.db.Delete(&systemdomain.FieldLabel{}, id).Error
}
