package sylvain

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"math/rand"
	"time"
)

type ServerDiscoverOption func(svr *Server)

//type ServerDiscoverHandler func(svr *Server)

type NamedServerDiscover struct {
	leaseID clientv3.LeaseID

	svr    *Server
	clt    *Client
	ticker *time.Ticker

	leaseKeepAlive <-chan *clientv3.LeaseKeepAliveResponse

	// handler 当服务数据变动时，会调用这个函数通知监听
	//handler ServerDiscoverHandler
}

func NewNamedServerDiscover(clt *Client, opts ...ServerDiscoverOption) *NamedServerDiscover {
	svr := &Server{ip: getHostIp()}
	for _, opt := range opts {
		opt(svr)
	}

	svr.Addr = svr.GetAddr()

	namedServerDiscover := &NamedServerDiscover{
		svr:    svr,
		clt:    clt,
		ticker: time.NewTicker(3 * time.Second),
	}
	namedServerDiscover.publishWithLeaseKey(5)

	return namedServerDiscover
}

func (discover *NamedServerDiscover) publishWithLeaseKey(lease int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grant, err := discover.clt.Grant(ctx, lease)
	if err != nil {
		panic(err)
	}

	discover.leaseID = grant.ID
	serverBuf, _ := discover.svr.MarshalBinary()
	_, err = discover.clt.Put(ctx, discover.svr.Name, string(serverBuf), clientv3.WithLease(grant.ID))
	if err != nil {
		panic(err)
	}

	go discover.keepAlive()
	go discover.monitor()
}

func (discover *NamedServerDiscover) keepAlive() {
	for {
		select {
		case <-discover.ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err := discover.clt.KeepAliveOnce(ctx, discover.leaseID)
			if err != nil {
				panic(err)
			}
			cancel()
		}
	}
}

func (discover *NamedServerDiscover) monitor() {
	for {
		select {
		case response, ok := <-discover.clt.Watch(context.Background(), "server", clientv3.WithPrefix()):
			if !ok {
				return
			}
			fmt.Println(response)
		}
	}
}

// ElectionServerEndpoint 选举服务，根据算法会返回一个最佳的服务器地址
func (discover *NamedServerDiscover) ElectionServerEndpoint(name string) string {
	response, err := discover.clt.Get(context.Background(), serverNamed(name))
	if err != nil {
		// ERROR
	}

	if response.Count == 0 {
		return ""
	}

	// 服务列表序列化，将所有的列表都序列化出来
	var servers = make([]*Server, response.Count)
	for i, value := range response.Kvs {
		var server = new(Server)
		_ = server.UnmarshalBinary(value.Value)
		servers[i] = server
	}

	// 生成随机数，随机因子是 response 里的 Count
	idx := rand.Intn(int(response.Count))
	if idx == 0 {
		idx++
	}

	// 随机选举，每一个服务的概率是一样的
	return servers[idx-1].Addr
}

func (discover *NamedServerDiscover) Close() {
	discover.ticker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if _, err := discover.clt.Revoke(ctx, discover.leaseID); err != nil {
		panic(err)
	}

	log.Println("Delete server", discover.svr.Name)
	if _, err := discover.clt.Delete(ctx, discover.svr.Name); err != nil {
		panic(err)
	}
	cancel()
	_ = discover.clt.Close()
}

func serverNamed(name string) string {
	return fmt.Sprintf("server/%s", name)
}

func WithServerNamed(name string) ServerDiscoverOption {
	return func(svr *Server) {
		svr.Name = serverNamed(name)
	}
}

func WithServerPort(port int) ServerDiscoverOption {
	return func(svr *Server) {
		svr.port = port
	}
}
