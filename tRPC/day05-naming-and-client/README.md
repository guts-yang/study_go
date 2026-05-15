# Day 5 · 名字服务 & 下游调用

**主题**：让一个服务（`gateway`）通过名字服务发现并调用另一个服务（`user`），理解 `Selector` 四件套，掌握"服务端同时也是客户端"的工程范式。

**核心目标**：

- 能解释 `Selector = Discovery + ServiceRouter + LoadBalance + CircuitBreaker`；
- 会用 `client.WithTarget("ip://...")` 直连，也会通过 yaml 配置 `callee` 完成寻址；
- 会自定义一个最小化 selector 插件，理解 `scheme://servicename` 协议；
- 知道 `client.WithTimeout` 与 `context.WithTimeout` 的优先级与组合策略。

## 1. 拓扑

```
Client ──▶ Gateway (port 8002, service A) ──▶ User (port 8001, service B, day02 风格)
```

Gateway 是个聚合服务：

- 对外提供 `GreetUser(id)` 方法；
- 内部用 `UserServiceClientProxy` 调下游 `User` 拿用户信息，再拼成 "Hello, Alice"。

## 2. Selector 四件套

```
client.NewProxy → 调用前 → Selector.Select(serviceName)
                                  ├─ Discovery       通过 service name 拿节点列表
                                  ├─ ServiceRouter   按规则过滤节点（机房/集群/灰度）
                                  ├─ LoadBalance     从剩余节点选一个
                                  └─ CircuitBreaker  熔断器判断该节点是否健康
                                  ↓
                              发起请求
                                  ↓
                          Selector.Report(node, cost, err)   反馈给熔断器
```

记住四件套各自的责任：

| 组件 | 职责 |
| --- | --- |
| Discovery | 拿到 `serviceName` 对应的全量节点 |
| ServiceRouter | 按规则过滤（同机房优先、灰度白名单等） |
| LoadBalance | 选一个（轮询、加权、一致性哈希、最小连接） |
| CircuitBreaker | 健康判断 + 上报反馈 |

`client.WithTarget("scheme://servicename")` 中 **scheme 决定用哪个 Selector**：

| scheme | 含义 |
| --- | --- |
| `ip://127.0.0.1:8001` | 直连（本地调试） |
| `dns://hostname:port` | DNS 解析 |
| `polaris://...` | 北极星（公司内部） |
| `consul://...` | Consul |
| `example://service-a` | 自定义（本 day 演示） |

## 3. 工程目录

```
day05-naming-and-client/
├── README.md
├── go.mod
├── proto/
│   └── user.proto         # 复用 day02
├── stub/
├── selector/
│   └── example.go         # 自定义 selector 插件
├── service-b-user/        # 下游：User 服务（端口 8001）
│   ├── main.go
│   └── service.go
├── service-a-gateway/     # 上游：Gateway 服务（端口 8002）
│   ├── main.go
│   └── service.go         # 在 handler 中调下游
├── client/
│   └── main.go            # 调 Gateway
└── trpc_go.yaml
```

## 4. 关键代码片段

### 4.1 自定义 selector

```go
// selector/example.go
package selector

import (
    "errors"
    "math/rand"
    "time"

    "trpc.group/trpc-go/trpc-go/naming/registry"
    "trpc.group/trpc-go/trpc-go/naming/selector"
)

// store 模拟一个最小化的服务注册表。
var store = map[string][]*registry.Node{
    "trpc.study.user.UserService": {
        {Address: "127.0.0.1:8001", Network: "tcp"},
    },
}

type exampleSelector struct{}

func (s *exampleSelector) Select(serviceName string, _ ...selector.Option) (*registry.Node, error) {
    list, ok := store[serviceName]
    if !ok || len(list) == 0 {
        return nil, errors.New("no available node for " + serviceName)
    }
    return list[rand.Intn(len(list))], nil
}

func (s *exampleSelector) Report(_ *registry.Node, _ time.Duration, _ error) error {
    return nil
}

func init() {
    selector.Register("example", &exampleSelector{})
}
```

启用方式：`client.WithTarget("example://trpc.study.user.UserService")`。

### 4.2 Gateway 调下游

