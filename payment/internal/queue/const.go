package queue

import (
	"errors"

	"gitee.com/ywengineer/smart-kit/payment/pkg/model"
	"github.com/bytedance/sonic"
	"github.com/hibiken/asynq"
)

// A list of task types.
const (
	TypePurchaseNotify = "purchase:notify"
)

var ErrNotInit = errors.New("queue is not initialized yet, please init queue use queue.InitQueue")

type PurchaseNotifyPayload struct {
	GameID        string `json:"game_id"`
	ServerId      string `json:"server_id"`
	TransactionId string `json:"transaction_id"`
	PlayerId      string `json:"player_id"`
	ProductId     string `json:"product_id"`
	Channel       string `json:"channel"`
}

//----------------------------------------------
// Write a function NewXXXTask to create a task.
// A task consists of a type and a payload.
//----------------------------------------------

func NewPurchaseNotifyTask(purchase model.Purchase) error {
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
	})
	if err != nil {
		return err
	}
	_, e := cli.Enqueue(
		asynq.NewTask(TypePurchaseNotify, payload),
		asynq.MaxRetry(0),
		asynq.Queue(TypePurchaseNotify),
	)
	return e
}
