namespace go passport

enum AccountType {
    Anonymous = 1, // 匿名
    Wx = 2,
    QQ = 3,
    Sina = 4,
    Facebook = 5,
    Google = 6,
    Apple = 7,
    Telegram = 8,
    GameCenter = 9,
    EMail = 10,
    Mobile = 11,
}

enum Gender {
    Male = 1,
    Female = 2,
    Unknown = 3,
}

// 游戏登陆: POST http://xxx/login application/json
// 返回：LoginResp
struct LoginReq {
    1: AccountType Type (api.body="type,required", api.vd="$==1||$==10||$==11"); // 账号类型, 邮件、手机登陆
    2: string AppBundleId (api.body="app_bundle_id,required"); // 应用唯一标识符
    3: string Id (api.body="id,required"); // 唯一标志
    4: string AccessToken (api.body="access_token,required"); // 访问Token. 密码/验证码
    5: map<string, string> DeviceInfo (api.body="device,required", api.vd="len($) < 15 && every($, 'device_model', 'ver', 'os', 'os_ver', 'net_type', 'lang', 'locale')"); // 客户端信息[JSON]. 必须包括但不限于
}

// 注册: http://xxx/reg
// 返回：LoginResp
struct RegisterReq {
    1: AccountType Type (api.body="type,required", api.vd="$>=1 && $<=11"); // 账号类型
    2: string AppBundleId (api.body="app_bundle_id,required" api.vd="len($) <= 50"); // 应用唯一标识符
    3: map<string, string> DeviceInfo (api.body="device,required", api.vd="len($) < 15 && every($, 'device_model', 'ver', 'os', 'os_ver', 'net_type', 'lang', 'locale')"); // 客户端信息[JSON]. 必须包括但不限于

    4: string Id (api.body="id,required", api.vd="len($) > 5 && len($) <= 50");       // 第三方平台产生的唯一标志
    5: string AccessToken (api.body="access_token");  // 第三方平台API访问Token
    6: string RefreshToken (api.body="refresh_token"); // 第三方平台访问Token
    7: string Name (api.body="name", api.vd="len($) <= 20");         // 第三方平台昵称
    8: Gender Gender (api.body="gender");       // 第三方平台性别
    9: string IconUrl (api.body="icon_url", api.vd="len($) <= 200");      // 第三方平台头像
    10: string DeviceId (api.body="device_id,required", api.vd="regexp('[A-Z0-9]{10,32}')") // 设备ID
    11: string Adid (api.body="adid,required", api.vd="len($) > 0 && len($) <= 50") // 设备广告ID
}

// 绑定游戏账号: http://xxx/bin  返回成功失败 LoginResp  只有code
// 当绑定账号类型为第三方平台，且返回成功。需要将本地保存的randomID替换为第三方平台的唯一标志uid
struct BindReq {
    1: AccountType Type (api.body="type,required", api.vd="$>=2&&$<=11");    // 账号类型
    3: string BindId (api.body="bind_id,required", api.vd="len($) < 32");       // 第三方平台产生的唯一标志
    4: string AccessToken (api.body="access_token,required", api.vd="len($) < 255");  // 第三方平台API访问Token
    5: string RefreshToken (api.body="refresh_token", api.vd="len($) < 255"); // 第三方平台访问Token
    6: string Name (api.body="name", api.vd="len($) < 32");         // 第三方平台昵称
    7: Gender Gender (api.body="gender,required" api.vd="$>=1&&$<=3");       // 第三方平台性别
    8: string IconUrl (api.body="icon_url", api.vd="len($) < 255");      // 第三方平台头像
    9: string AppBundleId (api.body="app_bundle_id,required", api.vd="len($) < 100");  // 应用唯一标识符
}

// 账号登陆/注册返回结果
struct LoginResp {
    1: i64 PassportId (api.body="passport_id");              // 账号ID
    2: string Token (api.body="token");                // 随机jwt token  需要客户端自己存一份
    3: set<AccountType> Bounds (api.body="bounds");     // 已绑定的平台账号
    4: bool BrandNew (api.body="brand_new");               // 是否是全新的账号
    5: i64 CreateTime (api.body="create_time");              // 账号创建时间
}

// 登陆
service LoginService {
    LoginResp Login(1: LoginReq req) (api.post="/login")
}
// 注册
service RegisterService {
    LoginResp Register(1: RegisterReq req) (api.post="/register")
}
// 绑定
service BindService {
    set<string> Bind(1: BindReq req) (api.post="/bind")
}