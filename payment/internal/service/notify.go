package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gitee.com/ywengineer/smart-kit/payment/internal/config"
	"gitee.com/ywengineer/smart-kit/payment/pkg/model"
	msg "gitee.com/ywengineer/smart-kit/payment/pkg/proto"
	"gitee.com/ywengineer/smart-kit/pkg/apps"
	"gitee.com/ywengineer/smart-kit/pkg/rpcs"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

func Notify(ctx context.Context, sCtx apps.SmartContext, purchase model.Purchase) error {
	serverInfo, ok := config.GetMeta().FindServer(purchase.GameId, purchase.ServerId)
	// 如果通知地址不存在
	if !ok || len(serverInfo.ApiUrl) == 0 {
		return errors.New("[NotifyPay] [参数错误] 服务器不存在或邮件通知地址未配置！")
	}
	// update notify times
	cnt, err := gorm.G[model.Purchase](sCtx.Rdb().WithContext(ctx)).
		Where("transaction_id = ?", purchase.TransactionId).
		Update(ctx, "notified_times", gorm.Expr("notified_times + 1"))
	if err != nil {
		return err
	} else if cnt <= 0 {
		return errors.New(fmt.Sprintf("[NotifyPay] order [%s] not found", purchase.TransactionId))
	}
	// build notify data
	pt, et := purchase.PurchaseDate.UnixMilli(), purchase.GetExpiredTime()
	notifyBody := msg.PayNotify{
		PlayerID:     &purchase.PlayerId,
		OrderID:      &purchase.TransactionId,
		ItemID:       &purchase.ProductId,
		PlatformCode: &purchase.Channel,
		PurchaseTime: &pt,
		ExpireTime:   &et,
	}
	//
	statusCode, resp, err := sCtx.Rpc().Post(ctx, rpcs.ContentTypeOctStream, serverInfo.GetApiMethodUrl("pay"), rpcs.ProtoBody{V: &notifyBody})
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
	now := time.Now()
	// update order data
	if updated, err := gorm.G[model.Purchase](sCtx.Rdb().WithContext(ctx)).
		Where("transaction_id = ?", notifyRet.GetOrderID()).
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
		return errors.New(fmt.Sprintf("[NotifyPay] order [%s] not found", notifyRet.String()))
	} else {
		return nil
	}
}
