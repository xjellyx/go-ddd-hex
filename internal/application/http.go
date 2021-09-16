package application

import (
	"context"
	"net/http"
)

// IController 这里是定义一个接口
type IController interface {
	// Router 这个传参就是脚手架主程
	Router(server XHttp, isGroup bool)
}

// XHttp  http 接口
type XHttp interface {
	Run()                                                                                    // http开启服务接口
	Init() XHttp                                                                             // 初始化http引擎和使用中间件
	Use(middlewares ...func(ctx context.Context)) XHttp                                      // 使用中间件
	GetMux() http.Handler                                                                    // 获取http handler
	Handle(httpMethod, relativePath string, f func(ctx context.Context), isGroup bool) XHttp // 开启http路由接口
	Route(...IController) XHttp
	RouterGroup(string, ...IController) XHttp
}

type UserHandler interface {
	SetService(serviceInterface UserServiceInterface) // 设置服务
	Get(ctx context.Context)
	Register(ctx context.Context)
	ChangePasswd(ctx context.Context)
	IController
}

type PostHandler interface {
	SetService(serviceInterface PostServiceInterface) // 设置服务
	GetByUserID(c context.Context)
	IController
}
