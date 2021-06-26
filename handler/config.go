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

// All 配置列表
func (srv *Config) All(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return client.Call(ctx, srv.ServiceName, "Configs.All", req, res)
}

// List 配置列表
func (srv *Config) List(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return client.Call(ctx, srv.ServiceName, "Configs.List", req, res)
}

// Get 获取配置
func (srv *Config) Get(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return client.Call(ctx, srv.ServiceName, "Configs.Get", req, res)
}

// Create 创建配置
func (srv *Config) Create(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return client.Call(ctx, srv.ServiceName, "Configs.Create", req, res)
}

// Update 更新配置
func (srv *Config) Update(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return client.Call(ctx, srv.ServiceName, "Configs.Update", req, res)
}

// Delete 删除配置
func (srv *Config) Delete(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return client.Call(ctx, srv.ServiceName, "Configs.Delete", req, res)
}
