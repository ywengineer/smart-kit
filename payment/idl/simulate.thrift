namespace go simulate

include "common.thrift"

struct SimulateReq {
    1: string Passport;
    2: string PlayerId;
    3: string PlayerName;
    4: string GameId;
    5: string ServerId;
    6: string PlatformId;
    7: string ProductId;
    8: string OrderId;
}

service SimulateService {
    common.ApiResult Simulate(1:SimulateReq req) (api.get="/simulate")
}
