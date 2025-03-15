namespace go mgr

struct MgrSignReq {
    1: string Account (api.body="account,required", api.vd="len($) <= 20");
    2: string Password (api.body="password,required", api.vd="len($) <= 20");
}

struct MgrSignRes {
    1: i64 Id
    2: string Act
    3: string Name
    4: i64 Dept
    5: string Title
    6: string Token
}

service MgrSignService {
    MgrSignRes Sign(1: MgrSignReq user) (api.post="/mgr/sign")
}

struct WhiteListReq {
    1: i64 id (api.query="id,required", api.vd="$ > 0")
}

service MgrWhiteListService {
    bool Add(1:WhiteListReq id) (api.get="/mgr/white-list/add")
    bool Remove(1:WhiteListReq id) (api.get="/mgr/white-list/rm")
}