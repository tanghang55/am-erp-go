package repository

import (
	"am-erp-go/internal/module/procurement/domain"

	"gorm.io/gorm"
)

type movementRepository struct {
	db *gorm.DB
}

func NewMovementRepository(db *gorm.DB) *movementRepository {
	return &movementRepository{db: db}
}

func (r *movementRepository) Create(movement *domain.InventoryMovement) error {
	if movement == nil {
		return nil
	}
	return r.db.Create(movement).Error
}
