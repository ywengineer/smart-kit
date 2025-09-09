package vk

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// PaymentCheckResponse -------------------------- 3. 支付检验相关结构体 --------------------------
// 支付检验顶层响应
type PaymentCheckResponse struct {
	Code      string       `json:"code"`
	Message   *string      `json:"message"`
	Body      *PaymentBody `json:"body"`
	Timestamp string       `json:"timestamp"`
}

// PaymentBody 支付检验核心数据（对应官方文档 body 字段）
type PaymentBody struct {
	InvoiceId        int64        `json:"invoiceId"`        // 账单ID（检验入参）
	InvoiceDate      string       `json:"invoiceDate"`      // 账单创建时间
	RefundDate       *string      `json:"refundDate"`       // 退款时间（仅 REFUNDED 有值）
	InvoiceStatus    string       `json:"invoiceStatus"`    // 账单状态（核心校验字段）
	DeveloperPayload *string      `json:"developerPayload"` // 自定义订单信息
	AppId            int64        `json:"appId"`            // 应用ID（需与配置匹配）
	OwnerCode        int64        `json:"ownerCode"`        // 应用所有者编码
	PaymentInfo      *PaymentInfo `json:"paymentInfo"`      // 支付详情（CREATED 状态为空）
	PurchaseId       string       `json:"purchaseId"`       // 唯一购买UUID
	Order            *OrderInfo   `json:"order"`            // 订单信息
}

// PaymentInfo 支付详情
type PaymentInfo struct {
	PaymentDate    *string `json:"paymentDate"`    // 支付时间
	MaskedPan      *string `json:"maskedPan"`      // 掩码卡号（如 **1111）
	PaymentSystem  *string `json:"paymentSystem"`  // 支付系统（如 Visa）
	PaymentWay     *string `json:"paymentWay"`     // 支付方式（如 SberPay）
	PaymentWayCode *string `json:"paymentWayCode"` // 支付方式编码
	BankName       *string `json:"bankName"`       // 发卡行名称
}

// OrderInfo 订单信息
type OrderInfo struct {
	OrderId       string  `json:"orderId"`       // 订单UUID
	OrderNumber   *string `json:"orderNumber"`   // 订单编号（可选）
	VisualName    string  `json:"visualName"`    // 操作名称
	AmountCreate  int64   `json:"amountCreate"`  // 创建时金额（最小货币单位，如 копейки）
	AmountCurrent int64   `json:"amountCurrent"` // 当前金额（含折扣）
	Currency      string  `json:"currency"`      // 货币代码（如 RUB）
	ItemCode      string  `json:"itemCode"`      // 产品编码（需与配置匹配）
	Description   string  `json:"description"`   // 订单描述
	Language      string  `json:"language"`      // 描述语言
}

// -------------------------- 4. 核心：支付订单检验逻辑 --------------------------
// PaymentChecker 支付检验器（整合令牌管理器）
type PaymentChecker struct {
	tokenManager *TokenManager
	config       RustoreConfig
}

// NewPaymentChecker 初始化支付检验器
func NewPaymentChecker(config RustoreConfig) (*PaymentChecker, error) {
	tm, err := NewTokenManager(config)
	if err != nil {
		return nil, fmt.Errorf("初始化令牌管理器失败：%w", err)
	}
	return &PaymentChecker{
		tokenManager: tm,
		config:       config,
	}, nil
}

