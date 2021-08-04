package service

import (
	"github.com/olongfen/go-ddd-hex/internal/domain/aggregate"
	"github.com/olongfen/go-ddd-hex/internal/domain/dependency"
)

type postService struct {
	repo        dependency.PostRepo
	postFactory *aggregate.UserPostFactory
	txImpl      dependency.Transaction
}

func NewPostService(txImpl dependency.Transaction, postRepo dependency.PostRepo,
	userRepo dependency.UserRepo) *postService {
	return &postService{repo: postRepo, txImpl: txImpl, postFactory: aggregate.NewUserPostFactory(postRepo, userRepo)}
}

func (p *postService) GetByUserID(userId string) (res *aggregate.QueryUserPostRes, err error) {
	return p.postFactory.UserPostQuery(userId)
}
