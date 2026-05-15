// Package selector 实现一个最小化的自定义 selector 插件 "example://"。
//
// 真实生产环境中，selector 通常对接：
//   - 北极星（公司内）
//   - Consul / etcd
//   - DNS
//   - k8s service
//
// 本 day 的 example selector 仅用一个写死的 map 演示 Selector 接口长什么样。
package selector

import (
	"errors"
	"math/rand"
	"time"

	"git.code.oa.com/trpc-go/trpc-go/naming/registry"
	"git.code.oa.com/trpc-go/trpc-go/naming/selector"
)

// store 模拟一个最小化注册表：service name → 节点列表。
// 真实场景里，这个 map 是异步从注册中心拉来的。
var store = map[string][]*registry.Node{
	"trpc.study.user.UserService": {
		{Address: "127.0.0.1:8001", Network: "tcp"},
	},
}

// exampleSelector 实现 selector.Selector 接口。
type exampleSelector struct{}

// Select 是核心方法：根据 service name 返回一个节点。
// 真实实现要在这里塞 ServiceRouter 过滤、LoadBalance 选择、CircuitBreaker 检查。
func (s *exampleSelector) Select(serviceName string, _ ...selector.Option) (*registry.Node, error) {
	list, ok := store[serviceName]
	if !ok || len(list) == 0 {
		return nil, errors.New("no available node for " + serviceName)
	}
	return list[rand.Intn(len(list))], nil
}

// Report 由框架在每次 RPC 完成后回调，是熔断器的输入信号源。
// 这里 demo 简单返回 nil，真实实现要更新节点的成功率/延迟统计。
func (s *exampleSelector) Report(_ *registry.Node, _ time.Duration, _ error) error {
	return nil
}

// init 把插件注册到框架。
// 业务侧 client.WithTarget("example://...") 即会路由到这里。
func init() {
	selector.Register("example", &exampleSelector{})
}
