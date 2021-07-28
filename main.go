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
	database := db.NewDatabase(&cfg.DBConfig)
	userDB := repository.NewUserDB(database)
	user := service.NewUserService(repository.NewNoopTransaction(), userDB.GetRepo)
	post := service.NewPostService(repository.NewNoopTransaction(), repository.NewPostDB(database).GetRepo, userDB.GetRepo)
	ctl := xgin.NewXGin(cfg.Ctx, cfg, map[string]interface{}{contanst.UserTag: user,
		contanst.PostTag: post})
	app := application.NewApplication(ctl, database)
	app.Database.Connect()
	app.Http.Inject().Register().Run()
}
