package main

import (
	_ "github.com/olongfen/go-ddd-hex/internal/adapter/repository" // 初始化存储库适配器
	_ "github.com/olongfen/go-ddd-hex/internal/adapter/xhttp/xgin" // 初始化http适配器
	"github.com/olongfen/go-ddd-hex/internal/application"
)

func main() {
	application.app.InjectServices().Run()
}
