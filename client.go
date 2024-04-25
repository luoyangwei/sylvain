package sylvain

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

// Client 包装自 etcd 的 client 客户端，负责调用 etcd 的 client 客户端
type Client struct {
	*clientv3.Client
}

func NewEtcdEndpoints(endpoints []string, timeout time.Duration) *Client {
	clt, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return &Client{Client: clt}
}
