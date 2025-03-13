namespace go mgr

struct MgrUser {
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
    MgrSignRes Sign(1: MgrUser user) (api.post="/mgr/sign")
}

service MgrWhiteListService {
    bool Add(1:i64 id) (api.get="/mgr/white-list/add")
    bool Remove(1:i64 id) (api.get="/mgr/white-list/rm")
}