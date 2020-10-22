package handler

import (
	"context"

	client "github.com/lecex/core/client"
	pb "github.com/lecex/socialite-api/proto/user"
)

// User 配置结构
type User struct {
	ServiceName string
}

// List 获取用户列表
func (srv *User) List(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return client.Call(ctx, srv.ServiceName, "Users.List", req, res)
}

// Get 根据 唯一 获取用户
func (srv *User) Get(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return client.Call(ctx, srv.ServiceName, "Users.Get", req, res)
}

// Create 创建用户
func (srv *User) Create(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return client.Call(ctx, srv.ServiceName, "Users.AuthURL", req, res)
}

// Update 更新用户
func (srv *User) Update(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return client.Call(ctx, srv.ServiceName, "Users.Update", req, res)
}

// Delete 删除用户
func (srv *User) Delete(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return client.Call(ctx, srv.ServiceName, "Users.Delete", req, res)
}

// SelfBind 绑定用户
func (srv *User) SelfBind(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return client.Call(ctx, srv.ServiceName, "Users.SelfBind", req, res)
}

// SelfUnbind 解除绑定
func (srv *User) SelfUnbind(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return client.Call(ctx, srv.ServiceName, "Users.SelfUnbind", req, res)
}
