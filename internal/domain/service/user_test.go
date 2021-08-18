package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/olongfen/go-ddd-hex/mock/user"
	"github.com/stretchr/testify/assert"
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
