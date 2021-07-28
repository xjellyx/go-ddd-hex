package dependency

import (
	"context"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/olongfen/go-ddd-hex/lib/query"
)

type UserRepo interface {
	Get(id string) (*entity.User, error)
	Find(cond map[string]interface{}, meta *query.Meta) ([]*entity.User, error)
	Create(user *entity.User) error
	Update(cond map[string]interface{}, change interface{}) error
	Delete(cond map[string]interface{}) error
}

type Transaction interface {
	Transaction(ctx context.Context, fn func(txCtx context.Context) error) error
}
