package main

import (
	"fmt"
	"os"

	"ptm-indonesia/config"
	"ptm-indonesia/model"

	"gorm.io/gorm/clause"
)

func main() {
	cfg, err := config.NewAppConfig()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	db, cleanup, err := config.NewDatabase(cfg)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "open database: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	adminUser := &model.User{
		Name:      "Super Admin",
		Email:     cfg.Admin.Email,
		Status:    model.UserStatusSuperAdmin,
		StatusRow: model.ActiveUserStatusRow(),
		CreatedBy: 0,
		UpdatedBy: 0,
	}

	if err := db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "email"}, {Name: "status_row"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name",
				"status",
				"status_row",
				"updated_by",
				"updated_at",
			}),
		}).
		Create(adminUser).Error; err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "seed admin user: %v\n", err)
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(os.Stdout, "seeded admin user with email=%s\n", cfg.Admin.Email)
}
