package vk

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
