package nacosprovider

import (
	"fmt"
	"gitee.com/ywengineer/smart-kit/pkg/logk"
	"github.com/bytedance/sonic"
	"log/slog"
	"net"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"

	"gitee.com/ywengineer/smart-kit/pkg/nacos"
	"gitee.com/ywengineer/smart-kit/pkg/nets"
	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/stretchr/testify/assert"
)

const serviceName = "my_service"
const groupName = "DEFAULT_GROUP"

func newClusterForTest(name string, addr string, cp cluster.ClusterProvider) *cluster.Cluster {
	host, _port, err := net.SplitHostPort(addr)
	if err != nil {
		panic(err)
	}
	host = utilk.DefaultIfEmpty(host, nets.GetDefaultIpv4())
	port, _ := strconv.Atoi(_port)
	remoteConfig := remote.Configure(host, port, remote.WithAdvertisedHost(host))
	lookup := disthash.New()
	config := cluster.Configure(name, cp, lookup, remoteConfig)
	// return cluster.NewForTest(system, config)

	system := actor.NewActorSystem(actor.WithLoggerFactory(func(system *actor.ActorSystem) *slog.Logger {
		return logk.GetDefaultSLogger().With("lib", "Proto.Actor").With("system", system.ID)
	}))
	c := cluster.New(system, config)

	// use for test without start remote
	c.ActorSystem.ProcessRegistry.Address = addr
	c.MemberList = cluster.NewMemberList(c)
	c.Remote = remote.NewRemote(c.ActorSystem, c.Config.RemoteConfig)
	return c
}

func newNacosProvider() cluster.ClusterProvider {
	conf := nacos.Nacos{
		Ip:          "127.0.0.1",
		Port:        8848,
		ContextPath: "/nacos",
		TimeoutMs:   20000,
		Namespace:   "public",
		User:        "nacos",
		Password:    "nacos",
		Group:       groupName,
	}
	nc, err := nacos.NewNamingClientWithConfig(conf, "debug")
	if err != nil {
		panic(err)
	}
	return New(nc, "public", groupName, WithServiceName(serviceName), WithRefreshTTL(5*time.Second), WithEphemeral())
}

func TestStartMember(t *testing.T) {
	if testing.Short() {
		return
	}
	a := assert.New(t)
	p := newNacosProvider()
	defer p.Shutdown(true)

	c := newClusterForTest(serviceName, "127.0.0.1:8000", p)
	eventstream := c.ActorSystem.EventStream
	eventstream.Subscribe(func(m interface{}) {
		t.Logf("[%s] %+v", reflect.TypeOf(m).String(), m)
	})

	err := p.StartMember(c)
	a.NoError(err)

	<-utilk.WatchQuitSignal()
}

func TestRegisterMultipleMembers(t *testing.T) {
	if testing.Short() {
		return
	}
	a := assert.New(t)

	members := []struct {
		cluster string
		host    string
		port    int
	}{
		{serviceName, "127.0.0.1", 8001},
		{serviceName, "127.0.0.1", 8002},
		{serviceName, "127.0.0.1", 8003},
	}

	p := newNacosProvider().(*Provider)
	defer p.Shutdown(true)
	//
	wg := sync.WaitGroup{}
	wg.Add(len(members))
	//
	for _, member := range members {
		mk := fmt.Sprintf("%s@%s:%d", member.cluster, member.host, member.port)
		go func(mk string) {
			defer wg.Done()
			addr := fmt.Sprintf("%s:%d", member.host, member.port)
			_p := newNacosProvider()
			c := newClusterForTest(member.cluster, addr, _p)
			eventstream := c.ActorSystem.EventStream
			eventstream.Subscribe(func(m interface{}) {
				if ct, ok := m.(*cluster.ClusterTopology); ok {
					topo, _ := sonic.Marshal(ct)
					t.Logf("[%s] %s", mk, topo)
				} else {
					t.Logf("[%s] [%s] %+v", mk, reflect.TypeOf(m).String(), m)
				}
			})
			err := _p.StartMember(c)
			a.NoError(err)
			t.Cleanup(func(__p cluster.ClusterProvider) func() {
				return func() {
					t.Logf("shutdown: %+v", __p.Shutdown(true))
				}
			}(_p))
			<-utilk.WatchQuitSignal()
			t.Logf("[%s], finished", mk)
		}(mk)
	}
	//
	time.Sleep(5 * time.Second)
	entries, err := p.client.GetService(vo.GetServiceParam{
		ServiceName: serviceName,
		GroupName:   groupName,
	})
	a.NoError(err)
	found := false
	for _, entry := range entries.Hosts {
		found = false
		for _, member := range members {
			if entry.Port == uint64(member.port) {
				found = true
			}
		}
		t.Logf("Member port [%v] - ExtensionID:%v Address: %v:%v, Metadata: %+v", found, entry.InstanceId, entry.Ip, entry.Port, entry.Metadata)
	}
	//
	wg.Wait()
	t.Log("TestRegisterMultipleMembers PASSED")
}

func TestUpdateTTL_DoesNotReregisterAfterShutdown(t *testing.T) {
	if testing.Short() {
		return
	}
	a := assert.New(t)

	p := newNacosProvider().(*Provider)
	c := newClusterForTest(serviceName, "127.0.0.1:8001", p)

	shutdownShouldHaveResolved := make(chan bool, 1)

	err := p.StartMember(c)
	a.NoError(err)

	time.Sleep(time.Second)
	found, _ := findService(t, p)
	a.True(found, "service was not registered in consul")

	go func() {
		// if after 5 seconds `Shutdown` did not resolve, assume that it will not resolve until `blockingUpdateTTL` resolves
		time.Sleep(5 * time.Second)
		shutdownShouldHaveResolved <- true
	}()

	err = p.Shutdown(true)
	a.NoError(err)
	shutdownShouldHaveResolved <- true

	// since `UpdateTTL` runs in a separate goroutine we need to wait until it is actually finished before checking the member's clusterstatus
	p.updateTTLWaitGroup.Wait()
	found, status := findService(t, p)
	a.Falsef(found, "service was still registered in consul after shutdown (service status: %s)", status)
}

func findService(t *testing.T, p *Provider) (found bool, status string) {
	service := p.cluster.Config.Name
	port := p.cluster.Config.RemoteConfig.Port
	entries, err := p.client.GetService(vo.GetServiceParam{
		ServiceName: service,
		GroupName:   groupName,
	})
	if err != nil {
		t.Error(err)
	}

	for _, entry := range entries.Hosts {
		if entry.Port == uint64(port) {
			return true, entry.InstanceId
		}
	}
	return false, ""
}
