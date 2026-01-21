package repository

import (
	"testing"

	systemdomain "am-erp-go/internal/module/system/domain"
)

func TestAuditLogRepositoryImplementsInterface(t *testing.T) {
	var _ systemdomain.AuditLogRepository = (*auditLogRepository)(nil)
}
