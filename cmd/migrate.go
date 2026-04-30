package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

const (
	migrationPath  = "internal/database"
	migrationTable = "schema_migrations"
	migrationLock  = int64(738459201)
)

var migrationFilePattern = regexp.MustCompile(`^([0-9]+)_.+\.(up|down)\.sql$`)

type migrationFile struct {
	Version int
	Dir     string
	Name    string
	SQL     string
}

type migrationRunner struct {
	conn       *sql.Conn
	migrations []migrationFile
}

type dirtyMigrationError struct {
	version int
}

func (e dirtyMigrationError) Error() string {
	return fmt.Sprintf("dirty database version %d. Fix the previous failed migration before running more migrations", e.version)
}

func MigrateUp() {
	runMigration("up")
}

func MigrateDown() {
	runMigration("down", "1")
}

func runMigration(command string, args ...string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

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

	config, err := pgx.ParseConfig(databaseURL)
	if err != nil {
		fmt.Printf("error parsing database connection: %v\n", err)
		return
	}
	config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	db := stdlib.OpenDB(*config)
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		fmt.Printf("error connecting to database: %v\n", err)
		return
	}

	conn, err := db.Conn(ctx)
	if err != nil {
		fmt.Printf("error reserving database connection: %v\n", err)
		return
	}
	defer conn.Close()

	if err := lockMigrations(ctx, conn); err != nil {
		fmt.Printf("error locking migrations: %v\n", err)
		return
	}
	defer unlockMigrations(conn)

	migrations, err := loadMigrationFiles(path)
	if err != nil {
		fmt.Printf("error loading migration files: %v\n", err)
		return
	}

	runner := migrationRunner{
		conn:       conn,
		migrations: migrations,
	}

	applied, err := runner.run(ctx, command, args...)
	if err != nil {
		fmt.Printf("migration %s failed: %v\n", command, err)
		return
	}

	if applied == 0 {
		fmt.Println("no migration changes detected")
		return
	}

	fmt.Printf("migration %s completed successfully (%d applied)\n", command, applied)
}

func lockMigrations(ctx context.Context, conn *sql.Conn) error {
	_, err := conn.ExecContext(ctx, "SELECT pg_advisory_lock($1)", migrationLock)
	return err
}

func unlockMigrations(conn *sql.Conn) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, _ = conn.ExecContext(ctx, "SELECT pg_advisory_unlock($1)", migrationLock)
}

func loadMigrationFiles(path string) ([]migrationFile, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	migrations := []migrationFile{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		matches := migrationFilePattern.FindStringSubmatch(entry.Name())
		if matches == nil {
			continue
		}

		version, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, fmt.Errorf("invalid migration version in %s: %w", entry.Name(), err)
		}

		migrationPath := filepath.Join(path, entry.Name())
		migrationSQL, err := os.ReadFile(migrationPath)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", entry.Name(), err)
		}

		migrations = append(migrations, migrationFile{
			Version: version,
			Dir:     matches[2],
			Name:    entry.Name(),
			SQL:     string(migrationSQL),
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		if migrations[i].Version == migrations[j].Version {
			return migrations[i].Dir < migrations[j].Dir
		}
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func (runner migrationRunner) run(ctx context.Context, command string, args ...string) (int, error) {
	if err := runner.ensureMigrationTable(ctx); err != nil {
		return 0, err
	}

	version, dirty, err := runner.currentVersion(ctx)
	if err != nil {
		return 0, err
	}

	if dirty {
		return 0, dirtyMigrationError{version: version}
	}

	switch command {
	case "up":
		return runner.up(ctx, version)
	case "down":
		steps, err := migrationSteps(args...)
		if err != nil {
			return 0, err
		}
		return runner.down(ctx, version, steps)
	default:
		return 0, fmt.Errorf("unsupported migration command %q", command)
	}
}

func (runner migrationRunner) ensureMigrationTable(ctx context.Context) error {
	_, err := runner.conn.ExecContext(ctx, fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (version bigint NOT NULL PRIMARY KEY, dirty boolean NOT NULL)",
		migrationTable,
	))
	return err
}

func (runner migrationRunner) currentVersion(ctx context.Context) (int, bool, error) {
	row := runner.conn.QueryRowContext(ctx, fmt.Sprintf("SELECT version, dirty FROM %s LIMIT 1", migrationTable))

	version := 0
	dirty := false
	if err := row.Scan(&version, &dirty); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, false, nil
		}
		return 0, false, err
	}

	return version, dirty, nil
}

func migrationSteps(args ...string) (int, error) {
	if len(args) == 0 {
		return 1, nil
	}

	steps, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("invalid migration step %q", args[0])
	}
	if steps < 1 {
		return 0, fmt.Errorf("migration step must be greater than 0")
	}

	return steps, nil
}

func (runner migrationRunner) up(ctx context.Context, currentVersion int) (int, error) {
	pending := []migrationFile{}
	for _, migration := range runner.migrations {
		if migration.Dir == "up" && migration.Version > currentVersion {
			pending = append(pending, migration)
		}
	}

	for i, migration := range pending {
		if err := runner.apply(ctx, migration, migration.Version); err != nil {
			return i, err
		}
	}

	return len(pending), nil
}

func (runner migrationRunner) down(ctx context.Context, currentVersion int, steps int) (int, error) {
	if currentVersion == 0 {
		return 0, nil
	}

	applied := 0
	for applied < steps && currentVersion > 0 {
		migration, ok := runner.findMigration(currentVersion, "down")
		if !ok {
			return applied, fmt.Errorf("missing down migration for version %d", currentVersion)
		}

		nextVersion := runner.previousVersion(currentVersion)
		if err := runner.apply(ctx, migration, nextVersion); err != nil {
			return applied, err
		}

		currentVersion = nextVersion
		applied++
	}

	return applied, nil
}

func (runner migrationRunner) findMigration(version int, dir string) (migrationFile, bool) {
	for _, migration := range runner.migrations {
		if migration.Version == version && migration.Dir == dir {
			return migration, true
		}
	}

	return migrationFile{}, false
}

func (runner migrationRunner) previousVersion(currentVersion int) int {
	previousVersion := 0
	for _, migration := range runner.migrations {
		if migration.Dir == "up" && migration.Version < currentVersion && migration.Version > previousVersion {
			previousVersion = migration.Version
		}
	}

	return previousVersion
}

func (runner migrationRunner) apply(ctx context.Context, migration migrationFile, targetVersion int) error {
	tx, err := runner.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if strings.TrimSpace(migration.SQL) != "" {
		if _, err := tx.ExecContext(ctx, migration.SQL); err != nil {
			return fmt.Errorf("%s: %w", migration.Name, err)
		}
	}

	if err := setMigrationVersion(ctx, tx, targetVersion); err != nil {
		return fmt.Errorf("updating migration version after %s: %w", migration.Name, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing %s: %w", migration.Name, err)
	}

	return nil
}

func setMigrationVersion(ctx context.Context, tx *sql.Tx, version int) error {
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s", migrationTable)); err != nil {
		return err
	}

	if version == 0 {
		return nil
	}

	_, err := tx.ExecContext(
		ctx,
		fmt.Sprintf("INSERT INTO %s (version, dirty) VALUES ($1, false)", migrationTable),
		version,
	)
	return err
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
