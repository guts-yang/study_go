---
name: trpc-go-7day-learning-plan
overview: 在 study_go/tRPC 目录下产出一份对齐 day01~day07 风格的 tRPC-Go 7 天工程级学习规划，包含从 0 到 1 上手 tRPC、用户服务（UserService）端到端 demo、以及拦截器/插件/监控/单测等工程实践，使用官方推荐的 trpc-cmdline 工具链。
todos:
  - id: root-readme
    content: 编写 tRPC/README.md：7 天总表、Windows+PowerShell 环境准备、版本选型与 GOPROXY、学习节奏、工程级面试复盘清单
    status: completed
  - id: day01-helloworld
    content: 使用 [subagent:code-explorer] 核对 v2 包路径，落地 day01-helloworld：proto + trpc create 生成 + server/client + yaml + README
    status: completed
    dependencies:
      - root-readme
  - id: day02-userservice-demo
    content: 落地核心 demo day02-userservice-demo：user.proto 设计、生成 stub、server（map+RWMutex 复用 day07）、client、yaml、详尽 README 与验证步骤
    status: completed
    dependencies:
      - day01-helloworld
  - id: day03-config-and-admin
    content: 落地 day03-config-and-admin：yaml 三段详解、log 标准接口、admin 端口与健康检查、热加载验证、README
    status: completed
    dependencies:
      - day02-userservice-demo
  - id: day04-filter-and-errs
    content: 落地 day04-filter-and-errs：服务端 auth filter、客户端 timing filter、errs 区间约定与示例、配置式 vs 代码式注册、README
    status: completed
    dependencies:
      - day03-config-and-admin
  - id: day05-naming-and-client
    content: 落地 day05-naming-and-client：service-a 调用 service-b 演示下游调用、WithTarget 寻址、自定义 selector、超时配置、README
    status: completed
    dependencies:
      - day04-filter-and-errs
  - id: day06-http-restful
    content: 落地 day06-http-restful：proto 加 google.api.http 注解、yaml 同时暴露 trpc+http、curl 与 trpc client 双调用、错误码映射表、README
    status: completed
    dependencies:
      - day05-naming-and-client
  - id: day07-test-and-best-practice
    content: 落地 day07-test-and-best-practice：mock 单测、testify/suite 集成测试、trpc-metrics-runtime、生产 checklist、与 day01–day07 对照总结、最终面试复盘
    status: completed
    dependencies:
      - day06-http-restful
---

## Product Overview

在 `study_go/tRPC/` 目录下整理一份 **7 天工程级 tRPC-Go 学习规划**，作为已完成的 day01–day07 Go 基础训练（语法/集合/OOP/接口/并发/sync/net-http）的进阶篇。整套规划严格对齐根目录 `README.md` 的"dayXX 一日一目录 + 独立 README + 可跑代码 + 面试复盘"风格，从 0 到 1 帮助使用者认识 tRPC-Go，并在 Day2 落地一个延续 day07 思路的 **UserService（CreateUser / GetUser）服务端 + 客户端连接 demo**。

## Core Features

- **7 天阶梯路线**：Day1 跑通 helloworld → Day2 端到端 UserService demo → Day3 配置与 admin → Day4 拦截器与错误码 → Day5 名字服务与下游调用 → Day6 HTTP/RESTful 多协议 → Day7 单测/监控/最佳实践复盘。
- **入口文档 `tRPC/README.md`**：包含 7 天总表、Windows + PowerShell 环境准备（Go 1.22、protoc、trpc-cmdline）、学习节奏建议、工程级面试复盘清单、v1(git.woa.com)/v2(trpc.tech)/开源(github.com/trpc-group)的版本选型说明。
- **每日子目录**：均含独立 `README.md`（主题 / 学习目标 / 知识点 / 可运行命令 / 当天验证标准 / 面试复盘）、可执行 Go 源码、必要的 `*.proto` 与 `trpc_go.yaml` 配置。
- **核心 Demo（Day2）**：`user.proto` 定义 CreateUser / GetUser，使用 `trpc create` 生成桩代码与 yaml 模板；服务端用内存 map + `sync.RWMutex` 复用 day07 思路；客户端用 `NewUserServiceClientProxy` 调用；提供 PowerShell 启动与请求示例。
- **进阶专题落地**：Day4 提供 `auth filter` + `errs.New` 错误码示例；Day5 演示 `client.WithTarget("ip://...")` 与服务端调用下游；Day6 通过 `google.api.http` 注解让同一服务同时暴露 tRPC + RESTful；Day7 演示 mock 桩代码 + `testify/suite` 集成测试 + `trpc-metrics-runtime` 接入。
- **不改动现有 day01–day07 内容**，仅在 `tRPC/` 下新增；所有命令以 PowerShell 形式给出，并附 GOPROXY、离线 fallback、版本兼容提示。

