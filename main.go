package main

import (
	"fmt"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"
	"github.com/yilei-pixocial/kubeapi/pkg/sysinit"
	"github.com/yilei-pixocial/kubeapi/router"
	"github.com/yilei-pixocial/kubeapi/router/middleware"
	"log"

	"github.com/iris-contrib/swagger/v12"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	_ "github.com/yilei-pixocial/kubeapi/docs"
)

var App *iris.Application

// @title kubeapi API
// @version 1.0
// @host localhost:8888
// @BasePath /
func main() {
	// 初始化App
	App = iris.New()

	// 初始化
	sysinit.InitConf()
	InitMiddleware()
	sysinit.InitLogger()
	sysinit.InitCron()
	router.SetRoutes(App)

	config := &swagger.Config{
		URL: fmt.Sprintf("http://localhost:%s/swagger/doc.json", sysinit.GCF.UString("server.port")), //The url pointing to API definition
	}
	// use swagger middleware to
	App.Get("/swagger/*any", swagger.CustomWrapHandler(config, swaggerFiles.Handler))

	// 启动
	run := App.Run(iris.Addr(":"+sysinit.GCF.UString("server.port")), iris.WithCharset("UTF-8"))
	log.Fatal(run)
}

func InitMiddleware() {
	// 设置未知异常捕获
	App.Use(recover.New())
	// 设置限流器
	middleware.InitHttpLimiter()
}
