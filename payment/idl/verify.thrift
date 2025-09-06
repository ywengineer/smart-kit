namespace go purchase

include "common.thrift"

// go.tag 中添加 Protobuf tag：wire_type=varint，tag=1，可选，名称 id
// 1: required i32 id (go.tag = 'protobuf:"varint,1,opt,name=id"')
// 类型映射正确：Thrift 类型与 Protobuf 类型的 wire_type 需正确对应（如 i32 → varint，string → bytes），参考下表：
// Thrift 类型	Protobuf 类型	wire_type
// i32	int32	varint
// i64	int64	varint
// string	string	bytes
// bool	bool	varint
// double	double	fixed64
// list<T>	repeated T	同 T 的类型
// map<K,V>	map<K,V>	bytes
struct VerifyReq {
    1: string Receipt (api.body="receipt,required", go.tag='protobuf:"bytes,1,req,name=receipt"'); // 订单数据
    2: string GameId (api.body="gameId,required", go.tag='protobuf:"bytes,2,req,name=gameId"');   // 游戏ID
    3: string ServerId (api.body="serverId,required", go.tag='protobuf:"bytes,3,req,name=serverId"');       // 服务器ID required
    4: string Passport (api.body="passport,required", go.tag='protobuf:"bytes,4,req,name=passport"');       // 账号ID required
    5: string PlayerId (api.body="playerId,required", go.tag='protobuf:"bytes,5,req,name=playerId"');       // 玩家ID required
    6: string PlayerName (api.body="playerName,required", go.tag='protobuf:"bytes,6,req,name=playerName"');     // 玩家名称 required
    7: string Channel (api.body="channel,required", go.tag='protobuf:"bytes,7,req,name=channel"');        // 支付渠道. google, apple required
    8: string Signature (api.body="signature", go.tag='protobuf:"bytes,8,opt,name=signature"');      // Google订单检验签名数据 optional
    9: string Locale (api.body="locale", go.tag='protobuf:"bytes,9,opt,name=locale"');      // 区域码（当前设置的区域） optional
    10: i64 CreateTime (api.body="createTime", go.tag='protobuf:"varint,10,opt,name=receipt"');      // 账号注册时间 optional
}

service VerifyService {
    common.ApiResult Verify(1:VerifyReq req) (api.get="/verify")
}
