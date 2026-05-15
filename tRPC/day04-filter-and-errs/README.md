# Day 4 · 拦截器（Filter） & 错误码（errs）

**主题**：把切面关注点（鉴权、耗时打点、限流）从业务代码中抽出来，统一走 filter 链；把"业务错误"和"框架错误"分开，按区间约定错误码。

**核心目标**：

- 能画出 server 与 client 双向 filter 的执行顺序图；
- 能写一个鉴权 server filter 与一个耗时 client filter；
- 知道 `errs` 错误码的区间约定，会用 `errs.New / errs.Wrap / errs.Code / errs.Msg`。

## 1. Filter 是什么

**Filter 就是 RPC 版本的中间件**，签名是"接收 ctx + req，调 next 拿到 rsp，自己加工"，与 day07 的 HTTP 中间件思路完全一致，只是工作在 RPC 层、感知到 service+method。

执行顺序示意（重要）：

```
[Client]
  filterA (前) ─→ filterB (前) ─→ selector → 网络 → ...
                                                   ↓
  filterA (后) ←─ filterB (后) ←─ selector ← ─────┘

[Server]
  收包 → filter1 (前) → filter2 (前) → handler → filter2 (后) → filter1 (后) → 发包
```

记住三条铁律：

1. **请求方向按数组顺序，响应方向按逆序**；
2. **所有 filter 共享同一个 `context`**，cancel/timeout 透明传递；
3. **filter 中只做 O(1)/O(log n) 的事**，绝不要做阻塞 I/O —— 否则会阻塞整条 RPC 链。

## 2. 工程目录

```
day04-filter-and-errs/
├── README.md
├── go.mod
├── filter/
│   ├── auth.go         # 服务端鉴权 filter
│   └── timing.go       # 客户端耗时 filter
├── server/
│   ├── main.go
│   └── service.go
├── client/
│   └── main.go
└── trpc_go.yaml
```

## 3. 服务端 filter：auth

业务诉求：调用方必须在 metadata 里带 `auth-token: demo`，否则拒绝服务。

```go
// filter/auth.go
package filter

import (
    "context"

    trpc "trpc.group/trpc-go/trpc-go"
    "trpc.group/trpc-go/trpc-go/errs"
    "trpc.group/trpc-go/trpc-go/filter"
    "trpc.group/trpc-go/trpc-go/log"
    "trpc.group/trpc-go/trpc-go/server"
)

const expectedToken = "demo"

// Auth 是服务端鉴权 filter。
// 注意签名：func(ctx, req, next) (rsp, err)，next 必须被调到，否则 handler 不会执行。
func Auth(ctx context.Context, req interface{}, next filter.ServerHandleFunc) (interface{}, error) {
    msg := trpc.Message(ctx)                          // 拿到本次 RPC 的 metadata 容器
    token := string(msg.ServerMetaData()["auth-token"])

    if token != expectedToken {
        log.WarnContextf(ctx, "auth failed: token=%q", token)
        // 401 是约定的"鉴权失败"业务错误码（区间 200-999）
        return nil, errs.New(401, "unauthorized")
    }

    return next(ctx, req)                             // 必须调；不调则请求挂死
}

// Register 在 init 中把 filter 注册进框架，yaml 里 filter: [auth] 才能找到它。
func Register() {
    server.RegisterFilter("auth", Auth)
}
```

## 4. 客户端 filter：timing

业务诉求：在客户端打点每次 RPC 的真实耗时（含序列化、网络、反序列化）。

```go
// filter/timing.go
package filter

import (
    "context"
    "time"

    "trpc.group/trpc-go/trpc-go/client"
    "trpc.group/trpc-go/trpc-go/filter"
    "trpc.group/trpc-go/trpc-go/log"
)

func Timing(ctx context.Context, req, rsp interface{}, next filter.ClientHandleFunc) error {
    start := time.Now()
    err := next(ctx, req, rsp)
    log.InfoContextf(ctx, "rpc cost=%s err=%v", time.Since(start), err)
    return err
}

func RegisterClient() {
    client.RegisterFilter("timing", Timing)
}
```

## 5. errs 错误码区间约定

