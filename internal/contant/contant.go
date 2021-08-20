package contant

const (
	UserTag = "USER"
	PostTag = "POST"
)

// GormSpanKey 包内静态变量
const GormSpanKey = "__gorm_span"
const (
	RepositoryMethodCtxTag = "repository_method"
	CallBackBeforeName     = "opentracing:before"
	CallBackAfterName      = "opentracing:after"
)
