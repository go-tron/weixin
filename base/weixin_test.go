package base

import (
	"github.com/go-resty/resty/v2"
	localTime "github.com/go-tron/local-time"
	"github.com/go-tron/logger"
	"github.com/go-tron/redis"
	"sync"
	"testing"
	"time"
)

var base = Weixin{
	Config: &Config{
		AppId:  "wx6c8124f1fbafb1f3",
		Secret: "868c1eadf49047effa89e63021819b35",
		Redis: redis.New(&redis.Config{
			Addr:     "127.0.0.1:6379",
			Password: "GBkrIO9bkOcWrdsC",
		}),
		Logger: logger.NewZap("weixin", "info"),
	},
}

func TestGetAccessToken(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	var i = 0
	for i < 1000 {
		i++

		go func() {
			result, err := base.GetAccessToken()
			if err != nil {
				t.Fatal(err)
			}
			t.Log(localTime.Now(), "result", result)
		}()

	}
	wg.Wait()
}

func TestAccessTokenValid(t *testing.T) {

	result, err := resty.New().R().
		SetQueryParams(map[string]string{
			"type":         "jsapi",
			"access_token": "37_dKtjG74cDuBMUYibQQj0-66lHoVfKN1zmEddGWYF3di4gmPIsO0iev5uW4qHMxSL7vfaxKHniPRKSADwrCU3fJWEe7suaODufGL5w4kYT_qmTnZBvKWAGjkMls3zb-oEKSkMB6Zq17rf4kdCZCAbADAGJH",
		}).
		Get("https://api.weixin.qq.com/cgi-bin/ticket/getticket")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", result)
}

func TestGetJsApiTicket(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	result, err := base.GetJsApiTicket()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", result)

	time.Sleep(time.Second * 6)
	base.jsApiTicket = nil

	time.Sleep(time.Second * 6)
	result, err = base.GetJsApiTicket()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", result)

	wg.Wait()
}
