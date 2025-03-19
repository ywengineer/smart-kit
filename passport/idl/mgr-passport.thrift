namespace go mgr.pst

struct PassportBindData {
    1: i64 Id
    2: i64 CreateAt
    3: i64 UpdateAt
    4: i64 DeleteAt
    5: string BindType
    6: string BindId
    7: string Token
    8: string SocialName
    9: i64 Gender
    10: string IconUrl
}

struct PassportData {
    1: i64 Id
    2: i64 CreateAt
    3: i64 UpdateAt
    4: i64 DeleteAt
    5: string DeviceId
    6: string Adid
    7: string SystemType
    8: string Locale
    9: list<PassportBindData> Bounds
    10: string Extra
}

struct MgrPassportDetailReq {
    1: i64 Id (api.query="id,required", api.vd="$ > 0")
}

struct MgrPassportPageReq {
    1: i32 PageNo (api.query="page,required", api.vd="$ > 0")
    2: i32 PageSize (api.query="page_size,required", api.vd="$ >= 50 && $ <= 100")
    3: i64 PassportId (api.query="passport_id")
}

struct MgrPassportPageRes {
    1: i32 Page
    2: i32 PageSize
    3: list<PassportData> Data
    4: i64 Total
    5: i32 MaxPage
}

service MgrPassportInfoService {
    PassportData Detail(1:MgrPassportDetailReq req) (api.get="/mgr/passport/detail")
}

service MgrPassportPageService {
    MgrPassportPageRes Page(1: MgrPassportPageReq req) (api.get="/mgr/passport/page")
}