package weixin

import (
	"crypto/sha1"
	"encoding/hex"
	localTime "github.com/go-tron/local-time"
	"github.com/go-tron/random"
	"strconv"
)

type JsApiConfig struct {
	Debug     bool   `json:"debug"`
	AppId     string `json:"appId"`
	Timestamp int64  `json:"timestamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
}

func (wx *Weixin) GetJsApiConfig(url string) (*JsApiConfig, error) {

	jsApiTicket, err := wx.GetJsApiTicket()
	if err != nil {
		return nil, err
	}

	jsApiConfig := &JsApiConfig{
		AppId:     wx.AppId,
		Timestamp: localTime.Now().Unix(),
		NonceStr:  random.String(10),
		Signature: "",
	}
	var signStr = "jsapi_ticket=" + jsApiTicket.Ticket + "&noncestr=" + jsApiConfig.NonceStr + "&timestamp=" + strconv.FormatInt(jsApiConfig.Timestamp, 10) + "&url=" + url
	hash := sha1.New()
	hash.Write([]byte(signStr))
	jsApiConfig.Signature = hex.EncodeToString(hash.Sum(nil))
	return jsApiConfig, nil
}
