package repository

import (
	"context"
	"github.com/olongfen/go-ddd-hex/config"
	"github.com/olongfen/go-ddd-hex/internal/application"
	"github.com/olongfen/go-ddd-hex/internal/domain/dependency"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/olongfen/go-ddd-hex/internal/infra/db"
	"github.com/olongfen/go-ddd-hex/lib/query"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type postDB struct {
	db application.Database
}

func (u *postDB) GetRepo(ctx context.Context) dependency.PostRepo {
	return newPostRepo(ctx, u.db)
}

func NewPostDB(db application.Database) *postDB {
	return &postDB{db: db}
}

type postRepo struct {
	db *gorm.DB
}

func init() {
	db.RegisterInjector(func(db *gorm.DB) {
		if config.GetConfig().AutoMigrate {
			err := db.AutoMigrate(&entity.Post{})
			if err != nil {
				logrus.Fatal(err)
			}
		}
	})
}

func newPostRepo(ctx context.Context, database application.Database) *postRepo {
	return &postRepo{db: database.DB(ctx).(*gorm.DB)}
}

func (u *postRepo) Get(id string) (res *entity.Post, err error) {
	var (
		data = new(entity.Post)
	)
	if err = u.db.Model(&entity.Post{}).Where("id = ?", id).First(data).Error; err != nil {
		return
	}

	res = data
	return
}

func (u *postRepo) Find(cond map[string]interface{}, meta *query.Meta) (res []*entity.Post, err error) {
	meta.WithOffsetLimit(u.db)
	if err = u.db.Where(cond).Find(&res).Error; err != nil {
		return
	}
	return nil, nil
}

func (u *postRepo) Create(Post *entity.Post) error {
	return u.db.Create(Post).Error
}

func (u *postRepo) Update(cond map[string]interface{}, change interface{}) error {
	if err := u.db.Model(&entity.Post{}).Where(cond).Updates(change).Error; err != nil {
		return err
	}
	return nil
}

func (u *postRepo) Delete(cond map[string]interface{}) error {
	return u.db.Model(&entity.Post{}).Delete(cond).Error
}