## 学习成果验收

学完后能独立回答：tRPC-Go 与 net/http 的本质差异、PB 桩代码生成原理、yaml 配置三段（server/client/plugins）的作用、filter 执行顺序、Selector 四件套、HTTP-RPC 互转机制、错误码区间约定、桩 mock 单测套路、生产环境监控与日志上报规范。

## Tech Stack

- **语言/框架**：Go 1.22（沿用现有 `go.mod`）+ tRPC-Go v2 主推，包路径优先使用 `trpc.tech/trpc-go/trpc-go/v2`；如无内网访问，提供 fallback 到开源版 `github.com/trpc-group/trpc-go`。
- **IDL 与桩代码**：Protocol Buffers（proto3）+ trpc-cmdline（`trpc create -p xxx.proto -o out --domain=trpc.tech --versionsuffix=v2`），桩代码默认带 mock。
- **配置**：`trpc_go.yaml`（server / client / plugins 三段，支持 fsnotify 热加载）。
- **存储**：Day2 demo 使用进程内 `map[uint64]*pb.User` + `sync.RWMutex`（与 day07 一致，保持轻量）。
- **测试**：`testing` + `github.com/stretchr/testify/{assert,require,suite}` + `github.com/golang/mock/gomock`。
- **监控/日志**：框架自带 `log`/`metrics` 标准接口；Day7 演示导入 `trpc.tech/trpc-go/trpc-metrics-runtime/v2` 上报 runtime 指标。
- **开发环境**：Windows + PowerShell（所有命令以 PS 形式给出，路径用反斜杠或正斜杠均可，强调 `setx GOPROXY https://goproxy.cn,direct`）。

## Implementation Approach

- **策略**：以官方推荐学习路径（helloworld → quick_start → user_guide → developer_guide → practice）为骨架，压缩为 7 天，每天一个独立可跑目录 + README，难度由浅入深、不重叠。
- **Day2 端到端 demo 设计**：先写 `proto/user.proto` → `trpc create` 生成到 `stub/` → 在 `server/main.go` 实现 `UserService` 接口 → 在 `client/main.go` 用 `pb.NewUserServiceClientProxy` 调用 → 通过 `trpc_go.yaml` 配置端口与协议；服务侧使用 day07 同款 `sync.RWMutex + map` 存储以建立心智迁移感。
- **关键决策与理由**：
- 选 v2（`trpc.tech/.../v2`）：知识库明确 v0.18.x 是 LTS，v2 是正在演进的主线，新建项目应贴 v2；同时给出开源版 fallback 兼顾外网。
- 用 `trpc create` 而非手写桩：贴官方推荐，且自动产出 yaml/mock，减少认知负担。
- 7 天严格对齐 day01–day07 风格：复用根 README 的体感（天数 / 目录 / 主题 / 核心目标 / 命令 / 复盘清单），降低切换成本。
- 不在 plan 阶段写文件：plan 模式禁止写盘，规划仅产出"目录树 + 文件清单 + 章节标题 + 关键片段示意 + 验证标准"。
- **性能与可靠性**：Day2 demo 在 hot path 用 `RWMutex` 而非 `Mutex`（读多写少），避免与 day07 重复踩坑；Day4 的 filter 中只放 O(1) 操作（取 metadata、写日志），强调 filter 不应做阻塞 I/O；Day5 强调 `client.WithTimeout` + `context.WithTimeout` 的组合，避免级联超时。

