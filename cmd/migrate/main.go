package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"ptm-indonesia/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.NewAppConfig()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	migrationConfig := config.NewMigrationRunner(cfg)

	migrator, err := migrate.New(migrationConfig.Source, migrationConfig.DatabaseURL)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "create migrator: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_, _ = migrator.Close()
	}()

	if len(os.Args) < 2 {
		_, _ = fmt.Fprintln(os.Stderr, "usage: go run ./cmd/migrate [up|down|version|force <version>]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "up":
		if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			_, _ = fmt.Fprintf(os.Stderr, "run migrate up: %v\n", err)
			os.Exit(1)
		}
	case "down":
		if err := migrator.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			_, _ = fmt.Fprintf(os.Stderr, "run migrate down: %v\n", err)
			os.Exit(1)
		}
	case "version":
		version, dirty, err := migrator.Version()
		if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
			_, _ = fmt.Fprintf(os.Stderr, "get migration version: %v\n", err)
			os.Exit(1)
		}

		if errors.Is(err, migrate.ErrNilVersion) {
			_, _ = fmt.Fprintln(os.Stdout, "version=0 dirty=false")
			return
		}

		_, _ = fmt.Fprintf(os.Stdout, "version=%d dirty=%t\n", version, dirty)
	case "force":
		if len(os.Args) < 3 {
			_, _ = fmt.Fprintln(os.Stderr, "usage: go run ./cmd/migrate force <version>")
			os.Exit(1)
		}

		version, convErr := strconv.Atoi(os.Args[2])
		if convErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "invalid force version: %v\n", convErr)
			os.Exit(1)
		}

		if err := migrator.Force(version); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "force migration version: %v\n", err)
			os.Exit(1)
		}
	default:
		_, _ = fmt.Fprintf(os.Stderr, "unknown migration command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
