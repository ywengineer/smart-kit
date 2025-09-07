namespace go simulate

include "common.thrift"

struct SimulateReq {
    1: string Passport (api.body="passport,required");
    2: string PlayerId (api.body="playerId,required");
    3: string PlayerName (api.body="playerName,required");
    4: string GameId (api.body="gameId,required");
    5: string ServerId (api.body="serverId,required");
    6: string PlatformId (api.body="platformId,required");
    7: string ProductId (api.body="productId,required");
    8: string OrderId (api.body="orderId,required");
}

service SimulateService {
    common.ApiResult Simulate(1:SimulateReq req) (api.post="/simulate")
}
