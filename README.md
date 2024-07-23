# Sylvain

sylvain 定位是提供给服务服务注册发现的能力, 支持快速结合 `etcd` 的 API 实现服务注册发现。

## 现在开始

> 请将 Go 升级到 ^1.17

```
go get -u github.com/luoyangwei/sylvain
```

## 例子

### provider

> 查看 example/provider

```golang
// etcd 节点
endpoints := sylvain.NewEtcdEndpoints([]string{"127.0.0.1:2379"}, 5*time.Second)
// 将 provider 注册到 etcd
_ = sylvain.NewNamedServerDiscover(endpoints, sylvain.WithServerNamed("provider"), sylvain.WithServerPort(8888))
```

### customer

> 查看 example/customer

```golang
// etcd 节点
endpoints := sylvain.NewEtcdEndpoints([]string{"127.0.0.1:2379"}, 5*time.Second)
// 将 customer 注册到 etcd
serverDiscover := sylvain.NewNamedServerDiscover(endpoints, sylvain.WithServerNamed("customer"), sylvain.WithServerPort(9999))
// 发现 provider
providerServer := serverDiscover.ElectionServerEndpoint("provider")
```
