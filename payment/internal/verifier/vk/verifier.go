package vk

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitee.com/ywengineer/smart-kit/payment/internal/config"
	"gitee.com/ywengineer/smart-kit/pkg/utilk"

	"gitee.com/ywengineer/smart-kit/payment/internal/verifier/inf"
	"gitee.com/ywengineer/smart-kit/payment/pkg/model"
	"gitee.com/ywengineer/smart-kit/pkg/rpcs"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/pkg/errors"
)

// verifyResp 支付检验顶层响应
type verifyResp struct {
	Code      string       `json:"code"`
	Message   string       `json:"message"`
	Body      *paymentBody `json:"body"`
	Timestamp string       `json:"timestamp"`
}

// paymentBody 支付检验核心数据（对应官方文档 body 字段）
// InvoiceStatus
//
//	CREATED - создан;
//	EXECUTED - 启动付款流程;
//	CONFIRMED — 成功付款的最终状态，现金从买方注销
//	CANCELLED - 用户在付款开始前取消;
//	REJECTED - 拒付（资金不足、CVC无效或其他原因）;
//	EXPIRED - 付款时间已过;
//	PAID - 购买消费商品，资金已成功保存，购买等待开发商确认;
//	REVERSED - 购买被消费商品：未收到持有确认请求，持有被取消，退款给买方;
//	REFUNDED - 已回收的资金将退还给买方.
//	REFUNDING - 退款中.
type paymentBody struct {
	InvoiceId        int64  `json:"invoiceId"`        // 账单ID（检验入参）
	InvoiceDate      string `json:"invoiceDate"`      // 账单创建时间
	RefundDate       string `json:"refundDate"`       // 退款时间（仅 REFUNDED 有值）
	InvoiceStatus    string `json:"invoiceStatus"`    // 账单状态（核心校验字段）
	DeveloperPayload string `json:"developerPayload"` // 自定义订单信息
	AppId            int64  `json:"appId"`            // 应用ID（需与配置匹配）
	OwnerCode        int64  `json:"ownerCode"`        // 应用所有者编码
	PurchaseId       string `json:"purchaseId"`       // 唯一购买UUID
	PaymentInfo      struct { // 支付详情（CREATED 状态为空）
		PaymentDate    string `json:"paymentDate"`    // 支付时间
		MaskedPan      string `json:"maskedPan"`      // 掩码卡号（如 **1111）
		PaymentSystem  string `json:"paymentSystem"`  // 支付系统（如 Visa）
		PaymentWay     string `json:"paymentWay"`     // 支付方式（如 SberPay）
		PaymentWayCode string `json:"paymentWayCode"` // 支付方式编码
		BankName       string `json:"bankName"`       // 发卡行名称
	} `json:"paymentInfo"`
	Order struct {
		OrderId       string  `json:"orderId"`       // 订单UUID
		OrderNumber   *string `json:"orderNumber"`   // 订单编号（可选）
		VisualName    string  `json:"visualName"`    // 操作名称
		AmountCreate  int64   `json:"amountCreate"`  // 创建时金额（最小货币单位，如 копейки）
		AmountCurrent int64   `json:"amountCurrent"` // 当前金额（含折扣）
		Currency      string  `json:"currency"`      // 货币代码（如 RUB）
		ItemCode      string  `json:"itemCode"`      // 产品编码（需与配置匹配）
		Description   string  `json:"description"`   // 订单描述
		Language      string  `json:"language"`      // 描述语言
	} `json:"order"` // 订单信息
}

// Rustore 支付检验器（整合令牌管理器）
type Rustore struct {
	tokenManager *TokenManager
	config       RustoreConfig
}

// New 初始化支付检验器
func New(cp config.ChannelProperty) (inf.Verifier, error) {
	c := RustoreConfig{
		ClientID:     cp.ClientID,
		ClientSecret: cp.ClientSecret,
		IsSandbox:    cp.Sandbox,
		Apps:         cp.Apps,
	}
	tm, err := NewTokenManager(c)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to initialize the token manager")
	}
	return &Rustore{
		tokenManager: tm,
		config:       c,
	}, nil
}

