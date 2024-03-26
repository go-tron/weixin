package weixin

import (
	"github.com/go-tron/logger"
	"github.com/go-tron/redis"
	"sync"
	"testing"
	"time"
)

var account = Weixin{
	Config: &Config{
		Username:         "h8ex8ug5yvsnbqz9",
		Password:         "f7d1168d49fafa6fb41edbc67a29ce0d",
		BaseUrl:          "https://weixin-config.eioos.com/config",
		AppId:            "wx6c8124f1fbafb1f3",
		Secret:           "868c1eadf49047effa89e63021819b35",
		OAuthRedirectUri: "https://weixin.eioos.com/oauth/return?uri=",
		Logger:           logger.NewZap("weixin", "info"),
		Redis: redis.New(&redis.Config{
			Addr:     "127.0.0.1:6379",
			Password: "GBkrIO9bkOcWrdsC",
		}),
	},
}

func TestGetAccessToken(t *testing.T) {
	result, err := account.GetAccessToken()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", result)
}

func TestGetJsApiTicket(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	result, err := account.GetJsApiTicket()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", result)

	time.Sleep(time.Second * 6)
	account.jsApiTicket = nil

	time.Sleep(time.Second * 6)
	result, err = account.GetJsApiTicket()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", result)

	wg.Wait()
}

func TestGetJsApiConfig(t *testing.T) {
	result, err := account.GetJsApiConfig("http://192.168.1.101:17000")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", result)
}

func TestGetOAuthCode(t *testing.T) {
	result, err := account.GetOAuthCode(&OAuthCodeReq{
		Uri:   "http://192.168.1.101:17000",
		Scope: ScopeUserInfo,
		State: "",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", result)
}

func TestGetOAuthAccessToken(t *testing.T) {
	result, err := account.GetOAuthAccessToken("0619PF0w3CBXZU2jYz2w3dMTwh19PF0Y")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", result)
}

func TestGetOAuthUserInfo(t *testing.T) {
	result, err := account.GetOAuthUserInfo(&OAuthUserInfoReq{
		AccessToken: "38_mSl_DWZJS3SV5HJ3EW_Y8XxRdguq7NXr4VFYgUNFAKWT2X8Ma_bfUk5BgTsaQLnkKg7mSYpMVejvVdsOOAMTDEizL4Rr2EGiEpxuqMXiZNU",
		OpenId:      "oasi95rPit953LHRYfaifGnTuqgs",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", result)
}

func TestGetUserInfo(t *testing.T) {
	result, err := account.GetUserInfo("oasi95rPit953LHRYfaifGnTuqgs")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", result)
}

func TestSendStartTemplate(t *testing.T) {
	err := account.SendTemplate(&TemplateReq{
		OpenId:     "oasi95rPit953LHRYfaifGnTuqgs",
		TemplateId: "0sWTgTRNs91psQ8PSkREh8-4h1ziHIQsvmdfyqTc6Qk",
		Url:        "https://app.eioos.com",
		Data: map[string]interface{}{
			"first": map[string]string{
				"value": "first",
			},
			"keyword1": map[string]string{
				"value": "keyword1",
			},
			"keyword2": map[string]string{
				"value": "keyword2",
			},
			"keyword3": map[string]string{
				"value": "keyword3",
			},
			"keyword4": map[string]string{
				"value": "keyword4",
			},
			"remark": map[string]string{
				"value": "remark",
				"color": "#173177",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", "success")
}

func TestSendFinishTemplate(t *testing.T) {
	err := account.SendTemplate(&TemplateReq{
		OpenId:     "opZow6IEeQkp1y03HWfjZW0njPUE",
		TemplateId: "je8cWU_cvGgDMlPMAK87DmJ9p4In11GnHRkaIW19IbU",
		Url:        "https://app.eioos.com",
		Data: map[string]interface{}{
			"first": map[string]string{
				"value": "first",
			},
			"keyword1": map[string]string{
				"value": "keyword1",
			},
			"keyword2": map[string]string{
				"value": "keyword2",
			},
			"keyword3": map[string]string{
				"value": "keyword3",
			},
			"keyword4": map[string]string{
				"value": "keyword4",
			},
			"remark": map[string]string{
				"value": "remark",
				"color": "#173177",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", "success")
}

func TestBatchGetMaterial(t *testing.T) {
	result, err := account.BatchGetMaterial(&BatchGetMaterialReq{
		Type:   "news",
		Offset: "0",
		Count:  "20",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", result)
}

func TestMenuCreate(t *testing.T) {
	err := account.MenuCreate(map[string]interface{}{
		"button": []map[string]interface{}{
			{
				"type": "view",
				"name": "name",
				"url":  "https://wx.nonsense.eioos.com",
			},
			{
				"name": "name",
				"sub_button": []map[string]interface{}{
					{
						"type":     "view_limited",
						"name":     "name",
						"media_id": "CP7F8oYVrBC4X55wK1Bnt_vbWo80GFmT-WRmntZ3hmU",
					},
					{
						"type":     "view_limited",
						"name":     "name",
						"media_id": "CP7F8oYVrBC4X55wK1BntxYZOPL7WCX6ofjKP4muqeQ",
					},
				},
			},
			{
				"name": "name",
				"sub_button": []map[string]interface{}{
					{
						"type":     "view_limited",
						"name":     "name",
						"media_id": "CP7F8oYVrBC4X55wK1Bnt-CuzcyCfFDrnHL3l0YALzw",
					},
					{
						"type":     "view_limited",
						"name":     "name",
						"media_id": "CP7F8oYVrBC4X55wK1BntwPyWyMHSUB1BckeNDH2txo",
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", "succeed")
}

func TestMenuDelete(t *testing.T) {
	err := account.MenuDelete()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", "succeed")
}

func TestSendUniformMessage(t *testing.T) {
	err := account.SendUniformMessage(&UniformMessageReq{
		OpenId:     "oasi95rPit953LHRYfaifGnTuqgs",
		TemplateId: "0sWTgTRNs91psQ8PSkREh8-4h1ziHIQsvmdfyqTc6Qk",
		Url:        "https://app.eioos.com",
		Data: map[string]interface{}{
			"first": map[string]string{
				"value": "first",
			},
			"keyword1": map[string]string{
				"value": "keyword1",
			},
			"keyword2": map[string]string{
				"value": "keyword2",
			},
			"keyword3": map[string]string{
				"value": "keyword3",
			},
			"keyword4": map[string]string{
				"value": "keyword4",
			},
			"remark": map[string]string{
				"value": "remark",
				"color": "#173177",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", "success")
}
