namespace go openapi

enum Channel {
    SelfHost = 1, // 自托管
    Wx = 2,
    QQ = 3,
    Huawei = 4,
    XiaoMi = 5,
    Google = 6,
    Apple = 7,
    RuStore = 8,
}

struct HealthResp {
    1: bool Healthy (api.body="healthy");
}

service HealthService {
    HealthResp Health() (api.get="/health")
}
