package repository

import (
	"am-erp-go/internal/module/identity/domain"
	"time"

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

func (r *userRepository) ListUsers(params *domain.UserListParams) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	query := r.db.Model(&domain.User{})
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(
			"username LIKE ? OR real_name LIKE ? OR email LIKE ? OR phone LIKE ?",
			keyword, keyword, keyword, keyword,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Order("id DESC").Offset(offset).Limit(params.PageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) GetUserByID(id uint64) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) UpdateUser(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) DisableUser(id uint64) error {
	return r.db.Model(&domain.User{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":       "DISABLED",
			"gmt_modified": time.Now(),
		}).Error
}

func (r *userRepository) ReplaceUserRoles(userID uint64, roleIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&domain.UserRole{}).Error; err != nil {
			return err
		}
		if len(roleIDs) == 0 {
			return nil
		}

		now := time.Now()
		rows := make([]domain.UserRole, 0, len(roleIDs))
		for _, roleID := range roleIDs {
			rows = append(rows, domain.UserRole{
				UserID:      userID,
				RoleID:      roleID,
				GmtCreate:   now,
				GmtModified: now,
			})
		}
		return tx.Create(&rows).Error
	})
}

func (r *userRepository) ListRoles() ([]domain.Role, error) {
	var roles []domain.Role
	if err := r.db.Order("id ASC").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *userRepository) ListPermissions() ([]domain.Permission, error) {
	var permissions []domain.Permission
	if err := r.db.Order("id ASC").Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}
