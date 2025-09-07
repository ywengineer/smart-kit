package model

import (
	"time"
)

type Purchase struct {
	BaseLog               `json:",inline"`
	Status                int        `json:"status" redis:"status"`                                                 // 订单状态
	ProductId             string     `json:"product_id" redis:"product_id"`                                         // 商品的标识
	Quantity              int        `json:"quantity" redis:"quantity"`                                             // 购买商品的数量
	TransactionId         string     `json:"transaction_id" redis:"transaction_id" gorm:"uniqueIndex:idx_order_id"` // 交易的标识
	OriginalTransactionId string     `json:"original_transaction_id" redis:"original_transaction_id"`               // 对于恢复的transaction对象，该键对应了原始的transaction标识
	PurchaseDate          time.Time  `json:"purchase_date" redis:"purchase_date"`                                   // 交易的日期(UTC)
	OriginalPurchaseDate  time.Time  `json:"original_purchase_date" redis:"original_purchase_date"`                 // 对于恢复的transaction对象，该键对应了原始的交易日期(UTC)
	ExpireDate            *time.Time `json:"expire_date" redis:"expire_date"`                                       // The expiration date for the subscription, expressed as the number of milliseconds since January 1, 1970, 00:00:00 GMT.
	FinishDate            *time.Time `json:"finish_date" redis:"finish_date"`                                       // 订单处理结束时间
	// GooglePlay 用于对给定商品和用户对进行唯一标识的令牌
	// App Store 用来标识程序的字符串。一个服务器可能需要支持多个server的支付功能，可以用这个标识来区分程序。链接sandbox用来测试的程序的不到这个值，因此该键不存在。
	AppItemId                  string `json:"app_item_id" redis:"app_item_id"`
	VersionExternalIdentifier  int    `json:"version_external_identifier" redis:"version_external_identifier"`   //  用来标识程序修订数。该键在sandbox环境下不存在
	OriginalApplicationVersion string `json:"original_application_version" redis:"original_application_version"` // App version
	BundleId                   string `json:"bundle_id" redis:"bundle_id"`                                       //  iPhone程序的bundle标识
	Notified                   bool   `json:"notified" redis:"notified"`                                         // 是否已通知
	NotifiedTimes              int    `json:"notified_times" redis:"notified_times"`                             // 通知次数
	Receipt                    string `json:"receipt" redis:"receipt" gorm:"type:text"`                          // 订单原始凭证
	ReceiptResult              string `json:"receipt_result" redis:"receipt_result" gorm:"type:text"`            // 订单验证结果
	Channel                    string `json:"channel" redis:"channel"`                                           // 支付渠道
	FreeTrail                  bool   `json:"free_trail" redis:"free_trail"`                                     // 试用订单
	TestOrder                  bool   `json:"test_order" redis:"test_order"`                                     // 测试订单
	//
	RankPoints int   `json:"rank_points" redis:"rank_points"` // 段位积分
	Price      int   `json:"price" redis:"price"`             // 价格(分)
	Credits    int64 `json:"credits" redis:"credits"`         // 点券剩余量
	Money      int64 `json:"money" redis:"money"`             // 钻石剩余量
	Coin       int64 `json:"coin" redis:"coin"`               // 当前金币
}

func (p Purchase) GetExpiredTime() int64 {
	if p.ExpireDate == nil {
		return 0
	}
	return p.ExpireDate.Unix()
}
