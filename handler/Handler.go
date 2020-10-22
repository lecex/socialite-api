package handler

import (
	"context"
	"time"

	"github.com/micro/go-micro/v2/util/log"

	server "github.com/micro/go-micro/v2/server"

	client "github.com/lecex/core/client"
	configPB "github.com/lecex/socialite-api/proto/config"
	socialitePB "github.com/lecex/socialite-api/proto/socialite"
	userPB "github.com/lecex/socialite-api/proto/user"

	"github.com/lecex/socialite-api/config"
	PB "github.com/lecex/user/proto/permission"
)

// Handler 注册方法
type Handler struct {
	Server server.Server
}

var Conf = config.Conf

// Register 注册
func (srv *Handler) Register() {
	configPB.RegisterConfigsHandler(srv.Server, &Config{Conf.Service["socialite"]})
	socialitePB.RegisterSocialitesHandler(srv.Server, &Socialite{Conf.Service["socialite"], Conf.Service["user"]})
	userPB.RegisterUsersHandler(srv.Server, &User{Conf.Service["socialite"]})

	go Sync() // 同步前端权限
}

// Sync 同步
func Sync() {
	time.Sleep(5 * time.Second)
	req := &PB.Request{
		Permissions: Conf.Permissions,
	}
	res := &PB.Response{}
	err := client.Call(context.TODO(), Conf.Service["user"], "Permissions.Sync", req, res)
	if err != nil {
		log.Log(err)
		Sync()
	}
}
