package application

import (
	"context"
	"github.com/olongfen/go-ddd-hex/internal/domain/aggregate"
	"github.com/olongfen/go-ddd-hex/internal/domain/service"
	"github.com/olongfen/go-ddd-hex/internal/domain/vo"
)

var (
	// 告诉编译器接口是否实现
	_ UserServiceInterface = (*service.UserService)(nil)
	_ PostServiceInterface = (*service.PostService)(nil)
)

// UserServiceInterface user 用户服务接口
type UserServiceInterface interface {
	Create(ctx context.Context, forms []*vo.UserVOForm) (res []*vo.UserVO, err error)
	ChangePassword(ctx context.Context, id string, oldPwd, newPwd string) error
	Get(ctx context.Context, unique vo.UserUnique) (res *vo.UserVO, err error)
	Register(ctx context.Context, f vo.RegisterForm) error
}

// PostServiceInterface post 服务接口
type PostServiceInterface interface {
	GetByUserID(ctx context.Context, userID string) (*aggregate.QueryUserPostRes, error)
	Create(ctx context.Context)
}

// Service service 服务接口
type Service interface {
}
