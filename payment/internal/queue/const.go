package queue

import (
	"errors"

	"gitee.com/ywengineer/smart-kit/payment/pkg/model"
	"github.com/bytedance/sonic"
	"github.com/hibiken/asynq"
)

type TaskType string

// A list of task types.
const (
	PurchaseNotify TaskType = "purchase:notify"
)

var ErrNotInit = errors.New("queue is not initialized yet, please init queue use queue.InitQueue")

type PurchaseNotifyPayload struct {
	GameID        string `json:"game_id"`
	ServerId      string `json:"server_id"`
	TransactionId string `json:"transaction_id"`
	PlayerId      string `json:"player_id"`
	ProductId     string `json:"product_id"`
	Channel       string `json:"channel"`
	PurchaseTime  int64  `json:"purchase_time"`
	ExpireTime    int64  `json:"expire_time"`
}

//----------------------------------------------
// Write a function NewXXXTask to create a task.
// A task consists of a type and a payload.
//----------------------------------------------

func PublishPurchaseNotify(purchase model.Purchase, options ...asynq.Option) error {
	if cli == nil {
		return ErrNotInit
	}
	payload, err := sonic.Marshal(PurchaseNotifyPayload{
		GameID:        purchase.GameId,
		ServerId:      purchase.ServerId,
		TransactionId: purchase.TransactionId,
		PlayerId:      purchase.PlayerId,
		ProductId:     purchase.ProductId,
		Channel:       purchase.Channel,
		PurchaseTime:  purchase.PurchaseDate.Unix(),
		ExpireTime:    purchase.GetExpiredTime(),
	})
	if err != nil {
		return err
	}
	var ops []asynq.Option
	ops = append(ops, options...)
	_, e := cli.Enqueue(asynq.NewTask(string(PurchaseNotify), payload), ops...)
	return e
}
