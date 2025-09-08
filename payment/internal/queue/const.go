package queue

import (
	"errors"
	"fmt"
	"time"

	"gitee.com/ywengineer/smart-kit/payment/pkg/model"
	"github.com/bytedance/sonic"
	"github.com/hibiken/asynq"
)

type TaskType string

// A list of task types.
const (
	PurchaseNotify TaskType = "queue:purchase:notify"
	Test           TaskType = "queue:test"
)

type QName string

const (
	High    QName = "high"
	Low     QName = "low"
	Default QName = "default"
)

func (name QName) String() string         { return fmt.Sprintf("Queue(%q)", string(name)) }
func (name QName) Type() asynq.OptionType { return asynq.QueueOpt }
func (name QName) Value() interface{}     { return string(name) }

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

type TestAsynqQueue struct {
	Time time.Time `bson:"time"`
}

//----------------------------------------------
// Write a function NewXXXTask to create a task.
// A task consists of a type and a payload.
//----------------------------------------------

func PublishPurchaseNotify(purchase model.Purchase, options ...asynq.Option) error {
	return publishTask(PurchaseNotify, PurchaseNotifyPayload{
		GameID:        purchase.GameId,
		ServerId:      purchase.ServerId,
		TransactionId: purchase.TransactionId,
		PlayerId:      purchase.PlayerId,
		ProductId:     purchase.ProductId,
		Channel:       purchase.Channel,
		PurchaseTime:  purchase.PurchaseDate.Unix(),
		ExpireTime:    purchase.GetExpiredTime(),
	}, options...)
}

func PublishTest() error {
	return publishTask(Test, TestAsynqQueue{Time: time.Now()})
}

func publishTask(task TaskType, payload interface{}, options ...asynq.Option) error {
	if cli == nil {
		return ErrNotInit
	}
	ps, err := sonic.Marshal(payload)
	if err != nil {
		return err
	}
	_, e := cli.Enqueue(asynq.NewTask(string(task), ps), options...)
	return e
}
