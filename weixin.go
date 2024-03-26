package weixin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/go-tron/config"
	"github.com/go-tron/logger"
	"github.com/go-tron/redis"
	"github.com/go-tron/weixin/base"
	"time"
)

func NewWithConfig(c *config.Config, redis *redis.Redis) *Weixin {
	return New(&Config{
		Username:         c.GetString("application.id"),
		Password:         c.GetString("application.secret"),
		BaseUrl:          c.GetString("weixin.baseUrl"),
		AppId:            c.GetString("weixin.appId"),
		Secret:           c.GetString("weixin.secret"),
		OAuthRedirectUri: c.GetString("weixin.oAuthRedirectUri"),
		Redis:            redis,
		Logger:           logger.NewZapWithConfig(c, "weixin", "info"),
	})
}

func New(c *Config) *Weixin {

	if c == nil {
		panic("config 必须设置")
	}
	if c.Username == "" {
		panic("Username 必须设置")
	}
	if c.Password == "" {
		panic("Password 必须设置")
	}
	if c.BaseUrl == "" {
		panic("BaseUrl 必须设置")
	}
	if c.Name == "" {
		panic("Name 必须设置")
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
	accessToken *AccessToken `json:"accessToken"`
	jsApiTicket *JsApiTicket `json:"jsApiTicket"`
}

type Config struct {
	Username         string        `json:"username"`
	Password         string        `json:"password"`
	BaseUrl          string        `json:"baseUrl"`
	Name             string        `json:"name"`
	AppId            string        `json:"appId"`
	Secret           string        `json:"secret"`
	Token            string        `json:"token"`
	SubscribeUrl     string        `json:"subscribeUrl"`
	OAuthRedirectUri string        `json:"oAuthRedirectUri"`
	Logger           logger.Logger `json:"logger"`
	Redis            *redis.Redis  `json:"redis"`
}

type AccessTokenRes struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    AccessToken `json:"data"`
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

			if wx.accessToken.ExpiresIn == 0 {
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

			if wx.jsApiTicket.ExpiresIn == 0 {
				wx.ClearJsApiTicket()
				break
			}
		}
	}()
}

func (wx *Weixin) GetAccessToken() (a *AccessToken, err error) {
	if wx.accessToken != nil && wx.accessToken.AccessToken != "" {
		wx.Logger.Debug("GetAccessToken from application", wx.Logger.Field("appId", wx.AppId))
		return wx.accessToken, nil
	}

	accessToken, err := wx.Redis.Get(context.Background(), base.AccessTokenPrefix+wx.AppId).Result()
	ttl, err := wx.Redis.TTL(context.Background(), base.AccessTokenPrefix+wx.AppId).Result()
	if accessToken != "" && ttl > 0 {
		wx.SetAccessToken(accessToken, int64(ttl/time.Second))
		wx.Logger.Debug("GetAccessToken from redis", wx.Logger.Field("appId", wx.AppId))
		return wx.accessToken, nil
	}

	resp, err := resty.New().R().
		SetBody(map[string]string{
			"appId":  wx.AppId,
			"secret": wx.Secret,
		}).
		SetBasicAuth(wx.Username, wx.Password).
		Post(wx.BaseUrl + "/token")
	if err != nil {
		return nil, err
	}

	wx.Logger.Debug("GetAccessToken", wx.Logger.Field("response", resp.Body()), wx.Logger.Field("appId", wx.AppId))

	var res = &AccessTokenRes{}
	if err := json.Unmarshal(resp.Body(), res); err != nil {
		return nil, err
	}

	if res.Code != "00" {
		if res.Message != "" {
			return nil, errors.New(fmt.Sprintf("(%s)%s", res.Code, res.Message))
		} else {
			return nil, errors.New("getAccessTokenError")
		}
	}

	wx.SetAccessToken(res.Data.AccessToken, res.Data.ExpiresIn)
	wx.Logger.Debug("GetAccessToken from request", wx.Logger.Field("appId", wx.AppId))

	return wx.accessToken, nil
}

type JsApiTicketRes struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    JsApiTicket `json:"data"`
}

func (wx *Weixin) GetJsApiTicket() (j *JsApiTicket, err error) {
	if wx.jsApiTicket != nil && wx.jsApiTicket.Ticket != "" {
		wx.Logger.Debug("GetJsApiTicket from application", wx.Logger.Field("appId", wx.AppId))
		return wx.jsApiTicket, nil
	}

	jsApiTicket, err := wx.Redis.Get(context.Background(), base.JsApiTicketPrefix+wx.AppId).Result()
	ttl, err := wx.Redis.TTL(context.Background(), base.JsApiTicketPrefix+wx.AppId).Result()
	if jsApiTicket != "" && ttl > 0 {
		wx.SetJsApiTicket(jsApiTicket, int64(ttl/time.Second))
		wx.Logger.Debug("GetJsApiTicket from redis", wx.Logger.Field("appId", wx.AppId))
		return wx.jsApiTicket, nil
	}

	resp, err := resty.New().R().
		SetBody(map[string]string{
			"appId":  wx.AppId,
			"secret": wx.Secret,
		}).
		SetBasicAuth(wx.Username, wx.Password).
		Post(wx.BaseUrl + "/ticket")
	if err != nil {
		return nil, err
	}
	wx.Logger.Debug("GetJsApiTicket", wx.Logger.Field("response", resp.Body()), wx.Logger.Field("appId", wx.AppId))

	var res = &JsApiTicketRes{}
	if err := json.Unmarshal(resp.Body(), res); err != nil {
		return nil, err
	}

	if res.Code != "00" {
		if res.Message != "" {
			return nil, errors.New(fmt.Sprintf("(%s)%s", res.Code, res.Message))
		} else {
			return nil, errors.New("getJsApiTicketError")
		}
	}

	wx.SetJsApiTicket(res.Data.Ticket, res.Data.ExpiresIn)

	wx.Logger.Debug("GetJsApiTicket from request", wx.Logger.Field("appId", wx.AppId))
	return wx.jsApiTicket, nil
}
