package base

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/go-tron/config"
	"github.com/go-tron/logger"
	"github.com/go-tron/redis"
	"sync"
	"time"
)

func NewWithConfig(c *config.Config, redis *redis.Redis) *Weixin {
	return New(&Config{
		AppId:  c.GetString("weixin.appId"),
		Secret: c.GetString("weixin.secret"),
		Redis:  redis,
		Logger: logger.NewZapWithConfig(c, "weixin-base", "info"),
	})
}

func New(c *Config) *Weixin {

	if c == nil {
		panic("config 必须设置")
	}
	if c.AppId == "" {
		panic("AppId 必须设置")
	}
	if c.Secret == "" {
		panic("Secret 必须设置")
	}
	if c.Logger == nil {
		panic("Logger 必须设置")
	}
	if c.Redis == nil {
		panic("Redis 必须设置")
	}

	return &Weixin{
		Config: c,
	}
}

type Config struct {
	AppId  string        `json:"appId"`
	Secret string        `json:"secret"`
	Logger logger.Logger `json:"logger"`
	Redis  *redis.Redis  `json:"redis"`
}

type AccessToken struct {
	AccessToken string       `json:"access_token"`
	ExpiresIn   int64        `json:"expires_in"`
	ticker      *time.Ticker `json:"-"`
}

type JsApiTicket struct {
	Ticket    string       `json:"ticket"`
	ExpiresIn int64        `json:"expires_in"`
	ticker    *time.Ticker `json:"-"`
}

type Weixin struct {
	*Config
	accessTokenLock sync.Mutex
	jsApiTicketLock sync.Mutex
	accessToken     *AccessToken `json:"accessToken"`
	jsApiTicket     *JsApiTicket `json:"jsApiTicket"`
}

const (
	AccessTokenPrefix = "wx-access-token:"
	JsApiTicketPrefix = "wx-jsapi-ticket:"
)

type AccessTokenRes struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (wx *Weixin) ClearAccessToken() {
	if wx.accessToken == nil {
		return
	}
	wx.accessToken.AccessToken = ""
	wx.accessToken.ExpiresIn = 0
	if wx.accessToken.ticker != nil {
		wx.accessToken.ticker.Stop()
	}
}

func (wx *Weixin) SetAccessToken(accessToken string, expiresIn int64) {

	wx.accessToken = &AccessToken{
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
		ticker:      time.NewTicker(time.Second),
	}

	go func() {
		for wx.accessToken.ExpiresIn > 0 {
			<-wx.accessToken.ticker.C

			if wx.accessToken == nil {
				break
			}

			wx.accessToken.ExpiresIn--
			wx.Logger.Debug(fmt.Sprintf("accessToken.ExpiresIn:%d", wx.accessToken.ExpiresIn), wx.Logger.Field("appId", wx.AppId))

			if wx.accessToken.ExpiresIn <= 0 {
				wx.ClearAccessToken()
				break
			}
		}
	}()
}

func (wx *Weixin) ClearJsApiTicket() {
	if wx.jsApiTicket == nil {
		return
	}
	wx.jsApiTicket.Ticket = ""
	wx.jsApiTicket.ExpiresIn = 0
	if wx.jsApiTicket.ticker != nil {
		wx.jsApiTicket.ticker.Stop()
	}
}

func (wx *Weixin) SetJsApiTicket(jsApiTicket string, expiresIn int64) {

	wx.jsApiTicket = &JsApiTicket{
		Ticket:    jsApiTicket,
		ExpiresIn: expiresIn,
		ticker:    time.NewTicker(time.Second),
	}

	go func() {
		for wx.jsApiTicket.ExpiresIn > 0 {
			<-wx.jsApiTicket.ticker.C
			if wx.jsApiTicket == nil {
				break
			}

			wx.jsApiTicket.ExpiresIn--
			wx.Logger.Debug(fmt.Sprintf("jsApiTicket.ExpiresIn:%d", wx.jsApiTicket.ExpiresIn), wx.Logger.Field("appId", wx.AppId))

			if wx.jsApiTicket.ExpiresIn <= 0 {
				wx.ClearJsApiTicket()
				break
			}
		}
	}()
}

