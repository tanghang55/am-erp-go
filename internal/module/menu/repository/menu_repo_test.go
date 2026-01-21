package repository

import (
    "testing"

    menudomain "am-erp-go/internal/module/menu/domain"
)

func TestMenuRepositoryImplementsInterface(t *testing.T) {
    var _ menudomain.MenuRepository = (*menuRepository)(nil)
}
