# Day 7 · 测试 · 监控 · 生产最佳实践

**主题**：把前 6 天积累的 demo 收敛到工程级别 —— 用桩代码 mock 写单测、用 `testify/suite` 写集成测试、接入 runtime metrics、列出生产 checklist。

**核心目标**：

- 会用 `trpc create` 自动生成的 mock 文件配合 gomock 写表驱动测试；
- 会用 `testify/suite` 写"启动 server → 真实 RPC → 断言 → 关 server"的集成测试；
- 会接入 `trpc-metrics-runtime`，从 admin `/metrics` 看到 `trpc.runtime.*` 指标；
- 能给出一份"上线前 8 项 checklist"。

## 1. 测试金字塔

```
        ┌─────┐
        │ E2E │   ~5%   全链路、依赖真实环境
        ├─────┤
        │ Int │  ~15%   多组件协作（含 trpc server）
        ├─────┤
        │ Unit│  ~80%   纯函数、单 service 用 mock 隔离下游
        └─────┘
```

tRPC-Go 对每一层都有"姿势"：

- **Unit**：`trpc create` 默认生成的 `*_mock.go`（gomock）+ `testify/assert`；
- **Int**：`testify/suite` 启一个内存 server，用 ip:// 直连；
- **E2E**：部署到测试环境跑 `trpc-cli` 或自己的 client。

## 2. 工程目录

```
day07-test-and-best-practice/
├── README.md
├── go.mod
├── proto/
│   └── user.proto                # 沿用 day02
├── stub/                         # trpc create --rpconly 生成（含 *_mock.go）
├── service/
│   ├── user_service.go           # 业务实现（依赖注入 downstream）
│   └── user_service_test.go      # gomock 单测
├── integration/
│   └── suite_test.go             # testify/suite 集成测试
├── server/
│   └── main.go                   # 接入 trpc-metrics-runtime 上报
└── trpc_go.yaml
```

## 3. Unit 测试：用桩代码 mock 写表驱动测试

`trpc create` 生成的 `stub/.../user_mock.go` 自动包含 `MockUserServiceClientProxy`。我们演示一个聚合 service 用 mock 隔离下游：

```go
// service/user_service.go
type AggregatorService struct {
    downstream pb.UserServiceClientProxy   // 接口而非实现，可被替换为 mock
}

func (a *AggregatorService) Greet(ctx context.Context, id uint64) (string, error) {
    rsp, err := a.downstream.GetUser(ctx, &pb.GetUserReq{Id: id})
    if err != nil {
        if errs.Code(err) == 404 {
            return "", errs.New(404, "user not found")
        }
        return "", errs.Wrap(err, 502, "downstream failed")
    }
    return "Hello, " + rsp.GetUser().GetName(), nil
}
```

```go
// service/user_service_test.go
func TestAggregator_Greet(t *testing.T) {
    cases := []struct {
        name    string
        prepare func(m *pb.MockUserServiceClientProxy)
        want    string
        wantCode int
    }{
        {
            name: "ok",
            prepare: func(m *pb.MockUserServiceClientProxy) {
                m.EXPECT().
                    GetUser(gomock.Any(), gomock.Any()).
                    Return(&pb.GetUserRsp{User: &pb.User{Id: 1, Name: "Alice"}}, nil)
            },
            want: "Hello, Alice",
        },
        {
            name: "downstream not found",
            prepare: func(m *pb.MockUserServiceClientProxy) {
                m.EXPECT().GetUser(gomock.Any(), gomock.Any()).
                    Return(nil, errs.New(404, "no user"))
            },
            wantCode: 404,
        },
        {
            name: "downstream timeout",
            prepare: func(m *pb.MockUserServiceClientProxy) {
                m.EXPECT().GetUser(gomock.Any(), gomock.Any()).
                    Return(nil, errs.NewFrameError(errs.RetClientTimeout, "timeout"))
            },
            wantCode: 502,
        },
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mock := pb.NewMockUserServiceClientProxy(ctrl)
            tc.prepare(mock)

            got, err := (&AggregatorService{downstream: mock}).Greet(context.Background(), 1)
            if tc.wantCode != 0 {
                require.Error(t, err)
                assert.Equal(t, tc.wantCode, errs.Code(err))
            } else {
                require.NoError(t, err)
                assert.Equal(t, tc.want, got)
            }
        })
    }
}
```

## 4. Integration 测试：testify/suite

```go
// integration/suite_test.go
type UserSuite struct {
    suite.Suite
    server *trpc.Server
    client pb.UserServiceClientProxy
}

func (s *UserSuite) SetupSuite() {
    // 启动一个内存 server（端口随机选 0，让系统分配）
    s.server = trpc.NewServer(
        server.WithServiceName("trpc.test.user.UserService"),
        server.WithAddress("127.0.0.1:0"),
        server.WithProtocol("trpc"),
    )
    pb.RegisterUserServiceService(s.server, newRealUserImpl())
    go s.server.Serve()                                    // 异步启动

    addr := s.server.Service("...").Address()              // 拿到实际监听地址
    s.client = pb.NewUserServiceClientProxy(client.WithTarget("ip://" + addr))
}

func (s *UserSuite) TearDownSuite() { s.server.Close(nil) }

func (s *UserSuite) TestCreateAndGet() {
    crsp, err := s.client.CreateUser(context.Background(), &pb.CreateUserReq{Name: "Alice"})
    s.Require().NoError(err)

    grsp, err := s.client.GetUser(context.Background(), &pb.GetUserReq{Id: crsp.GetUser().GetId()})
    s.Require().NoError(err)
    s.Equal("Alice", grsp.GetUser().GetName())
}

func TestUserSuite(t *testing.T) { suite.Run(t, new(UserSuite)) }
```

