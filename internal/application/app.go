package application

import (
	"context"
	"github.com/olongfen/go-ddd-hex/internal/domain/aggregate"
	"github.com/olongfen/go-ddd-hex/internal/domain/vo"
)

type UserInterface interface {
	ChangePassword(ctx context.Context, id string, oldPwd, newPwd string) error
	Get(ctx context.Context, id string) (res *vo.UserRes, err error)
}

type PostInterface interface {
	GetByUserID(ctx context.Context, userID string) (*aggregate.QueryUserPostRes, error)
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
	Http     XHttp
	Database Database
}

func NewApplication(h XHttp, db Database) *Application {
	return &Application{
		Http:     h,
		Database: db,
	}
}
