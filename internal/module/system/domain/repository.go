package domain

type FieldLabelRepository interface {
	GetAll() ([]*FieldLabel, error)
	List(page, pageSize int, keyword string) ([]*FieldLabel, int64, error)
	GetByID(id uint64) (*FieldLabel, error)
	GetByKey(key string) (*FieldLabel, error)
	Create(label *FieldLabel) error
	Update(label *FieldLabel) error
	Delete(id uint64) error
}

type AuditLogRepository interface {
	List(params AuditLogListParams) ([]*AuditLog, int64, error)
	Create(log *AuditLog) error
}
