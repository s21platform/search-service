package user_create

import (
	"context"
)

type Elastic interface {
	SetUser(ctx context.Context, data map[string]interface{}) error
}
