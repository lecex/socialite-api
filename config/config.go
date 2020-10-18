package config

import (
	"github.com/lecex/core/config"
	"github.com/lecex/core/env"
	PB "github.com/lecex/user/proto/permission"
)

// 	Conf 配置
// 	Service // 服务名称
//	Method // 方法
//	Auth // 是否认证授权
//	Policy // 是否认证权限
//	Name // 权限名称
//	Description // 权限解释
var Conf config.Config = config.Config{
	Name:    env.Getenv("MICRO_API_NAMESPACE", "go.micro.api.") + "socialite-api",
	Version: "latest",
	Service: map[string]string{
		"user":      env.Getenv("USER_SERVICE", "go.micro.srv.user"),
		"socialite": env.Getenv("SOCIALITE_SERVICE", "go.micro.srv.socialite"),
	},
	Permissions: []*PB.Permission{
		{Service: "socialite-api", Method: "Configs.Get", Auth: true, Policy: false, Name: "获取社会登录配置", Description: "获取社会登录配置。"},
		{Service: "socialite-api", Method: "Configs.Update", Auth: true, Policy: false, Name: "更新社会登录配置", Description: "更新社会登录配置。"},

		{Service: "socialite-api", Method: "Socialites.Auth", Auth: false, Policy: false, Name: "社会登录授权", Description: "社会登录授权。"},
		{Service: "socialite-api", Method: "Socialites.Register", Auth: false, Policy: false, Name: "社会登录注册", Description: "社会登录注册。"},
		{Service: "socialite-api", Method: "Socialites.AuthURL", Auth: false, Policy: false, Name: "登录URL获取", Description: "登录URL获取。"},

		{Service: "socialite-api", Method: "Users.List", Auth: true, Policy: true, Name: "社会登录用户列表", Description: "社会登录用户列表。"},
		{Service: "socialite-api", Method: "Users.Get", Auth: true, Policy: true, Name: "社会登录获取用户", Description: "社会登录获取用户。"},
		{Service: "socialite-api", Method: "Users.Create", Auth: true, Policy: true, Name: "社会登录创建用户", Description: "社会登录创建用户。"},
		{Service: "socialite-api", Method: "Users.Update", Auth: true, Policy: true, Name: "社会登录用户灯芯", Description: "社会登录用户灯芯。"},
		{Service: "socialite-api", Method: "Users.Delete", Auth: true, Policy: true, Name: "社会登录删除用户", Description: "社会登录删除用户。"},
		{Service: "socialite-api", Method: "Users.SelfBind", Auth: true, Policy: false, Name: "社会登录绑定用户", Description: "社会登录绑定用户。"},
		{Service: "socialite-api", Method: "Users.SelfUnbind", Auth: true, Policy: false, Name: "社会登录接触绑定", Description: "社会登录接触绑定。"},
	},
}
