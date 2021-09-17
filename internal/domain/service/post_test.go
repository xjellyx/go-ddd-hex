package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/olongfen/go-ddd-hex/internal/domain/vo"
	"github.com/olongfen/go-ddd-hex/lib/query"
	"github.com/olongfen/go-ddd-hex/mock/post"
	"github.com/olongfen/go-ddd-hex/mock/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPostService_GetByUserID(t *testing.T) {
	var (
		ctl      = gomock.NewController(t)
		postRepo *post.MockPostRepo
		userRepo *user.MockUserRepo
	)
	defer ctl.Finish()

	postRepo = post.NewMockPostRepo(ctl)
	ctx := context.Background()
	user_ := &entity.User{Username: "test", UUID: "1"}
	postRepo.EXPECT().Find(ctx, map[string]interface{}{
		"user_uuid": user_.UUID,
	}, &query.Meta{}).Return([]*entity.Post{
		{Title: "test_title"},
	}, nil)
	userRepo = user.NewMockUserRepo(ctl)
	userRepo.EXPECT().Get(ctx, gomock.Eq(vo.UserUnique{ID: "1"})).Return(user_, nil)
	p := NewPostService(postRepo, userRepo)
	if data, err := p.GetByUserID(ctx, "1"); err != nil {
		t.Fatal(err)
	} else {
		assert.Equal(t, len(data.Posts), 1)
		assert.Equal(t, data.User.Username, "test")
	}
}
