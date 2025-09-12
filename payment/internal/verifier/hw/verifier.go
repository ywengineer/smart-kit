package hw

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"gitee.com/ywengineer/smart-kit/payment/internal/config"
	"gitee.com/ywengineer/smart-kit/payment/internal/verifier/inf"
	"gitee.com/ywengineer/smart-kit/payment/pkg/model"
	"gitee.com/ywengineer/smart-kit/pkg/rpcs"
	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type huawei struct {
	tokenManager *TokenManager
	config       Config
	key          *rsa.PublicKey
}

func (r *huawei) Verify(ctx context.Context, receipt string) (*model.Purchase, error) {
	var rp IapPurchaseDetails
	if err := sonic.Unmarshal(utilk.S2b(receipt), &rp); err != nil {
		return nil, errors.WithMessage(err, "Failed to unmarshal the huawei purchase receipt")
	}
	// get token
	token, err := r.tokenManager.getToken()
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to obtain the api access token")
	}
	//
	statusCode, resp, err := rpcs.GetDefaultRpc().Post(ctx, rpcs.ContentTypeJSON, r.config.ApiRoot+"/applications/purchases/tokens/verify", http.Header{
		"Authorization": []string{token},
	}, rpcs.JsonBody{V: map[string]string{"purchaseToken": rp.PurchaseToken, "productId": rp.ProductId}})

	if err != nil {
		return nil, errors.WithMessage(err, "Failed to post huawei purchase token verify with error")
	} else if statusCode != consts.StatusOK {
		return nil, errors.New(fmt.Sprintf("Failed to post huawei purchase token verify with status: %d, resp = %s", statusCode, string(resp)))
	}
	//
	var vr verifyResp
	if err = sonic.Unmarshal(resp, &vr); err != nil {
		return nil, errors.WithMessagef(err, "failed to parse huawei purchase verify api response json: %s", string(resp))
	} else if vr.ResponseCode != "0" || len(vr.PurchaseTokenData) == 0 {
		return nil, errors.New(fmt.Sprintf("huawei purchase verify api error: %s", string(resp)))
	}
	//
	var hp IapPurchaseDetails
	if err = sonic.Unmarshal(utilk.S2b(vr.PurchaseTokenData), &hp); err != nil {
		return nil, errors.WithMessagef(err, "failed to parse huawei purchase detail json: %s", string(resp))
	}
	// 校验账单状态
	if hp.PurchaseState != Confirmed {
		return nil, inf.IncompletePurchase
	}
	// 所有校验通过，返回支付核心数据
	p, err := r.convert(&hp)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to convert huawei purchase model.Purchase")
	}
	//
	p.ReceiptResult, p.Receipt = string(resp), receipt
	p.TestOrder = hp.PurchaseType != nil && *hp.PurchaseType == 0
	return p, nil
}

func (r *huawei) verifyRsaSign(content string, sign string) error {
	hashed := sha256.Sum256([]byte(content))
	signature, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(r.key, crypto.SHA256, hashed[:], signature)
}

func (r *huawei) parseSignKey() error {
	publicKeyByte, err := base64.StdEncoding.DecodeString(r.config.ClientSecret)
	if err != nil {
		return err
	}
	pub, err := x509.ParsePKIXPublicKey(publicKeyByte)
	if err != nil {
		return err
	}
	r.key = pub.(*rsa.PublicKey)
	return nil
}

// convert rustore payment 2 model.Purchase
// need reset data after conver
// - TestOrder
// - FreeTrail
func (r *huawei) convert(data *IapPurchaseDetails) (*model.Purchase, error) {
	var payload developerPayload
	p := &model.Purchase{}
	_ = sonic.Unmarshal(utilk.S2b(data.DeveloperPayload), &payload)
	p.SystemType = payload.SystemType
	// ------------------------------------------------------------------------------------------
	p.TransactionId = data.OrderId
	p.TestOrder, p.FreeTrail = false, false
	p.OriginalTransactionId = p.TransactionId                  // 同transaction_id
	p.PurchaseDate = time.UnixMilli(data.PurchaseTime).Local() // 商品的购买时间（从新纪年（1970 年 1 月 1 日）开始计算的毫秒数）。
	p.OriginalPurchaseDate = p.PurchaseDate                    // 对于恢复的transaction对象，该键对应了原始的交易日期
	p.ExpireDate = nil                                         // 普通消耗类型,没有过期时间
	p.OriginalApplicationVersion = payload.AppVersion          // 开发者指定的字符串，包含订单的补充信息。您可以在发起 getBuyIntent 请求时为此字段指定一个值。
	p.AppItemId = data.PurchaseToken                           // 用于对给定商品和用户对进行唯一标识的令牌。
	p.VersionExternalIdentifier = 0                            // 用来标识程序修订数。该键在sandbox环境下不存在
	// ------------------------------------------------------------------------------------------
	p.Quantity = 1                // 购买商品的数量
	p.Status = 0                  // 订单的购买状态。可能的值为 0（已购买）、1（已取消）或者 2（已退款）
	p.ProductId = data.ProductId  // 商品的商品 ID。每种商品都有一个商品 ID，您必须通过 Google Play Developer Console 在应用的商品列表中指定此 ID。
	p.BundleId = data.PackageName //  Android程序的bundle标识
	return p, nil
}

// New 初始化
func New(config config.ChannelProperty) (inf.Verifier, error) {
	hw := huawei{config: Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		IsSandbox:    config.Sandbox,
		ApiRoot:      utilk.DefaultIfEmpty(config.ApiRoot, "https://orders-drru.iap.cloud.huawei.ru"),
	}}
	if err := hw.parseSignKey(); err != nil {
		return nil, err
	}
	return &hw, nil
}
