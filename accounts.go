package weixin

type accounts interface {
	GetAccountById(string) (*Weixin, error)
}

type Accounts struct {
	Accounts accounts
}

func (u *Accounts) GetName(appId string) (string, error) {
	account, err := u.Accounts.GetAccountById(appId)
	if err != nil {
		return "", err
	}
	return account.Name, nil
}

func (u *Accounts) GetAccessToken(appId string) (*AccessToken, error) {
	account, err := u.Accounts.GetAccountById(appId)
	if err != nil {
		return nil, err
	}
	return account.GetAccessToken()
}

func (u *Accounts) GetJsApiTicket(appId string) (*JsApiTicket, error) {
	account, err := u.Accounts.GetAccountById(appId)
	if err != nil {
		return nil, err
	}
	return account.GetJsApiTicket()
}

func (u *Accounts) GetJsApiConfig(appId string, url string) (*JsApiConfig, error) {
	account, err := u.Accounts.GetAccountById(appId)
	if err != nil {
		return nil, err
	}
	return account.GetJsApiConfig(url)
}

func (u *Accounts) GetOAuthCode(appId string, params *OAuthCodeReq) (string, error) {
	account, err := u.Accounts.GetAccountById(appId)
	if err != nil {
		return "", err
	}
	return account.GetOAuthCode(params)
}

func (u *Accounts) GetOAuthAccessToken(appId string, code string) (*OAuthAccessTokenRes, error) {
	account, err := u.Accounts.GetAccountById(appId)
	if err != nil {
		return nil, err
	}
	return account.GetOAuthAccessToken(code)
}

func (u *Accounts) GetOAuthUserInfo(appId string, params *OAuthUserInfoReq) (*OAuthUserInfoRes, error) {
	account, err := u.Accounts.GetAccountById(appId)
	if err != nil {
		return nil, err
	}
	return account.GetOAuthUserInfo(params)
}

func (u *Accounts) GetOAuthUserInfoFromCode(appId string, code string) (*OAuthUserInfoRes, error) {
	account, err := u.Accounts.GetAccountById(appId)
	if err != nil {
		return nil, err
	}
	return account.GetOAuthUserInfoFromCode(code)
}

func (u *Accounts) VerifySignature(appId string, params *SignatureReq) error {
	account, err := u.Accounts.GetAccountById(appId)
	if err != nil {
		return err
	}
	return account.VerifySignature(params)
}

func (u *Accounts) GetUserInfo(appId string, openId string) (*GetUserInfoRes, error) {
	account, err := u.Accounts.GetAccountById(appId)
	if err != nil {
		return nil, err
	}
	return account.GetUserInfo(openId)
}

func (u *Accounts) GetSubscribeUrl(appId string) (string, error) {
	account, err := u.Accounts.GetAccountById(appId)
	if err != nil {
		return "", err
	}
	return account.SubscribeUrl, nil
}

func (u *Accounts) SendTemplate(appId string, template *TemplateReq) error {
	account, err := u.Accounts.GetAccountById(appId)
	if err != nil {
		return err
	}
	return account.SendTemplate(template)
}

func (u *Accounts) GetUserPhoneNumber(appId string, code string) (*UserPhoneNumber, error) {
	account, err := u.Accounts.GetAccountById(appId)
	if err != nil {
		return nil, err
	}
	return account.GetUserPhoneNumber(code)
}
