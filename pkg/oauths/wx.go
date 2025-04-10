package oauths

import (
	"context"
	"errors"
	"github.com/bytedance/sonic"
	"github.com/ywengineer/smart-kit/pkg/nets"
	"github.com/ywengineer/smart-kit/pkg/utilk"
	"net/url"
	"strconv"
)

func NewWxAuth(appId, appSecret string) AuthFacade {
	return &wxAuth{appId: appId, appSecret: appSecret, baseUrl: "https://api.weixin.qq.com"}
}

type wxAuth struct {
	baseUrl   string
	appId     string
	appSecret string
}

func (w *wxAuth) Validate(metadata string) (AuthFacade, error) {
	if len(w.appSecret) == 0 || len(w.appId) == 0 {
		return nil, errors.New("missing prop [app-id/app-secret] for weixin auth: " + metadata)
	}
	return w, nil
}

func (w *wxAuth) GetToken(code string) (*AccessToken, error) {
	p := url.Values{}
	p.Set("appid", w.appId)
	p.Set("secret", w.appSecret)
	p.Set("code", code)
	p.Set("grant_type", "authorization_code")
	sc, body, err := cli.Get(context.Background(), w.baseUrl+"/sns/oauth2/access_token?"+p.Encode())
	if err != nil {
		return nil, err
	} else if !nets.Is2xx(sc) {
		return nil, errors.New("bad status code: " + strconv.Itoa(sc))
	} else {
		at := &AccessToken{}
		err = sonic.Unmarshal(body, at)
		if err != nil {
			return nil, err
		} else if at.ErrCode != 0 {
			return nil, errors.New(at.ErrMsg + ": " + strconv.Itoa(at.ErrCode))
		}
		return at, err
	}
}

func (w *wxAuth) RefreshToken(refreshToken string) (*RefreshToken, error) {
	p := url.Values{}
	p.Set("appid", w.appId)
	p.Set("refresh_token", refreshToken)
	p.Set("grant_type", "refresh_token")
	sc, body, err := cli.Get(context.Background(), w.baseUrl+"/sns/oauth2/refresh_token?"+p.Encode())
	if err != nil {
		return nil, err
	} else if !nets.Is2xx(sc) {
		return nil, errors.New("bad status code: " + strconv.Itoa(sc))
	} else {
		at, err := utilk.UnmarshalJSON(body, &RefreshToken{})
		if err != nil {
			return nil, err
		} else if at.ErrCode != 0 {
			return nil, errors.New(at.ErrMsg + ": " + strconv.Itoa(at.ErrCode))
		}
		return at, err
	}
}

func (w *wxAuth) GetUserInfo(openid string, accessToken string) (*UserInfo, error) {
	p := url.Values{}
	p.Set("openid", openid)
	p.Set("access_token", accessToken)
	sc, body, err := cli.Get(context.Background(), w.baseUrl+"/sns/userinfo?"+p.Encode())
	if err != nil {
		return nil, err
	} else if !nets.Is2xx(sc) {
		return nil, errors.New("bad status code: " + strconv.Itoa(sc))
	} else {
		at, err := utilk.UnmarshalJSON(body, &UserInfo{})
		if err != nil {
			return nil, err
		} else if at.ErrCode != 0 {
			return nil, errors.New(at.ErrMsg + ": " + strconv.Itoa(at.ErrCode))
		}
		return at, err
	}
}
