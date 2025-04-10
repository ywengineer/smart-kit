package oauths

import (
	"context"
	"errors"
	"fmt"
	"github.com/ywengineer/smart-kit/pkg/nets"
	"github.com/ywengineer/smart-kit/pkg/utilk"
	"net/url"
	"strconv"
)

func NewFacebookAuth(appId, appSecret, redirectUrl string) AuthFacade {
	return &fbAuth{
		appId:       appId,
		appSecret:   appSecret,
		baseUrl:     "https://graph.facebook.com/v22.0",
		redirectURL: redirectUrl,
	}
}

type fbAuth struct {
	baseUrl     string
	appId       string
	appSecret   string
	redirectURL string
}

func (f *fbAuth) Validate(metadata string) (AuthFacade, error) {
	if len(f.appSecret) == 0 || len(f.appId) == 0 || len(f.redirectURL) == 0 {
		return nil, errors.New("missing prop [app-id/app-secret/redirect-url] for facebook auth: " + metadata)
	}
	return f, nil
}

func (f *fbAuth) GetToken(code string) (*AccessToken, error) {
	params := url.Values{}
	params.Add("client_id", f.appId)
	params.Add("client_secret", f.appSecret)
	params.Add("redirect_uri", f.redirectURL)
	params.Add("code", code)

	sc, body, err := cli.Get(context.Background(), f.baseUrl+"/oauth/access_token?"+params.Encode())
	if err != nil {
		return nil, err
	} else if !nets.Is2xx(sc) {
		return nil, errors.New("bad status code: " + strconv.Itoa(sc))
	} else {
		return utilk.UnmarshalJSON(body, &AccessToken{})
	}
}

func (f *fbAuth) RefreshToken(refreshToken string) (*RefreshToken, error) {
	return nil, errors.New("unsupported")
}

func (f *fbAuth) GetUserInfo(_ string, accessToken string) (*UserInfo, error) {
	sc, body, err := cli.Get(context.Background(), f.baseUrl+"/me?access_token="+accessToken)
	if err != nil {
		return nil, err
	} else if !nets.Is2xx(sc) {
		return nil, errors.New("bad status code: " + strconv.Itoa(sc))
	} else if mui, err := utilk.UnmarshalJSON(body, &MetaUserInfo{}); err != nil {
		return nil, err
	} else if mui.Error != nil {
		return nil, errors.New(fmt.Sprintf("bad data: code = %d, type = %s, msg = %s", mui.Error.Code, mui.Error.Type, mui.Error.Message))
	} else {
		return &UserInfo{
			Openid:   mui.Id,
			Nickname: mui.Name,
		}, nil
	}
}

type MetaUserInfo struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Picture   struct {
		Data struct {
			Url          string `json:"url"`
			IsSilhouette bool   `json:"is_silhouette"`
		} `json:"data"`
	} `json:"picture"`
	Error *struct {
		Message   string `json:"message"`
		Type      string `json:"type"`
		Code      int    `json:"code"`
		FbtraceId string `json:"fbtrace_id"`
	} `json:"error,omitempty"`
}
