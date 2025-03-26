package oauths

import (
	"errors"
	"github.com/google/uuid"
	"strings"
)

type anoAuth struct {
	baseUrl     string
	appId       string
	appSecret   string
	redirectURL string
}

func (q *anoAuth) Validate(_ string) (AuthFacade, error) {
	return q, nil
}

func (q *anoAuth) GetToken(code string) (*AccessToken, error) {
	return &AccessToken{
		AccessToken: uuid.New().String(),
		Openid:      strings.ToLower(strings.ReplaceAll(uuid.New().String(), "-", "")),
	}, nil
}

func (q *anoAuth) RefreshToken(refreshToken string) (*RefreshToken, error) {
	return nil, errors.New("unsupported")
}

func (q *anoAuth) GetUserInfo(openid string, _ string) (*UserInfo, error) {
	// get openid
	return &UserInfo{
		Openid: openid,
	}, nil
}
