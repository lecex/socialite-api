package handler

import (
	"context"

	client "github.com/lecex/core/client"
	pb "github.com/lecex/socialite-api/proto/config"
)

// Config 配置结构
type Config struct {
	ServiceName string
}

// Get 获取配置
func (srv *Config) Get(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return client.Call(ctx, srv.ServiceName, "Configs.Get", req, res)
}

// Update 更新配置
func (srv *Config) Update(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return client.Call(ctx, srv.ServiceName, "Configs.Update", req, res)
}
