package service

import (
	"github.com/golang/mock/gomock"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
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
		txmlRepo *user.MockTransaction
	)
	defer ctl.Finish()

	postRepo = post.NewMockPostRepo(ctl)
	postRepo.EXPECT().Find(map[string]interface{}{
		"user_uuid": "1",
	}, &query.Meta{PageNum: 1, PageSize: 10}).Return([]*entity.Post{
		{Title: "test_title"},
	}, nil)
	userRepo = user.NewMockUserRepo(ctl)
	userRepo.EXPECT().Get("1").Return(&entity.User{Username: "test"}, nil)
	txmlRepo = user.NewMockTransaction(ctl)
	p := NewPostService(txmlRepo, postRepo, userRepo)
	if data, err := p.GetByUserID("1"); err != nil {
		t.Fatal(err)
	} else {
		assert.Equal(t, len(data.Posts), 1)
		assert.Equal(t, data.User.Username, "test")
	}
}
