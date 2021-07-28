package xgin

import "github.com/olongfen/go-ddd-hex/internal/application"

type PostCtl struct {
	PostQueryCtl
}

func NewPostCtl(domain application.PostInterface) *PostCtl {
	return &PostCtl{PostQueryCtl{domain: domain}}
}
