package dependency

import (
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/olongfen/go-ddd-hex/lib/query"
)

type PostRepo interface {
	Get(id string) (*entity.Post, error)
	Find(cond map[string]interface{}, meta *query.Meta) ([]*entity.Post, error)
	Create(post *entity.Post) error
	Update(cond map[string]interface{}, change interface{}) error
	Delete(cond map[string]interface{}) error
}
