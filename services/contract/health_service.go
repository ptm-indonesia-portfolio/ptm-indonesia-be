package contract

import (
	"context"

	"ptm-indonesia/model"
)

type HealthService interface {
	Check(ctx context.Context) *model.HealthResponse
}
