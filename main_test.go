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
			Code:   "0919Noll2zVnQ54Utill2hwJjN19Nol0",
		},
	}
	res := &socialitePB.Response{}
	h := handler.Socialite{Conf.Service["socialite"], Conf.Service["user"]}
	err := h.Auth(context.TODO(), req, res)
	fmt.Println("----Auth----", res, err)
}

func TestSocialiteInfo(t *testing.T) {
}