| 错误码 | 类型 | 用途 | 来源 |
| --- | --- | --- | --- |
| 1 ~ 199 | 框架错误（`ErrorTypeFramework`） | 框架层的网络/编解码/超时 | tRPC-Go 内置 |
| 200 ~ 999 | 业务错误（`ErrorTypeBusiness`） | 业务自定义 | 业务方 |
| 1 ~ 199 | 下游框架错误（`ErrorTypeCalleeFramework`） | 调下游时下游回的框架错 | 框架自动包装 |

> ⚠️ **业务错误码上限 999**。Day2 中我们在 demo 里写的 `errs.New(404, ...)` 是合规的；但如果业务想表达更细的语义（如 `1001 用户名重复`），按规范要么挤进 `200-999` 区间（如 `409` 表冲突），要么用 `errs.Wrap` 在 msg 中编码。

四种正确姿势：

```go
// 1) 创建业务错误（最常用）
return nil, errs.New(404, "user not found")

// 2) 包装下游错误，保留原始 chain
if err := proxy.Call(ctx, ...); err != nil {
    return nil, errs.Wrap(err, 500, "call downstream failed")
}

// 3) 提取错误码（client 侧最常用）
if errs.Code(err) == 404 {
    // 用户不存在，按业务降级处理
}

// 4) 提取错误消息
fmt.Println(errs.Msg(err))
```

## 6. 配置 filter：代码式 vs 配置式

**代码式**（在 `trpc.NewServer()` 之前用 Option）：

```go
s := trpc.NewServer(server.WithFilter(filter.Auth))
```

**配置式**（推荐，运行期可调）：

```yaml
server:
  service:
    - name: trpc.study.user.UserService
      filter: [auth]                # 这里的 "auth" 必须事先 server.RegisterFilter 过
```

经验：

- **共享逻辑（recovery、metrics、tracing）走配置式** —— 改 yaml 即可调整顺序；
- **临时调试逻辑走代码式** —— 比如某个 service 临时加个调试 filter，不污染 yaml。

## 7. 跑起来

```powershell
cd .\tRPC\day04-filter-and-errs
trpc create -p ..\day02-userservice-demo\proto\user.proto -o . --rpconly
go mod tidy

# 窗口 1
go run .\server\

# 窗口 2
go run .\client\
```

期望输出：

```
[client] CreateUser ok: id=1 name=Alice cost=8.2ms     # timing filter 打点
[client] GetUser   ok: id=1 cost=2.1ms
[client] 不带 token 调用 → code=401 msg=unauthorized
```

服务端窗口同时会打印 `auth failed: token=""`。

## 8. 验证标准

- [ ] 客户端日志能看到每次 RPC 的耗时；
- [ ] 客户端故意不传 token 时，服务端打印鉴权失败日志，客户端拿到 `code=401`；
- [ ] 把 yaml 中的 `filter: [auth]` 删掉，服务端不再鉴权；
- [ ] 在 server filter 中故意 `time.Sleep(2 * time.Second)`，观察客户端因超时拿到 `code=101`（框架超时码）。

## 9. 面试复盘

1. **filter 与 day07 的 `logRequest` 中间件本质区别？** Day07 中间件作用于 HTTP `Handler`，只能感知 method+path；filter 作用于 RPC `service+method`，能拿到强类型 req/rsp、metadata、错误码。
2. **client filter 中的"选择器过滤器"为什么必须在最后一个？** 因为它要做"已知 service 名 → 选出具体节点"这个动作 —— 必须在所有自定义 filter 处理完前置逻辑后才执行，否则负载均衡决策会基于不完整的上下文。
3. **`errs.New(404, "user not found")` 和 HTTP 协议下的 404 是同一个吗？** 在 tRPC 私有协议下，404 是业务错误码、跨进程透传；在 HTTP 协议下，框架会按映射表把它转成 HTTP 状态码 404。**值看起来一样，是巧合也是惯例**。
4. **`errs.Wrap(err, code, msg)` 和 day04 标准库 `fmt.Errorf("...: %w", err)` 区别？** `errs.Wrap` 保留 chain 的同时**重置 code** —— 让外层接口能用自己的语义码报错；`%w` 只 wrap 不能改 code，因为标准库 error 没有 code 概念。
5. **filter 里能改 `req` 吗？** 技术上能（指针修改），但**强烈不推荐**。filter 应当是"观察者"或"短路终止者"，不要悄悄改业务参数 —— 排障时会让人怀疑人生。
