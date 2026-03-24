package seed

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"
)

var minimalSeedTables = []string{"permission", "role", "role_permission", "menu", "finance_exchange_rate"}

func ExportMinimalSeedSQL(db *sql.DB) (string, error) {
	if err := ValidateMinimalSeedSource(db); err != nil {
		return "", err
	}

	parts := []string{
		"SET NAMES utf8mb4;",
		"SET FOREIGN_KEY_CHECKS = 0;",
		"",
		"DELETE FROM `user_role`;",
		"DELETE FROM `role_permission`;",
		"DELETE FROM `user`;",
		"DELETE FROM `finance_exchange_rate`;",
		"DELETE FROM `menu`;",
		"DELETE FROM `permission`;",
		"DELETE FROM `role`;",
		"",
	}

	for _, table := range minimalSeedTables {
		sqlText, err := exportTable(db, table)
		if err != nil {
			return "", err
		}
		if sqlText != "" {
			parts = append(parts, sqlText)
		}
	}

	parts = append(parts, "", "SET FOREIGN_KEY_CHECKS = 1;", "")
	return strings.Join(parts, "\n"), nil
}

func exportTable(db *sql.DB, table string) (string, error) {
	columns, err := loadColumns(db, table)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("SELECT * FROM `%s` ORDER BY `%s`", table, columns[0])
	rows, err := db.Query(query)
	if err != nil {
		return "", fmt.Errorf("query table %s: %w", table, err)
	}
	defer rows.Close()

	valuesSQL := make([]string, 0)
	for rows.Next() {
		rowValues := make([]any, len(columns))
		scanTargets := make([]any, len(columns))
		for i := range rowValues {
			scanTargets[i] = &rowValues[i]
		}
		if err := rows.Scan(scanTargets...); err != nil {
			return "", fmt.Errorf("scan table %s: %w", table, err)
		}

		serialized := make([]string, 0, len(rowValues))
		for _, value := range rowValues {
			serialized = append(serialized, sqlLiteral(value))
		}
		valuesSQL = append(valuesSQL, fmt.Sprintf("(%s)", strings.Join(serialized, ",")))
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("iterate table %s: %w", table, err)
	}
	if len(valuesSQL) == 0 {
		return "", nil
	}

	return fmt.Sprintf("INSERT INTO `%s` VALUES %s;", table, strings.Join(valuesSQL, ",")), nil
}

func loadColumns(db *sql.DB, table string) ([]string, error) {
	rows, err := db.Query(fmt.Sprintf("SHOW COLUMNS FROM `%s`", table))
	if err != nil {
		return nil, fmt.Errorf("show columns for %s: %w", table, err)
	}
	defer rows.Close()

	columns := make([]string, 0)
	for rows.Next() {
		var field string
		var typeName string
		var nullable string
		var key sql.NullString
		var defaultValue sql.NullString
		var extra string
		if err := rows.Scan(&field, &typeName, &nullable, &key, &defaultValue, &extra); err != nil {
			return nil, fmt.Errorf("scan columns for %s: %w", table, err)
		}
		columns = append(columns, field)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate columns for %s: %w", table, err)
	}
	if len(columns) == 0 {
		return nil, fmt.Errorf("table %s has no columns", table)
	}
	return columns, nil
}

func sqlLiteral(value any) string {
	if value == nil {
		return "NULL"
	}

	switch v := value.(type) {
	case []byte:
		return "'" + escapeSQLString(string(v)) + "'"
	case string:
		return "'" + escapeSQLString(v) + "'"
	case time.Time:
		return "'" + v.Format("2006-01-02 15:04:05") + "'"
	case bool:
		if v {
			return "1"
		}
		return "0"
	default:
		return "'" + escapeSQLString(fmt.Sprint(v)) + "'"
	}
}

func escapeSQLString(value string) string {
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		"'", "\\'",
		"\n", "\\n",
		"\r", "\\r",
		"\x00", "\\0",
	)
	return replacer.Replace(value)
}

func MinimalSeedTables() []string {
	result := make([]string, len(minimalSeedTables))
	copy(result, minimalSeedTables)
	sort.Strings(result)
	return result
}