// Verify 按 invoiceId 检验支付状态（核心方法）
// 返回值：
// - *model.Purchase：支付成功且校验通过时返回核心数据
// - error：校验失败（网络错误、接口错误、状态非法等）
func (rustore *Rustore) Verify(ctx context.Context, invoiceId string) (*model.Purchase, error) {
	// 获取有效 JWE 令牌
	token, err := rustore.tokenManager.getToken()
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to obtain the token")
	}
	//
	var verifyURL string
	if rustore.config.IsSandbox {
		verifyURL = sandboxVerifyURL + invoiceId
	} else {
		verifyURL = prodVerifyURL + invoiceId
	}
	//
	statusCode, resp, err := rpcs.GetDefaultRpc().Get(ctx, verifyURL, http.Header{
		"Public-Token": []string{token},
	})
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to get rustore purchase with error")
	} else if statusCode != consts.StatusOK {
		return nil, errors.New(fmt.Sprintf("Failed to get rustore purchase with status: %d, resp = %s", statusCode, string(resp)))
	}
	//
	var vr verifyResp
	if err = sonic.Unmarshal(resp, &vr); err != nil || !strings.EqualFold(vr.Code, "OK") || vr.Body == nil {
		return nil, errors.WithMessagef(err, "failed to parse rustore purchase json: %s", string(resp))
	}
	// 校验账单状态（仅 CONFIRMED 为支付成功）
	if !strings.EqualFold(vr.Body.InvoiceStatus, "CONFIRMED") {
		return nil, inf.IncompletePurchase
	}
	// 校验应用ID（可选，确保请求对应正确应用）
	if !rustore.config.IsValidApp(strconv.FormatInt(vr.Body.AppId, 10)) {
		return nil, inf.OtherAppPurchase
	}
	// 校验产品编码（可选，确保购买的是正确产品）
	// expectedItemCode := "your-item-code" // 示例值，需修改！
	// if paymentBody.Order != nil && paymentBody.Order.ItemCode != expectedItemCode {
	// 	return nil, fmt.Errorf("产品编码不匹配：实际 %s，期望 %s", paymentBody.Order.ItemCode, expectedItemCode)
	// }
	// 校验金额（可选，确保支付金额与订单一致）
	// expectedAmount := int64(1000) // 示例：1000 копейки = 10 RUB，需修改！
	// if paymentBody.Order != nil && paymentBody.Order.AmountCurrent != expectedAmount {
	// 	return nil, fmt.Errorf("支付金额不匹配：实际 %d，期望 %d", paymentBody.Order.AmountCurrent, expectedAmount)
	// }
	// 所有校验通过，返回支付核心数据
	p, err := rustore.convert(vr.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to convert rustore payment 2 model.Purchase")
	}
	//
	p.ReceiptResult = string(resp)
	p.TestOrder = rustore.config.IsSandbox
	return p, nil
}

type developerPayload struct {
	SystemType string `json:"systemType"`
	AppVersion string `json:"appVersion"`
	BundleId   string `json:"bundleId"`
}

// convert rustore payment 2 model.Purchase
// need reset data after conver
// - TestOrder
// - FreeTrail
func (rustore *Rustore) convert(data *paymentBody) (*model.Purchase, error) {
	var payload developerPayload
	var err error
	p := &model.Purchase{}
	_ = sonic.Unmarshal(utilk.S2b(data.DeveloperPayload), &payload)
	p.SystemType = payload.SystemType
	// ------------------------------------------------------------------------------------------
	p.TransactionId = strconv.FormatInt(data.InvoiceId, 64)
	p.TestOrder, p.FreeTrail = false, false
	p.OriginalTransactionId = p.TransactionId                                    // 同transaction_id
	p.PurchaseDate, err = time.Parse(time.RFC3339, data.PaymentInfo.PaymentDate) // 商品的购买时间（从新纪年（1970 年 1 月 1 日）开始计算的毫秒数）。
	if err != nil {
		return nil, err
	}
	p.OriginalPurchaseDate = p.PurchaseDate           // 对于恢复的transaction对象，该键对应了原始的交易日期
	p.ExpireDate = nil                                // 普通消耗类型,没有过期时间
	p.OriginalApplicationVersion = payload.AppVersion // 开发者指定的字符串，包含订单的补充信息。您可以在发起 getBuyIntent 请求时为此字段指定一个值。
	p.AppItemId = data.Order.ItemCode                 // 用于对给定商品和用户对进行唯一标识的令牌。
	p.VersionExternalIdentifier = 0                   // 用来标识程序修订数。该键在sandbox环境下不存在
	// ------------------------------------------------------------------------------------------
	p.Quantity = 1                    // 购买商品的数量
	p.Status = 0                      // 订单的购买状态。可能的值为 0（已购买）、1（已取消）或者 2（已退款）
	p.ProductId = data.Order.ItemCode // 商品的商品 ID。每种商品都有一个商品 ID，您必须通过 Google Play Developer Console 在应用的商品列表中指定此 ID。
	p.BundleId = payload.BundleId     //  Android程序的bundle标识
	return p, nil
}
