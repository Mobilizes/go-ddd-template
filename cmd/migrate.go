package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const migrationPath = "internal/database"

func MigrateUp() {
	runMigration("up")
}

func MigrateDown() {
	runMigration("down", "1")
}

func runMigration(command string, args ...string) {
	path, err := filepath.Abs(migrationPath)
	if err != nil {
		fmt.Printf("error resolving migration path: %v\n", err)
		return
	}

	databaseURL, err := buildMigrationDatabaseURL()
	if err != nil {
		fmt.Printf("error building migration database URL: %v\n", err)
		return
	}

	commandArgs := []string{"-path", path, "-database", databaseURL, command}
	commandArgs = append(commandArgs, args...)

	migrateCommand := exec.Command("migrate", commandArgs...)
	output, err := migrateCommand.CombinedOutput()
	outputText := strings.TrimSpace(string(output))

	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			fmt.Println("migrate CLI is not available in PATH")
			return
		}

		if strings.Contains(outputText, "no change") {
			fmt.Println("no migration changes detected")
			return
		}

		if outputText != "" {
			fmt.Printf("migration %s failed: %s\n", command, outputText)
			return
		}

		fmt.Printf("migration %s failed: %v\n", command, err)
		return
	}

	if outputText != "" {
		fmt.Println(outputText)
	}

	fmt.Printf("migration %s completed successfully\n", command)
}

func buildMigrationDatabaseURL() (string, error) {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	missingEnv := []string{}
	if dbUser == "" {
		missingEnv = append(missingEnv, "DB_USER")
	}
	if dbHost == "" {
		missingEnv = append(missingEnv, "DB_HOST")
	}
	if dbName == "" {
		missingEnv = append(missingEnv, "DB_NAME")
	}
	if dbPort == "" {
		missingEnv = append(missingEnv, "DB_PORT")
	}

	if len(missingEnv) > 0 {
		return "", fmt.Errorf("missing required env vars: %s", strings.Join(missingEnv, ", "))
	}

	if dbSSLMode == "" {
		dbSSLMode = "require"
	}

	query := url.Values{}
	query.Set("sslmode", dbSSLMode)

	databaseURL := &url.URL{
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%s", dbHost, dbPort),
		Path:     "/" + dbName,
		RawQuery: query.Encode(),
	}

	if dbPass == "" {
		databaseURL.User = url.User(dbUser)
	} else {
		databaseURL.User = url.UserPassword(dbUser, dbPass)
	}

	return databaseURL.String(), nil
}
