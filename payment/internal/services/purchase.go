package services

import (
	"context"
	"errors"

	"gitee.com/ywengineer/smart-kit/payment/internal/config"
	"gitee.com/ywengineer/smart-kit/payment/pkg/api"
	"gitee.com/ywengineer/smart-kit/payment/pkg/model"
	"gitee.com/ywengineer/smart-kit/pkg/apps"
)

var ErrDuplicateOrder = errors.New(api.DuplicateOrder)

func OnPurchase(ctx context.Context, sCtx apps.SmartContext, gameId, serverId, passport, playerId, playerName string, purchaseLog *model.Purchase, channel config.Channel, product config.Product) error {
	// 订单数据补充
	purchaseLog.GameId = gameId
	purchaseLog.ServerId = serverId
	purchaseLog.Passport = passport
	purchaseLog.PlayerId = playerId
	purchaseLog.PlayerName = playerName
	purchaseLog.Channel = channel.Code
	purchaseLog.Notified = false
	// 查看订单是否已处理
	var old model.Purchase
	ret := sCtx.Rdb().
		WithContext(ctx).
		Select("id", "expire_date").
		First(&old, "transaction_id = ?", purchaseLog.TransactionId)
	if ret.Error != nil {
		return ret.Error
	}
	// 如果是订阅且已处理, 不再进行后续处理
	if old.ID > 0 && old.ExpireDate != nil {
		return nil
	}
	// 如果已处理
	if old.ID > 0 {
		return ErrDuplicateOrder
	}
	// 如果是试用或者测试订单, 价格为0. 否则为填写金额
	if purchaseLog.FreeTrail || purchaseLog.TestOrder {
		purchaseLog.Price = 0
	} else {
		purchaseLog.Price = 0
	}
	purchaseLog.Price = int(product.Money * 100)
	// 记录
	return sCtx.Rdb().WithContext(ctx).Create(purchaseLog).Error
}
