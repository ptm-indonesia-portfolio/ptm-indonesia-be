package contract

import "context"

type SystemRepository interface {
	Ping(ctx context.Context) error
}
