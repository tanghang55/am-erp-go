package domain

type UserListParams struct {
	Page     int
	PageSize int
	Status   string
	Keyword  string
}
