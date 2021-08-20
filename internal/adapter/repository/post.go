package repository

import (
	"context"
	"github.com/olongfen/go-ddd-hex/internal/application"
	"github.com/olongfen/go-ddd-hex/internal/contant"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/olongfen/go-ddd-hex/lib/query"
	"gorm.io/gorm"
)

type postRepo struct {
	db *gorm.DB
}

func init() {
	application.App.AppendRepo(NewPostRepo(application.App.GetDB()))
}

func NewPostRepo(database application.Database) *postRepo {
	return &postRepo{db: database.DB().(*gorm.DB)}
}

func (u *postRepo) Get(ctx context.Context, id string) (res *entity.Post, err error) {
	var (
		data = new(entity.Post)
	)
	ctx = context.WithValue(ctx, contant.RepositoryMethodCtxTag, "postRepo-Get")
	if err = u.db.WithContext(ctx).Model(&entity.Post{}).Where("id = ?", id).First(data).Error; err != nil {
		return
	}

	res = data
	return
}

func (u *postRepo) Find(ctx context.Context, cond map[string]interface{}, meta *query.Meta) (res []*entity.Post, err error) {
	ctx = context.WithValue(ctx, contant.RepositoryMethodCtxTag, "postRepo-Find")
	withContext := u.db.WithContext(ctx)
	if err = withContext.Where(cond).Offset(meta.Offset()).Limit(meta.Limit()).Find(&res).Error; err != nil {
		return
	}
	return
}

func (u *postRepo) Create(ctx context.Context, Posts []*entity.Post) error {
	ctx = context.WithValue(ctx, contant.RepositoryMethodCtxTag, "postRepo-Create")
	return u.db.WithContext(ctx).Create(Posts).Error
}

func (u *postRepo) Update(ctx context.Context, cond map[string]interface{}, change interface{}) error {
	ctx = context.WithValue(ctx, contant.RepositoryMethodCtxTag, "postRepo-ChangePasswd")
	if err := u.db.WithContext(ctx).Model(&entity.Post{}).Where(cond).Updates(change).Error; err != nil {
		return err
	}
	return nil
}

func (u *postRepo) Delete(ctx context.Context, cond map[string]interface{}) error {
	ctx = context.WithValue(ctx, contant.RepositoryMethodCtxTag, "postRepo-Delete")
	return u.db.WithContext(ctx).Model(&entity.Post{}).Delete(cond).Error
}
