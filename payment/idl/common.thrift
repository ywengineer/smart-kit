namespace go common

struct ApiResult {
    1: string Code (api.body="code", go.tag='protobuf:"bytes,1,req,name=code"');
    2: string Message (api.body="msg", go.tag='protobuf:"bytes,2,opt,name=msg"');
    3: string ErrCode (api.body="err_code", go.tag='protobuf:"bytes,3,opt,name=errorCode"');
}