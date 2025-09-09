package config

import (
	"context"
	"fmt"
	"net/http"

	"gitee.com/ywengineer/smart-kit/pkg/rpcs"
	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/robfig/cron/v3"
)

type Channel struct {
	Id   int64  `json:"id" yaml:"id" redis:"id"`
	Code string `json:"code" yaml:"code" redis:"code"`
	Name string `json:"name" yaml:"name" redis:"name"`
}

func (c Channel) String() string {
	return fmt.Sprintf("%d:%s:%s", c.Id, c.Code, c.Name)
}

type Product struct {
	Id          int64   `json:"id" yaml:"id" redis:"id"`
	GameId      string  `json:"game_id" yaml:"game_id" redis:"game_id"`
	ServerId    string  `json:"server_id" yaml:"server_id" redis:"server_id"`
	ProductId   string  `json:"product_id" yaml:"product_id" redis:"product_id"`
	Service     string  `json:"service" yaml:"service" redis:"service"`
	ServiceDays int     `json:"serviceDays" yaml:"serviceDays" redis:"serviceDays"`
	Money       float64 `json:"money" yaml:"money" redis:"money"`
	Feed        int32   `json:"feed" yaml:"feed" redis:"feed"` // 兑换代币数量
	PlatformId  int64   `json:"platformId" yaml:"platformId" redis:"platformId"`
	//Time        time.Time `json:"time" yaml:"time" redis:"time"`
}

type GameServerInfo struct {
	Id       int64  `json:"id" yaml:"id" redis:"id"`
	GameId   string `json:"gameID" yaml:"gameID" redis:"gameID"`
	ServerId string `json:"serverID" yaml:"serverID" redis:"serverID"`
	Name     string `json:"name" yaml:"name" redis:"name"`
	GameIP   string `json:"gameIP" yaml:"gameIP" redis:"gameIP"`
	GamePort int    `json:"gamePort" yaml:"gamePort" redis:"gamePort"`
	Status   int    `json:"status" yaml:"status" redis:"status"`
	ApiUrl   string `json:"apiUrl" yaml:"apiUrl" redis:"apiUrl"`
	Metadata string `json:"metadata" yaml:"metadata" redis:"metadata"`
	md       map[string]interface{}
}

func (g GameServerInfo) IsFuncOpen(funcName string) bool {
	if len(g.Metadata) == 0 {
		return false
	}
	if err := sonic.Unmarshal([]byte(g.Metadata), &g.md); err != nil {
		return false
	}
	o, ok := g.md[funcName]
	return ok && o.(bool)
}

func (g GameServerInfo) GetApiMethodUrl(method string) string {
	return g.ApiUrl + "?m=" + method
}

type Metadata struct {
	channelMap    map[string]Channel
	productMap    map[uint64]Product
	gameServerMap map[uint64]GameServerInfo
}

func (g *Metadata) refresh(ctx context.Context) {
	rpc := rpcs.GetDefaultRpc()
	rpc.GetAsync(ctx, p.RemoteUrl.Product, nil, func(statusCode int, body []byte, err error) {
		if err != nil || statusCode != http.StatusOK {
			hlog.CtxErrorf(ctx, "failed to retrieve product data: [status = %d] [body = %s], [err = %v]", statusCode, body, err)
		} else {
			var data []Product
			if err = sonic.Unmarshal(body, &data); err != nil {
				hlog.CtxErrorf(ctx, "failed to decode product data: [status = %d] [body = %s], [err = %v]", statusCode, body, err)
			} else {
				for _, v := range data {
					g.productMap[utilk.Hash(v.ProductId, v.PlatformId)] = v
				}
			}
		}

	})
	rpc.GetAsync(ctx, p.RemoteUrl.Platform, nil, func(statusCode int, body []byte, err error) {
		if err != nil || statusCode != http.StatusOK {
			hlog.CtxErrorf(ctx, "failed to retrieve channel data: [status = %d] [body = %s], [err = %v]", statusCode, body, err)
		} else {
			var data []Channel
			if err = sonic.Unmarshal(body, &data); err != nil {
				hlog.CtxErrorf(ctx, "failed to decode channel data: [status = %d] [body = %s], [err = %v]", statusCode, body, err)
			} else {
				for _, v := range data {
					g.channelMap[v.Code] = v
				}
			}
		}
	})
	rpc.GetAsync(ctx, p.RemoteUrl.GameServer, nil, func(statusCode int, body []byte, err error) {
		if err != nil || statusCode != http.StatusOK {
			hlog.CtxErrorf(ctx, "failed to retrieve game server data: [status = %d] [body = %s], [err = %v]", statusCode, body, err)
		} else {
			var data []GameServerInfo
			if err = sonic.Unmarshal(body, &data); err != nil {
				hlog.CtxErrorf(ctx, "failed to decode game server data: [status = %d] [body = %s], [err = %v]", statusCode, body, err)
			} else {
				for _, v := range data {
					g.gameServerMap[utilk.Hash(v.GameId, v.ServerId)] = v
				}
			}
		}
	})
}

func (g *Metadata) FindServer(gameId, serverId string) (r GameServerInfo, ok bool) {
	r, ok = g.gameServerMap[utilk.Hash(gameId, serverId)]
	return
}

func (g *Metadata) FindProduct(productId string, platformId int64) (r Product, ok bool) {
	r, ok = g.productMap[utilk.Hash(productId, platformId)]
	return
}

func (g *Metadata) FindChannel(code string) (r Channel, ok bool) {
	r, ok = g.channelMap[code]
	return
}

type metaUpdateJob struct {
	ctx      context.Context
	executed bool
}

func (m *metaUpdateJob) Run() {
	if !p.RemoteUrl.EnableUpdate {
		if !m.executed {
			// at least execute at a time
			m.executed = true
			mt.refresh(m.ctx)
		}
		return
	}
	mt.refresh(m.ctx)
}

func MetaUpdateJob(ctx context.Context) cron.Job {
	return &metaUpdateJob{ctx: ctx}
}
