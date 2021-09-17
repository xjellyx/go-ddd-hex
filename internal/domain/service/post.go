package service

import (
	"context"
	"github.com/olongfen/go-ddd-hex/internal/domain/aggregate"
	"github.com/olongfen/go-ddd-hex/internal/domain/dependency"
	"github.com/olongfen/go-ddd-hex/internal/domain/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type PostService struct {
	repo        dependency.PostRepo
	postFactory *aggregate.UserPostFactory
}

func NewPostService(postRepo dependency.PostRepo,
	userRepo dependency.UserRepo) *PostService {
	return &PostService{repo: postRepo, postFactory: aggregate.NewUserPostFactory(postRepo, userRepo)}
}

func (p *PostService) GetByUserID(ctx context.Context, userId string) (res *aggregate.QueryUserPostRes, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "PostService-GetByUserID")
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

func (p *PostService) Create(ctx context.Context) {
	p.repo.Create(ctx, []*entity.Post{{Title: "dsafadsfadsf", Content: "dsafjldsajglkjhglkhjlkhjiuhfkghsdakhg"}})
}
