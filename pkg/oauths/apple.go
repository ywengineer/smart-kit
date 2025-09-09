package oauths

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"gitee.com/ywengineer/smart-kit/pkg/nets"
	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"github.com/bytedance/sonic"
	"github.com/golang-jwt/jwt/v4"
	"math/big"
	"strconv"
	"time"
)

func NewAppleAuth(appId string) AuthFacade {
	return &appleAuth{appId: appId, basePath: "https://appleid.apple.com", keys: make(map[string]*rsa.PublicKey)}
}

type appleAuth struct {
	appId    string
	basePath string
	keys     map[string]*rsa.PublicKey
}

func (a *appleAuth) Validate(metadata string) (AuthFacade, error) {
	if len(a.appId) == 0 {
		return nil, errors.New("missing prop [app-id/app-secret] for apple auth: " + metadata)
	}
	if err := a.initPublicKeys(); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *appleAuth) GetToken(code string) (*AccessToken, error) {
	return &AccessToken{
		AccessToken: code,
	}, nil
}

func (a *appleAuth) RefreshToken(refreshToken string) (*RefreshToken, error) {
	return nil, errors.New("unsupported")
}

func (a *appleAuth) GetUserInfo(_ string, accessToken string) (*UserInfo, error) {
	/*
			{
			    "header": {
			        "alg": "RS256",
			        "kid": "ABC123"
			    },
			    "claims": {
			        "iss": "https://appleid.apple.com",
			        "aud": "com.example.app",
			        "exp": 1699450800,
			        "iat": 1699447200,
			        "sub": "001234.abcdef1234567890.5678",
			        "at_hash": "abcdefghijklmnopqrstuvwxyz123456",
			        "email": "example_private@privaterelay.appleid.com",
			        "email_verified": "true",
			        "is_private_email": "true",
			        "auth_time": 1699447200
			    },
			    "signature": "abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdef"
			}
		1. header（头部）
			alg：指定 JWT 使用的签名算法，通常 Apple ID 登录使用的是 RS256（RSA 256 位签名算法）。
			kid：公钥 ID，用于标识验证签名时所使用的公钥，可从 Apple 的公钥端点（https://appleid.apple.com/auth/keys）获取对应的公钥。
		2. claims（声明）
			iss：令牌的颁发者，固定为 "https://appleid.apple.com"，表明该令牌是由 Apple 颁发的。
			aud：受众，指的是该令牌的目标接收者，一般是你的应用的客户端 ID，用于确保令牌是发给你的应用的。
			exp：令牌的过期时间，以 Unix 时间戳表示，超过此时间后令牌将失效。
			iat：令牌的颁发时间，以 Unix 时间戳表示，用于验证令牌的时效性。
			sub：用户的唯一标识符，在 Apple 生态系统中唯一标识一个用户，可用于在你的系统中关联和识别用户。
			at_hash：访问令牌的哈希值，用于验证 ID Token 与访问令牌是否匹配。
			email：用户的电子邮件地址。如果用户选择了使用 Apple 的私人电子邮件转发服务，这里会显示一个以 @privaterelay.appleid.com 结尾的临时邮箱地址。
			email_verified：表示用户的电子邮件地址是否已经过验证，值为 "true" 或 "false"。
			is_private_email：指示用户是否使用了 Apple 的私人电子邮件转发服务，值为 "true" 或 "false"。
			auth_time：用户进行身份验证的时间，以 Unix 时间戳表示。
		3. signature（签名）
			用于验证令牌的完整性和真实性，确保令牌在传输过程中没有被篡改。
	*/
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		} else if kid, ok := token.Header["kid"].(string); !ok {
			return nil, errors.New("kid not found in token header")
		} else if pubKey, ok := a.keys[kid]; !ok {
			return nil, errors.New(fmt.Sprintf("kid [%s] not found in public keys", kid))
		} else {
			return pubKey, nil
		}
	})

	if err != nil {
		return nil, err
	}
	//
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	//
	if !claims.VerifyAudience(a.appId, true) {
		return nil, errors.New("invalid app from token")
	}
	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, errors.New("token expired")
	}
	return &UserInfo{
		Openid:     utilk.ToString(claims["sub"]),
		Nickname:   utilk.ToString(claims["email"]),
		Sex:        0,
		Province:   "",
		City:       "",
		Country:    "",
		HeadImgUrl: "",
		Privilege:  nil,
		UnionId:    utilk.ToString(claims["sub"]),
	}, nil
}

type AppleJWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// AppleJWKS 定义 Apple 的 JWKS 结构体
type AppleJWKS struct {
	Keys []AppleJWK `json:"keys"`
}

// GetApplePublicKeys 获取 Apple 的公钥
func (a *appleAuth) initPublicKeys() error {
	sc, body, err := cli.Get(context.Background(), a.basePath+"/auth/keys", nil)
	if err != nil {
		return err
	} else if !nets.Is2xx(sc) {
		return errors.New("failed to GetApplePublicKeys: " + strconv.Itoa(sc))
	} else {
		var jwks AppleJWKS
		if err = sonic.Unmarshal(body, &jwks); err != nil {
			return err
		}
		//
		for _, key := range jwks.Keys {
			a.keys[key.Kid], err = buildPublicKeyByKid(&key)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// buildPublicKeyByKid 根据 KID 获取对应的公钥
func buildPublicKeyByKid(jwk *AppleJWK) (*rsa.PublicKey, error) {
	// 解码 Base64 编码的模数 (n)
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, err
	}
	n := new(big.Int).SetBytes(nBytes)

	// 解码 Base64 编码的指数 (e)
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, err
	}
	e := int(new(big.Int).SetBytes(eBytes).Int64())

	pubKey := &rsa.PublicKey{
		N: n,
		E: e,
	}
	return pubKey, nil
}