要点：

- **`SetupSuite` 一次性、`SetupTest` 每用例一次**，重资源进 SetupSuite；
- `127.0.0.1:0` 让系统分配空闲端口，**避免并行测试端口冲突**；
- 集成测试不应依赖外部服务（DB、北极星），下游一律 mock 或起内嵌进程。

## 5. Metrics：接入 trpc-metrics-runtime

```go
// server/main.go
import (
    _ "trpc.group/trpc-go/trpc-metrics-runtime"   // 空 import 触发上报
)
```

启动后：

```powershell
curl http://127.0.0.1:11014/metrics
```

会看到：

```
trpc.runtime.gc.num_gc          12
trpc.runtime.goroutine.num      27
trpc.runtime.heap.alloc_bytes   12345678
trpc.runtime.cpu.usage_percent  3.2
...
```

业务自定义指标：

```go
import "trpc.group/trpc-go/trpc-go/metrics"

metrics.IncrCounter("biz.user.create.success", 1)
metrics.HistogramObserve("biz.user.create.latency_ms", 23.5)
```

## 6. 生产环境 Checklist（上线前必看 8 项）

1. **超时**：每跳 `timeout` 显式配置；调用方 ctx 携带 deadline；级联预算 = 总预算 × 0.6；
2. **重试**：仅对**幂等接口**开启 `client.WithFilter(retry)`；最大 2 次；带指数退避；
3. **熔断**：通过 selector 内置或自定义 filter 实现；阈值 50% 错误率 / 30s 窗口；
4. **限流**：服务端入口 filter（令牌桶），按 caller 分桶；保护下游用 `rate.Limiter`；
5. **日志**：禁用 `fmt.Printf`，全部走 `log.InfoContextf`；线上级别 info，慎开 debug；
6. **错误码**：业务错误码统一在常量包中定义（如 `const ErrUserNotFound = 404`），文档化；
7. **监控**：runtime + 业务自定义指标双管齐下；P95/P99 延迟、错误率、QPS 至少要有；
8. **配置**：敏感信息（DB 密码、token）走配置中心 + 加密，**绝不**进 yaml；admin 端口绑内网。

## 7. 与 day01-day07 的对照总结

| `study_go` 基础天 | 在 tRPC-Go 里的对应 |
| --- | --- |
| day01-basics（Go 语法） | tRPC stub 仍是普通 Go struct/interface |
| day02-collections（slice/map） | UserStore 仍用 map + 锁 |
| day03-oop（struct + 方法） | 业务实现 = struct 实现 stub interface |
| day04-interface-error（隐式接口、error） | `errs.Error` 实现 `error`；`errs.Code/Msg` 跨进程透传 |
| day05-concurrency（goroutine） | server 每请求一个 goroutine，handler 必须并发安全 |
| day06-advanced-sync（sync/Context） | `context.WithTimeout` + RWMutex 全部用上 |
| day07-web-server（net/http） | tRPC RESTful 协议 = "带 IDL 契约的 net/http" |

> **核心结论**：tRPC-Go 没有否定 Go 后端的任何基础知识，它只是把"工程脚手架"做完了 —— 让你少写胶水代码、多写业务代码。

## 8. 工程级面试复盘清单（最终版）

学完 7 天后的"自检题"：

1. 描述一个 RPC 请求从客户端到服务端的完整生命周期（含 filter、selector、编解码）。
2. tRPC 私有协议帧结构有哪几段？为什么要有"流 ID"？
3. `errs` 错误码区间为什么必须遵守？混用的代价是什么？
4. 一个进程里同时有 server 和 client，怎么共享 context、监控、日志？
5. server filter 抛 panic 会发生什么？为什么 `recovery` 必须放第一个？
6. 客户端 filter 中"选择器过滤器"为什么必须是最后一个？如果不是会出什么问题？
7. `context.WithTimeout` 透传到下游时，下游拿到的 deadline 是绝对时间还是相对时间？
8. 桩代码自动生成的 mock 文件里 `MockXxxClientProxy` 是怎么实现的？为什么不能直接 mock 业务接口？
9. RESTful 协议里 `body: "*"` 与 `body: "user"` 的差别，以及对前端联调的影响？
10. 一个 service 突然 P99 飙到秒级，admin 端口的哪几个接口能帮你定位？
11. yaml 热加载哪些字段会即时生效？哪些必须重启？为什么？
12. 业务码 404 在 tRPC 协议下和 RESTful 协议下表现的差异？
13. `polaris://servicename` 与 `dns://hostname:port` 在故障容忍上有什么差别？
14. `client.WithTarget` 与 yaml `client.service[].target` 同时存在时谁生效？
15. 你新设计一个 `OrderService`，proto 里至少要做哪 5 个决定？

能把这 15 题答出来，你已经是工程级 tRPC-Go 玩家。

## 9. 推荐继续阅读

- 官方文档 `docs/practice/`：性能调优、故障排查、灰度发布；
- `examples/features/` 下的 stream 示例：流式 RPC（双向 / 服务端流 / 客户端流）；
- 北极星 / Polaris 名字服务的真实接入；
- gRPC 互通：`protocol: grpc` 让 tRPC 服务能被 grpc-go 客户端调用。

---

**你已经走完了 14 天的学习路径**：前 7 天打 Go 后端基础，后 7 天上 tRPC 工程框架。下一步建议：找一个真实业务场景（订单、消息、推送），用这套框架落地一个 5-7 个 service 的微服务系统，把今天列的 8 项 checklist 全部跑一遍。
