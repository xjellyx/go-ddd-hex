package service

import (
	"context"
	"errors"
	"github.com/olongfen/go-ddd-hex/internal/domain/dependency"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/olongfen/go-ddd-hex/internal/domain/vo"
	"github.com/olongfen/go-ddd-hex/lib/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo dependency.UserRepo
}

func NewUserService(repo dependency.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) Create(ctx context.Context, forms []*vo.UserVOForm) (res []*vo.UserVO, err error) {
	var (
		data []*entity.User
	)
	for _, v := range forms {
		data = append(data, vo.UserVOForm2Entity(v))
	}

	span, _ := opentracing.StartSpanFromContext(ctx, "UserService-Create")
	defer func() {
		if err != nil {
			span.LogFields(log.Error(err))
		}
		span.Finish()
	}()
	if err = u.repo.Create(ctx, data); err != nil {
		return
	}

	for _, v := range data {
		res = append(res, vo.UserEntity2VO(v))
	}

	return
}

func (u *UserService) Get(ctx context.Context, unique vo.UserUnique) (res *vo.UserVO, err error) {
	var (
		data *entity.User
	)
	span, _ := opentracing.StartSpanFromContext(ctx, "UserService-Get")
	span.SetTag("cond", unique)
	defer func() {
		if err != nil {
			span.LogFields(log.Error(err))
		}
		span.Finish()
	}()
	if data, err = u.repo.Get(ctx, unique); err != nil {
		return
	}
	res = vo.UserEntity2VO(data)
	return
}

func (u *UserService) ChangePassword(ctx context.Context, id string, oldPwd, newPwd string) (err error) {
	var (
		data *entity.User
	)
	span, _ := opentracing.StartSpanFromContext(ctx, "UserService-ChangePasswd")
	defer func() {
		if err != nil {
			span.LogFields(log.Error(err))
		}
		span.Finish()
	}()
	if data, err = u.repo.Get(ctx, vo.UserUnique{ID: id}); err != nil {
		return
	}
	if data.Password.Ptr() != nil {
		if err = bcrypt.CompareHashAndPassword([]byte(data.Password.String), []byte(oldPwd)); err != nil {
			return
		}
	}
	//_n, _err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	//if _err != nil {
	//	err = _err
	//	return
	//}
	//data.Password.SetValid(string(_n))
	data.SetPassword(newPwd)
	if err = u.repo.Update(ctx, data.QueryCond(), map[string]interface{}{"password": data.Password.String}); err != nil {
		return
	}

	return
}

func (u *UserService) Register(ctx context.Context, f vo.RegisterForm) (err error) {
	if len(f.Phone) == 0 {
		err = errors.New("phone must be send")
		return
	}

	if len(f.Password) == 0 {
		err = errors.New("password must be send")
		return
	}
	user := entity.NewUser(utils.RandString(16)).SetPhone(f.Phone).SetPassword(f.Password)
	if err = u.repo.Create(ctx, []*entity.User{user}); err != nil {
		return
	}
	return
}
