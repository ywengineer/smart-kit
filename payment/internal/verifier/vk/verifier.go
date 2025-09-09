package vk

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"

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
	InvoiceId        int64    `json:"invoiceId"`        // 账单ID（检验入参）
	InvoiceDate      string   `json:"invoiceDate"`      // 账单创建时间
	RefundDate       string   `json:"refundDate"`       // 退款时间（仅 REFUNDED 有值）
	InvoiceStatus    string   `json:"invoiceStatus"`    // 账单状态（核心校验字段）
	DeveloperPayload string   `json:"developerPayload"` // 自定义订单信息
	AppId            int64    `json:"appId"`            // 应用ID（需与配置匹配）
	OwnerCode        int64    `json:"ownerCode"`        // 应用所有者编码
	PurchaseId       string   `json:"purchaseId"`       // 唯一购买UUID
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

// NewRustore 初始化支付检验器
func NewRustore(config RustoreConfig) (*Rustore, error) {
	tm, err := NewTokenManager(config)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to initialize the token manager")
	}
	return &Rustore{
		tokenManager: tm,
		config:       config,
	}, nil
}

// Verify 按 invoiceId 检验支付状态（核心方法）
// 返回值：
// - *model.Purchase：支付成功且校验通过时返回核心数据
// - error：校验失败（网络错误、接口错误、状态非法等）
func (pc *Rustore) Verify(ctx context.Context, invoiceId string) (*model.Purchase, error) {
	// 获取有效 JWE 令牌
	token, err := pc.tokenManager.getToken()
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to obtain the token")
	}
	//
	var verifyURL string
	if pc.config.IsSandbox {
		verifyURL = path.Join(sandboxVerifyURL, invoiceId)
	} else {
		verifyURL = path.Join(prodVerifyURL, invoiceId)
	}
	//
	statusCode, resp, err := rpcs.GetDefaultRpc().Get(context.Background(), verifyURL, http.Header{
		"Public-Token": []string{token},
	})
	if statusCode != consts.StatusOK {
		return nil, errors.New(fmt.Sprintf("Failed to get rustore purchase with status: %d", statusCode))
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
	if !pc.config.IsValidApp(strconv.FormatInt(vr.Body.AppId, 10)) {
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
	p := &model.Purchase{}
	return p, nil
}

//func main() {
//	// 1. 配置参数（需替换为你的实际值！）
//	config := RustoreConfig{
//		ClientID:     "your-client-id",     // 从控制台获取
//		ClientSecret: "your-client-secret", // 从控制台获取（严格保密）
//		IsSandbox:    true,                 // 测试用 true，生产用 false
//	}
//
//	// 2. 初始化支付检验器
//	checker, err := NewPaymentChecker(config)
//	if err != nil {
//		fmt.Printf("初始化支付检验器失败：%v\n", err)
//		return
//	}
//
//	// 3. 待检验的 invoiceId（从 Pay SDK 回调或订单记录获取）
//	invoiceId := int64(123456789) // 示例值，需修改！
//
//	// 4. 执行支付检验
//	paymentData, err := checker.CheckPaymentByInvoiceId(invoiceId)
//	if err != nil {
//		fmt.Printf("支付检验失败：%v\n", err)
//		return
//	}
//
//	// 5. 检验成功，处理业务逻辑（如更新订单状态、发放商品等）
//	fmt.Println("=== 支付检验成功 ===")
//	// 打印核心信息（生产环境建议用日志记录，而非直接打印）
//	fmt.Printf("Invoice ID: %d\n", paymentData.InvoiceId)
//	fmt.Printf("支付状态: %s\n", paymentData.InvoiceStatus)
//	fmt.Printf("支付时间: %s\n", *paymentData.PaymentInfo.PaymentDate)
//	fmt.Printf("订单金额: %d %s（%0.2f 元）\n",
//		paymentData.Order.AmountCurrent,
//		paymentData.Order.Currency,
//		float64(paymentData.Order.AmountCurrent)/100) // 转换为元（假设最小单位是分）
//	if paymentData.DeveloperPayload != nil {
//		fmt.Printf("自定义订单信息: %s\n", *paymentData.DeveloperPayload)
//	}
//
//	// 后续业务逻辑：更新数据库订单状态、调用业务接口发放权益等
//	// ...
//}
