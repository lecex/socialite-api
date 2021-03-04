package handler

import (
	context "context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	client "github.com/lecex/core/client"
	authSrvPB "github.com/lecex/user/proto/auth"
	userSrvPB "github.com/lecex/user/proto/user"
	"github.com/micro/go-micro/v2/metadata"

	pb "github.com/lecex/socialite-api/proto/socialite"

	"github.com/go-redis/redis"
)

// Socialite 配置结构
type Socialite struct {
	ServiceName string
	UserService string
	Redis       *redis.Client
}

func (srv *Socialite) getCache(code string, res *pb.Response) (err error) {
	// 读取缓存
	r, err := srv.Redis.Get("SocialiteCode_" + code).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil
		}
		return err
	}
	return json.Unmarshal([]byte(r), res)
}

func (srv *Socialite) setCache(code string, res *pb.Response, t time.Duration) (err error) {
	// 缓存30秒 防止重复请求
	value, _ := json.Marshal(res)
	return srv.Redis.Set("SocialiteCode_"+code, value, t*time.Second).Err()
}
func (srv *Socialite) getUsers(ctx context.Context, res *pb.Response) (err error) {
	// 获取关联用户token
	for _, user := range res.SocialiteUser.Users {
		reqAuthSrv := &authSrvPB.Request{
			User: &authSrvPB.User{
				Id: user.Id,
			},
		}
		resAuthSrv := &authSrvPB.Response{}
		err = client.Call(ctx, srv.UserService, "Auth.AuthById", reqAuthSrv, resAuthSrv)
		if err != nil {
			return err
		}
		user.Token = resAuthSrv.Token
	}
	return
}

// Auth 获取授权
func (srv *Socialite) Auth(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	res.SocialiteUser = &pb.SocialiteUser{
		Users: []*pb.User{},
	}
	err = srv.getCache(req.Socialite.Code, res) // 读取缓存
	if err != nil {
		return err
	}
	if res.SocialiteUser.Id == "" {
		err = client.Call(ctx, srv.ServiceName, "Socialites.Auth", req, res)
		if err != nil {
			return err
		}
		err = srv.setCache(req.Socialite.Code, res, 30) // 设置缓存 30秒
		if err != nil {
			return err
		}
		res.SocialiteUser.Content = ""
		res.SocialiteUser.OauthId = ""
	}
	if len(res.SocialiteUser.Users) > 0 {
		err = srv.getUsers(ctx, res) // 设置缓存 30秒
		if err != nil {
			return err
		}
		res.Valid = true
	} else {
		res.Valid = false
	}
	return err
}

// AuthURL 授权网址
func (srv *Socialite) AuthURL(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return client.Call(ctx, srv.ServiceName, "Socialites.AuthURL", req, res)
}

// RegisterUser 注册用户
func (srv *Socialite) RegisterUser(ctx context.Context, user *pb.User, captcha string) (res *userSrvPB.User, err error) {
	if user.Mobile != "" { // 验证手机
		err = srv.VerifyCaptcha(user.Mobile, captcha)
		if err != nil {
			return nil, err
		}
	}
	if user.Email != "" { // 验证邮箱
		err = srv.VerifyCaptcha(user.Email, captcha)
		if err != nil {
			return nil, err
		}
	}
	if user.Password == "" {
		user.Password = srv.getRandomString(8)
	}
	meta, _ := metadata.FromContext(ctx) // debug 无法获取 meta
	// 无用户先通过用户服务创建用户
	reqUserSrv := &userSrvPB.Request{
		User: &userSrvPB.User{
			Username: user.Username,
			Mobile:   user.Mobile, // 绑定手机必须是后端通过验证的
			Email:    user.Email,
			Password: user.Password,
			Name:     user.Name,
			Avatar:   user.Avatar,
			Origin:   meta["Service"],
		},
	}
	resUserSrv := &userSrvPB.Response{}
	err = client.Call(ctx, srv.UserService, "Users.Exist", reqUserSrv, resUserSrv)
	if resUserSrv.Valid {
		err = client.Call(ctx, srv.UserService, "Users.Get", reqUserSrv, resUserSrv)
	} else {
		err = client.Call(ctx, srv.UserService, "Users.Create", reqUserSrv, resUserSrv)
	}
	if err != nil {
		return res, err
	}
	return resUserSrv.User, err
}

// Register 授权后注册【可用于增加新账号】
func (srv *Socialite) Register(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	resAuth := &pb.Response{}
	err = srv.Auth(ctx, req, resAuth)
	if err != nil {
		return err
	}
	if resAuth.SocialiteUser.Id != "" {
		user, err := srv.RegisterUser(ctx, req.SocialiteUser.Users[0], req.Captcha)
		if err != nil {
			return err
		}
		req.SocialiteUser.Id = resAuth.SocialiteUser.Id
		req.SocialiteUser.Users[0].Id = user.Id
		err = client.Call(ctx, srv.ServiceName, "Socialites.BuildUser", req, res)
		if err != nil {
			return err
		}
	}
	return err
}

// getRandomString 生成随机字符串
func (srv *Socialite) getRandomString(length int64) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; int64(i) < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// VerifyCaptcha 校验验证码
func (srv *Socialite) VerifyCaptcha(addressee string, captcha string) (err error) {
	r, err := srv.Redis.Get("Captcha_" + addressee).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return fmt.Errorf("验证码已超时")
		}
		return err
	}
	if r != captcha {
		return fmt.Errorf("验证码错误")
	}
	srv.Redis.Set("Captcha_"+addressee, r, 1*time.Second).Err() // 1秒后自动过期
	return nil

}
