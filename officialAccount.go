package weixin

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"sort"
	"strings"
)

type SignatureReq struct {
	Signature string
	Nonce     string
	Timestamp string
}

func (wx *Weixin) VerifySignature(params *SignatureReq) error {
	signArr := []string{params.Nonce, params.Timestamp, wx.Token}
	sort.Strings(signArr)
	signStr := strings.Join(signArr, "")
	hash := sha1.New()
	hash.Write([]byte(signStr))
	signature := hex.EncodeToString(hash.Sum(nil))
	if signature != params.Signature {
		return errors.New("signature invalid")
	}
	return nil
}

type GetUserInfoRes struct {
	ErrCode    int      `json:"errcode"`
	ErrMsg     string   `json:"errmsg"`
	Subscribe  int      `json:"subscribe"`
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

func (wx *Weixin) GetUserInfo(openId string) (*GetUserInfoRes, error) {

	accessToken, err := wx.GetAccessToken()
	if err != nil {
		return nil, err
	}

	resp, err := resty.New().R().
		SetQueryParams(map[string]string{
			"access_token": accessToken.AccessToken,
			"openid":       openId,
			"lang":         "zh_CN",
		}).
		Get("https://api.weixin.qq.com/cgi-bin/user/info")
	if err != nil {
		return nil, err
	}

	var res = &GetUserInfoRes{}
	if err := json.Unmarshal(resp.Body(), res); err != nil {
		return nil, err
	}
	if res.ErrCode != 0 {
		if res.ErrMsg != "" {
			return nil, errors.New(fmt.Sprintf("(%d)%s", res.ErrCode, res.ErrMsg))
		} else {
			return nil, errors.New("GetUserInfo")
		}
	}
	return res, nil
}

type TemplateReq struct {
	OpenId     string                 `json:"openId"`
	TemplateId string                 `json:"templateId"`
	Url        string                 `json:"url"`
	Data       map[string]interface{} `json:"data"`
}

type TemplateRes struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	MsgId   string `json:"msgid"`
}

func (wx *Weixin) SendTemplate(template *TemplateReq) error {
	accessToken, err := wx.GetAccessToken()
	if err != nil {
		return err
	}

	resp, err := resty.New().R().
		SetBody(map[string]interface{}{
			"touser":      template.OpenId,
			"template_id": template.TemplateId,
			"url":         template.Url,
			"data":        template.Data,
		}).
		Post("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + accessToken.AccessToken)
	if err != nil {
		return err
	}

	var res = &OAuthAccessTokenRes{}
	if err := json.Unmarshal(resp.Body(), res); err != nil {
		return err
	}

	if res.ErrCode != 0 {
		if res.ErrMsg != "" {
			return errors.New(fmt.Sprintf("(%d)%s", res.ErrCode, res.ErrMsg))
		} else {
			return errors.New("SendTemplate")
		}
	}
	return nil
}

type BatchGetMaterialReq struct {
	Type   string `json:"type"`
	Offset string `json:"offset"`
	Count  string `json:"count"`
}

func (wx *Weixin) BatchGetMaterial(params *BatchGetMaterialReq) (map[string]interface{}, error) {

	accessToken, err := wx.GetAccessToken()
	if err != nil {
		return nil, err
	}

	resp, err := resty.New().R().
		SetQueryParams(map[string]string{
			"access_token": accessToken.AccessToken,
		}).
		SetBody(params).
		Post("https://api.weixin.qq.com/cgi-bin/material/batchget_material")
	if err != nil {
		return nil, err
	}

	var res = make(map[string]interface{})
	if err := json.Unmarshal(resp.Body(), res); err != nil {
		return nil, err
	}

	if res["errcode"] != nil && res["errcode"] != 0 {
		if res["errmsg"] != "" {
			return nil, errors.New(fmt.Sprintf("(%d)%s", res["errcode"], res["errmsg"]))
		} else {
			return nil, errors.New("BatchGetMaterial")
		}
	}

	return res, nil
}

func (wx *Weixin) MenuCreate(params map[string]interface{}) error {

	accessToken, err := wx.GetAccessToken()
	if err != nil {
		return err
	}

	resp, err := resty.New().R().
		SetQueryParams(map[string]string{
			"access_token": accessToken.AccessToken,
		}).
		SetBody(params).
		Post("https://api.weixin.qq.com/cgi-bin/menu/create")
	if err != nil {
		return err
	}

	var res = make(map[string]interface{})
	if err := json.Unmarshal(resp.Body(), &res); err != nil {
		return err
	}

	if res["errcode"] != nil && int(res["errcode"].(float64)) != 0 {
		if res["errmsg"] != "" {
			return errors.New(fmt.Sprintf("(%d)%s", int(res["errcode"].(float64)), res["errmsg"]))
		} else {
			return errors.New("MenuCreate")
		}
	}
	return nil
}

func (wx *Weixin) MenuDelete() error {

	accessToken, err := wx.GetAccessToken()
	if err != nil {
		return err
	}

	resp, err := resty.New().R().
		SetQueryParams(map[string]string{
			"access_token": accessToken.AccessToken,
		}).
		Post("https://api.weixin.qq.com/cgi-bin/menu/delete")
	if err != nil {
		return err
	}

	var res = make(map[string]interface{})
	if err := json.Unmarshal(resp.Body(), &res); err != nil {
		return err
	}

	if res["errcode"] != nil && int(res["errcode"].(float64)) != 0 {
		if res["errmsg"] != "" {
			return errors.New(fmt.Sprintf("(%d)%s", int(res["errcode"].(float64)), res["errmsg"]))
		} else {
			return errors.New("MenuDelete")
		}
	}
	return nil
}

type UniformMessageReq struct {
	OpenId     string                 `json:"openId"`
	TemplateId string                 `json:"templateId"`
	Url        string                 `json:"url"`
	Data       map[string]interface{} `json:"data"`
}

type UniformMessageRes struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	MsgId   string `json:"msgid"`
}

func (wx *Weixin) SendUniformMessage(template *UniformMessageReq) error {
	accessToken, err := wx.GetAccessToken()
	if err != nil {
		return err
	}

	resp, err := resty.New().R().
		SetBody(map[string]interface{}{
			"touser": template.OpenId,
			"mp_template_msg": map[string]interface{}{
				"appid":       wx.AppId,
				"template_id": template.TemplateId,
				"url":         template.Url,
				"data":        template.Data,
			},
		}).
		Post("https://api.weixin.qq.com/cgi-bin/message/wxopen/template/uniform_send?access_token=" + accessToken.AccessToken)
	if err != nil {
		return err
	}

	var res = &OAuthAccessTokenRes{}
	if err := json.Unmarshal(resp.Body(), res); err != nil {
		return err
	}

	if res.ErrCode != 0 {
		if res.ErrMsg != "" {
			return errors.New(fmt.Sprintf("(%d)%s", res.ErrCode, res.ErrMsg))
		} else {
			return errors.New("SendUniformMessage")
		}
	}
	return nil
}
