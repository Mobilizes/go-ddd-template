package cmd

import (
	"fmt"
	"mob/ddd-template/cmd/seed"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

const defaultSeedPath = "seeds"

type seedFunc func(db *gorm.DB, path string) error

type pairSeed struct {
	key string
	fn  seedFunc
}

var seeders = []pairSeed{
	{"users", seed.SeedUsers},
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
	seedersByName := map[string]seedFunc{}
	for _, seeder := range seeders {
		seedersByName[seeder.key] = seeder.fn

		path, ok := seedFilesByName[seeder.key]
		if !ok {
			continue
		}

		fmt.Printf("Seeding %s from %s...\n", seeder.key, path)
		if err := seeder.fn(db, path); err != nil {
			return fmt.Errorf("seed %s: %w", seeder.key, err)
		}
		seeded[seeder.key] = true
	}

	for name, path := range seedFilesByName {
		if seeded[name] {
			continue
		}

		seeder, ok := seedersByName[name]
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
