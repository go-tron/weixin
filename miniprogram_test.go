package weixin

import (
	"github.com/go-tron/logger"
	"github.com/go-tron/redis"
	"testing"
)

var maccount = Weixin{
	Config: &Config{
		Username:         "h8ex8ug5yvsnbqz9",
		Password:         "f7d1168d49fafa6fb41edbc67a29ce0d",
		BaseUrl:          "http://localhost:7118/config",
		AppId:            "wx8f3c9c583c35ebd5",
		Secret:           "2b41a03dfed156b35280be3d4592ba50",
		OAuthRedirectUri: "https://weixin.eioos.com/oauth/return?uri=",
		Logger:           logger.NewZap("weixin", "info"),
		Redis: redis.New(&redis.Config{
			Addr:     "127.0.0.1:6379",
			Password: "GBkrIO9bkOcWrdsC",
		}),
	},
}

func TestGetUserPhoneNumber(t *testing.T) {
	result, err := maccount.GetUserPhoneNumber("7785a45674df481fa13d895525e6d2d87c8b33f3f69e6751a755fb199af0f924")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", result)
}
