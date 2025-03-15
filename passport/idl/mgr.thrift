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

struct WhiteListData {
    1: i64 Id
    2: i64 CreateAt
    3: i64 UpdateAt
    4: i64 DeleteAt
    5: i64 PassportId
}

struct WhiteListPageReq {
    1: i32 PageNo (api.query="page,required", api.vd="$ > 0")
    2: i32 PageSize (api.query="page_size,required", api.vd="$ >= 50 && $ <= 100")
}

struct WhiteListPageRes {
    1: i32 Page
    2: i32 PageSize
    3: list<WhiteListData> Data
    4: i64 Total
    5: i32 MaxPage
}

service MgrWhiteListService {
    bool Add(1:WhiteListReq id) (api.get="/mgr/white-list/add")
    bool Remove(1:WhiteListReq id) (api.get="/mgr/white-list/rm")
}

service MgrWhiteListPageService {
    WhiteListPageRes Page(1:WhiteListPageReq req) (api.get="/mgr/white-list/page")
}