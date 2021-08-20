package handler

import (
	"github.com/olongfen/go-ddd-hex/internal/application"
	log "github.com/sirupsen/logrus"
)

var (
	_ application.PostHandler = (*postCtl)(nil)
)

func init() {
	application.App.AppendHTTPGroupHandler(&postCtl{})
}

type postCtl struct {
	postQueryCtl
}

func (p *postCtl) SetService(domain application.PostServiceInterface) {
	p.domain = domain
}

func (p *postCtl) Router(xhttp application.XHttp, isGroup bool) {
	if p.domain == nil {
		log.Fatal("post controller domain service not init")
	}
	const post = "posts"
	xhttp.Handle("GET", post+"/:userId", p.GetByUserID, isGroup)
}
