package repository

import (
	"context"
	"github.com/olongfen/go-ddd-hex/internal/application"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/olongfen/go-ddd-hex/lib/query"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func init() {
	application.App.AppendRepo(NewUserRepo(application.App.Ctx, application.App.Database))
}

func NewUserRepo(ctx context.Context, database application.Database) *userRepo {
	return &userRepo{db: database.DB(ctx).(*gorm.DB)}
}

func (u *userRepo) Get(ctx context.Context, id string) (res *entity.User, err error) {
	var (
		data = new(entity.User)
	)
	if err = u.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).First(data).Error; err != nil {
		return
	}

	res = data
	return
}

func (u *userRepo) Find(ctx context.Context, cond map[string]interface{}, meta *query.Meta) (res []*entity.User, err error) {
	withContext := u.db.WithContext(ctx)
	meta.WithOffsetLimit(withContext)
	if err = withContext.WithContext(ctx).Where(cond).Find(&res).Error; err != nil {
		return
	}
	return
}

func (u *userRepo) Create(ctx context.Context, user *entity.User) error {
	return u.db.WithContext(ctx).Create(user).Error
}

func (u *userRepo) Update(ctx context.Context, cond map[string]interface{}, change interface{}) error {
	if err := u.db.WithContext(ctx).Model(&entity.User{}).Where(cond).Updates(change).Error; err != nil {
		return err
	}
	return nil
}

func (u *userRepo) Delete(ctx context.Context, cond map[string]interface{}) error {
	return u.db.WithContext(ctx).Model(&entity.User{}).Delete(cond).Error
}