// CheckPaymentByInvoiceId 按 invoiceId 检验支付状态（核心方法）
// 返回值：
// - *PaymentBody：支付成功且校验通过时返回核心数据
// - error：校验失败（网络错误、接口错误、状态非法等）
func (pc *PaymentChecker) CheckPaymentByInvoiceId(invoiceId int64) (*PaymentBody, error) {
	// 1. 获取有效 JWE 令牌
	token, err := pc.tokenManager.getToken()
	if err != nil {
		return nil, fmt.Errorf("获取令牌失败：%w", err)
	}

	// 2. 构造支付检验接口地址
	checkURL := fmt.Sprintf(prodCheckURL, invoiceId)
	if pc.config.IsSandbox {
		checkURL = fmt.Sprintf(sandboxCheckURL, invoiceId)
	}

	// 3. 发送 GET 请求（带令牌授权）
	req, err := http.NewRequest("GET", checkURL, nil)
	if err != nil {
		return nil, fmt.Errorf("检验请求创建失败：%w", err)
	}
	req.Header.Set("Public-Token", token) // 关键：令牌放入请求头
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("检验请求发送失败：%w", err)
	}
	defer resp.Body.Close()

	// 4. 解析检验响应
	var checkResp PaymentCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&checkResp); err != nil {
		return nil, fmt.Errorf("检验响应解析失败：%w", err)
	}

	// 5. 校验接口响应码（仅 OK 为成功）
	if strings.ToUpper(checkResp.Code) != "OK" {
		msg := "未知错误"
		if checkResp.Message != nil {
			msg = *checkResp.Message
		}
		return nil, fmt.Errorf("检验接口返回错误：code=%s, message=%s", checkResp.Code, msg)
	}

	// 6. 校验响应体非空
	if checkResp.Body == nil {
		return nil, fmt.Errorf("检验响应体为空（code=OK但body=null）")
	}
	paymentBody := checkResp.Body

	// 7. 核心业务校验（根据实际需求调整）
	// 7.1 校验账单状态（仅 CONFIRMED 为支付成功）
	if paymentBody.InvoiceStatus != "CONFIRMED" {
		return nil, fmt.Errorf("支付状态非法：当前状态为 %s（需为 CONFIRMED）", paymentBody.InvoiceStatus)
	}

	// 7.2 校验应用ID（可选，确保请求对应正确应用）
	// 请替换为你的实际 AppId（从 Rustore 控制台获取）
	expectedAppId := int64(123456) // 示例值，需修改！
	if paymentBody.AppId != expectedAppId {
		return nil, fmt.Errorf("应用ID不匹配：实际 %d，期望 %d", paymentBody.AppId, expectedAppId)
	}

	// 7.3 校验产品编码（可选，确保购买的是正确产品）
	// expectedItemCode := "your-item-code" // 示例值，需修改！
	// if paymentBody.Order != nil && paymentBody.Order.ItemCode != expectedItemCode {
	// 	return nil, fmt.Errorf("产品编码不匹配：实际 %s，期望 %s", paymentBody.Order.ItemCode, expectedItemCode)
	// }

	// 7.4 校验金额（可选，确保支付金额与订单一致）
	// expectedAmount := int64(1000) // 示例：1000 копейки = 10 RUB，需修改！
	// if paymentBody.Order != nil && paymentBody.Order.AmountCurrent != expectedAmount {
	// 	return nil, fmt.Errorf("支付金额不匹配：实际 %d，期望 %d", paymentBody.Order.AmountCurrent, expectedAmount)
	// }

	// 8. 所有校验通过，返回支付核心数据
	return paymentBody, nil
}

// -------------------------- 5. 示例：调用支付检验 --------------------------
func main() {
	// 1. 配置参数（需替换为你的实际值！）
	config := RustoreConfig{
		ClientID:     "your-client-id",     // 从控制台获取
		ClientSecret: "your-client-secret", // 从控制台获取（严格保密）
		IsSandbox:    true,                 // 测试用 true，生产用 false
	}

	// 2. 初始化支付检验器
	checker, err := NewPaymentChecker(config)
	if err != nil {
		fmt.Printf("初始化支付检验器失败：%v\n", err)
		return
	}

	// 3. 待检验的 invoiceId（从 Pay SDK 回调或订单记录获取）
	invoiceId := int64(123456789) // 示例值，需修改！

	// 4. 执行支付检验
	paymentData, err := checker.CheckPaymentByInvoiceId(invoiceId)
	if err != nil {
		fmt.Printf("支付检验失败：%v\n", err)
		return
	}

	// 5. 检验成功，处理业务逻辑（如更新订单状态、发放商品等）
	fmt.Println("=== 支付检验成功 ===")
	// 打印核心信息（生产环境建议用日志记录，而非直接打印）
	fmt.Printf("Invoice ID: %d\n", paymentData.InvoiceId)
	fmt.Printf("支付状态: %s\n", paymentData.InvoiceStatus)
	fmt.Printf("支付时间: %s\n", *paymentData.PaymentInfo.PaymentDate)
	fmt.Printf("订单金额: %d %s（%0.2f 元）\n",
		paymentData.Order.AmountCurrent,
		paymentData.Order.Currency,
		float64(paymentData.Order.AmountCurrent)/100) // 转换为元（假设最小单位是分）
	if paymentData.DeveloperPayload != nil {
		fmt.Printf("自定义订单信息: %s\n", *paymentData.DeveloperPayload)
	}

	// 后续业务逻辑：更新数据库订单状态、调用业务接口发放权益等
	// ...
}
