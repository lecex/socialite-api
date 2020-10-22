package handler

import (
	"context"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"

	client "github.com/lecex/core/client"
	authSrvPB "github.com/lecex/user/proto/auth"
	userSrvPB "github.com/lecex/user/proto/user"

	pb "github.com/lecex/socialite-api/proto/socialite"
	"github.com/lecex/socialite-api/providers/redis"
)

// Socialite 配置结构
type Socialite struct {
	ServiceName string
	UserService string
}

// Auth 获取授权
func (srv *Socialite) Auth(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	res.SocialiteUser = &pb.SocialiteUser{
		Users: []*pb.User{},
	}
	err = client.Call(ctx, srv.ServiceName, "Socialites.Auth", req, res)
	// 获取关联用户token
	for _, user := range res.SocialiteUser.Users {
		reqAuthSrv := &authSrvPB.Request{
			User: &authSrvPB.User{
				Id: user.Id,
			},
		}
		resAuthSrv := &authSrvPB.Response{}
		err = client.Call(context.TODO(), srv.UserService, "Auth.AuthById", reqAuthSrv, resAuthSrv)
		if err != nil {
			return err
		}
		res.SocialiteUser.Users = append(res.SocialiteUser.Users, &pb.User{
			Id:    user.Id,
			Name:  resAuthSrv.User.Name,
			Token: resAuthSrv.Token,
		})
	}
	if len(res.SocialiteUser.Users) == 0 {
		key := uuid.NewV4().String()
		redis := redis.NewClient()
		// 过期时间默认30分钟
		err = redis.Set(key, res.SocialiteUser, 30*time.Minute).Err()
		if err != nil {
			return err
		}
		res.Uuid = key
	}
	fmt.Println(res)
	return err
}

// AuthURL 授权网址
func (srv *Socialite) AuthURL(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return client.Call(ctx, srv.ServiceName, "Socialites.AuthURL", req, res)
}

// Register 授权后注册【可用于增加新账号】
func (srv *Socialite) Register(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	err = client.Call(ctx, srv.ServiceName, "Users.Get", req, res)
	if err != nil {
		return err
	}
	if len(req.SocialiteUser.Users) > 0 {
		for _, user := range req.SocialiteUser.Users {
			// 无用户先通过用户服务创建用户
			reqUserSrv := &userSrvPB.Request{
				User: &userSrvPB.User{
					Username: user.Username,
					Mobile:   user.Mobile,
					Email:    user.Email,
					Password: user.Password,
					Name:     user.Name,
					Avatar:   user.Avatar,
				},
			}
			resUserSrv := &userSrvPB.Response{}
			err = client.Call(context.TODO(), srv.ServiceName, "Users.Create", reqUserSrv, resUserSrv)
			if err != nil {
				return err
			}
			// if resUserSrv.Valid {
			// 	u.Users = append(u.Users, &userPB.User{
			// 		Id: resUserSrv.User.Id,
			// 	})
			// }
		}
	} else {
		err = fmt.Errorf("未收到用户注册信息")
	}
	// u.CreatedAt = ""
	// u.UpdatedAt = ""
	// _, err = srv.Repo.Update(u)
	// fmt.Println("---Register---", u)
	return err
}
