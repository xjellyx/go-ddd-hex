package service

import (
	"context"
	"github.com/olongfen/go-ddd-hex/internal/domain/aggregate"
	"github.com/olongfen/go-ddd-hex/internal/domain/dependency"
)

type postService struct {
	repoFn      func(ctx context.Context) dependency.PostRepo
	postFactory *aggregate.UserPostFactory
	txImpl      dependency.Transaction
}

func NewPostService(txImpl dependency.Transaction, repoFn func(ctx context.Context) dependency.PostRepo,
	userRepoFn func(ctx context.Context) dependency.UserRepo) *postService {
	return &postService{repoFn: repoFn, txImpl: txImpl, postFactory: aggregate.NewUserPostFactory(repoFn, userRepoFn)}
}

func (p *postService) GetByUserID(ctx context.Context, userId string) (res *aggregate.QueryUserPostRes, err error) {
	return p.postFactory.UserPostQuery(ctx, userId)
}
