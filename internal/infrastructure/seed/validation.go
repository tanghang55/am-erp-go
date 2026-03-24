package seed

import (
	"database/sql"
	"fmt"
	"strings"
)

type permissionSeedRow struct {
	ID     uint64
	Code   string
	Status string
}

type roleSeedRow struct {
	ID   uint64
	Name string
}

type rolePermissionSeedRow struct {
	RoleID       uint64
	PermissionID uint64
}

type menuSeedRow struct {
	ID             uint64
	Code           string
	ParentID       sql.NullInt64
	PermissionCode sql.NullString
	Status         string
}

func ValidateMinimalSeedSource(db *sql.DB) error {
	permissions, err := loadPermissionSeedRows(db)
	if err != nil {
		return err
	}
	roles, err := loadRoleSeedRows(db)
	if err != nil {
		return err
	}
	rolePermissions, err := loadRolePermissionSeedRows(db)
	if err != nil {
		return err
	}
	menus, err := loadMenuSeedRows(db)
	if err != nil {
		return err
	}

	return validateMinimalSeedRows(permissions, roles, rolePermissions, menus)
}

func validateMinimalSeedRows(permissions []permissionSeedRow, roles []roleSeedRow, rolePermissions []rolePermissionSeedRow, menus []menuSeedRow) error {
	issues := make([]string, 0)
	permissionByID := make(map[uint64]permissionSeedRow, len(permissions))
	permissionByCode := make(map[string]permissionSeedRow, len(permissions))
	for _, item := range permissions {
		code := strings.TrimSpace(item.Code)
		if code == "" {
			issues = append(issues, fmt.Sprintf("permission id=%d code is empty", item.ID))
			continue
		}
		if _, exists := permissionByCode[code]; exists {
			issues = append(issues, fmt.Sprintf("permission code duplicated: %s", code))
			continue
		}
		permissionByID[item.ID] = item
		permissionByCode[code] = item
	}
	if len(permissionByID) == 0 {
		issues = append(issues, "permission table is empty")
	}

	roleByID := make(map[uint64]roleSeedRow, len(roles))
	adminRoleCount := 0
	for _, item := range roles {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			issues = append(issues, fmt.Sprintf("role id=%d name is empty", item.ID))
			continue
		}
		roleByID[item.ID] = item
		if name == "admin" {
			adminRoleCount++
		}
	}
	if adminRoleCount != 1 {
		issues = append(issues, fmt.Sprintf("expected exactly one admin role, got %d", adminRoleCount))
	}

	rolePermissionPairs := make(map[string]struct{}, len(rolePermissions))
	for _, item := range rolePermissions {
		if _, exists := roleByID[item.RoleID]; !exists {
			issues = append(issues, fmt.Sprintf("role_permission references missing role_id=%d", item.RoleID))
		}
		if _, exists := permissionByID[item.PermissionID]; !exists {
			issues = append(issues, fmt.Sprintf("role_permission references missing permission_id=%d", item.PermissionID))
		}
		pairKey := fmt.Sprintf("%d:%d", item.RoleID, item.PermissionID)
		if _, exists := rolePermissionPairs[pairKey]; exists {
			issues = append(issues, fmt.Sprintf("role_permission duplicated pair %s", pairKey))
			continue
		}
		rolePermissionPairs[pairKey] = struct{}{}
	}

	menuByID := make(map[uint64]menuSeedRow, len(menus))
	menuByCode := make(map[string]menuSeedRow, len(menus))
	for _, item := range menus {
		code := strings.TrimSpace(item.Code)
		if code == "" {
			issues = append(issues, fmt.Sprintf("menu id=%d code is empty", item.ID))
		} else {
			if _, exists := menuByCode[code]; exists {
				issues = append(issues, fmt.Sprintf("menu code duplicated: %s", code))
			}
			menuByCode[code] = item
		}
		menuByID[item.ID] = item
	}
	if len(menuByID) == 0 {
		issues = append(issues, "menu table is empty")
	}

	for _, item := range menus {
		if item.ParentID.Valid {
			if _, exists := menuByID[uint64(item.ParentID.Int64)]; !exists {
				issues = append(issues, fmt.Sprintf("menu code=%s references missing parent_id=%d", strings.TrimSpace(item.Code), item.ParentID.Int64))
			}
		}

		permissionCode := strings.TrimSpace(item.PermissionCode.String)
		if permissionCode == "" {
			issues = append(issues, fmt.Sprintf("menu code=%s permission_code is empty", strings.TrimSpace(item.Code)))
			continue
		}
		if _, exists := permissionByCode[permissionCode]; !exists {
			issues = append(issues, fmt.Sprintf("menu code=%s references missing permission_code=%s", strings.TrimSpace(item.Code), permissionCode))
		}
	}

	if len(issues) > 0 {
		return fmt.Errorf("minimal seed validation failed: %s", strings.Join(issues, "; "))
	}
	return nil
}

func loadPermissionSeedRows(db *sql.DB) ([]permissionSeedRow, error) {
	rows, err := db.Query("SELECT id, code, status FROM permission ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("query permission seed rows: %w", err)
	}
	defer rows.Close()

	result := make([]permissionSeedRow, 0)
	for rows.Next() {
		var item permissionSeedRow
		if err := rows.Scan(&item.ID, &item.Code, &item.Status); err != nil {
			return nil, fmt.Errorf("scan permission seed row: %w", err)
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func loadRoleSeedRows(db *sql.DB) ([]roleSeedRow, error) {
	rows, err := db.Query("SELECT id, name FROM role ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("query role seed rows: %w", err)
	}
	defer rows.Close()

	result := make([]roleSeedRow, 0)
	for rows.Next() {
		var item roleSeedRow
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, fmt.Errorf("scan role seed row: %w", err)
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func loadRolePermissionSeedRows(db *sql.DB) ([]rolePermissionSeedRow, error) {
	rows, err := db.Query("SELECT role_id, permission_id FROM role_permission ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("query role_permission seed rows: %w", err)
	}
	defer rows.Close()

	result := make([]rolePermissionSeedRow, 0)
	for rows.Next() {
		var item rolePermissionSeedRow
		if err := rows.Scan(&item.RoleID, &item.PermissionID); err != nil {
			return nil, fmt.Errorf("scan role_permission seed row: %w", err)
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func loadMenuSeedRows(db *sql.DB) ([]menuSeedRow, error) {
	rows, err := db.Query("SELECT id, code, parent_id, permission_code, status FROM menu ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("query menu seed rows: %w", err)
	}
	defer rows.Close()

	result := make([]menuSeedRow, 0)
	for rows.Next() {
		var item menuSeedRow
		if err := rows.Scan(&item.ID, &item.Code, &item.ParentID, &item.PermissionCode, &item.Status); err != nil {
			return nil, fmt.Errorf("scan menu seed row: %w", err)
		}
		result = append(result, item)
	}
	return result, rows.Err()
}
