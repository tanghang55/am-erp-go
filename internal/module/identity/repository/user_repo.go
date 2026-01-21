package repository

import (
	"am-erp-go/internal/module/identity/domain"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ? AND status = 'ACTIVE'", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint64) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ? AND status = 'ACTIVE'", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserRoles(userID uint64) ([]domain.Role, error) {
	var roles []domain.Role
	err := r.db.Table("role").
		Joins("JOIN user_role ON user_role.role_id = role.id").
		Where("user_role.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

func (r *userRepository) GetUserPermissions(userID uint64) ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := r.db.Table("permission").
		Joins("JOIN role_permission ON role_permission.permission_id = permission.id").
		Joins("JOIN user_role ON user_role.role_id = role_permission.role_id").
		Where("user_role.user_id = ? AND permission.status = 'ACTIVE'", userID).
		Distinct().
		Find(&permissions).Error
	return permissions, err
}
