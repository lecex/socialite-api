package handler

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
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
	if res.SocialiteUser.Id != "" && len(res.SocialiteUser.Users) == 0 {
		session := uuid.NewV4().String()

		redis := redis.NewClient()
		value, _ := json.Marshal(res.SocialiteUser)
		// 过期时间默认 30 分钟
		err = redis.Set(session, value, 30*time.Minute).Err()

		if err != nil {
			return err
		}
		res.Session = session
	}
	res.SocialiteUser.Content = ""
	res.SocialiteUser.OauthId = ""
	return err
}

// AuthURL 授权网址
func (srv *Socialite) AuthURL(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	return client.Call(ctx, srv.ServiceName, "Socialites.AuthURL", req, res)
}

// Register 授权后注册【可用于增加新账号】
func (srv *Socialite) Register(ctx context.Context, req *pb.Request, res *pb.Response) (err error) {
	mobile := ""
	// 过期时间默认 30 分钟
	redis := redis.NewClient()
	socialiteUser, err := redis.Get(req.Session).Result()
	u := &pb.SocialiteUser{}
	err = json.Unmarshal([]byte(socialiteUser), u)
	if u.Id == "" && err != err {
		return fmt.Errorf("Session 未查询到相关信息")
	}

	if u.Origin == "miniprogram_"+req.Miniprogram.Type {
		if req.Miniprogram.Type == "wechat" {
			mobile, err = srv.getWechatMobile(u, req.Miniprogram)
			if err != nil {
				return err
			}
		}
	}
	if len(req.SocialiteUser.Users) > 0 { // 前端传入的用户数据
		user := req.SocialiteUser.Users[0]
		// 禁止直接传入手机邮箱
		user.Email = ""
		// 无用户先通过用户服务创建用户
		reqUserSrv := &userSrvPB.Request{
			User: &userSrvPB.User{
				Username: user.Username,
				Mobile:   mobile, // 绑定手机必须是后端通过验证的
				Email:    user.Email,
				Password: user.Password,
				Name:     user.Name,
				Avatar:   user.Avatar,
			},
		}
		resUserSrv := &userSrvPB.Response{}
		err = client.Call(context.TODO(), srv.UserService, "Users.Create", reqUserSrv, resUserSrv)
		if err != nil {
			return err
		}
		// if resUserSrv.Valid {
		// 	u.Users = append(u.Users, &userPB.User{
		// 		Id: resUserSrv.User.Id,
		// 	})
		// }
	} else {
		err = fmt.Errorf("未收到用户注册信息")
	}

	// err = client.Call(ctx, srv.ServiceName, "Users.Get", req, res)
	// if err != nil {
	// 	return err
	// }
	// if len(req.SocialiteUser.Users) > 0 {
	// 	for _, user := range req.SocialiteUser.Users {
	// 		// 无用户先通过用户服务创建用户
	// 		reqUserSrv := &userSrvPB.Request{
	// 			User: &userSrvPB.User{
	// 				Username: user.Username,
	// 				Mobile:   user.Mobile,
	// 				Email:    user.Email,
	// 				Password: user.Password,
	// 				Name:     user.Name,
	// 				Avatar:   user.Avatar,
	// 			},
	// 		}
	// 		resUserSrv := &userSrvPB.Response{}
	// 		err = client.Call(context.TODO(), srv.ServiceName, "Users.Create", reqUserSrv, resUserSrv)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		// if resUserSrv.Valid {
	// 		// 	u.Users = append(u.Users, &userPB.User{
	// 		// 		Id: resUserSrv.User.Id,
	// 		// 	})
	// 		// }
	// 	}
	// } else {
	// 	err = fmt.Errorf("未收到用户注册信息")
	// }
	// u.CreatedAt = ""
	// u.UpdatedAt = ""
	// _, err = srv.Repo.Update(u)
	// fmt.Println("---Register---", u)
	return err
}

// getWechatMobile 获取微信手机
func (srv *Socialite) getWechatMobile(u *pb.SocialiteUser, m *pb.Miniprogram) (mobile string, err error) {
	c := map[string]string{}
	err = json.Unmarshal([]byte(u.Content), c)
	if err != err {
		return "", fmt.Errorf("微信配置信息解析错误")
	}
	fmt.Println(1, m.EncryptedData, c["session_key"], m.Iv)
	info, err := srv.sessionInfo(m.EncryptedData, c["session_key"], m.Iv)
	mobile = info["phoneNumber"].(string)
	if err != nil {
		return
	}
	return
}

// sessionInfo 解密小程序会话加密信息
func (srv *Socialite) sessionInfo(encryptedData, sessionKey, iv string) (info map[string]interface{}, err error) {
	cipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return
	}
	aesKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return
	}
	aesIv, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return
	}

	const (
		BLOCK_SIZE = 32             // PKCS#7
		BLOCK_MASK = BLOCK_SIZE - 1 // BLOCK_SIZE 为 2^n 时, 可以用 mask 获取针对 BLOCK_SIZE 的余数
	)
	if len(cipherText) < BLOCK_SIZE {
		err = fmt.Errorf("the length of ciphertext too short: %d", len(cipherText))
		return
	}
	plaintext := make([]byte, len(cipherText)) // len(plaintext) >= BLOCK_SIZE
	// 解密
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		panic(err)
	}
	mode := cipher.NewCBCDecrypter(block, aesIv)
	mode.CryptBlocks(plaintext, cipherText)
	// PKCS#7 去除补位
	amountToPad := int(plaintext[len(plaintext)-1])
	if amountToPad < 1 || amountToPad > BLOCK_SIZE {
		err = fmt.Errorf("the amount to pad is incorrect: %d", amountToPad)
		return
	}
	plaintext = plaintext[:len(plaintext)-amountToPad]
	// 反拼接
	// len(plaintext) == 16+4+len(rawXMLMsg)+len(appId)
	if len(plaintext) <= 20 {
		err = fmt.Errorf("plaintext too short, the length is %d", len(plaintext))
		return
	}
	if err != nil {
		return
	}
	if err = json.Unmarshal(plaintext, &info); err != nil {
		return
	}
	return
}
