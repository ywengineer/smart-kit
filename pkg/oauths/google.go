package oauths

import (
	"context"
	"errors"
	"fmt"
	"gitee.com/ywengineer/smart-kit/pkg/logk"
	"gitee.com/ywengineer/smart-kit/pkg/rpcs"
	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"google.golang.org/api/idtoken"
)

func NewGoogleAuth(clientId string) AuthFacade {
	va, err := idtoken.NewValidator(context.Background(), idtoken.WithHTTPClient(rpcs.StandardUseHertz()))
	if err != nil {
		logk.Fatalf("Failed to create google idtoken validator: %v", err)
		return nil
	}
	return &googleAuth{appId: clientId, validator: va}
}

type googleAuth struct {
	appId     string
	appSecret string
	validator *idtoken.Validator
}

func (g *googleAuth) Validate(metadata string) (AuthFacade, error) {
	if len(g.appId) == 0 {
		return nil, errors.New("missing prop [app-id/app-secret] for google auth: " + metadata)
	}
	return g, nil
}

func (g *googleAuth) GetToken(code string) (*AccessToken, error) {
	return &AccessToken{
		AccessToken: code,
	}, nil
}

func (g *googleAuth) RefreshToken(refreshToken string) (*RefreshToken, error) {
	return nil, errors.New("unsupported")
}

func (g *googleAuth) GetUserInfo(_ string, accessToken string) (*UserInfo, error) { // 创建验证器
	/*
		{
		    "iss": "https://accounts.google.com",
		    "azp": "YOUR_CLIENT_ID",
		    "aud": "YOUR_CLIENT_ID",
		    "sub": "123456789012345678901",
		    "email": "example@example.com",
		    "email_verified": true,
		    "at_hash": "abcdefghijklmnopqrstuvwxyz123456",
		    "name": "John Doe",
		    "picture": "https://lh3.googleusercontent.com/a-/AOh14Gi71234567890abcdefghijklmnopqrstuvwxyz",
		    "given_name": "John",
		    "family_name": "Doe",
		    "locale": "en-US",
		    "iat": 1630435200,
		    "exp": 1630438800
		}
		各字段说明：
		iss：令牌的颁发者，对于 Google ID Token 来说，固定为 "https://accounts.google.com"。
		azp：授权方，通常是客户端 ID。
		aud：受众，指定这个 ID Token 是为哪个客户端颁发的，通常也是客户端 ID。
		sub：用户的唯一标识符，是一个字符串，在 Google 账户系统中唯一标识一个用户。
		email：用户的电子邮件地址。
		email_verified：一个布尔值，表示用户的电子邮件地址是否已经过验证。
		at_hash：访问令牌的哈希值，用于验证 ID Token 与访问令牌是否匹配。
		name：用户的全名。
		picture：用户的头像 URL。
		given_name：用户的名字。
		family_name：用户的姓氏。
		locale：用户的语言和地区设置，例如 "en-US" 表示美式英语。
		iat：令牌的颁发时间，以 Unix 时间戳表示。
		exp：令牌的过期时间，以 Unix 时间戳表示。
	*/
	payload, err := g.validator.Validate(context.Background(), accessToken, g.appId)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %v", err)
	}
	return &UserInfo{
		Openid:     payload.Subject,
		Nickname:   fmt.Sprintf("%s (%s)", utilk.ToString(payload.Claims["name"]), utilk.ToString(payload.Claims["email"])),
		Sex:        0,
		Province:   "",
		City:       "",
		Country:    utilk.ToString(payload.Claims["locale"]),
		HeadImgUrl: utilk.ToString(payload.Claims["picture"]),
		Privilege:  nil,
		UnionId:    payload.Subject,
	}, nil
}
