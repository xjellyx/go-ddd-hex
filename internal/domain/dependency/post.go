package dependency

import (
	"context"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/olongfen/go-ddd-hex/lib/query"
)

type PostRepo interface {
	Get(ctx context.Context, id string) (*entity.Post, error)
	Find(ctx context.Context, cond map[string]interface{}, meta *query.Meta) ([]*entity.Post, error)
	Create(ctx context.Context, post *entity.Post) error
	Update(ctx context.Context, cond map[string]interface{}, change interface{}) error
	Delete(ctx context.Context, cond map[string]interface{}) error
}
