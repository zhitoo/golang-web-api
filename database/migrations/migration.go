package migrations

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

func getMigrationsDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}

func getSourceUrl() string {
	sourceURL := "file://./src/data/migrations"
	return sourceURL
}

func getMigrator(gormDb *gorm.DB) (*migrate.Migrate, error) {
	sqlDb, err := gormDb.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDb, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		getSourceUrl(),
		"postgres",
		driver,
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// RunMigrations executes database migrations
func Up(gormDb *gorm.DB) error {
	m, err := getMigrator(gormDb)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Migration UP executed successfully.")
	return nil
}

func Down(gormDb *gorm.DB, steps int) error {
	m, err := getMigrator(gormDb)
	if err != nil {
		return err
	}

	if steps > 0 {
		steps *= -1
	}

	// m.Down() همه را پاک میکند. برای برگشت فقط یک مرحله از m.Steps(-1) استفاده کنید
	if err := m.Steps(steps); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Printf("Migration DOWN (%d step) executed successfully.\n", steps*-1)
	return nil
}

func Create(name string) error {
	if name == "" {
		return fmt.Errorf("migration name cannot be empty")
	}

	timestamp := time.Now().Unix()
	basePath := getMigrationsDir()

	upName := fmt.Sprintf("%s/%d_%s.up.sql", basePath, timestamp, name)
	downName := fmt.Sprintf("%s/%d_%s.down.sql", basePath, timestamp, name)

	// ساخت فایل up
	if err := os.WriteFile(upName, []byte(""), 0644); err != nil {
		return err
	}
	// ساخت فایل down
	if err := os.WriteFile(downName, []byte(""), 0644); err != nil {
		return err
	}

	log.Printf("Migration files created:\n%s\n%s\n", upName, downName)
	return nil
}

func Force(gormDb *gorm.DB, version int) error {
	m, err := getMigrator(gormDb)
	if err != nil {
		return err
	}
	return m.Force(version)
}
