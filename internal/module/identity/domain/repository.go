package domain

type UserRepository interface {
	FindByUsername(username string) (*User, error)
	FindByID(id uint64) (*User, error)
	GetUserRoles(userID uint64) ([]Role, error)
	GetUserPermissions(userID uint64) ([]Permission, error)
}
