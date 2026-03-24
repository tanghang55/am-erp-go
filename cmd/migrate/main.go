package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"am-erp-go/internal/infrastructure/config"
	"am-erp-go/internal/infrastructure/db"
	"am-erp-go/internal/infrastructure/migration"

	"github.com/joho/godotenv"
)

func main() {
	baseline := flag.Bool("baseline", false, "record existing migration files without executing SQL")
	baselineFile := flag.String("baseline-file", "", "record manifest versions into schema_migration without executing SQL")
	envFile := flag.String("env-file", "", "env file path override")
	flag.Parse()

	if *envFile != "" {
		if err := godotenv.Overload(*envFile); err != nil {
			log.Fatalf("load env file failed: %v", err)
		}
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

	migrator := migration.NewMigrator(sqlDB, filepath.Clean("migrations"))
	var applied []migration.MigrationFile
	if *baseline && *baselineFile != "" {
		log.Fatalf("apply migrations failed: -baseline and -baseline-file cannot be used together")
	}
	if *baselineFile != "" {
		versions, err := migration.LoadVersionManifest(*baselineFile)
		if err != nil {
			log.Fatalf("load baseline file failed: %v", err)
		}
		applied, err = migrator.BaselineVersions(versions)
	} else if *baseline {
		applied, err = migrator.BaselineAll()
	} else {
		applied, err = migrator.ApplyAll()
	}
	if err != nil {
		if errors.Is(err, migration.ErrBaselineRequired) {
			log.Fatalf("apply migrations failed: %v. Run: go run .\\cmd\\migrate\\main.go -baseline", err)
		}
		log.Fatalf("apply migrations failed: %v", err)
	}

	if len(applied) == 0 {
		if *baseline {
			fmt.Println("no pending migrations to baseline")
		} else {
			fmt.Println("no pending migrations")
		}
		return
	}

	if *baseline || *baselineFile != "" {
		fmt.Printf("baselined %d migrations:\n", len(applied))
	} else {
		fmt.Printf("applied %d migrations:\n", len(applied))
	}
	for _, item := range applied {
		fmt.Printf("- %s\n", item.Version)
	}
}
