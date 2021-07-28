package aggregate

import (
	"context"
	"github.com/olongfen/go-ddd-hex/internal/domain/dependency"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/olongfen/go-ddd-hex/internal/domain/vo"
	"github.com/olongfen/go-ddd-hex/lib/query"
)

type UserPostFactory struct {
	UserRepo func(ctx context.Context) dependency.UserRepo
	PostRepo func(ctx context.Context) dependency.PostRepo
}

type QueryUserPostRes struct {
	User  vo.UserRes   `json:"user"`
	Posts []vo.PostRes `json:"posts"`
}

func NewUserPostFactory(postRepo func(ctx context.Context) dependency.PostRepo,
	userRepoFn func(ctx context.Context) dependency.UserRepo) *UserPostFactory {
	return &UserPostFactory{userRepoFn, postRepo}
}

func (f *UserPostFactory) UserPostQuery(ctx context.Context, userId string) (res *QueryUserPostRes, err error) {
	var (
		data     *entity.User
		dataPost []*entity.Post
	)

	if data, err = f.UserRepo(ctx).Get(userId); err != nil {
		return nil, err
	}
	if dataPost, err = f.PostRepo(ctx).Find(map[string]interface{}{
		"user_uuid": userId,
	}, &query.Meta{PageNum: 1, PageSize: 10}); err != nil {
		return
	}

	res = new(QueryUserPostRes)
	res.User.Username = data.Username
	res.User.Nickname = data.Nickname.String
	res.User.IsAdmin = data.IsAdmin.Bool
	res.User.CreatedAt = data.CreatedAt
	res.User.UpdatedAt = data.UpdatedAt
	for _, v := range dataPost {
		d := vo.PostRes{
			BaseRes: vo.BaseRes{
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
			},
			Title:   v.Title,
			Content: v.Content,
		}
		res.Posts = append(res.Posts, d)

	}
	return
}