func (wx *Weixin) GetAccessToken() (a *AccessToken, err error) {

	wx.accessTokenLock.Lock()
	defer func() {
		if err != nil {
			wx.Logger.Error("GetAccessToken", wx.Logger.Field("error", err), wx.Logger.Field("appId", wx.AppId))
		}
		wx.accessTokenLock.Unlock()
	}()

	if wx.accessToken != nil && wx.accessToken.AccessToken != "" {
		wx.Logger.Debug("GetAccessToken from application", wx.Logger.Field("appId", wx.AppId))
		return wx.accessToken, nil
	}

	accessToken, err := wx.Redis.Get(context.Background(), AccessTokenPrefix+wx.AppId).Result()
	ttl, err := wx.Redis.TTL(context.Background(), AccessTokenPrefix+wx.AppId).Result()
	if accessToken != "" && ttl > 0 {
		wx.SetAccessToken(accessToken, int64(ttl/time.Second))
		wx.Logger.Debug("GetAccessToken from redis", wx.Logger.Field("appId", wx.AppId))
		return wx.accessToken, nil
	}

	resp, err := resty.New().R().
		SetQueryParams(map[string]string{
			"grant_type": "client_credential",
			"appid":      wx.AppId,
			"secret":     wx.Secret,
		}).
		Get("https://api.weixin.qq.com/cgi-bin/token")

	if err != nil {
		return nil, err
	}

	wx.Logger.Debug("GetAccessToken", wx.Logger.Field("response", resp.Body()), wx.Logger.Field("appId", wx.AppId))

	var res = &AccessTokenRes{}
	if err := json.Unmarshal(resp.Body(), res); err != nil {
		return nil, err
	}

	if res.AccessToken == "" || res.ExpiresIn == 0 {
		if res.ErrMsg != "" {
			return nil, errors.New(fmt.Sprintf("(%d)%s", res.ErrCode, res.ErrMsg))
		} else {
			return nil, errors.New("request failed")
		}
	}

	expireIn := time.Second * time.Duration(res.ExpiresIn)

	wx.SetAccessToken(res.AccessToken, res.ExpiresIn)

	wx.Redis.Set(context.Background(), AccessTokenPrefix+wx.AppId, res.AccessToken, expireIn).Result()

	wx.Logger.Debug("GetAccessToken from request", wx.Logger.Field("appId", wx.AppId))
	return wx.accessToken, nil
}

type JsApiTicketRes struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"`
}

func (wx *Weixin) GetJsApiTicket() (j *JsApiTicket, err error) {

	wx.jsApiTicketLock.Lock()
	defer func() {
		if err != nil {
			wx.Logger.Error("GetJsApiTicket", wx.Logger.Field("error", err), wx.Logger.Field("appId", wx.AppId))
		}
		wx.jsApiTicketLock.Unlock()
	}()

	if wx.jsApiTicket != nil && wx.jsApiTicket.Ticket != "" {
		wx.Logger.Debug("GetJsApiTicket from application", wx.Logger.Field("appId", wx.AppId))
		return wx.jsApiTicket, nil
	}

	jsApiTicket, err := wx.Redis.Get(context.Background(), JsApiTicketPrefix+wx.AppId).Result()
	ttl, err := wx.Redis.TTL(context.Background(), JsApiTicketPrefix+wx.AppId).Result()
	if jsApiTicket != "" && ttl > 0 {
		wx.SetJsApiTicket(jsApiTicket, int64(ttl/time.Second))
		wx.Logger.Debug("GetJsApiTicket from redis", wx.Logger.Field("appId", wx.AppId))
		return wx.jsApiTicket, nil
	}

	accessToken, err := wx.GetAccessToken()
	if err != nil {
		return nil, err
	}

	resp, err := resty.New().R().
		SetQueryParams(map[string]string{
			"type":         "jsapi",
			"access_token": accessToken.AccessToken,
		}).
		Get("https://api.weixin.qq.com/cgi-bin/ticket/getticket")

	if err != nil {
		return nil, err
	}

	wx.Logger.Debug("GetJsApiTicket", wx.Logger.Field("response", resp.Body()), wx.Logger.Field("appId", wx.AppId))

	var res = &JsApiTicketRes{}
	if err := json.Unmarshal(resp.Body(), res); err != nil {
		return nil, err
	}

	if res.Ticket == "" || res.ExpiresIn == 0 {
		if res.ErrCode == 40001 || res.ErrCode == 42001 {
			wx.ClearAccessToken()
			_, err := wx.Redis.Del(context.Background(), AccessTokenPrefix+wx.AppId).Result()
			if err != nil {
				return nil, err
			}
			return wx.GetJsApiTicket()
		}

		if res.ErrMsg != "" {
			return nil, errors.New(fmt.Sprintf("(%d)%s", res.ErrCode, res.ErrMsg))
		} else {
			return nil, errors.New("request failed")
		}
	}

	expireIn := time.Second * time.Duration(res.ExpiresIn)

	wx.SetJsApiTicket(res.Ticket, res.ExpiresIn)

	wx.Redis.Set(context.Background(), JsApiTicketPrefix+wx.AppId, res.Ticket, expireIn).Result()

	wx.Logger.Debug("GetJsApiTicket from request", wx.Logger.Field("appId", wx.AppId))
	return wx.jsApiTicket, nil
}