## Implementation Notes

- **Grounded**：路径/包名/命令必须与知识库一致（`trpc.NewServer()` / `pb.RegisterUserServiceService` / `errs.New(code, msg)` / `client.WithTarget("ip://127.0.0.1:8000")`）。fallback 命名空间为 `github.com/trpc-group/trpc-go`，仅在 README 中作为可选提示，不混入示例代码。
- **版本统一**：所有示例 import 统一使用 v2（`trpc.tech/trpc-go/trpc-go/v2`），README 顶部说明若选择开源版需要将 import 路径整体替换。
- **错误处理**：所有业务错误必须 `errs.New(code, msg)`；错误码遵循 `200~999` 业务区间；Day4 给出 `ErrUserNotFound = 1001` 之类的常量约定（注意：业务错误码上限 999，规划中给出 200–999 区间表，避免越界）。
- **日志**：Day3 起一律用 `trpc.Log` / `log.Infof`，禁用 `fmt.Printf`，README 中显式提示原因（无法上报远程日志中心）。
- **Windows/PowerShell**：所有命令都给 PS 版本，路径示例使用 `./...`；`go install` 装 trpc-cmdline 后强调 `$env:Path` 是否包含 `$env:GOPATH\bin`。
- **Blast radius**：本规划只新增 `tRPC/` 子目录，不动 day01–day07；`go.mod` 的 `module study_go` 保持不变，子目录用 `study_go/tRPC/dayXX-...` 作为内部包路径，外部依赖（trpc-go、testify、mock、protobuf）由 `go get` 增量加入 `go.sum`。

## Architecture Design

整体是"一份入口 README + 7 个独立日目录"的横向并列结构，Day2 是核心垂直 demo，其它日围绕 Day2 中已生成的 UserService 做正交扩展（同一 proto，不同主题）。

```mermaid
graph TD
    Root[tRPC/README.md<br/>7 天总表 + 环境准备 + 复盘清单]
    Root --> D1[day01-helloworld<br/>跑通官方 helloworld]
    Root --> D2[day02-userservice-demo<br/>核心 demo: server + client]
    Root --> D3[day03-config-and-admin<br/>yaml + log + admin]
    Root --> D4[day04-filter-and-errs<br/>拦截器 + 错误码]
    Root --> D5[day05-naming-and-client<br/>Selector + 下游调用]
    Root --> D6[day06-http-restful<br/>HTTP/RESTful 多协议]
    Root --> D7[day07-test-and-best-practice<br/>mock 单测 + metrics + 复盘]

    D2 -. 复用 user.proto .-> D3
    D2 -. 复用 user.proto .-> D4
    D2 -. 复用 user.proto .-> D5
    D2 -. 复用 user.proto .-> D6
    D2 -. 复用 user.proto .-> D7
```

## Directory Structure

### Directory Structure Summary

本规划仅在 `tRPC/` 目录下新增内容，不修改 day01–day07。每日目录都自包含一个可独立 `go run` 的迷你项目（`server/`、`client/`、`proto/`、`stub/`、`trpc_go.yaml`），共享 Day2 沉淀下来的 `user.proto` 心智模型；每日 README 严格按照"主题 / 知识点 / 命令 / 验证标准 / 面试复盘"五段式书写，与根 README 风格一致。

