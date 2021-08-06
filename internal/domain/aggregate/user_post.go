package aggregate

import (
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
	User  vo.UserRes   `json:"user"`
	Posts []vo.PostRes `json:"posts"`
}

func NewUserPostFactory(postRepo dependency.PostRepo,
	userRepo dependency.UserRepo) *UserPostFactory {
	return &UserPostFactory{userRepo, postRepo}
}

func (f *UserPostFactory) UserPostQuery(userId string) (res *QueryUserPostRes, err error) {
	var (
		data     *entity.User
		dataPost []*entity.Post
	)

	if data, err = f.UserRepo.Get(userId); err != nil {
		return nil, err
	}
	if dataPost, err = f.PostRepo.Find(map[string]interface{}{
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
