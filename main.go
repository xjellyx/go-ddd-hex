package main

import (
	"github.com/olongfen/go-ddd-hex/config"
	"github.com/olongfen/go-ddd-hex/internal/adapter/repository"
	"github.com/olongfen/go-ddd-hex/internal/adapter/xhttp/xgin"
	"github.com/olongfen/go-ddd-hex/internal/application"
	"github.com/olongfen/go-ddd-hex/internal/contanst"
	service "github.com/olongfen/go-ddd-hex/internal/domain/service"
	"github.com/olongfen/go-ddd-hex/internal/infra/db"
)

func main() {
	cfg := config.GetConfig()
	app := application.NewApplication()
	database := db.NewDatabase(&cfg.DBConfig)
	app.SetDatabase(database).Connect()
	userResp := repository.NewUserRepo(cfg.Ctx, database)
	user := service.NewUserService(repository.NewNoopTransaction(), userResp)
	post := service.NewPostService(repository.NewNoopTransaction(), repository.NewPostRepo(cfg.Ctx, database), userResp)
	ctl := xgin.NewXGin(cfg.Ctx, cfg, map[string]interface{}{contanst.UserTag: user,
		contanst.PostTag: post})
	app.SetXHttp(ctl).Inject().Register().Run()
}
