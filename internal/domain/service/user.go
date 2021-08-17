package service

import (
	"context"
	"github.com/olongfen/go-ddd-hex/internal/domain/dependency"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/olongfen/go-ddd-hex/internal/domain/vo"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo   dependency.UserRepo
	txImpl dependency.Transaction
}

func NewUserService(txImpl dependency.Transaction, repo dependency.UserRepo) *userService {
	return &userService{repo: repo, txImpl: txImpl}
}

func (u *userService) Get(ctx context.Context, id string) (res *vo.UserVO, err error) {
	var (
		data *entity.User
	)
	if data, err = u.repo.Get(ctx, id); err != nil {
		return
	}
	res = vo.UserEntity2VO(data)
	return
}

func (u *userService) ChangePassword(ctx context.Context, id string, oldPwd, newPwd string) (err error) {
	var (
		data *entity.User
	)
	span, _ := opentracing.StartSpanFromContext(ctx, "userService-ChangePassword")
	defer func() {
		if err != nil {
			span.LogFields(log.Error(err))
		}
		span.Finish()
	}()
	if data, err = u.repo.Get(ctx, id); err != nil {
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(data.Password.String), []byte(oldPwd)); err != nil {
		return
	}
	_n, _err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if _err != nil {
		err = _err
		return
	}
	data.Password.SetValid(string(_n))
	if err = u.repo.Update(ctx, data.QueryCond(), map[string]interface{}{"password": data.Password.String}); err != nil {
		return
	}

	return
}
