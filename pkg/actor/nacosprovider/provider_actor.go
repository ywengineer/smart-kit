package nacosprovider

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/scheduler"
)

type providerActor struct {
	*Provider
	actor.Behavior
	refreshCanceller scheduler.CancelFunc
}

type (
	RegisterService   struct{}
	UpdateTTL         struct{}
	MemberListUpdated struct {
		members []*cluster.Member
	}
)

func (pa *providerActor) Receive(ctx actor.Context) {
	pa.Behavior.Receive(ctx)
}

func newProviderActor(provider *Provider) actor.Actor {
	pa := &providerActor{
		Behavior: actor.NewBehavior(),
		Provider: provider,
	}
	pa.Become(pa.init)
	return pa
}

func (pa *providerActor) init(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		ctx.Send(ctx.Self(), &RegisterService{})
	case *RegisterService:
		if err := pa.registerService(); err != nil {
			ctx.Logger().Error("Failed to register service to nacos, will retry", slog.Any("error", err))
			ctx.Send(ctx.Self(), &RegisterService{})
		} else {
			ctx.Logger().Info("Registered service to nacos")
			pa.Become(pa.running)
			refreshScheduler := scheduler.NewTimerScheduler(ctx)
			pa.refreshCanceller = refreshScheduler.SendRepeatedly(0, pa.refreshTTL, ctx.Self(), &UpdateTTL{})
			if err := pa.doSubscribe(ctx); err != nil {
				ctx.Logger().Error("Failed to subscribe nacos service", slog.String("service", pa.serviceName), slog.Any("error", err))
			}
		}
	}
}

func (pa *providerActor) running(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *UpdateTTL:
		if err := blockingUpdateTTL(pa.Provider); err != nil {
			ctx.Logger().Warn("Failed to update TTL", slog.Any("error", err))
		}
	case *MemberListUpdated:
		pa.cluster.MemberList.UpdateClusterTopology(msg.members)
	case *actor.Stopping:
		pa.refreshCanceller()
		if err := pa.deregisterService(); err != nil {
			ctx.Logger().Error("Failed to deregister service from nacos", slog.Any("error", err))
		} else {
			ctx.Logger().Info("De-registered service from nacos")
		}
	}
}

func (pa *providerActor) doSubscribe(ctx actor.Context) error {
	err := pa.client.Subscribe(&vo.SubscribeParam{
		Clusters:    []string{pa.clusterName},
		GroupName:   pa.groupName,
		ServiceName: pa.serviceName,
		SubscribeCallback: func(services []model.Instance, err error) {
			pa.processNacosUpdate(services, err, ctx)
		},
	})
	if err != nil {
		ctx.Logger().Error("Failed to subscribe nacos service", slog.Any("error", err))
		return err
	}
	return nil
}

func (pa *providerActor) processNacosUpdate(services []model.Instance, err error, ctx actor.Context) {
	if err != nil {
		ctx.Logger().Error("Didn't get expected data from nacos subscription")
		return
	}
	ctx.Logger().Info("Subs trigger, service topology changed")
	var members []*cluster.Member
	for _, v := range services {
		if v.Enable {
			memberId := v.Metadata["id"]
			if memberId == "" {
				memberId = fmt.Sprintf("%v[%s]@%v:%v", pa.clusterName, v.InstanceId, v.Ip, v.Port)
				ctx.Logger().Info("metadata['id'] was empty, fixed", slog.String("id", memberId))
			}
			members = append(members, &cluster.Member{
				Id:    memberId,
				Host:  v.Ip,
				Port:  int32(v.Port),
				Kinds: strings.Split(v.Metadata["tags"], ","),
			})
		}
	}

	// delay the fist update until there is at least one member
	if len(members) > 0 {
		ctx.Send(ctx.Self(), &MemberListUpdated{
			members: members,
		})
	}
}