```
study_go/
└── tRPC/
    ├── README.md                                # [NEW] tRPC 学习入口文档：7 天总表（仿根 README 表格）、Windows+PowerShell 环境准备（Go 1.22 / protoc / trpc-cmdline 安装）、版本选型（v2 主推 + 开源版 fallback + GOPROXY 设置）、学习节奏（先读 README -> 看代码 -> 跑 -> 改 -> 复述）、工程级面试复盘清单（约 15 条：tRPC vs net/http、PB 桩代码、yaml 三段、filter 顺序、Selector 四件套、HTTP-RPC 互转、errs 区间、mock 套路、log/metrics 标准接口等）。
    │
    ├── day01-helloworld/
    │   ├── README.md                            # [NEW] 章节：1)tRPC 是什么 2)与 net/http 的差异 3)trpc-cmdline 安装与原理 4)proto3 速览 5)PowerShell 跑通命令 6)桩代码目录解读 7)验证标准（curl/客户端打印 "Hello, world"）8)面试复盘（5 条）。
    │   ├── proto/helloworld.proto               # [NEW] 官方风格 Greeter 服务，含 SayHello(HelloRequest)->HelloReply。
    │   ├── stub/                                # [NEW] 由 `trpc create -p proto/helloworld.proto -o stub --rpconly` 生成（README 给命令，不预写文件内容）。
    │   ├── server/main.go                       # [NEW] 最小服务端：trpc.NewServer() + pb.RegisterGreeterService + Serve()，业务实现拼接 "Hello, " + req.Msg。
    │   ├── client/main.go                       # [NEW] 最小客户端：proxy := pb.NewGreeterClientProxy(client.WithTarget("ip://127.0.0.1:8000"))；调用 SayHello 并打印响应。
    │   └── trpc_go.yaml                         # [NEW] 最小化 server 段（service name=trpc.helloworld.Greeter, network=tcp, port=8000, protocol=trpc）。
    │
    ├── day02-userservice-demo/                  # 核心 demo：从 day07 net/http 用户服务迁移到 tRPC
    │   ├── README.md                            # [NEW] 章节：1)需求与 day07 对照 2)proto 设计要点（包名/服务名/字段编号） 3)trpc create 命令完整步骤 4)生成产物解读 5)server 实现讲解（map+RWMutex 复用 day07） 6)client 调用讲解 7)PowerShell 双终端运行步骤 8)验证标准（CreateUser 后 GetUser 能拿回数据） 9)面试复盘（PB 字段编号为何不能改、Proxy 为何线程安全、与 net/http 的迁移代价）。
    │   ├── proto/user.proto                     # [NEW] syntax=proto3; package trpc.study.user; service UserService { rpc CreateUser(CreateUserReq) returns (CreateUserRsp); rpc GetUser(GetUserReq) returns (GetUserRsp); }; message User{uint64 id=1; string name=2; int64 created_at=3;}; 含 CreateUserReq{string name=1;}, CreateUserRsp{User user=1;}, GetUserReq{uint64 id=1;}, GetUserRsp{User user=1;}。
    │   ├── stub/                                # [NEW] trpc create 生成（README 给命令）。
    │   ├── server/
    │   │   ├── main.go                          # [NEW] 注册服务：s := trpc.NewServer(); pb.RegisterUserServiceService(s, &userImpl{store: newStore()}); s.Serve()。
    │   │   └── service.go                       # [NEW] userImpl：内嵌 store；CreateUser 自增 id 后写入 map；GetUser 读不存在时返回 errs.New(404, "user not found")，提示后续 Day4 会改为业务错误码。
    │   ├── client/main.go                       # [NEW] 顺序调用 CreateUser -> GetUser，打印两次响应；演示 client.WithTimeout(time.Second)。
    │   └── trpc_go.yaml                         # [NEW] server.service[0].name=trpc.study.user.UserService, port=8001；client 段示例性留空，由 day05 启用。
    │
    ├── day03-config-and-admin/
    │   ├── README.md                            # [NEW] 章节：1)yaml 三段（global/server/client/plugins）逐字段解读 2)log 插件配置（控制台 + 文件 + 级别） 3)admin 端口启用与默认接口（/cmds /healthz /metrics /pprof） 4)框架 log API 用法（log.Infof/log.WithFields） 5)fsnotify 热加载验证 6)PowerShell 命令 7)验证标准（改 yaml 不重启日志级别生效；curl admin 健康检查 200） 8)面试复盘（为什么不能用 fmt.Printf）。
    │   ├── server/main.go                       # [NEW] 复用 Day2 的 UserService 实现（import 路径指向 day02 的 stub），重点演示 log.Infof 与 admin 自定义 HandleFunc。
    │   ├── trpc_go.yaml                         # [NEW] 含 admin: { ip: 127.0.0.1, port: 11014 }；plugins.log.default 配置 console+file+level=debug。
    │   └── (无 client，复用 Day2 client)
    │
    ├── day04-filter-and-errs/
    │   ├── README.md                            # [NEW] 章节：1)filter 执行顺序图（请求顺序/响应逆序） 2)errs 区间约定（框架 1-199、业务 200-999） 3)自定义 server filter（鉴权：检查 trpc-auth-token） 4)自定义 client filter（耗时打点） 5)errs.New / errs.Wrap / errs.Code 用法 6)PowerShell 命令 7)验证标准（无 token 返回 errs.RetServerAuthFail；CreateUser 重名返回业务错误 1001 → 在区间内调整为 401） 8)面试复盘（filter 与中间件区别、何时用 WithNamedFilter）。
    │   ├── filter/auth.go                       # [NEW] 实现 ServerFilter：从 trpc.GetMetaData(ctx, "auth-token") 取 token，缺失或不匹配返回 errs.New(errs.RetServerAuthFail, "unauthorized")。
    │   ├── filter/timing.go                     # [NEW] 实现 ClientFilter：start := time.Now(); err := next(...); log.Infof("rpc cost=%s", time.Since(start))。
    │   ├── server/main.go                       # [NEW] trpc.NewServer(server.WithFilter(authFilter)) 注册到 UserService。
    │   ├── client/main.go                       # [NEW] proxy := pb.NewUserServiceClientProxy(client.WithFilter(timingFilter), client.WithMetaData("auth-token", "demo"))。
    │   └── trpc_go.yaml                         # [NEW] 同 Day3，新增 server.filter: [auth] 演示配置式注册（与代码式择一）。
    │
    ├── day05-naming-and-client/
    │   ├── README.md                            # [NEW] 章节：1)Selector 四件套（Discovery/ServiceRouter/LoadBalance/CircuitBreaker） 2)WithTarget 寻址协议（ip:// / dns:// / 自定义 scheme） 3)在服务端调用下游（service A 同时是 service B 的客户端） 4)client.WithTimeout vs context.WithTimeout 5)简单自定义 selector 注册（example://）6)PowerShell 命令 7)验证标准（A 收到请求后调 B，把 B 的结果加工返回） 8)面试复盘（为什么 trpc 推崇按 service name 寻址而非 ip）。
    │   ├── proto/user.proto                     # [NEW] 复用 Day2 user.proto，新增 Greet 方法（演示聚合调用）。
    │   ├── stub/                                # [NEW] trpc create 重新生成。
    │   ├── service-a/main.go                    # [NEW] 实现 Greet：内部用 NewUserServiceClientProxy(WithTarget("ip://127.0.0.1:8001")) 调 Day2 的 server，再拼接返回。
    │   ├── service-b/main.go                    # [NEW] 复用 Day2 server 角色，端口 8001。
    │   ├── selector/example.go                  # [NEW] 演示 selector.Register("example", &exampleSelector{})，store 内置两个节点做随机选择。
    │   └── trpc_go.yaml                         # [NEW] 两个 service：a(8002), b(8001)；client 段配 callee=trpc.study.user.UserService target=ip://127.0.0.1:8001 timeout=1000。
    │
    ├── day06-http-restful/
    │   ├── README.md                            # [NEW] 章节：1)tRPC-Go 三种 HTTP 服务对比表（泛 HTTP 标准 / 泛 HTTP RPC / RESTful） 2)google.api.http 注解写法 3)同一 service 多协议暴露（trpc + http） 4)tRPC 错误码到 HTTP 状态码映射表（404/401/429/504...） 5)PowerShell + curl 命令（同时 curl 和 trpc client 调用） 6)验证标准（curl POST /v1/users 与 trpc client CreateUser 等价） 7)面试复盘（RESTful 路由前缀树、何时该选 RESTful vs 泛 HTTP）。
    │   ├── proto/user.proto                     # [NEW] 在 Day2 user.proto 上加 import "google/api/annotations.proto"；CreateUser 加 option (google.api.http) = { post: "/v1/users" body: "*" }；GetUser 加 { get: "/v1/users/{id}" }。
    │   ├── stub/                                # [NEW] trpc create 重新生成（含 RESTful 路由）。
    │   ├── server/main.go                       # [NEW] 注册一次 service 实现，但通过 yaml 暴露 trpc + http 两个 service。
    │   ├── client/
    │   │   ├── trpc_main.go                     # [NEW] 通过 trpc client 调用。
    │   │   └── http_main.go                     # [NEW] 通过 net/http 直接 POST /v1/users。
    │   └── trpc_go.yaml                         # [NEW] server.service: [{name: ..., protocol: trpc, port: 8003}, {name: ..., protocol: http, port: 8004}]。
    │
    ├── day07-test-and-best-practice/
    │   ├── README.md                            # [NEW] 章节：1)测试金字塔（单测 80/集成 15/系统 5） 2)用 trpc create 自动生成的 mock 写单测 3)testify/suite 集成测试模板（SetupSuite/SetupTest/TearDown） 4)trpc-metrics-runtime 接入与可视化指标 5)生产最佳实践 checklist（超时/重试/熔断/日志/错误码/配置中心） 6)与 day01-day07 的对照总结 7)PowerShell 命令（go test ./... -coverprofile=coverage.out） 8)工程级面试复盘清单（10 条）。
    │   ├── service/user_service.go              # [NEW] 抽出可被 mock 的下游 client 接口，演示依赖注入便于测试。
    │   ├── service/user_service_test.go         # [NEW] 用 gomock 生成的 MockUserServiceClientProxy 写表驱动测试，覆盖 CreateUser 成功 / GetUser 不存在 / 下游超时三种用例。
    │   ├── integration/suite_test.go            # [NEW] testify/suite：SetupSuite 启动内嵌 trpc server，TestCreateAndGet 端到端断言，TearDownSuite 优雅关停。
    │   ├── metrics/main.go                      # [NEW] import _ "trpc.tech/trpc-go/trpc-metrics-runtime/v2"；展示 admin /metrics 输出 trpc.runtime.* 指标。
    │   └── trpc_go.yaml                         # [NEW] 含 admin + plugins.log + plugins.metrics 完整配置示例。
    │
    └── (不新增其它根级文件；go.mod/go.sum 在实施阶段由 go get 增量更新)
```

