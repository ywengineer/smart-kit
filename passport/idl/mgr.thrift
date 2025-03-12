namespace go mgr

struct MgrUser {
    1: string Account (api.body="account,required", api.vd="len($) <= 20");
    2: string Password (api.body="password,required", api.vd="len($) <= 20");
}

service MgrSignService {
    string Sign(1: MgrUser user) (api.post="/mgr/sign")
}

service MgrWhiteListService {
    bool Add(1:i64 id) (api.get="/mgr/white-list/add")
    bool Remove(1:i64 id) (api.get="/mgr/white-list/rm")
}