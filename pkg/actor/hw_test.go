package actor

import (
	"gitee.com/ywengineer/smart-kit/pkg/nets"
	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	"log/slog"
	"os"
	"testing"
)

type (
	hello      struct{ Who string }
	helloActor struct{}
)

func (state *helloActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *hello:
		context.Logger().Info("Hello ", slog.String("who", msg.Who))
	}
}

func TestHelloActor(t *testing.T) {
	system := actor.NewActorSystem(actor.WithLoggerFactory(func(system *actor.ActorSystem) *slog.Logger {
		return slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}))
	ip := nets.GetDefaultIpv4()
	r := remote.NewRemote(system, remote.Configure(ip, 0, remote.WithAdvertisedHost(ip)))
	r.Register("network", actor.PropsFromProducer(func() actor.Actor { return &helloActor{} }))
	r.Start()
	//pid := r.Spawn() system.Root.Spawn(props)
	//system.Root.Send(pid, &hello{Who: "Roger"})
	//r.SendMessage()
	<-utilk.WatchQuitSignal()
	r.Shutdown(true)
}
