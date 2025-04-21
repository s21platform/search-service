package elsearch

import (
	"context"

	"github.com/s21platform/search-service/internal/model"
)

type Elastic interface {
	Update(ctx context.Context, id string, doc model.UserInfo) error
}
