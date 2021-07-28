package main

import (
	"github.com/olongfen/go-ddd-hex/config"
	"github.com/olongfen/go-ddd-hex/internal/adapter/repository"
	"github.com/olongfen/go-ddd-hex/internal/contanst"
	service2 "github.com/olongfen/go-ddd-hex/internal/domain/service"

	"github.com/olongfen/go-ddd-hex/internal/adapter/xhttp/xgin"
	"github.com/olongfen/go-ddd-hex/internal/application"
	"github.com/olongfen/go-ddd-hex/internal/infra/db"
)

func main() {
	cfg := config.GetConfig()
	db := db.NewDatabase(&cfg.DBConfig)
	userDB := repository.NewUserDB(db)
	u := service2.NewUserService(repository.NewNoopTransaction(), userDB.GetRepo)
	p := service2.NewPostService(repository.NewNoopTransaction(), repository.NewPostDB(db).GetRepo, userDB.GetRepo)
	g := xgin.NewXGin(cfg.Ctx, cfg, map[string]interface{}{contanst.UserTag: u,
		contanst.PostTag: p})
	app := application.NewApplication(g, db)
	app.Database.Connect()
	app.Http.Run()
}
