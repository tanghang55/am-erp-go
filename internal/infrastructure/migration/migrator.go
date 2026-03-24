package migration

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

var migrationFilePattern = regexp.MustCompile(`^\d{8,}.*\.sql$`)

var ErrBaselineRequired = errors.New("existing schema detected without migration history; run baseline first")

type MigrationFile struct {
	Version string
	Path    string
}

func LoadVersionManifest(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open version manifest: %w", err)
	}
	defer file.Close()

	versions := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		versions = append(versions, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan version manifest: %w", err)
	}
	return versions, nil
}

type Migrator struct {
	db  *sql.DB
	dir string
}

func NewMigrator(db *sql.DB, dir string) *Migrator {
	return &Migrator{db: db, dir: dir}
}

func ListPendingMigrations(dir string, applied map[string]struct{}) ([]MigrationFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read migration dir: %w", err)
	}

	files := make([]MigrationFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !migrationFilePattern.MatchString(name) {
			continue
		}
		if _, ok := applied[name]; ok {
			continue
		}
		files = append(files, MigrationFile{
			Version: name,
			Path:    filepath.Join(dir, name),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Version < files[j].Version
	})
	return files, nil
}

func ListMissingAppliedMigrations(dir string, applied map[string]struct{}) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read migration dir: %w", err)
	}

	existing := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !migrationFilePattern.MatchString(name) {
			continue
		}
		existing[name] = struct{}{}
	}

	missing := make([]string, 0)
	for version := range applied {
		if !migrationFilePattern.MatchString(version) {
			continue
		}
		if _, ok := existing[version]; ok {
			continue
		}
		missing = append(missing, version)
	}
	sort.Strings(missing)
	return missing, nil
}

func (m *Migrator) ApplyAll() ([]MigrationFile, error) {
	if err := m.ensureSchemaMigrationTable(); err != nil {
		return nil, err
	}
	applied, err := m.loadAppliedVersions()
	if err != nil {
		return nil, err
	}
	pending, err := ListPendingMigrations(m.dir, applied)
	if err != nil {
		return nil, err
	}
	tableCount, err := m.countBusinessTables()
	if err != nil {
		return nil, err
	}
	if RequireBaseline(len(applied), len(pending), tableCount) {
		return nil, ErrBaselineRequired
	}

	appliedNow := make([]MigrationFile, 0, len(pending))
	for _, file := range pending {
		if err := m.applyFile(file); err != nil {
			return appliedNow, err
		}
		appliedNow = append(appliedNow, file)
	}
	return appliedNow, nil
}

func (m *Migrator) BaselineAll() ([]MigrationFile, error) {
	if err := m.ensureSchemaMigrationTable(); err != nil {
		return nil, err
	}
	applied, err := m.loadAppliedVersions()
	if err != nil {
		return nil, err
	}
	pending, err := ListPendingMigrations(m.dir, applied)
	if err != nil {
		return nil, err
	}

	for _, file := range pending {
		if _, err := m.db.Exec(
			"INSERT INTO schema_migration (version, executed_at) VALUES (?, ?)",
			file.Version,
			time.Now(),
		); err != nil {
			return nil, fmt.Errorf("record baseline migration %s: %w", file.Version, err)
		}
	}
	return pending, nil
}

func (m *Migrator) BaselineVersions(versions []string) ([]MigrationFile, error) {
	if err := m.ensureSchemaMigrationTable(); err != nil {
		return nil, err
	}

	appliedNow := make([]MigrationFile, 0, len(versions))
	for _, version := range versions {
		version = strings.TrimSpace(version)
		if version == "" {
			continue
		}
		if _, err := m.db.Exec(
			"INSERT IGNORE INTO schema_migration (version, executed_at) VALUES (?, ?)",
			version,
			time.Now(),
		); err != nil {
			return nil, fmt.Errorf("record baseline migration %s: %w", version, err)
		}
		appliedNow = append(appliedNow, MigrationFile{Version: version})
	}
	return appliedNow, nil
}

func (m *Migrator) ensureSchemaMigrationTable() error {
	const ddl = `
CREATE TABLE IF NOT EXISTS schema_migration (
  id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  version varchar(255) NOT NULL COMMENT '迁移版本文件名',
  checksum varchar(64) DEFAULT NULL COMMENT '文件校验',
  executed_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '执行时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_version (version)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据库迁移执行记录表';`
	_, err := m.db.Exec(ddl)
	return err
}

func (m *Migrator) loadAppliedVersions() (map[string]struct{}, error) {
	rows, err := m.db.Query("SELECT version FROM schema_migration")
	if err != nil {
		return nil, fmt.Errorf("query applied migrations: %w", err)
	}
	defer rows.Close()

	result := make(map[string]struct{})
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		result[version] = struct{}{}
	}
	return result, rows.Err()
}

func (m *Migrator) countBusinessTables() (int, error) {
	row := m.db.QueryRow(`
SELECT COUNT(*)
FROM information_schema.tables
WHERE table_schema = DATABASE()
  AND table_type = 'BASE TABLE'
  AND table_name <> 'schema_migration'`)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("count business tables: %w", err)
	}
	return count, nil
}

func RequireBaseline(appliedCount int, pendingCount int, businessTableCount int) bool {
	return appliedCount == 0 && pendingCount > 0 && businessTableCount > 0
}

func (m *Migrator) applyFile(file MigrationFile) error {
	content, err := os.ReadFile(file.Path)
	if err != nil {
		return fmt.Errorf("read migration %s: %w", file.Version, err)
	}

	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("begin migration tx %s: %w", file.Version, err)
	}

	statements := splitSQLStatements(string(content))
	for _, statement := range statements {
		if _, err := tx.Exec(statement); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("execute migration %s: %w", file.Version, err)
		}
	}

	if _, err := tx.Exec(
		"INSERT INTO schema_migration (version, executed_at) VALUES (?, ?)",
		file.Version,
		time.Now(),
	); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("record migration %s: %w", file.Version, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %s: %w", file.Version, err)
	}
	return nil
}

func splitSQLStatements(sqlText string) []string {
	statements := make([]string, 0)
	var current strings.Builder
	var quote rune
	escaped := false

	flush := func() {
		statement := strings.TrimSpace(current.String())
		if statement != "" {
			statements = append(statements, statement)
		}
		current.Reset()
	}

	for _, ch := range sqlText {
		switch {
		case quote != 0:
			current.WriteRune(ch)
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' && quote != '`' {
				escaped = true
				continue
			}
			if ch == quote {
				quote = 0
			}
		case ch == '\'' || ch == '"' || ch == '`':
			quote = ch
			current.WriteRune(ch)
		case ch == ';':
			flush()
		default:
			current.WriteRune(ch)
		}
	}
	flush()
	return statements
}
