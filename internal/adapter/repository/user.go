package repository

import (
	"context"
	"github.com/olongfen/go-ddd-hex/internal/application"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/olongfen/go-ddd-hex/internal/infra/db"
	"github.com/olongfen/go-ddd-hex/lib/query"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func init() {
	application.App.AppendRepo(NewUserRepo(application.App.GetDB()))
}

func NewUserRepo(database application.Database) *userRepo {
	return &userRepo{db: database.DB().(*gorm.DB)}
}

func (u *userRepo) Get(ctx context.Context, id string) (res *entity.User, err error) {
	var (
		data = new(entity.User)
	)
	ctx = context.WithValue(ctx, db.RepositoryMethodCtxTag, "userRepo-Get")
	if err = u.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).First(data).Error; err != nil {
		return
	}

	res = data
	return
}

func (u *userRepo) Find(ctx context.Context, cond map[string]interface{}, meta *query.Meta) (res []*entity.User, err error) {
	withContext := u.db.WithContext(ctx)
	ctx = context.WithValue(ctx, db.RepositoryMethodCtxTag, "userRepo-Find")
	if err = withContext.WithContext(ctx).Where(cond).Offset(meta.Offset()).Limit(meta.Limit()).Find(&res).Error; err != nil {
		return
	}
	return
}

func (u *userRepo) Create(ctx context.Context, users []*entity.User) error {
	ctx = context.WithValue(ctx, db.RepositoryMethodCtxTag, "userRepo-Create")
	return u.db.WithContext(ctx).Create(users).Error
}

func (u *userRepo) Update(ctx context.Context, cond map[string]interface{}, change interface{}) error {
	ctx = context.WithValue(ctx, db.RepositoryMethodCtxTag, "userRepo-ChangePasswd")
	if err := u.db.WithContext(ctx).Model(&entity.User{}).Where(cond).Updates(change).Error; err != nil {
		return err
	}
	return nil
}

func (u *userRepo) Delete(ctx context.Context, cond map[string]interface{}) error {
	ctx = context.WithValue(ctx, db.RepositoryMethodCtxTag, "userRepo-Delete")
	return u.db.WithContext(ctx).Model(&entity.User{}).Delete(cond).Error
}
