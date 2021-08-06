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

type userRepo struct {
	db *gorm.DB
}

func init() {
	application.app.AppendRepo(NewUserRepo(application.app.Ctx, application.app.Database))
	db.RegisterInjector(func(db *gorm.DB) {
		if config.GetConfig().AutoMigrate {
			err := db.AutoMigrate(&entity.User{})
			if err != nil {
				logrus.Fatal(err)
			}
		}
	})
}

func NewUserRepo(ctx context.Context, database application.Database) *userRepo {
	return &userRepo{db: database.DB(ctx).(*gorm.DB)}
}

func (u *userRepo) Get(id string) (res *entity.User, err error) {
	var (
		data = new(entity.User)
	)
	if err = u.db.Model(&entity.User{}).Where("id = ?", id).First(data).Error; err != nil {
		return
	}

	res = data
	return
}

func (u *userRepo) Find(cond map[string]interface{}, meta *query.Meta) (res []*entity.User, err error) {
	meta.WithOffsetLimit(u.db)
	if err = u.db.Where(cond).Find(&res).Error; err != nil {
		return
	}
	return nil, nil
}

func (u *userRepo) Create(user *entity.User) error {
	return u.db.Create(user).Error
}

func (u *userRepo) Update(cond map[string]interface{}, change interface{}) error {
	if err := u.db.Model(&entity.User{}).Where(cond).Updates(change).Error; err != nil {
		return err
	}
	return nil
}

func (u *userRepo) Delete(cond map[string]interface{}) error {
	return u.db.Model(&entity.User{}).Delete(cond).Error
}