```go
// service-a-gateway/service.go
type gatewayImpl struct {
    userCli pb.UserServiceClientProxy
}

func newGatewayImpl() *gatewayImpl {
    return &gatewayImpl{
        userCli: pb.NewUserServiceClientProxy(
            client.WithTarget("example://trpc.study.user.UserService"),
            client.WithTimeout(500 * time.Millisecond),
        ),
    }
}

func (g *gatewayImpl) GreetUser(ctx context.Context, req *pb.GreetUserReq) (*pb.GreetUserRsp, error) {
    // 给下游传递 deadline；context 携带 trace_id 自动透传
    rsp, err := g.userCli.GetUser(ctx, &pb.GetUserReq{Id: req.GetId()})
    if err != nil {
        return nil, errs.Wrap(err, 502, "downstream user service failed")
    }
    return &pb.GreetUserRsp{Greeting: "Hello, " + rsp.GetUser().GetName()}, nil
}
```

要点：

- `pb.NewUserServiceClientProxy` 在 Gateway 启动时**一次构造**，handler 里复用 —— **proxy 是并发安全的**；
- `errs.Wrap(err, 502, "...")` 把"下游错误"包装成"网关业务错"，错误码区间自洽；
- 下游 `GetUser` 失败时 `err` 链中保留了原始 errs，排障可用 `errs.Code(errs.Unwrap(err))` 拿到下游码。

## 5. 超时组合策略

```
[Outer ctx 1500ms] ─→ Gateway.WithTimeout(500ms) ─→ User
                                ↑
                           取最小值生效
```

- **`client.WithTimeout(d)`** 是该 callee 的"上限"；
- **`context.WithTimeout(ctx, d)`** 是这次请求的"上限"；
- **谁先到 deadline 谁生效**。生产规则：链路总预算 = 用户感知超时；每跳 timeout = 预算 / 跳数 × 安全系数（通常 0.6）。

## 6. proto 扩展

在 day02 `user.proto` 基础上新增一个 service：

```protobuf
service GatewayService {
  rpc GreetUser(GreetUserReq) returns (GreetUserRsp);
}

message GreetUserReq { uint64 id       = 1; }
message GreetUserRsp { string greeting = 1; }
```

> 一个 proto 文件可以定义多个 service。生成的 stub 会同时给出 `UserService*` 和 `GatewayService*` 的全套类型。

## 7. 跑起来（三窗口）

```powershell
cd .\tRPC\day05-naming-and-client
trpc create -p .\proto\user.proto -o . --rpconly
go mod tidy

# 窗口 1：下游 User
go run .\service-b-user\

# 窗口 2：上游 Gateway
go run .\service-a-gateway\

# 窗口 3：客户端
go run .\client\
```

期望客户端输出：

```
CreateUser via User      → id=1
GreetUser  via Gateway   → "Hello, Alice"
```

## 8. 验证标准

- [ ] 把 `service-b-user` 杀掉再调 Gateway，得到 `code=502 msg=downstream user service failed`；
- [ ] 把 `client.WithTimeout` 改成 `1ms`，验证 Gateway 立即超时（框架码 101）；
- [ ] 把 `selector` 中的端口故意改成 `8009`（不存在），验证 selector 报错；
- [ ] 把 `WithTarget` 改成 `ip://127.0.0.1:8001`，绕过 selector 也能跑通。

## 9. 面试复盘

1. **`ip://` 与 `polaris://` 的根本区别？** 前者是写死的端点，后者是"问名字服务要节点"；前者无熔断、无路由、无健康检查，后者四件套俱全。
2. **服务端如何在 handler 里调下游？** 在 service struct 里持有 `proxy`，构造时初始化、handler 里直接 `proxy.Method(ctx, req)`。**绝不要**在 handler 内动态 new proxy（连接池没法复用）。
3. **`context.WithTimeout` 透传到下游的 deadline 是"剩余时间"还是"原始时间"？** **剩余时间**。框架在 metadata 里写入 `deadline`，下游自己减去当前时间得到剩余预算。
4. **为什么 Selector 的 `Report` 方法很重要？** 没有 Report 就没有"正/负反馈"，熔断器无法判断节点健康度。生产里跳过 Report 等于禁用熔断。
5. **下游变更（IP 漂移、扩容）时，上游 proxy 需要重启吗？** 不需要。Selector 内部的 Discovery 会异步刷新节点列表（北极星默认 1s 一次），proxy 引用的是 Selector，对它而言节点变更是透明的。
