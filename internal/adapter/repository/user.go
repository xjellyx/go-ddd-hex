package repository

import (
	"context"
	"github.com/olongfen/go-ddd-hex/internal/application"
	"github.com/olongfen/go-ddd-hex/internal/contant"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/olongfen/go-ddd-hex/internal/domain/vo"
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

func (u *userRepo) Get(ctx context.Context, unique vo.UserUnique) (res *entity.User, err error) {
	var (
		data = new(entity.User)
	)
	ctx = context.WithValue(ctx, contant.RepositoryMethodCtxTag, "userRepo-Get")
	db := u.db.WithContext(ctx).Model(&entity.User{})
	switch {
	case len(unique.ID) > 0:
		if err = db.Where("id = ?", unique.ID).First(data).Error; err != nil {
			return
		}
	case len(unique.Phone) > 0:
		if err = db.Where("phone = ?", unique.Phone).First(data).Error; err != nil {
			return
		}
	case len(unique.UUID) > 0:
		if err = db.Where("uuid = ?", unique.UUID).First(data).Error; err != nil {
			return
		}
	case len(unique.Username) > 0:
		if err = db.Where("username = ?", unique.Username).First(data).Error; err != nil {
			return
		}
	}

	res = data
	return
}

func (u *userRepo) Find(ctx context.Context, cond map[string]interface{}, meta *query.Meta) (res []*entity.User, err error) {
	withContext := u.db.WithContext(ctx)
	ctx = context.WithValue(ctx, contant.RepositoryMethodCtxTag, "userRepo-Find")
	if err = withContext.WithContext(ctx).Where(cond).Offset(meta.Offset()).Limit(meta.Limit()).Find(&res).Error; err != nil {
		return
	}
	return
}

func (u *userRepo) Create(ctx context.Context, users []*entity.User) error {
	ctx = context.WithValue(ctx, contant.RepositoryMethodCtxTag, "userRepo-Create")
	return u.db.WithContext(ctx).Create(users).Error
}

func (u *userRepo) Update(ctx context.Context, cond map[string]interface{}, change interface{}) error {
	ctx = context.WithValue(ctx, contant.RepositoryMethodCtxTag, "userRepo-ChangePasswd")
	if err := u.db.WithContext(ctx).Model(&entity.User{}).Where(cond).Updates(change).Error; err != nil {
		return err
	}
	return nil
}

func (u *userRepo) Delete(ctx context.Context, cond map[string]interface{}) error {
	ctx = context.WithValue(ctx, contant.RepositoryMethodCtxTag, "userRepo-Delete")
	return u.db.WithContext(ctx).Model(&entity.User{}).Delete(cond).Error
}
