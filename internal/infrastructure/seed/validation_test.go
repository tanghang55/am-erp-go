package seed

import (
	"database/sql"
	"strings"
	"testing"
)

func TestValidateMinimalSeedRowsPassesForHealthyData(t *testing.T) {
	err := validateMinimalSeedRows(
		[]permissionSeedRow{
			{ID: 1, Code: "system.manage", Status: "ACTIVE"},
			{ID: 2, Code: "product.manage", Status: "ACTIVE"},
		},
		[]roleSeedRow{
			{ID: 10, Name: "admin"},
		},
		[]rolePermissionSeedRow{
			{RoleID: 10, PermissionID: 1},
			{RoleID: 10, PermissionID: 2},
		},
		[]menuSeedRow{
			{ID: 100, Code: "SYSTEM", PermissionCode: sql.NullString{String: "system.manage", Valid: true}, Status: "ACTIVE"},
			{ID: 101, Code: "PRODUCT", ParentID: sql.NullInt64{Int64: 100, Valid: true}, PermissionCode: sql.NullString{String: "product.manage", Valid: true}, Status: "ACTIVE"},
		},
	)
	if err != nil {
		t.Fatalf("expected validation to pass, got %v", err)
	}
}

func TestValidateMinimalSeedRowsFailsForMissingAdminRole(t *testing.T) {
	err := validateMinimalSeedRows(
		[]permissionSeedRow{{ID: 1, Code: "system.manage", Status: "ACTIVE"}},
		[]roleSeedRow{{ID: 10, Name: "operator"}},
		nil,
		[]menuSeedRow{{ID: 100, Code: "SYSTEM", PermissionCode: sql.NullString{String: "system.manage", Valid: true}, Status: "ACTIVE"}},
	)
	if err == nil || !strings.Contains(err.Error(), "expected exactly one admin role") {
		t.Fatalf("expected missing admin role error, got %v", err)
	}
}

func TestValidateMinimalSeedRowsFailsForMissingMenuPermission(t *testing.T) {
	err := validateMinimalSeedRows(
		[]permissionSeedRow{{ID: 1, Code: "system.manage", Status: "ACTIVE"}},
		[]roleSeedRow{{ID: 10, Name: "admin"}},
		[]rolePermissionSeedRow{{RoleID: 10, PermissionID: 1}},
		[]menuSeedRow{{ID: 100, Code: "SYSTEM", PermissionCode: sql.NullString{String: "missing.permission", Valid: true}, Status: "ACTIVE"}},
	)
	if err == nil || !strings.Contains(err.Error(), "references missing permission_code") {
		t.Fatalf("expected missing permission_code error, got %v", err)
	}
}

func TestValidateMinimalSeedRowsFailsForMissingParentMenu(t *testing.T) {
	err := validateMinimalSeedRows(
		[]permissionSeedRow{{ID: 1, Code: "system.manage", Status: "ACTIVE"}},
		[]roleSeedRow{{ID: 10, Name: "admin"}},
		[]rolePermissionSeedRow{{RoleID: 10, PermissionID: 1}},
		[]menuSeedRow{{ID: 100, Code: "SYSTEM_CHILD", ParentID: sql.NullInt64{Int64: 999, Valid: true}, PermissionCode: sql.NullString{String: "system.manage", Valid: true}, Status: "ACTIVE"}},
	)
	if err == nil || !strings.Contains(err.Error(), "references missing parent_id") {
		t.Fatalf("expected missing parent_id error, got %v", err)
	}
}
