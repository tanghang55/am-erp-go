package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"am-erp-go/internal/infrastructure/config"
	"am-erp-go/internal/infrastructure/db"
	"am-erp-go/internal/infrastructure/seed"

	"github.com/joho/godotenv"
)

func main() {
	envFile := flag.String("env-file", ".env", "env file path")
	outputFile := flag.String("output-file", filepath.Clean("baseline/minimal/minimal_seed.sql"), "minimal seed output file")
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

	sqlText, err := seed.ExportMinimalSeedSQL(sqlDB)
	if err != nil {
		log.Fatalf("export minimal seed failed: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(*outputFile), 0o755); err != nil {
		log.Fatalf("create output directory failed: %v", err)
	}
	if err := os.WriteFile(*outputFile, []byte(sqlText), 0o644); err != nil {
		log.Fatalf("write output file failed: %v", err)
	}

	fmt.Printf("minimal seed exported successfully\noutput file: %s\n", *outputFile)
}
