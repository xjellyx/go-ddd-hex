package aggregate

import (
	"context"
	"github.com/olongfen/go-ddd-hex/internal/domain/dependency"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/olongfen/go-ddd-hex/internal/domain/vo"
	"github.com/olongfen/go-ddd-hex/lib/query"
)

// UserPostFactory 帖子与用户逻辑聚合工厂模式
type UserPostFactory struct {
	UserRepo dependency.UserRepo // 用户存储库
	PostRepo dependency.PostRepo // 帖子存储库
}

type QueryUserPostRes struct {
	User  vo.UserVO   `json:"user"`
	Posts []vo.PostVO `json:"posts"`
}

func NewUserPostFactory(postRepo dependency.PostRepo,
	userRepo dependency.UserRepo) *UserPostFactory {
	return &UserPostFactory{userRepo, postRepo}
}

func (f *UserPostFactory) UserPostQuery(ctx context.Context, userId string) (res *QueryUserPostRes, err error) {
	var (
		data     *entity.User
		dataPost []*entity.Post
	)

	if data, err = f.UserRepo.Get(ctx, userId); err != nil {
		return nil, err
	}
	if dataPost, err = f.PostRepo.Find(ctx, map[string]interface{}{
		"user_uuid": data.UUID,
	}, &query.Meta{}); err != nil {
		return
	}

	res = new(QueryUserPostRes)
	res.User = *vo.UserEntity2VO(data)
	for _, v := range dataPost {
		d := vo.PostEntity2VO(v)
		res.Posts = append(res.Posts, *d)

	}
	return
}
