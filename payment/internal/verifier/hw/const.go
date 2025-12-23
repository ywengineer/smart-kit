package hw

import (
	"encoding/base64"
	"fmt"
)

const (
	// 令牌接口地址
	prodTokenURL = "https://oauth-login.cloud.huawei.com/oauth2/v3/token"
	// 固定参数
	tokenExpiryBuffer = 30 // 令牌刷新缓冲时间（提前30秒刷新，避免过期）
)

// Config config
type Config struct {
	ClientID     string // 控制台获取的 Client ID
	ClientSecret string // 控制台获取的 application public key, base64 encode
	IsSandbox    bool   // 是否启用沙箱环境
	ApiRoot      string
}

type PurchaseKind int

func (p PurchaseKind) Code() int {
	return int(p)
}

const (
	Consumable PurchaseKind = iota
	NonConsumable
	Subscription
)

type PurchaseState int

func (p PurchaseState) Code() int {
	return int(p)
}

const (
	Init      = PurchaseState(-1)
	Confirmed = PurchaseState(0)
	Canceled  = PurchaseState(1)
	Refunded  = PurchaseState(2)
	Expired   = PurchaseState(3)
)

// IapPurchaseDetails 华为HMS支付API返回的用户购买详情数据模型，包含消耗型、非消耗型及订阅型商品信息
type IapPurchaseDetails struct {
	// 必选参数 - 所有商品类型通用
	ApplicationId int64         `json:"applicationId"` // 应用ID
	AutoRenewing  bool          `json:"autoRenewing"`  // 自动续订状态：消耗/非消耗型固定为false；订阅型true=活动且自动续订（含宽限期），false=已取消
	OrderId       string        `json:"orderId"`       // 订单ID，唯一标识收费收据，新收据/订阅续期生成新ID
	Kind          PurchaseKind  `json:"kind"`          // 商品类别：0=消耗型，1=非消耗型，2=订阅型
	ProductId     string        `json:"productId"`     // 商品唯一ID（PMS维护或购买时传入，验签后需校验）
	PurchaseState PurchaseState `json:"purchaseState"` // 订单交易状态：-1=初始化，0=已购买，1=已取消，2=已撤销/已退款，3=待处理
	PurchaseToken string        `json:"purchaseToken"` // 购买令牌（唯一标识商品-用户关系，订阅续期不变；建议加密存储，预留128位长度）
	LastOrderId   string        `json:"lastOrderId"`   // 上次续期订单ID（仅订阅场景），首次购买与orderId相同

	// 可选参数 - 所有商品类型通用
	PackageName        string `json:"packageName,omitempty"`        // 应用安装包名
	ProductName        string `json:"productName,omitempty"`        // 商品名称
	PurchaseTime       int64  `json:"purchaseTime,omitempty"`       // 购买时间（UTC毫秒时间戳）
	PurchaseTimeMillis int64  `json:"purchaseTimeMillis,omitempty"` // 历史兼容字段，同PurchaseTime（新接入无需关注）
	DeveloperPayload   string `json:"developerPayload,omitempty"`   // 商户侧保留信息（禁止传入个人敏感信息）
	DeveloperChallenge string `json:"developerChallenge,omitempty"` // 消耗请求自定义挑战字（仅一次性商品）
	ConsumptionState   int    `json:"consumptionState,omitempty"`   // 消耗状态（仅一次性商品）：0=未消耗，1=已消耗
	Confirmed          int    `json:"confirmed,omitempty"`          // 确认状态：0=未确认，1=已确认（无值表示无需确认，仅兼容用）
	PurchaseType       *int   `json:"purchaseType,omitempty"`       // 购买类型：0=沙盒环境，1=促销（暂不支持，正式购买不返回）
	Currency           string `json:"currency,omitempty"`           // 定价货币（ISO 4217标准，验签后需校验）
	Price              int64  `json:"price,omitempty"`              // 实际价格（单位：分，即原价*100，如501=5.01元，验签后需校验）
	Country            string `json:"country,omitempty"`            // 国家/地区码（ISO 3166标准）
	PayType            string `json:"payType,omitempty"`            // 支付方式（取值参考PayType枚举说明）
	PayOrderId         string `json:"payOrderId,omitempty"`         // 交易单号（用户支付成功后生成）
	Quantity           int    `json:"quantity,omitempty"`           // 购买数量
	AppInfo            string `json:"appInfo,omitempty"`            // App信息（预留字段）

	// 可选参数 - 仅订阅场景返回
	ProductGroup         string `json:"productGroup,omitempty"`         // 订阅型商品所属订阅组ID
	OriPurchaseTime      int64  `json:"oriPurchaseTime,omitempty"`      // 原购买时间（订阅首次成功收费时间，UTC毫秒时间戳）
	SubscriptionId       string `json:"subscriptionId,omitempty"`       // 订阅ID（用户-商品唯一对应，订阅续期不变）
	OriSubscriptionId    string `json:"oriSubscriptionId,omitempty"`    // 原订阅ID（当前订阅从其他商品切换时，关联原订阅信息）
	DaysLasted           int64  `json:"daysLasted,omitempty"`           // 已付费订阅天数（不含免费试用、促销期）
	NumOfPeriods         int64  `json:"numOfPeriods,omitempty"`         // 标准续期成功期数（0=未续期）
	NumOfDiscount        int64  `json:"numOfDiscount,omitempty"`        // 促销续期成功期数
	ExpirationDate       int64  `json:"expirationDate,omitempty"`       // 订阅过期时间（UTC毫秒时间戳，过去时间表示已过期）
	ExpirationIntent     int    `json:"expirationIntent,omitempty"`     // 订阅过期原因（仅过期订阅）：1=用户取消，2=商品不可用，3=签约异常，4=Billing错误，5=未同意涨价，6=未知（优先级1>2>3...）
	RetryFlag            int    `json:"retryFlag,omitempty"`            // 续期重试状态（仅过期订阅）：0=终止尝试，1=仍在尝试
	IntroductoryFlag     int    `json:"introductoryFlag,omitempty"`     // 是否处于促销价续期周期：1=是，0=否
	TrialFlag            int    `json:"trialFlag,omitempty"`            // 是否处于免费试用周期：1=是，0=否
	CancelTime           int64  `json:"cancelTime,omitempty"`           // 订阅撤销时间（UTC毫秒时间戳，撤销后等同于未购买）
	CancelReason         int    `json:"cancelReason,omitempty"`         // 取消原因：3=商户发起退款/撤销，2=用户升级/跨级，1=用户因App问题取消，0=其他（cancelTime空时3表示返还费用）
	NotifyClosed         int    `json:"notifyClosed,omitempty"`         // 是否关闭订阅通知（仅查询订阅关系返回）：1=是，0=否
	RenewStatus          int    `json:"renewStatus,omitempty"`          // 续期状态（仅自动续期订阅）：1=到期自动续期，0=用户停止续期
	PriceConsentStatus   int    `json:"priceConsentStatus,omitempty"`   // 商品提价用户意见：1=已同意，0=未动作（超期后订阅失效）
	RenewPrice           int64  `json:"renewPrice,omitempty"`           // 下次续期价格（单位：分，供客户端提示用户）
	SubIsvalid           bool   `json:"subIsvalid,omitempty"`           // 订阅有效性：true=已收费且未过期/在宽限期，false=未完成购买/已过期/已退款（取消后未过期仍为true）
	DeferFlag            int    `json:"deferFlag,omitempty"`            // 是否延迟结算：1=是，其他=否
	CancelWay            int    `json:"cancelWay,omitempty"`            // 取消订阅途径：0=用户，1=商户，2=华为
	CancellationTime     int64  `json:"cancellationTime,omitempty"`     // 取消续期时间（UTC毫秒时间戳，仅停止续期，不涉及退款）
	CancelledSubKeepDays int    `json:"cancelledSubKeepDays,omitempty"` // 用户取消后订阅关系保留天数（不表示已取消）
	ResumeTime           int64  `json:"resumeTime,omitempty"`           // 暂停订阅的恢复时间（UTC毫秒时间戳）
	SurveyReason         int    `json:"surveyReason,omitempty"`         // 顾客取消原因：0=其他，1=费用过高，2=技术问题，5=欺诈，7=商品切换，9=较少使用，10=有更好应用
	SurveyDetails        string `json:"surveyDetails,omitempty"`        // 顾客自定义取消原因（仅SurveyReason=0时返回）
	GraceExpirationTime  int64  `json:"graceExpirationTime,omitempty"`  // 订阅宽限期过期时间（UTC毫秒时间戳）
}

type verifyResp struct {
	ResponseCode       string `json:"responseCode"`       // 返回码。 0：成功。 其他：失败，具体请参见错误码。
	ResponseMessage    string `json:"responseMessage"`    // 响应描述。
	PurchaseTokenData  string `json:"purchaseTokenData"`  // 包含购买数据的JSON字符串，具体请参见表 IapPurchaseDetails。 该字段原样参与签名。
	DataSignature      string `json:"dataSignature"`      // purchaseTokenData基于应用RSA IAP私钥的签名信息，签名算法为signatureAlgorithm。应用请参见对返回结果验签使用IAP公钥对PurchaseTokenData的JSON字符串进行验签。
	SignatureAlgorithm string `json:"signatureAlgorithm"` // 签名算法。
}

func encodeAccessToken(t string) string {
	oriString := fmt.Sprintf("APPAT:%s", t)
	var authString = base64.StdEncoding.EncodeToString([]byte(oriString))
	var authHeaderString = fmt.Sprintf("Basic %s", authString)
	return authHeaderString
}

type developerPayload struct {
	SystemType string `json:"systemType"`
	AppVersion string `json:"appVersion"`
	GiftBoxId  int64  `json:"giftBoxId"`
}
