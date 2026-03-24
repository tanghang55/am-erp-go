package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"am-erp-go/internal/infrastructure/config"
	"am-erp-go/internal/infrastructure/db"
	"am-erp-go/internal/infrastructure/seed"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	envFile := flag.String("env-file", ".env", "env file path")
	seedFile := flag.String("seed-file", filepath.Clean("baseline/minimal/minimal_seed.sql"), "minimal seed SQL path")
	credentialOut := flag.String("credential-out", "admin_credentials.txt", "admin credential output file")
	adminUsername := flag.String("admin-username", "admin", "default admin username")
	adminRealName := flag.String("admin-real-name", "系统管理员", "default admin real name")
	passwordLength := flag.Int("password-length", 20, "generated admin password length")
	flag.Parse()

	if err := godotenv.Overload(*envFile); err != nil {
		log.Fatalf("load env file failed: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	gormDB, err := db.NewMySQL(&cfg.Database)
	if err != nil {
		log.Fatalf("connect database failed: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("get sql db failed: %v", err)
	}

	seedContent, err := os.ReadFile(*seedFile)
	if err != nil {
		log.Fatalf("read seed file failed: %v", err)
	}

	password, err := seed.GeneratePassword(*passwordLength)
	if err != nil {
		log.Fatalf("generate password failed: %v", err)
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("hash password failed: %v", err)
	}

	if err := applyMinimalSeed(sqlDB, string(seedContent), *adminUsername, *adminRealName, string(passwordHash)); err != nil {
		log.Fatalf("apply minimal seed failed: %v", err)
	}

	credentialText := fmt.Sprintf("username=%s\npassword=%s\n", *adminUsername, password)
	if err := os.MkdirAll(filepath.Dir(*credentialOut), 0o755); err != nil {
		log.Fatalf("create credential directory failed: %v", err)
	}
	if err := os.WriteFile(*credentialOut, []byte(credentialText), 0o600); err != nil {
		log.Fatalf("write credential file failed: %v", err)
	}

	fmt.Printf("minimal seed initialized successfully\ncredential file: %s\n", *credentialOut)
}

func applyMinimalSeed(db *sql.DB, seedSQL string, adminUsername string, adminRealName string, passwordHash string) error {
	adminUsername = strings.TrimSpace(adminUsername)
	if adminUsername == "" {
		return fmt.Errorf("admin username is required")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for _, stmt := range seed.SplitStatements(seedSQL) {
		if _, err := tx.Exec(stmt); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("exec seed statement failed: %w", err)
		}
	}

	if _, err := tx.Exec(
		"DELETE ur FROM user_role ur INNER JOIN user u ON u.id = ur.user_id WHERE u.username = ?",
		adminUsername,
	); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("cleanup admin user_role failed: %w", err)
	}
	if _, err := tx.Exec("DELETE FROM user WHERE username = ?", adminUsername); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("cleanup admin user failed: %w", err)
	}

	result, err := tx.Exec(
		"INSERT INTO user (username, password_hash, real_name, status) VALUES (?, ?, ?, 'ACTIVE')",
		adminUsername,
		passwordHash,
		adminRealName,
	)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("insert admin user failed: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("get admin user id failed: %w", err)
	}

	var roleID uint64
	if err := tx.QueryRow("SELECT id FROM role WHERE name = 'admin' LIMIT 1").Scan(&roleID); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("load admin role failed: %w", err)
	}

	if _, err := tx.Exec("INSERT INTO user_role (user_id, role_id) VALUES (?, ?)", userID, roleID); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("bind admin role failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
