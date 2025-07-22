package oauths

import (
	"context"
	"errors"
	"gitee.com/ywengineer/smart-kit/pkg/nets"
	"github.com/bytedance/sonic"
	"net/url"
	"strconv"
	"strings"
)

func NewSteamWebAuth(appId, appSecret string) AuthFacade {
	return &steamAuth{appId: appId, appSecret: appSecret, baseUrl: "https://partner.steam-api.com"}
}

type steamAuth struct {
	baseUrl   string
	appId     string
	appSecret string
}

func (steam *steamAuth) Validate(metadata string) (AuthFacade, error) {
	if len(steam.appSecret) == 0 || len(steam.appId) == 0 {
		return nil, errors.New("missing prop [app-id/app-secret] for steam auth: " + metadata)
	}
	return steam, nil
}

// GetToken code = ticket_identity or ticket
func (steam *steamAuth) GetToken(code string) (*AccessToken, error) {
	return &AccessToken{AccessToken: code}, nil
}

func (steam *steamAuth) RefreshToken(refreshToken string) (*RefreshToken, error) {
	return nil, errors.New("unsupported")
}

// GetUserInfo token = ticket_identity or ticket
func (steam *steamAuth) GetUserInfo(_ string, token string) (*UserInfo, error) {
	ticket, identity := "", ""
	underline := strings.IndexRune(token, '_')
	if underline > -1 {
		ticket, identity = token[:underline], token[underline+1:]
	} else {
		ticket = token
	}
	p := url.Values{}
	p.Set("appid", steam.appId)
	p.Set("key", steam.appSecret)
	p.Set("ticket", ticket)
	if len(identity) > 0 {
		p.Set("identity", identity)
	}
	sc, body, err := cli.Get(context.Background(), steam.baseUrl+"/ISteamUserAuth/AuthenticateUserTicket/v1/?"+p.Encode())
	if err != nil {
		return nil, err
	} else if !nets.Is2xx(sc) {
		return nil, errors.New("bad status code: " + strconv.Itoa(sc))
	} else {
		at := &SteamAuthResponse{}
		err = sonic.Unmarshal(body, at)
		if err != nil {
			return nil, err
		} else if !strings.EqualFold(at.Response.Params.Result, "OK") {
			return nil, errors.New(at.Response.Params.Result + ": " + strconv.Itoa(at.Response.Error.Code) + ", " + at.Response.Error.Desc)
		} else if at.Response.Params.VacBanned || at.Response.Params.PublisherBans {
			return nil, errors.New("banned")
		}
		return &UserInfo{
			UnionId:  at.Response.Params.SteamID,
			Openid:   at.Response.Params.SteamID,
			Nickname: at.Response.Params.DisplayName,
			Sex:      0,
		}, nil
	}
}

// SteamAuthResponse 定义 Steam 验证响应的结构体
type SteamAuthResponse struct {
	Response struct {
		Error struct {
			Code int    `json:"errorcode"`
			Desc string `json:"errordesc"`
		} `json:"error"`
		Params struct {
			Result        string `json:"result"`        // "OK" 表示验证成功，其他值（如 "ExpiredTicket"）表示失败。
			SteamID       string `json:"steamid"`       // 用户的64位SteamID
			OwnerSteamID  string `json:"ownersteamid"`  // 通常与steamid一致
			VacBanned     bool   `json:"vacbanned"`     // 是否被VAC封禁
			PublisherBans bool   `json:"publisherbans"` // 是否被开发者封禁
			DisplayName   string `json:"displayname"`   // 用户昵称
			SessionFlags  int    `json:"sessionflags"`  // 会话标志（1=普通用户，2=受限用户（如未消费），4=机器人。）
			TicketType    int    `json:"tickettype"`    // Ticket类型（0=普通会话，1=持久会话）
		} `json:"params"`
	} `json:"response"`
}
