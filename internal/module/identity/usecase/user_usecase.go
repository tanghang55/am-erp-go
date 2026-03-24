package usecase

import (
	"fmt"
	"strings"

	"am-erp-go/internal/module/identity/domain"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

type CreateUserInput struct {
	Username string
	Password string
	RealName string
	Email    string
	Phone    string
	Status   string
}

type UpdateUserInput struct {
	Password *string
	RealName *string
	Email    *string
	Phone    *string
	Status   *string
}

func (u *UserUsecase) ListUsers(params *domain.UserListParams) ([]domain.User, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	return u.userRepo.ListUsers(params)
}

func (u *UserUsecase) GetUserDetail(userID uint64) (*domain.User, []domain.Role, []domain.Permission, error) {
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, nil, nil, err
	}
	roles, err := u.userRepo.GetUserRoles(userID)
	if err != nil {
		return nil, nil, nil, err
	}
	permissions, err := u.userRepo.GetUserPermissions(userID)
	if err != nil {
		return nil, nil, nil, err
	}
	return user, roles, permissions, nil
}

func (u *UserUsecase) CreateUser(input *CreateUserInput) (*domain.User, error) {
	username := strings.TrimSpace(input.Username)
	if len(username) < 4 {
		return nil, fmt.Errorf("username length must be >= 4")
	}
	password := strings.TrimSpace(input.Password)
	if len(password) < 8 {
		return nil, fmt.Errorf("password length must be >= 8")
	}

	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = "ACTIVE"
	}
	if status != "ACTIVE" && status != "DISABLED" {
		return nil, fmt.Errorf("invalid status")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:     username,
		PasswordHash: string(passwordHash),
		RealName:     strings.TrimSpace(input.RealName),
		Email:        strings.TrimSpace(input.Email),
		Phone:        strings.TrimSpace(input.Phone),
		Status:       status,
	}

	if err := u.userRepo.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) UpdateUser(userID uint64, input *UpdateUserInput) (*domain.User, error) {
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if input.Password != nil {
		password := strings.TrimSpace(*input.Password)
		if password != "" {
			if len(password) < 8 {
				return nil, fmt.Errorf("password length must be >= 8")
			}
			passwordHash, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if hashErr != nil {
				return nil, hashErr
			}
			user.PasswordHash = string(passwordHash)
		}
	}
	if input.RealName != nil {
		user.RealName = strings.TrimSpace(*input.RealName)
	}
	if input.Email != nil {
		user.Email = strings.TrimSpace(*input.Email)
	}
	if input.Phone != nil {
		user.Phone = strings.TrimSpace(*input.Phone)
	}
	if input.Status != nil {
		status := strings.TrimSpace(*input.Status)
		if status != "ACTIVE" && status != "DISABLED" {
			return nil, fmt.Errorf("invalid status")
		}
		user.Status = status
	}

	if err := u.userRepo.UpdateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) DisableUser(userID uint64) error {
	return u.userRepo.DisableUser(userID)
}

func (u *UserUsecase) AssignRoles(userID uint64, roleIDs []uint64) error {
	if _, err := u.userRepo.GetUserByID(userID); err != nil {
		return err
	}
	return u.userRepo.ReplaceUserRoles(userID, roleIDs)
}

func (u *UserUsecase) ListRoles() ([]domain.Role, error) {
	return u.userRepo.ListRoles()
}

func (u *UserUsecase) ListPermissions() ([]domain.Permission, error) {
	return u.userRepo.ListPermissions()
}
