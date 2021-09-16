package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/guregu/null"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/olongfen/go-ddd-hex/internal/domain/vo"
	"github.com/olongfen/go-ddd-hex/mock/user"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestUserService_Get(t *testing.T) {
	var (
		repo *user.MockUserRepo
	)
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	repo = user.NewMockUserRepo(ctl)
	s := NewUserService(repo)
	ctx := context.Background()
	repo.EXPECT().Get(ctx, "1").Return(&entity.User{Username: "test1"}, nil)
	if u, err := s.Get(ctx, "1"); err != nil {
		t.Fatal(err)
	} else {
		assert.Equal(t, u.Username, "test1")
	}

}

func TestUserService_Register(t *testing.T) {
	var (
		repo *user.MockUserRepo
	)
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	repo = user.NewMockUserRepo(ctl)
	s := NewUserService(repo)
	ctx := context.Background()
	repo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
	if err := s.Register(ctx, vo.RegisterForm{Phone: "11111111111", Password: "aaaaaaaaaa"}); err != nil {
		t.Fatal(err)
	}
}

func TestService_ChangePassword(t *testing.T) {
	var (
		repo *user.MockUserRepo
		ctl  = gomock.NewController(t)
		id   = "1"
	)
	defer ctl.Finish()
	repo = user.NewMockUserRepo(ctl)
	s := NewUserService(repo)
	ctx := context.Background()
	p := null.String{}
	passwd, _ := bcrypt.GenerateFromPassword([]byte("111111"), bcrypt.DefaultCost)
	p.SetValid(string(passwd))
	repo.EXPECT().Get(ctx, id).Return(&entity.User{Password: p}, nil)
	repo.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).Return(nil)
	if err := s.ChangePassword(ctx, id, "111111", "123456"); err != nil {
		t.Fatal(err)
	}

}
