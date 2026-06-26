package services

import (
	"context"
	"time"

	"ptm-indonesia/config"
	"ptm-indonesia/model"
	repositoryContract "ptm-indonesia/repository/contract"
)

type HealthService struct {
	cfg              *config.AppConfig
	systemRepository repositoryContract.SystemRepository
}

func NewHealthService(
	cfg *config.AppConfig,
	systemRepository repositoryContract.SystemRepository,
) *HealthService {
	return &HealthService{
		cfg:              cfg,
		systemRepository: systemRepository,
	}
}

func (s *HealthService) Check(ctx context.Context) *model.HealthResponse {
	databaseStatus := "up"
	if err := s.systemRepository.Ping(ctx); err != nil {
		databaseStatus = "down"
	}

	location, err := time.LoadLocation(s.cfg.App.Timezone)
	if err != nil {
		location = time.UTC
	}

	return &model.HealthResponse{
		Name:               s.cfg.App.Name,
		Environment:        s.cfg.App.Environment,
		Database:           databaseStatus,
		DefaultLanguage:    s.cfg.App.DefaultLanguage,
		SupportedLanguages: s.cfg.App.SupportedLanguage,
		Timestamp:          time.Now().In(location).Format(time.RFC3339),
	}
}
