package weixin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type GetUserPhoneNumberRes struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	PhoneInfo struct {
		PhoneNumber     string `json:"phoneNumber"`
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     string `json:"countryCode"`
	} `json:"phone_info"`
}

type UserPhoneNumber struct {
	PhoneNumber string `json:"phoneNumber"`
}

func (wx *Weixin) GetUserPhoneNumber(code string) (*UserPhoneNumber, error) {
	accessToken, err := wx.GetAccessToken()
	if err != nil {
		return nil, err
	}
	resp, err := resty.New().R().
		SetQueryParams(map[string]string{
			"access_token": accessToken.AccessToken,
		}).
		SetBody(map[string]string{
			"code": code,
		}).
		Post("https://api.weixin.qq.com/wxa/business/getuserphonenumber")
	if err != nil {
		return nil, err
	}

	var res = &GetUserPhoneNumberRes{}
	if err := json.Unmarshal(resp.Body(), res); err != nil {
		return nil, err
	}
	if res.ErrCode != 0 {
		if res.ErrMsg != "" {
			return nil, errors.New(fmt.Sprintf("(%d)%s", res.ErrCode, res.ErrMsg))
		} else {
			return nil, errors.New("GetUserPhoneNumber Failed")
		}
	}

	if res.PhoneInfo.PhoneNumber == "" {
		return nil, errors.New("获取微信绑定手机号失败")
	}
	return &UserPhoneNumber{
		PhoneNumber: res.PhoneInfo.PhoneNumber,
	}, nil
}
