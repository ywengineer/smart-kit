package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gitee.com/ywengineer/smart-kit/payment/internal/queue"
	"gitee.com/ywengineer/smart-kit/pkg/apps"

	"gitee.com/ywengineer/smart-kit/payment/internal/config"
	"gitee.com/ywengineer/smart-kit/payment/pkg/model"
	msg "gitee.com/ywengineer/smart-kit/payment/pkg/proto"
	"gitee.com/ywengineer/smart-kit/pkg/rpcs"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

func Notify(ctx context.Context, data queue.PurchaseNotifyPayload) error {
	serverInfo, ok := config.FindServer(data.GameID, data.ServerId)
	// 如果通知地址不存在
	if !ok || len(serverInfo.ApiUrl) == 0 {
		return errors.New("[NotifyPay] [参数错误] 服务器不存在或邮件通知地址未配置！")
	}
	sCtx := apps.GetContext(ctx)
	// update notify times
	cnt, err := gorm.G[model.Purchase](sCtx.Rdb().WithContext(ctx)).
		Where("transaction_id = ? and notified = 0", data.TransactionId).
		Update(ctx, "notified_times", gorm.Expr("notified_times + 1"))
	if err != nil {
		return err
	} else if cnt <= 0 {
		return nil // 已通知, 无需重复通知
	}
	// build notify data
	notifyBody := msg.PayNotify{
		PlayerID:     &data.PlayerId,
		OrderID:      &data.TransactionId,
		ItemID:       &data.ProductId,
		PlatformCode: &data.Channel,
		PurchaseTime: &data.PurchaseTime,
		ExpireTime:   &data.ExpireTime,
	}
	//
	statusCode, resp, err := sCtx.Rpc().Post(ctx, rpcs.ContentTypeOctStream, serverInfo.GetApiMethodUrl("pay"), nil, rpcs.ProtoBody{V: &notifyBody})
	if err != nil {
		return err
	} else if statusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("[NotifyPay] api error: status = %d, body = %s", statusCode, string(resp)))
	}
	var ret msg.ManageRet
	if err = proto.Unmarshal(resp, &ret); err != nil {
		return err
	} else if ret.GetCode() != 0 {
		return errors.New(fmt.Sprintf("[NotifyPay] api failed: code = %d, message = %s", ret.GetCode(), ret.GetMsg()))
	}
	//
	var notifyRet msg.PayNotify_Ret
	if err = proto.Unmarshal(ret.GetData(), &notifyRet); err != nil {
		return err
	} else if len(notifyRet.GetOrderID()) == 0 {
		return errors.New(fmt.Sprintf("[NotifyPay] result missing order [%s]", notifyRet.String()))
	}
	now := time.Now().Local()
	// update order data
	if updated, err := gorm.G[model.Purchase](sCtx.Rdb().WithContext(ctx)).
		Where("transaction_id = ? AND notified = 0", notifyRet.GetOrderID()).
		Select("RankPoints", "Credits", "Money", "Coin", "Notified", "FinishDate").
		Updates(ctx, model.Purchase{
			RankPoints: int(notifyRet.GetRankPoints()),
			Credits:    notifyRet.GetCredits(),
			Money:      notifyRet.GetMoney(),
			Coin:       notifyRet.GetCoin(),
			Notified:   true,
			FinishDate: &now,
		}); err != nil {
		return err
	} else if updated <= 0 {
		return nil // 已通知, 无需重复通知
	} else {
		return nil
	}
}
