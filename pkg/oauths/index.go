package oauths

import (
	"gitee.com/ywengineer/smart-kit/pkg/rpcs"
	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"time"
)

var facadeMap = make(map[string]AuthFacade)

var cli, _ = rpcs.NewHertzRpc(nil, rpcs.RpcClientInfo{
	ClientName:     "smart-oauth-client",
	MaxRetry:       2,
	Delay:          time.Millisecond * 100,
	MaxConnPerHost: 512,
})

type AuthError struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type AccessToken struct {
	*AuthError
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"` // seconds
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionId      string `json:"unionid"`
}

type RefreshToken struct {
	*AuthError
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}

type UserInfo struct {
	*AuthError
	Openid     string   `json:"openid"`     // 普通用户的标识，对当前开发者账号唯一
	Nickname   string   `json:"nickname"`   // 普通用户昵称
	Sex        int      `json:"sex"`        // 普通用户性别，1 为男性，2 为女性
	Province   string   `json:"province"`   // 普通用户个人资料填写的省份
	City       string   `json:"city"`       // 普通用户个人资料填写的城市
	Country    string   `json:"country"`    // 国家，如中国为 CN
	HeadImgUrl string   `json:"headimgurl"` // 用户头像，最后一个数值代表正方形头像大小（有 0、46、64、96、132 数值可选，0 代表 640*640 正方形头像），用户没有头像时该项为空
	Privilege  []string `json:"privilege"`  // 用户特权信息，json 数组，如微信沃卡用户为（chinaunicom）
	UnionId    string   `json:"unionid"`    // 用户统一标识。针对一个微信开放平台账号下的应用，同一用户的 unionid 是唯一的。
}

func (u *UserInfo) UniqueId() string {
	return utilk.DefaultIfEmpty(u.UnionId, u.Openid)
}

type AuthFacade interface {
	Validate(metadata string) (AuthFacade, error)
	GetToken(code string) (*AccessToken, error)
	RefreshToken(refreshToken string) (*RefreshToken, error)
	GetUserInfo(openid string, accessToken string) (*UserInfo, error)
}
