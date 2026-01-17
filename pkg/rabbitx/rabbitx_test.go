package rabbitx

import (
	"context"
	"strconv"
	"testing"
	"time"

	"gitee.com/ywengineer/smart-kit/pkg/loaders"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

func TestRabbitx(t *testing.T) {
	var cfg RabbitMQConfig
	err := loaders.NewLocalLoader("./amqp.yaml").Load(&cfg)
	assert.Nil(t, err)
	// 创建客户端实例
	client, err := NewRabbitX(cfg)
	if err != nil {
		t.Fatalf("创建 RabbitX 客户端失败: %v", err)
	}
	defer client.Shutdown()
	// 启动消费者：自定义消息处理逻辑
	_ = client.ConsumeMsg(map[string]Processor{
		"log-queue": ProcessorFunc(func(ctx context.Context, msg *amqp.Delivery) error {
			t.Logf("✅ 收到消息:  routingKey = %s, body = %+v", msg.RoutingKey, string(msg.Body))
			return nil
		}),
	})
	// 6. 启动生产者：循环发送测试消息
	go func() {
		i := 0
		for {
			i++
			msg := map[string]interface{}{
				"msg":  "test message " + strconv.Itoa(i),
				"time": time.Now().Format(time.RFC3339Nano),
			}
			_ = client.PublishMsg(context.Background(), "rabbitx-exchange", "rabbitx.log.info", msg)
			time.Sleep(2 * time.Second)
		}
	}()

	// 阻塞主线程
	select {}
}

func TestChannel(t *testing.T) {

	single := make(chan struct{})

	go func() {
		for range single {
			t.Logf("✅ recv single")
		}
		t.Logf("✅ 信道已关闭")
	}()

	go func() {
		single <- struct{}{}
		time.Sleep(5 * time.Second)
		close(single)
	}()

	select {}
}
