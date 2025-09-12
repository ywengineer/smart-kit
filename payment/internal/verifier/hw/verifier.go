package hw

import (
	"context"

	"gitee.com/ywengineer/smart-kit/payment/internal/config"
	"gitee.com/ywengineer/smart-kit/payment/internal/verifier/inf"
	"gitee.com/ywengineer/smart-kit/payment/pkg/model"
)

type huawei struct {
	config config.ChannelProperty
}

func (r *huawei) Verify(ctx context.Context, receipt string) (*model.Purchase, error) {
	//bodyMap := map[string]string{"purchaseToken": purchaseToken, "productId": productId}
	//url := getOrderUrl(accountFlag)+ "/applications/purchases/tokens/verify"
	//bodyBytes, err := SendRequest(url, bodyMap)
	//if err != nil {
	//	log.Printf("err is %s", err)
	//}
	//// TODO: display the response as string in console, you can replace it with your business logic.
	//log.Printf("%s", bodyBytes)
	//TODO implement me
	panic("implement me")
}

// New 初始化
func New(config config.ChannelProperty) (inf.Verifier, error) {
	hw := huawei{config: config}
	return &hw, nil
}
