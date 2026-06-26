package config

type MigrationRunner struct {
	Source      string
	DatabaseURL string
}

func NewMigrationRunner(cfg *AppConfig) *MigrationRunner {
	return &MigrationRunner{
		Source:      cfg.Migration.Source,
		DatabaseURL: cfg.MigrationDatabaseURL(),
	}
}
