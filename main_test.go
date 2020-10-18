package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/lecex/socialite-api/handler"
)

func TestSocialiteAuth(t *testing.T) {
	req := &socialPB.Request{
		Socialite: &socialPB.Socialite{
			Driver: "miniprogram_wechat",
			Code:   "0914TD000B1drK1M73000dW94N34TD02",
		},
	}
	res := &socialPB.Response{}
	h := handler.Socialite{
	}
	err := h.Auth(context.TODO(), req, res)
	fmt.Println("----Auth----", res, err)
}

func TestSocialiteInfo(t *testing.T) {
}
