package service

import (
	"context"
	"github.com/olongfen/go-ddd-hex/internal/domain/dependency"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/olongfen/go-ddd-hex/internal/domain/vo"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repoFn func(ctx context.Context) dependency.UserRepo
	txImpl dependency.Transaction
}

func NewUserService(txImpl dependency.Transaction, repoFn func(ctx context.Context) dependency.UserRepo) *userService {
	return &userService{repoFn: repoFn, txImpl: txImpl}
}

func (u *userService) Get(ctx context.Context, id string) (res *vo.UserRes, err error) {
	var (
		data *entity.User
	)
	if data, err = u.repoFn(ctx).Get(id); err != nil {
		return
	}
	res = new(vo.UserRes)
	res.Username = data.Username
	res.Nickname = data.Nickname.String
	res.IsAdmin = data.IsAdmin.Bool
	res.CreatedAt = data.CreatedAt
	res.UpdatedAt = data.UpdatedAt
	return
}

func (u *userService) ChangePassword(ctx context.Context, id string, oldPwd, newPwd string) (err error) {
	var (
		data *entity.User
	)
	if data, err = u.repoFn(ctx).Get(id); err != nil {
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
	if err = u.repoFn(ctx).Update(data.QueryCond(), map[string]interface{}{"password": data.Password.String}); err != nil {
		return
	}

	return
}
