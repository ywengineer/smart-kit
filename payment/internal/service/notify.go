package service

import (
	"context"

	"gitee.com/ywengineer/smart-kit/payment/internal/config"
	"gitee.com/ywengineer/smart-kit/payment/pkg/model"
	"gitee.com/ywengineer/smart-kit/pkg/apps"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func Notify(ctx context.Context, sCtx apps.SmartContext, purchase model.Purchase) {
	 serverInfo,ok := config.GetMeta().FindServer(purchase.GameId, purchase.ServerId);
	// 如果通知地址不存在
	if !ok   || len(serverInfo.ApiUrl) ==0 {
		hlog.CtxErrorf(ctx, "[NotifyPay] notify purchase error: [参数错误] 服务器不存在或邮件通知地址未配置！")
		return
	};
	// 通知次数
	if (purchaseLogMapper.incrementNotifyTimes(purchaseLog.getTransaction_id()) == 0) {
		hlog.CtxErrorf(ctx, "[NotifyPay] order [%s] not found",purchase.TransactionId)
		return
	}
	final PayNotify notifyBody = PayNotify.newBuilder()
	.setPlayerID(purchaseLog.getPlayer_id())
	.setOrderID(purchaseLog.getTransaction_id())
	.setItemID(purchaseLog.getProduct_id())
	.setPlatformCode(purchaseLog.getChannel())
	.setPurchaseTime(purchaseLog.getPurchase_date().getTime())
	.setExpireTime(purchaseLog.expiredTime())
	.build();
	//
	ManageRet ret = postProtoBuf(serverInfo.apiUrl(ApiMethod.PayNotify), notifyBody.toByteArray(), Collections.emptyMap());
	//
	if (ret == null) {
		log.error("[NotifyPayFailed], {}", notifyBody);
	} else {
		//
		if (ret.getCode() != ManageHandlerCodes.SUC.getId()) return ApiResultUtils.WebApiResult.fail(MessageFormatter.format("code = {}, message = {}", ret.getCode(), ret.getMsg()).getMessage());
		//
		PayNotify.Ret notifyRet = PayNotify.Ret.parseFrom(ret.getData());
		//
		if (notifyRet.hasOrderID() && !StringUtils.isEmpty(notifyRet.getOrderID())) {
		//
			PurchaseLog order = purchaseLogMapper.queryByTransaction(notifyRet.getOrderID());
			//
			if (order != null) {
				// 结果
				order.setRankPoints(notifyRet.getRankPoints());
				order.setCredits(notifyRet.getCredits());
				order.setMoney(notifyRet.getMoney());
				order.setCoin(notifyRet.getCoin());
				// 通知成功
				order.setNotified(true);
				order.setFinishTime(new Date());
				// 更新订单状态
				purchaseLogMapper.updateNotify(order);
				// 通知Kafka
				producerService.send(KafkaTopics.Topic_Purchase, JSON.toJSONBytes(order));
			} else {
				log.error("[NotifyResult], {}", notifyRet);
			}
		} else {
			log.error("[NotifyResult], missing orderId, {}", notifyRet);
		}
	}
}
