package repository

import (
	"context"
	"github.com/olongfen/go-ddd-hex/config"
	"github.com/olongfen/go-ddd-hex/internal/application"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/olongfen/go-ddd-hex/internal/infra/db"
	"github.com/olongfen/go-ddd-hex/lib/query"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type postRepo struct {
	db *gorm.DB
}

func init() {
	application.App.AppendRepo(NewPostRepo(application.App.Ctx, application.App.Database))
	db.RegisterInjector(func(db *gorm.DB) {
		if config.GetConfig().AutoMigrate {
			err := db.AutoMigrate(&entity.Post{})
			if err != nil {
				logrus.Fatal(err)
			}
		}
	})
}

func NewPostRepo(ctx context.Context, database application.Database) *postRepo {
	return &postRepo{db: database.DB(ctx).(*gorm.DB)}
}

func (u *postRepo) Get(ctx context.Context, id string) (res *entity.Post, err error) {
	var (
		data = new(entity.Post)
	)
	if err = u.db.WithContext(ctx).Model(&entity.Post{}).Where("id = ?", id).First(data).Error; err != nil {
		return
	}

	res = data
	return
}

func (u *postRepo) Find(ctx context.Context, cond map[string]interface{}, meta *query.Meta) (res []*entity.Post, err error) {
	withContext := u.db.WithContext(ctx)
	meta.WithOffsetLimit(withContext)
	if err = withContext.Where(cond).Find(&res).Error; err != nil {
		return
	}
	return
}

func (u *postRepo) Create(ctx context.Context, Post *entity.Post) error {
	return u.db.WithContext(ctx).Create(Post).Error
}

func (u *postRepo) Update(ctx context.Context, cond map[string]interface{}, change interface{}) error {
	if err := u.db.WithContext(ctx).Model(&entity.Post{}).Where(cond).Updates(change).Error; err != nil {
		return err
	}
	return nil
}

func (u *postRepo) Delete(ctx context.Context, cond map[string]interface{}) error {
	return u.db.WithContext(ctx).Model(&entity.Post{}).Delete(cond).Error
}
