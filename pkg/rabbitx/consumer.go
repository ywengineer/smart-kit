package rabbitx

import (
	"context"
	"fmt"
	"sync/atomic"

	"gitee.com/ywengineer/smart-kit/pkg/logk"
	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	amqp "github.com/rabbitmq/amqp091-go"
	"gopkg.in/errgo.v2/fmt/errors"
)

type Processor interface {
	Process(ctx context.Context, msg *amqp.Delivery) error
}

type ProcessorFunc func(ctx context.Context, msg *amqp.Delivery) error

func (p ProcessorFunc) Process(ctx context.Context, msg *amqp.Delivery) error {
	return p(ctx, msg)
}

type xConsumer struct {
	name   string
	conf   *ConsumerConfig
	ch     *amqp.Channel
	runner Processor
	seq    uint64
}

func (c *xConsumer) run() *xConsumer {
	c.conf.Size = utilk.Min(utilk.Max(1, c.conf.Size), 10)
	//
	for range c.conf.Size {
		cName := fmt.Sprintf("%s-%d", c.name, c.nextSeq())
		if q, err := c.ch.Consume( // start consume if processor exist
			c.conf.Queue,     // queue
			cName,            // consumer
			c.conf.AutoAck,   // auto ack
			c.conf.Exclusive, // exclusive
			false,            // no local
			c.conf.NoWait,    // no wait
			c.conf.Arguments, // args
		); err != nil { // consume failed
			logk.Errorf("[%s] consume queue [%s] failed %+v", cName, c.conf.Queue, err)
		} else { // consume success
			go c.consume(q, cName, c.runner)
		}
	}
	return c
}

func (c *xConsumer) consume(q <-chan amqp.Delivery, consumer string, processor Processor) {
	logk.Infof("✅ Consumer [%s] has been started and is waiting to receive messages", consumer)
	for msg := range q {
		var err error
		if err = processor.Process(context.Background(), &msg); err == nil {
			err = msg.Ack(false)
		} else {
			err = errors.Because(err, msg.Reject(true), "process message failed")
		}
		if err != nil {
			logk.Errorf("[%s] consume message err: %s", consumer, err.Error())
		}
	}
	logk.Infof("✅ Consumer [%s] has been stopped", consumer)
}

func (c *xConsumer) nextSeq() uint64 {
	return atomic.AddUint64(&c.seq, 1)
}

func (c *xConsumer) Close() error {
	return c.ch.Close()
}
