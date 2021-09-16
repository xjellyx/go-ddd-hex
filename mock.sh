
mockgen -source internal/domain/dependency/post.go -destination ./mock/post/post_mock.go -package post
mockgen -source internal/domain/dependency/user.go -destination ./mock/user/user_mock.go -package user