## Key Code Structures

仅给出 Day2 demo 的核心 proto 与服务端 handler 签名（实施时 import 路径以生成出的 stub 为准；接口签名由 trpc-cmdline 决定，规划仅描述意图，不写实现体）：

```
// tRPC/day02-userservice-demo/proto/user.proto
syntax = "proto3";
package trpc.study.user;
option go_package = "study_go/tRPC/day02-userservice-demo/stub/trpc/study/user";

service UserService {
  rpc CreateUser(CreateUserReq) returns (CreateUserRsp);
  rpc GetUser(GetUserReq) returns (GetUserRsp);
}

message User {
  uint64 id = 1;
  string name = 2;
  int64  created_at = 3;
}
message CreateUserReq { string name = 1; }
message CreateUserRsp { User user = 1; }
message GetUserReq    { uint64 id = 1; }
message GetUserRsp    { User user = 1; }
```

```
// 服务端 handler 意图（实际方法签名以 trpc create 生成的 UserServiceService 接口为准）
type userImpl struct {
    mu    sync.RWMutex
    seq   uint64
    store map[uint64]*pb.User
}
// CreateUser(ctx, *pb.CreateUserReq) (*pb.CreateUserRsp, error)
// GetUser   (ctx, *pb.GetUserReq)    (*pb.GetUserRsp, error)   // 不存在时 return nil, errs.New(404, "user not found")
```

## Agent Extensions

### SubAgent

- **code-explorer**
- Purpose: 在 build 阶段批量探查 trpc-cmdline 实际生成的 stub 目录结构、`pb.RegisterUserServiceService` / `NewUserServiceClientProxy` 的精确签名，以及 `errs`、`server.WithFilter`、`client.WithTarget` 在用户最终选定的 v2 / 开源版包路径下的真实导出符号。
- Expected outcome: 在每个 day 的 `server/main.go` / `client/main.go` 落地前，确认 import 路径、函数签名、yaml 字段名 100% 与所选版本一致，避免 v1/v2/开源版混用导致编译失败。译失败。