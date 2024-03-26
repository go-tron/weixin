package weixin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
)

type Scope string

const (
	ScopeBase     Scope = "snsapi_base"
	ScopeUserInfo Scope = "snsapi_userinfo"
)

type OAuthCodeReq struct {
	Uri            string `json:"uri"`
	Scope          Scope  `json:"scope"`
	State          string `json:"state"`
	UseRedirectUri bool   `json:"useRedirectUri"`
}

type OAuthCodeQuery struct {
	AppId        string `url:"appid"`
	RedirectUri  string `url:"redirect_uri"`
	ResponseType string `url:"response_type"`
	Scope        string `url:"scope"`
	State        string `url:"state"`
}

func (wx *Weixin) GetOAuthCode(params *OAuthCodeReq) (string, error) {
	req := OAuthCodeQuery{
		AppId:        wx.AppId,
		RedirectUri:  params.Uri,
		ResponseType: "code",
		Scope:        string(params.Scope),
		State:        params.State,
	}

	if wx.OAuthRedirectUri != "" {
		req.RedirectUri = wx.OAuthRedirectUri + req.RedirectUri
	}

	v, err := query.Values(req)
	if err != nil {
		return "", err
	}
	return "https://open.weixin.qq.com/connect/oauth2/authorize?" + v.Encode() + "#wechat_redirect", nil
}

type OAuthAccessTokenRes struct {
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
}

func (wx *Weixin) GetOAuthAccessToken(code string) (*OAuthAccessTokenRes, error) {
	resp, err := resty.New().R().
		SetQueryParams(map[string]string{
			"appid":      wx.AppId,
			"secret":     wx.Secret,
			"code":       code,
			"grant_type": "authorization_code",
		}).
		Get("https://api.weixin.qq.com/sns/oauth2/access_token")
	if err != nil {
		return nil, err
	}

	var res = &OAuthAccessTokenRes{}
	if err := json.Unmarshal(resp.Body(), res); err != nil {
		return nil, err
	}

	if res.ErrCode != 0 {
		if res.ErrMsg != "" {
			return nil, errors.New(fmt.Sprintf("(%d)%s", res.ErrCode, res.ErrMsg))
		} else {
			return nil, errors.New("GetOAuthAccessToken")
		}
	}
	return res, nil
}

type OAuthUserInfoReq struct {
	OpenId      string `json:"openid"`
	AccessToken string `json:"access_token"`
}

type OAuthUserInfoRes struct {
	ErrCode    int      `json:"errcode"`
	ErrMsg     string   `json:"errmsg"`
	OpenId     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Language   string   `json:"language"`
	Province   string   `json:"province"`
	Country    string   `json:"country"`
	City       string   `json:"city"`
	HeadImgUrl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	UnionId    string   `json:"unionid"`
}

func (wx *Weixin) GetOAuthUserInfo(params *OAuthUserInfoReq) (*OAuthUserInfoRes, error) {
	resp, err := resty.New().R().
		SetQueryParams(map[string]string{
			"access_token": params.AccessToken,
			"openid":       params.OpenId,
			"lang":         "zh_CN",
		}).
		Get("https://api.weixin.qq.com/sns/userinfo")
	if err != nil {
		return nil, err
	}

	var res = &OAuthUserInfoRes{}
	if err := json.Unmarshal(resp.Body(), res); err != nil {
		return nil, err
	}
	if res.ErrCode != 0 {
		if res.ErrMsg != "" {
			return nil, errors.New(fmt.Sprintf("(%d)%s", res.ErrCode, res.ErrMsg))
		} else {
			return nil, errors.New("GetOAuthUserInfo")
		}
	}
	return res, nil
}

func (wx *Weixin) GetOAuthUserInfoFromCode(code string) (*OAuthUserInfoRes, error) {
	res, err := wx.GetOAuthAccessToken(code)
	if err != nil {
		return nil, err
	}
	return wx.GetOAuthUserInfo(&OAuthUserInfoReq{
		OpenId:      res.OpenId,
		AccessToken: res.AccessToken,
	})
}
