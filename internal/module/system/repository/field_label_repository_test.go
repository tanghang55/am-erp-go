package repository

import (
	"testing"

	systemdomain "am-erp-go/internal/module/system/domain"
)

func TestFieldLabelRepositoryImplementsInterface(t *testing.T) {
	var _ systemdomain.FieldLabelRepository = (*fieldLabelRepository)(nil)
}
