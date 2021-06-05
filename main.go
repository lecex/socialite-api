package main

import (
	// 公共引入

	_ "github.com/lecex/core/plugins"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/util/log"

	"github.com/lecex/socialite-api/config"
	"github.com/lecex/socialite-api/handler"
)

func main() {
	var Conf = config.Conf
	service := micro.NewService(
		micro.Name(Conf.Name),
		micro.Version(Conf.Version),
		micro.WrapHandler(Conf.Middleware().Wrapper), //验证权限
	)
	service.Init()
	// 注册服务
	h := handler.Handler{
		Server: service.Server(),
	}
	h.Register()
	// Run the server
	log.Fatal("serviser run ... Version:" + Conf.Version)
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
