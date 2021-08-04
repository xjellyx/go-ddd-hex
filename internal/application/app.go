package application

import (
	"context"
	"github.com/olongfen/go-ddd-hex/internal/domain/aggregate"
	"github.com/olongfen/go-ddd-hex/internal/domain/vo"
)

type UserInterface interface {
	ChangePassword(id string, oldPwd, newPwd string) error
	Get(id string) (res *vo.UserRes, err error)
}

type PostInterface interface {
	GetByUserID(userID string) (*aggregate.QueryUserPostRes, error)
}

type XHttp interface {
	Run()
	Register() XHttp
	Inject() XHttp
}

type Database interface {
	Connect()
	DB(ctx context.Context) interface{}
}

type Application struct {
	XHttp
	Database
}

func NewApplication() *Application {
	return &Application{}
}
func (a *Application) SetXHttp(x XHttp) *Application {
	a.XHttp = x
	return a
}

func (a *Application) SetDatabase(d Database) *Application {
	a.Database = d
	return a
}
