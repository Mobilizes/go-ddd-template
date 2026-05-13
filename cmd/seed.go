package cmd

import (
	"encoding/json"
	"fmt"
	"mob/ddd-template/internal/domain/entity"
	"mob/ddd-template/internal/infra/security"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const defaultSeedPath = "seeds"

type seedFunc func(db *gorm.DB, path string) error

var seeders = map[string]seedFunc{
	"users": seedUsers,
	"files": seedFiles,
}

var seedOrder = []string{
	"users",
	"files",
}

func Seed() {
	db := SetUpDatabaseConnectionOrFail()
	fmt.Println("Database connection established for seeding.")

	seedPath := os.Getenv("SEED_PATH")
	if seedPath == "" {
		seedPath = defaultSeedPath
	}

	if err := runSeeders(db, seedPath); err != nil {
		fmt.Printf("Error seeding database: %v\n", err)
		return
	}

	fmt.Println("Seeding completed successfully.")
}

func runSeeders(db *gorm.DB, seedPath string) error {
	entries, err := os.ReadDir(seedPath)
	if err != nil {
		return fmt.Errorf("read seed directory %q: %w", seedPath, err)
	}

	seedFilesByName := map[string]string{}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		seedFilesByName[name] = filepath.Join(seedPath, entry.Name())
	}

	seeded := map[string]bool{}
	for _, name := range seedOrder {
		path, ok := seedFilesByName[name]
		if !ok {
			continue
		}

		seeder, ok := seeders[name]
		if !ok {
			return fmt.Errorf("configured seed order includes unregistered seeder %q", name)
		}

		fmt.Printf("Seeding %s from %s...\n", name, path)
		if err := seeder(db, path); err != nil {
			return fmt.Errorf("seed %s: %w", name, err)
		}
		seeded[name] = true
	}

	for name, path := range seedFilesByName {
		if seeded[name] {
			continue
		}

		seeder, ok := seeders[name]
		if !ok {
			fmt.Printf("Skipping %s: no registered seeder named %q.\n", filepath.Base(path), name)
			continue
		}

		fmt.Printf("Seeding %s from %s...\n", name, path)
		if err := seeder(db, path); err != nil {
			return fmt.Errorf("seed %s: %w", name, err)
		}
	}

	return nil
}

type userSeed struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordHash string `json:"password_hash"`
}

func seedUsers(db *gorm.DB, path string) error {
	var rows []userSeed
	if err := readSeedJSON(path, &rows); err != nil {
		return err
	}

	hasher := security.NewHasher()
	users := make([]*entity.User, 0, len(rows))
	for _, row := range rows {
		if row.ID == "" {
			row.ID = uuid.NewString()
		}

		password := row.PasswordHash
		if password == "" {
			if row.Password == "" {
				return fmt.Errorf("user %q must have password or password_hash", row.Email)
			}

			hashedPassword, err := hasher.RandomHash(row.Password)
			if err != nil {
				return err
			}
			password = hashedPassword
		}

		users = append(users, entity.NewUser(row.ID, row.Name, row.Email, password))
	}

	if len(users) == 0 {
		return nil
	}

	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"password",
			"updated_at",
			"deleted_at",
		}),
	}).Create(&users).Error
}

type fileSeed struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	StorageKey string `json:"storage_key"`
	MimeType   string `json:"mime_type"`
	SizeBytes  int64  `json:"size_bytes"`
}

func seedFiles(db *gorm.DB, path string) error {
	var rows []fileSeed
	if err := readSeedJSON(path, &rows); err != nil {
		return err
	}

	now := time.Now()
	files := make([]*entity.File, 0, len(rows))
	for _, row := range rows {
		if row.ID == "" {
			row.ID = uuid.NewString()
		}

		file := entity.NewFile(row.ID, row.UserID, row.Name, row.StorageKey, row.MimeType, row.SizeBytes)
		file.CreatedAt = now
		file.UpdatedAt = now
		files = append(files, file)
	}

	if len(files) == 0 {
		return nil
	}

	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"user_id",
			"name",
			"storage_key",
			"mime_type",
			"size_bytes",
			"updated_at",
			"deleted_at",
		}),
	}).Create(&files).Error
}

func readSeedJSON(path string, dst any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}

	return nil
}
