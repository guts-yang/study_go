# tRPC-Go 7 天工程级学习规划

这是 [`study_go`](../README.md) 项目的进阶篇。在你完成了 `day01-basics` 到 `day07-web-server` 这 7 天 Go 基础训练（语法、并发、`net/http`、`context`、`sync`、`go test`）之后，这一份规划带你从 0 到 1 熟悉 **tRPC-Go**，并在 Day2 用它重写 day07 的用户服务，建立"从 HTTP 框架走向 RPC 框架"的工程直觉。

> tRPC 是腾讯开源的多语言 RPC 框架，tRPC-Go 是它的 Go 实现。你可以把它理解为：在 `net/http` 的并发模型之上，叠加了 IDL（proto）契约、服务发现、拦截器、监控、配置、多协议（tRPC/HTTP/HTTP2/gRPC）等一整套工程能力。

## 7 天路线

| 天数 | 目录 | 主题 | 核心目标 |
| --- | --- | --- | --- |
| Day 1 | [`day01-helloworld`](./day01-helloworld) | 环境准备 & helloworld | 装好 `trpc-cmdline`，跑通官方 Greeter |
| Day 2 | [`day02-userservice-demo`](./day02-userservice-demo) | **核心 demo：UserService 端到端** | 用 tRPC 重写 day07 的用户服务，跑通 server+client |
| Day 3 | [`day03-config-and-admin`](./day03-config-and-admin) | 配置文件 & admin 端口 & log | 看懂 `trpc_go.yaml` 三段，会用框架 log 与 admin |
| Day 4 | [`day04-filter-and-errs`](./day04-filter-and-errs) | 拦截器 & 错误码 | 写一个鉴权 filter，理解 `errs` 区间约定 |
| Day 5 | [`day05-naming-and-client`](./day05-naming-and-client) | 名字服务 & 下游调用 | 服务端调用下游服务，理解 `Selector` 四件套 |
| Day 6 | [`day06-http-restful`](./day06-http-restful) | HTTP/RESTful 多协议 | 同一 service 同时暴露 tRPC + HTTP/RESTful |
| Day 7 | [`day07-test-and-best-practice`](./day07-test-and-best-practice) | 单测 & 监控 & 最佳实践 | 用桩代码 mock 写单测，接入 metrics，做生产级复盘 |

## 环境准备（Windows + PowerShell）

### 1. Go 1.22 及以上

沿用 `study_go` 根目录的 `go.mod`：

```powershell
go version
# 期望：go version go1.22.x windows/amd64
```

### 2. Protocol Buffers 编译器

```powershell
# 下载 https://github.com/protocolbuffers/protobuf/releases 的 protoc-*-win64.zip
# 解压后把 bin/protoc.exe 放到 $env:Path 里
protoc --version
# 期望：libprotoc 3.x 或 4.x
```

### 3. trpc-cmdline（脚手架）

```powershell
# 推荐：开源版（外网可访问，零内网依赖）
go install trpc.group/trpc-go/trpc-cmdline/trpc@latest

# 内网（腾讯）版：
# go install trpc.tech/trpc-go/trpc-cmdline/trpc@latest

# 安装后让命令可用
$env:Path += ";$(go env GOPATH)\bin"
trpc version
```

如果 `trpc version` 报命令找不到，把 `$(go env GOPATH)\bin` 永久加到 PATH：

```powershell
[Environment]::SetEnvironmentVariable(
    "Path",
    [Environment]::GetEnvironmentVariable("Path", "User") + ";$(go env GOPATH)\bin",
    "User"
)
```

### 4. GOPROXY（强烈建议）

```powershell
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=sum.golang.org
```

### 5. tRPC-Go 版本选型

|     | 包路径 | 适用场景 |
| --- | --- | --- |
| 开源版 | `trpc.group/trpc-go/trpc-go`（推荐） | 个人学习、外网环境、本仓库默认 |
| 内网 v2 | `trpc.tech/trpc-go/trpc-go/v2` | 腾讯内网，正在演进的主线 |
| 内网 v1（LTS） | `git.woa.com/trpc-go/trpc-go` | 仅历史项目维护用，新项目不建议 |

> 本规划所有示例默认使用 **开源版** `trpc.group/trpc-go/trpc-go`。如果你在腾讯内网，按上表把 import 路径整体替换为 v2 即可，API 是兼容的。

## 学习方式

每天按这个节奏（与 `study_go` 根 README 一致）：

1. 先读当天目录的 `README.md`，看懂主题和知识点；
2. 再读 `*.proto`、`*.go`、`trpc_go.yaml`，把代码当讲义；
3. 按 README 里的命令跑起来，观察输出；
4. 改一改：换字段、加方法、改端口、加 filter；
5. 用自己的话回答当天的"面试复盘"问题。

> tRPC 的学习曲线不在 API 多复杂，而在"它把多少工程问题前置到了配置和约定里"。多读 yaml 比多读代码更重要。

## 常用命令

```powershell
# 进入某天目录
cd .\tRPC\day01-helloworld

# 第一次跑（每个 day 目录都是独立 module）
go mod tidy

# 生成桩代码（具体命令在每个 day 的 README）
trpc create -p .\proto\helloworld.proto -o . --rpconly

# 启动服务端（PowerShell 第一个窗口）
go run .\server\

# 调用客户端（PowerShell 第二个窗口）
go run .\client\

# 跑测试
go test ./...
```

## 工程级面试复盘清单

学完 7 天后，你应当能用自己的话清晰回答以下问题：

- tRPC-Go 与 `net/http` 在请求处理生命周期上有哪些根本差异？
- 为什么 RPC 框架普遍使用 IDL（proto）？字段编号修改会带来什么后果？
- `trpc_go.yaml` 的 `global / server / client / plugins` 四段分别解决什么问题？
- 服务端 filter 的执行顺序是怎样的？请求阶段与响应阶段为何方向相反？
- 客户端 filter 链中"选择器过滤器"在哪个位置？为什么必须在最后？
- `errs.New(code, msg)` 与 `errors.New(msg)` 的区别？错误码区间为什么要约定？
- Selector 的四个组件 `Discovery / ServiceRouter / LoadBalance / CircuitBreaker` 各自负责什么？
- `client.WithTarget("ip://...")` 与 `client.WithTarget("polaris://...")` 的差异点在哪？
- `client.WithTimeout` 与 `context.WithTimeout` 同时存在时谁优先？为什么级联超时要用 context？
- tRPC 的 `protocol: trpc` 私有协议比 HTTP/JSON 强在哪？什么场景反而该用 HTTP？
- RESTful 服务在框架内部怎么做 HTTP-RPC 互转？前缀树路由解决了什么？
- 为什么不能用 `fmt.Printf` 打日志？框架 `log` 接口背后做了什么？
- admin 端口默认提供哪些接口？为什么生产环境要把它绑到内网 IP？
- `trpc create` 默认生成的 mock 文件怎么用 `gomock` 写单测？
- `trpc-metrics-runtime` 上报了哪些指标？怎么在 admin 的 `/metrics` 看到？
- 一个生产级 tRPC 服务上线前，配置层面至少要 review 哪 8 项？

## 推荐扩展资料

- 官方源码与示例：<https://github.com/trpc-group/trpc-go>（外网）/ <https://git.woa.com/trpc-go/trpc-go>（内网）
- examples/helloworld：最简 RPC 入门
- examples/features：70+ 主题示例（filter、selector、stream、restful、metric、tracing 等）
- docs/quick_start.zh_CN.md：官方快速开始
- docs/user_guide：分主题深度文档
