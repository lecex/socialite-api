package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/lecex/socialite-api/config"
	"github.com/lecex/socialite-api/providers/redis"

	"github.com/lecex/socialite-api/handler"
	socialitePB "github.com/lecex/socialite-api/proto/socialite"
)

var Conf = config.Conf

func TestSocialiteAuth(t *testing.T) {
	// redis := redis.NewClient()
	// req := &socialitePB.Request{
	// 	Socialite: &socialitePB.Socialite{
	// 		Driver: "miniprogram_wechat",
	// 		Code:   "041Gzu000q09dL1Ovc200ilHic2Gzu00",
	// 	},
	// }
	// res := &socialitePB.Response{}
	// h := handler.Socialite{Conf.Service["socialite"], Conf.Service["user"], redis}
	// err := h.Auth(context.TODO(), req, res)
	// fmt.Println(res, err)
	// t.Log("----Auth----", res, err)
}

func TestSocialiteRegister(t *testing.T) {
	redis := redis.NewClient()
	req := &socialitePB.Request{
		Socialite: &socialitePB.Socialite{
			Driver: "miniprogram_wechat",
			Code:   "031CUvFa1P6txA0zEXFa1Q3DxD2CUvF1",
		},
		SocialiteUser: &socialitePB.SocialiteUser{
			Users: []*socialitePB.User{
				&socialitePB.User{
					Mobile: "13954386521",
					Name:   "🤩 测试网名123*-🥰",
					Avatar: "https://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83ep1m5aI7y3WJAP6XIXN4e39124xvcjJoI9AM8QXjB9jN6VJpl3C32VNeXELnB71EWk8sE7zp32n4A/132",
				},
			},
		},
		Verify: "123456",
	}
	res := &socialitePB.Response{}
	fmt.Println("----Register----", res, res)
	h := handler.Socialite{Conf.Service["socialite"], Conf.Service["user"], redis}
	err := h.Register(context.TODO(), req, res)
	fmt.Println("----Register----", res, err)
}
