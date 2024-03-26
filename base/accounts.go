package base

type accounts interface {
	GetAccountById(string) (*Weixin, error)
}

type Accounts struct {
	Accounts accounts
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
