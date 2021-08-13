package dependency

import (
	"context"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/olongfen/go-ddd-hex/lib/query"
)

type UserRepo interface {
	Get(ctx context.Context, id string) (*entity.User, error)
	Find(ctx context.Context, cond map[string]interface{}, meta *query.Meta) ([]*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, cond map[string]interface{}, change interface{}) error
	Delete(ctx context.Context, cond map[string]interface{}) error
}

type Transaction interface {
	Transaction(ctx context.Context, fn func(txCtx context.Context) error) error
}
