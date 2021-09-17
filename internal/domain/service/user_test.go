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
	"gorm.io/gorm"
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
	f := vo.UserUnique{ID: "1"}
	repo.EXPECT().Get(ctx, f).Return(&entity.User{Username: "test1"}, nil)
	if u, err := s.Get(ctx, f); err != nil {
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
	user := &entity.User{Model: gorm.Model{ID: 1}, Password: p}
	repo.EXPECT().Get(ctx, gomock.Eq(vo.UserUnique{ID: id})).Return(user, nil)
	repo.EXPECT().Update(ctx, map[string]interface{}{"id": user.ID}, gomock.Any()).Return(nil)
	if err := s.ChangePassword(ctx, id, "111111", "123456"); err != nil {
		t.Fatal(err)
	}

}
