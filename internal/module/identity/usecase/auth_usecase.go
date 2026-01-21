package usecase

import (
	"errors"

	"am-erp-go/internal/infrastructure/auth"
	"am-erp-go/internal/module/identity/domain"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserDisabled       = errors.New("user account is disabled")
)

type AuthUsecase struct {
	userRepo   domain.UserRepository
	jwtManager *auth.JWTManager
}

func NewAuthUsecase(
	userRepo domain.UserRepository,
	jwtManager *auth.JWTManager,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

type LoginResponse struct {
	User        *domain.User        `json:"user"`
	Roles       []domain.Role       `json:"roles"`
	Permissions []domain.Permission `json:"permissions"`
	AccessToken string              `json:"access_token"`
}

func (uc *AuthUsecase) Login(username, password string) (*LoginResponse, error) {
	user, err := uc.userRepo.FindByUsername(username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive() {
		return nil, ErrUserDisabled
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	roles, err := uc.userRepo.GetUserRoles(user.ID)
	if err != nil {
		return nil, err
	}

	permissions, err := uc.userRepo.GetUserPermissions(user.ID)
	if err != nil {
		return nil, err
	}

	token, err := uc.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:        user,
		Roles:       roles,
		Permissions: permissions,
		AccessToken: token,
	}, nil
}

func (uc *AuthUsecase) GetCurrentUser(userID uint64) (*domain.User, []domain.Role, []domain.Permission, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, nil, nil, err
	}

	roles, err := uc.userRepo.GetUserRoles(userID)
	if err != nil {
		return nil, nil, nil, err
	}

	permissions, err := uc.userRepo.GetUserPermissions(userID)
	if err != nil {
		return nil, nil, nil, err
	}

	return user, roles, permissions, nil
}
