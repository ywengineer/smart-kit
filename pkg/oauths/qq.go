package oauths

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/ywengineer/smart-kit/pkg/nets"
	"github.com/ywengineer/smart-kit/pkg/utilk"
	"net/url"
	"strconv"
	"strings"
)

func NewQQAuth(appId, appSecret, redirectUrl string) AuthFacade {
	return &qqAuth{
		appId:       appId,
		appSecret:   appSecret,
		baseUrl:     "https://graph.qq.com",
		redirectURL: redirectUrl,
	}
}

type qqAuth struct {
	baseUrl     string
	appId       string
	appSecret   string
	redirectURL string
}

func (q *qqAuth) Validate(metadata string) (AuthFacade, error) {
	if len(q.appSecret) == 0 || len(q.appId) == 0 || len(q.redirectURL) == 0 {
		return nil, errors.New("missing prop [app-id/app-secret/redirect-url] for qq auth: " + metadata)
	}
	return q, nil
}

func (q *qqAuth) GetToken(code string) (*AccessToken, error) {
	// 获取 access_token
	params := url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("client_id", q.appId)
	params.Add("client_secret", q.appSecret)
	params.Add("code", code)
	params.Add("redirect_uri", q.redirectURL)

	tokenURL := fmt.Sprintf("%s/oauth2.0/token?%s", q.baseUrl, params.Encode())
	sc, body, err := cli.Get(context.Background(), tokenURL)
	if err != nil {
		return nil, err
	} else if !nets.Is2xx(sc) {
		return nil, errors.New("bad status code: " + strconv.Itoa(sc))
	} else {
		// decode access_token
		token := string(body)
		if !strings.Contains(token, "callback") {
			return nil, errors.New("error token response: " + token)
		}
		posL, posR := strings.IndexRune(token, '('), strings.IndexByte(token, ')')
		if posL < 0 || posR < 0 || posR >= posL {
			return nil, errors.New("error token response: " + token)
		}
		token = token[posL+1 : posR]
		var tokenMap map[string]interface{}
		err = sonic.Unmarshal(body, &tokenMap)
		if err != nil {
			return nil, err
		} else if errCode, ok := tokenMap["error"]; ok {
			return nil, errors.New(fmt.Sprintf("%v: %v", errCode, tokenMap["error_description"]))
		} else {
			return &AccessToken{
				AccessToken: tokenMap["access_token"].(string),
			}, nil
		}
	}
}

func (q *qqAuth) RefreshToken(refreshToken string) (*RefreshToken, error) {
	return nil, errors.New("unsupported")
}

func (q *qqAuth) GetUserInfo(_ string, accessToken string) (*UserInfo, error) {
	// get openid
	openidURL := fmt.Sprintf("%s/oauth2.0/me?access_token=%s", q.baseUrl, accessToken)
	sc, body, err := cli.Get(context.Background(), openidURL)
	if err != nil {
		return nil, err
	} else if !nets.Is2xx(sc) {
		return nil, errors.New("bad status code: " + strconv.Itoa(sc))
	} else {
		// decode openid
		openidStr := string(body)
		if !strings.Contains(openidStr, "callback") {
			return nil, errors.New("error openid response: " + openidStr)
		}
		posL, posR := strings.IndexRune(openidStr, '('), strings.IndexByte(openidStr, ')')
		if posL < 0 || posR < 0 || posR >= posL {
			return nil, errors.New("error openid response: " + openidStr)
		}
		openidStr = openidStr[posL+1 : posR]
		var openIdMap map[string]interface{}
		err = sonic.Unmarshal(body, &openIdMap)
		if err != nil {
			return nil, err
		} else if errCode, ok := openIdMap["error"]; ok {
			return nil, errors.New(fmt.Sprintf("%v: %v", errCode, openIdMap["error_description"]))
		} else if at, err := utilk.UnmarshalJSON(body, &UserInfo{}); err != nil {
			return nil, err
		} else {
			return at, nil
		}
	}
}
