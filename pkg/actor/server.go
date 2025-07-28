package actor

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
)

type EchoActor struct{}

func (e *EchoActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case string:
		log.Printf("收到消息: %s", msg)

		// 模拟20%概率处理失败
		if time.Now().UnixNano()%5 == 0 {
			log.Println("模拟处理失败")
			return // 不回复，触发客户端超时
		}

		ctx.Respond("reply: " + msg)
	}
}

func main() {
	system := actor.NewActorSystem()

	// 配置远程服务（添加错误日志）
	remoteConfig := remote.Configure("localhost", 8080)
	//remoteConfig.WithErrorHandler(func(err error) {
	//	log.Printf("远程服务错误: %v", err)
	//})

	rs := remote.NewRemote(system, remoteConfig)

	// 注册服务
	rs.Register("EchoService", actor.PropsFromProducer(func() actor.Actor {
		return &EchoActor{}
	}))

	// 启动远程服务
	rs.Start()

	log.Println("服务端已启动，监听端口 8080...")

	// 注册健康检查
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			log.Println("服务端运行中...")
		}
	}()

	// 优雅关闭
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	log.Println("服务端正在关闭...")
	rs.Shutdown(true)
	log.Println("服务端已关闭")
}
