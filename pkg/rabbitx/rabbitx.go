package rabbitx

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gitee.com/ywengineer/smart-kit/pkg/logk"
	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/samber/lo"
)

var ErrNotInit = errors.New("rabbitmq client is not initialized yet, please init rabbitmq client use NewRabbitX")
var ErrPublisherClosed = errors.New("rabbitmq publisher is closed")

// RabbitX 定义全局的RabbitMQ客户端结构体，封装所有连接信息和重连逻辑
type RabbitX struct {
	conn        *amqp.Connection // MQ连接
	ch          *amqp.Channel    // MQ信道
	cfg         *RabbitMQConfig  // 配置信息
	isConnected *atomic.Bool     // 连接状态标记

	reconnectChan chan struct{} // 重连信号通道
	reconnectOnce *sync.Once

	consumers          []*xConsumer
	consumerProcessors map[string]Processor
}

// NewRabbitX 创建RabbitMQ客户端实例
func NewRabbitX(config RabbitMQConfig) (*RabbitX, error) {
	rmc := &RabbitX{
		reconnectOnce: &sync.Once{},
		reconnectChan: make(chan struct{}, 1),
		cfg:           &config,
		isConnected:   &atomic.Bool{},
	}
	return rmc, rmc.connect()
}

func (r *RabbitX) newChannel(name string) (*amqp.Channel, error) {
	//  创建信道
	ch, err := r.conn.Channel()
	if err != nil {
		r.Close()
		return nil, fmt.Errorf("rabbit create channel failed: %w", err)
	}
	// 3. 设置QoS
	if err = ch.Qos(r.cfg.QoS.PrefetchCount, r.cfg.QoS.PrefetchSize, r.cfg.QoS.Global); err != nil {
		_ = ch.Close()
		r.Close()
		return nil, fmt.Errorf("rabbit set channel qos failed: %w", err)
	}
	// 8. 开启监听：信道断开则触发重连信号
	r.listenChanClose(name, ch)
	//
	return ch, nil
}

// 核心：建立MQ连接 + 创建信道 + 声明交换机/队列/绑定关系
// 所有异常都会返回error，外部统一处理重连
func (r *RabbitX) connect() error {
	var tlsConfig *tls.Config
	nameMainCh := "publisher"
	if strings.HasPrefix(r.cfg.Addr, "amqps://") {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		}
	}
	//
	conn, err := amqp.DialTLS(r.cfg.Addr, tlsConfig)
	// 1. 建立TCP连接
	if err != nil {
		return fmt.Errorf("rabbit dial failed: %w", err)
	}
	//
	r.conn = conn
	// 2. 创建信道
	ch, err := r.newChannel(nameMainCh)
	if err != nil {
		return fmt.Errorf("rabbit create channel failed: %w", err)
	}
	// 3. 声明交换机 (幂等操作，重复声明无影响)
	err = r.ensureExchanges(ch, r.cfg.Exchanges)
	if err != nil {
		r.Close()
		return fmt.Errorf("rabbit declare exchange failed: %w", err)
	}
	// 4. 声明队列 (幂等操作)
	err = r.ensureQueues(ch, r.cfg.Queues)
	if err != nil {
		r.Close()
		return fmt.Errorf("rabbit declare queue failed: %w", err)
	}
	// 5. 绑定队列到交换机
	err = r.ensureBindings(ch, r.cfg.Bindings)
	if err != nil {
		r.Close()
		return fmt.Errorf("rabbit bind queue failed: %w", err)
	}
	// 6. 赋值连接和信道，标记连接状态
	r.conn = conn
	r.ch = ch
	r.isConnected.Store(true)
	//
	logk.Infof("✅ RabbitMQ connection has been established, resources have been initialized")
	//
	r.reconnectOnce.Do(func() {
		go r.reconnect()
	})
	// 7. 开启监听：连接断开则触发重连信号
	go r.listenConnClose()
	// 8. start consumers
	_ = r.ConsumeMsg(r.consumerProcessors)
	//
	return nil
}

// 监听连接关闭事件，连接断开时发送重连信号
func (r *RabbitX) listenConnClose() {
	err := <-r.conn.NotifyClose(make(chan *amqp.Error))
	if r.isConnected.CompareAndSwap(true, false) {
		logk.Errorf("❌ RabbitMQ connection has been closed, preparing to reconnect...: %v", err)
		r.reconnectChan <- struct{}{}
	}
}

// 监听信道关闭事件，信道断开时发送重连信号
func (r *RabbitX) listenChanClose(name string, ch *amqp.Channel) {
	go func() {
		err := <-ch.NotifyClose(make(chan *amqp.Error))
		if r.isConnected.CompareAndSwap(true, false) {
			logk.Errorf("❌ RabbitMQ Channel [%s] has been closed, preparing to reconnect...: %v", name, err)
			r.reconnectChan <- struct{}{}
		}
	}()
}

