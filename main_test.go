package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/lecex/socialite-api/config"

	"github.com/lecex/socialite-api/handler"
	socialitePB "github.com/lecex/socialite-api/proto/socialite"
)

var Conf = config.Conf

func TestSocialiteAuth(t *testing.T) {
	req := &socialitePB.Request{
		Socialite: &socialitePB.Socialite{
			Driver: "miniprogram_wechat",
			Code:   "051P42ll2rGHQ54VwJll2ChWci0P42lM",
		},
	}
	res := &socialitePB.Response{}
	h := handler.Socialite{Conf.Service["socialite"], Conf.Service["user"]}
	err := h.Auth(context.TODO(), req, res)
	t.Log("----Auth----", res, err)
}

func TestSocialiteRegister(t *testing.T) {
	req := &socialitePB.Request{
		Session: "eceb29e7-2851-4b3b-ace4-35af789c5304",
		SocialiteUser: &socialitePB.SocialiteUser{
			Users: []*socialitePB.User{
				&socialitePB.User{
					Username: "bvbv011",
					Mobile:   "19054386521",
					Email:    "bigrocs1@qq.com",
					Password: "123456",
					Name:     "BigRocs",
					Avatar:   "https://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83ep1m5aI7y3WJAP6XIXN4e39124xvcjJoI9AM8QXjB9jN6VJpl3C32VNeXELnB71EWk8sE7zp32n4A/132",
				},
			},
		},
		Miniprogram: &socialitePB.Miniprogram{
			EncryptedData: "eBtIKhQhqwhMyZVpOvjIcKLJYWsK4iXUqAzwCzYk3sUQ4vARXrs4W2m+w5GT4Zu5wjYbBn2Vcmg2YVG8PuiTMUJjMlv/2i1KnlNlJW6LpMzdz505kbVKSrBijqsrJqzVzA1CeLDuB73WkPugiBdcia6vMhQhTGGlUaC4CYWWC7a+P8YUMOWe6VUSumhgH/PGuIBXpRT4Nc3Tdp6ffX8yDg==",
			Iv:            "MDxsuRIKYDLLtt0qUefbzA==",
			Type:          "wechat",
		},
	}
	res := &socialitePB.Response{}
	fmt.Println("----Register----", res, res)
	h := handler.Socialite{Conf.Service["socialite"], Conf.Service["user"]}
	err := h.Register(context.TODO(), req, res)
	fmt.Println("----Register----", res, err)
}
