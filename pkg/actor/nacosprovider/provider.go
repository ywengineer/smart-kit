package nacosprovider

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/asynkron/protoactor-go/cluster"
)

var ProviderShuttingDownError = fmt.Errorf("nacos cluster provider is shutting down")
var ClusterDisabledError = fmt.Errorf("nacos cluster is disabled")

type Provider struct {
	cluster            *cluster.Cluster
	deregistered       bool
	shutdown           bool
	id                 string
	serviceName        string
	address            string
	port               int
	knownKinds         []string
	ttl                time.Duration
	refreshTTL         time.Duration
	updateTTLWaitGroup sync.WaitGroup
	deregisterCritical time.Duration
	blockingWaitTime   time.Duration
	clusterError       error
	pid                *actor.PID
	client             naming_client.INamingClient
	//
	weight      int
	metadata    map[string]string
	clusterName string
	groupName   string
	namespace   string
}

func New(client naming_client.INamingClient, opts ...Option) *Provider {
	p := &Provider{
		ttl:                3 * time.Second,
		refreshTTL:         1 * time.Second,
		deregisterCritical: 60 * time.Second,
		blockingWaitTime:   20 * time.Second,
		client:             client,
		weight:             10,
		groupName:          "DEFAULT_GROUP",
		clusterName:        "DEFAULT",
		namespace:          "public",
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *Provider) init(c *cluster.Cluster) error {
	knownKinds := c.GetClusterKinds()
	clusterName := c.Config.Name
	memberId := c.ActorSystem.ID

	host, port, err := c.ActorSystem.GetHostPort()
	if err != nil {
		return err
	}

	p.cluster = c
	p.id = memberId
	p.clusterName = clusterName
	p.address = host
	p.port = port
	p.knownKinds = knownKinds
	return nil
}

func (p *Provider) StartMember(c *cluster.Cluster) error {
	err := p.init(c)
	if err != nil {
		return err
	}

	p.pid, err = c.ActorSystem.Root.SpawnNamed(actor.PropsFromProducer(func() actor.Actor {
		return newProviderActor(p)
	}), "nacos-provider")
	if err != nil {
		p.cluster.Logger().Error("Failed to start consul-provider actor", slog.Any("error", err))
		return err
	}

	return nil
}

func (p *Provider) StartClient(c *cluster.Cluster) error {
	if err := p.init(c); err != nil {
		return err
	}
	p.blockingStatusChange()
	p.monitorMemberStatusChanges()
	return nil
}

func (p *Provider) DeregisterMember() error {
	err := p.deregisterService()
	if err != nil {
		fmt.Println(err)
		return err
	}
	p.deregistered = true
	return nil
}

func (p *Provider) Shutdown(graceful bool) error {
	if p.shutdown {
		return nil
	}
	p.shutdown = true
	if p.pid != nil {
		if err := p.cluster.ActorSystem.Root.StopFuture(p.pid).Wait(); err != nil {
			p.cluster.Logger().Error("Failed to stop nacos-provider actor", slog.Any("error", err))
		}
		p.pid = nil
	}

	return nil
}

func blockingUpdateTTL(p *Provider) error {
	if !p.client.ServerHealthy() {
		p.clusterError = ClusterDisabledError
	} else {
		p.clusterError = nil
	}
	return p.clusterError
}

func (p *Provider) registerService() error {
	s := vo.RegisterInstanceParam{
		ClusterName: p.clusterName,
		GroupName:   p.groupName,
		ServiceName: p.serviceName,
		Weight:      float64(p.weight),
		Ip:          p.address,
		Port:        uint64(p.port),
		Healthy:     true,
		Enable:      true,
		Metadata: map[string]string{
			"id":   p.id,
			"tags": strings.Join(p.knownKinds, ","),
		},
	}
	_, err := p.client.RegisterInstance(s)
	return err
}

func (p *Provider) deregisterService() error {
	s := vo.DeregisterInstanceParam{
		Ip:          p.address,
		Port:        uint64(p.port),
		Cluster:     p.clusterName,
		ServiceName: p.serviceName,
		GroupName:   p.groupName,
	}
	_, err := p.client.DeregisterInstance(s)
	return err
}

// call this directly after registering the service
func (p *Provider) blockingStatusChange() {
	p.notifyStatuses()
}

func (p *Provider) notifyStatuses() {
	service, err := p.client.GetService(vo.GetServiceParam{
		GroupName:   p.groupName,
		ServiceName: p.serviceName,
	})
	if err != nil {
		p.cluster.Logger().Error("notifyStatues", slog.Any("error", err))
		return
	}
	var members []*cluster.Member
	for _, v := range service.Hosts {
		if v.Enable && v.Healthy {
			memberId := v.Metadata["id"]
			if memberId == "" {
				memberId = fmt.Sprintf("%v[%s]@%v:%v", p.clusterName, v.InstanceId, v.Ip, v.Port)
				p.cluster.Logger().Info("metadata['id'] was empty, fixeds", slog.String("id", memberId))
			}
			members = append(members, &cluster.Member{
				Id:    memberId,
				Host:  v.Ip,
				Port:  int32(v.Port),
				Kinds: strings.Split(v.Metadata["tags"], ","),
			})
		}
	}
	// the reason why we want this in a batch and not as individual messages is that
	// if we have an atomic batch, we can calculate what nodes have left the cluster
	// passing events one by one, we can't know if someone left or just haven't changed status for a long time

	// publish the current cluster topology onto the event stream
	p.cluster.MemberList.UpdateClusterTopology(members)
}

func (p *Provider) monitorMemberStatusChanges() {
	go func() {
		for !p.shutdown {
			p.notifyStatuses()
		}
	}()
}
