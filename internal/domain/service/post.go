package service

import (
	"context"
	"github.com/olongfen/go-ddd-hex/internal/domain/aggregate"
	"github.com/olongfen/go-ddd-hex/internal/domain/dependency"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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

func (p *postService) GetByUserID(ctx context.Context, userId string) (res *aggregate.QueryUserPostRes, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "postService-GetByUserID")
	defer func() {
		if err != nil {
			span.LogFields(log.Error(err))
		}
		span.Finish()
	}()
	if res, err = p.postFactory.UserPostQuery(ctx, userId); err != nil {
		return
	}
	return
}