// reconnect 核心：自动重连方法，带指数退避策略，无限重试直到重连成功
func (r *RabbitX) reconnect() {
	for range r.reconnectChan {
		if r.isConnected.Load() {
			continue
		}
		// 尝试重连
		_backoff := 0
		for {
			err := r.connect()
			if err == nil {
				break // 重连成功，退出循环
			}
			// 随机指数退避休眠：1s → 2s →4s →8s →16s，最大16s，防止CPU飙升
			backoff := time.Duration(1<<_backoff) * time.Second
			logk.Errorf("⏳ Retry reconnecting RabbitMQ after %v ...", backoff)
			_backoff = utilk.Min(4, _backoff+1)
			time.Sleep(backoff)
		}
	}
}

// PublishMsg 生产者：发送消息，断线自动重连后可继续发送，消息持久化
func (r *RabbitX) PublishMsg(ctx context.Context, exchange string, routeKey string, payload interface{}, options ...Option) error {
	if !r.isConnected.Load() {
		return ErrNotInit
	}
	ch := r.ch
	if ch == nil || ch.IsClosed() {
		return ErrPublisherClosed
	}

	var msg *amqp.Publishing
	if p, ok := payload.(amqp.Publishing); ok {
		msg = &p
	} else {
		dataBytes, _err := sonic.Marshal(payload)
		if _err != nil {
			return _err
		}
		msg = &amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  consts.MIMEApplicationJSONUTF8,
			Body:         dataBytes,
		}
	}
	for _, option := range options {
		option(msg)
	}
	//
	if msg.MessageId == "" {
		msg.MessageId = lo.RandomString(6, lo.AlphanumericCharset)
	}
	//
	err := ch.PublishWithContext(
		ctx,
		exchange, // exchange
		routeKey, // routing key
		false,    // mandatory
		false,    // immediate
		*msg,
	)
	logk.Debugf("[RabbitX] publish msg. Queue[%s], Route[%s], Payload[%+v], Err[%+v]", exchange, routeKey, payload, err)
	//
	return err
}

// ConsumeMsg 消费者：消费消息，断线重连后自动恢复消费，自动ACK确认
func (r *RabbitX) ConsumeMsg(processors map[string]Processor) error {
	if !r.isConnected.Load() {
		return ErrNotInit
	} else if len(processors) == 0 {
		return nil
	}
	r.consumers = make([]*xConsumer, 0, len(processors))
	cmd := os.Args[0]
	for _, consumer := range r.cfg.Consumers {
		//
		var ch *amqp.Channel
		var err error
		consumerName := uniqueConsumerTag(cmd + "-" + consumer.Queue)
		if ch, err = r.newChannel(consumerName); err != nil {
			logk.Errorf("❌ Consumer [%s] failed to create channel: %v", consumerName, err)
			continue
		}
		r.consumers = append(r.consumers, (&xConsumer{
			name:   consumerName,
			conf:   &consumer,
			ch:     ch,
			runner: processors[consumer.Queue],
		}).run())
	}
	r.consumerProcessors = processors
	return nil
}

// Close 关闭连接和信道，优雅退出
func (r *RabbitX) Close() {
	for _, consumer := range r.consumers {
		_ = consumer.Close()
	}
	if r.ch != nil {
		_ = r.ch.Close()
	}
	if r.conn != nil {
		_ = r.conn.Close()
	}
	r.conn, r.ch = nil, nil
	r.consumers = make([]*xConsumer, 0)
	r.isConnected.Store(false)
	logk.Infof("✅ RabbitMQ connection has been gracefully closed")
}

func (r *RabbitX) Shutdown() error {
	r.Close()
	close(r.reconnectChan)
	logk.Infof("✅ RabbitMQ has been shutdown and stop reconnect")
	return nil
}

func (r *RabbitX) ensureBindings(ch *amqp.Channel, bindings []BindingConfig) error {
	for _, binding := range bindings {
		if err := ch.QueueBind(
			binding.Queue,      // queue name
			binding.RoutingKey, // routing key
			binding.Exchange,   // exchange
			binding.NoWait,     // no-wait
			binding.Arguments,  // arguments
		); err != nil {
			return err
		}
	}
	return nil
}

func (r *RabbitX) ensureQueues(ch *amqp.Channel, queues []QueueConfig) error {
	for _, queue := range queues {
		if qi, err := ch.QueueDeclare(
			queue.Name,       // name
			queue.Durable,    // durable
			queue.AutoDelete, // auto-deleted
			queue.Exclusive,  // exclusive
			queue.NoWait,     // no-wait
			queue.Arguments,  // arguments
		); err != nil {
			return err
		} else {
			logk.Infof("queue %s declared, count of messages not awaiting acknowledgment: %d, consumers: %d", queue.Name, qi.Messages, qi.Consumers)
		}
	}
	return nil
}

func (r *RabbitX) ensureExchanges(ch *amqp.Channel, exchanges []ExchangeConfig) error {
	for _, e := range exchanges {
		if err := ch.ExchangeDeclare(
			e.Name,       // name
			e.Kind,       // type
			e.Durable,    // durable
			e.AutoDelete, // auto-deleted
			e.Internal,   // internal
			e.NoWait,     // no-wait
			e.Arguments,  // arguments
		); err != nil {
			return err
		}
	}
	return nil
}
