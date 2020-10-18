package handler

import (
	"context"

	client "github.com/lecex/core/client"
	pb "github.com/lecex/socialite-api/proto/socialite"
)

// Socialite 配置结构
type Socialite struct {
	ServiceName string
}

// Auth 获取授权
func (srv *Socialite) Auth(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return client.Call(ctx, srv.ServiceName, "Socialites.Auth", req, res)
}

// Register 授权后注册【可用于增加新账号】
func (srv *Socialite) Register(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return client.Call(ctx, srv.ServiceName, "Socialites.Register", req, res)
}

// AuthURL 授权网址
func (srv *Socialite) AuthURL(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return client.Call(ctx, srv.ServiceName, "Socialites.AuthURL", req, res)
}
