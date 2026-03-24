package domain

type UserRepository interface {
	FindByUsername(username string) (*User, error)
	FindByID(id uint64) (*User, error)
	GetUserRoles(userID uint64) ([]Role, error)
	GetUserPermissions(userID uint64) ([]Permission, error)
	ListUsers(params *UserListParams) ([]User, int64, error)
	GetUserByID(id uint64) (*User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	DisableUser(id uint64) error
	ReplaceUserRoles(userID uint64, roleIDs []uint64) error
	ListRoles() ([]Role, error)
	ListPermissions() ([]Permission, error)
}
