package service

import (
	"github.com/golang/mock/gomock"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/olongfen/go-ddd-hex/mock/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserService_Get(t *testing.T) {
	var (
		txrepo *user.MockTransaction
		repo   *user.MockUserRepo
	)
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	repo = user.NewMockUserRepo(ctl)
	txrepo = user.NewMockTransaction(ctl)
	s := NewUserService(txrepo, repo)
	repo.EXPECT().Get("1").Return(&entity.User{Username: "test1"}, nil)
	if u, err := s.Get("1"); err != nil {
		t.Fatal(err)
	} else {
		assert.Equal(t, u.Username, "test1")
	}

}
