# Day 3 · 配置文件 & admin 端口 & 标准日志

**主题**：吃透 `trpc_go.yaml` 的四段结构（global / server / client / plugins），学会用框架 `log` 接口和 admin 端口。

**核心目标**：

- 看到一份生产级 yaml，能逐段解释每个字段；
- 知道为什么不能用 `fmt.Printf`；
- 会用 admin 端口做健康检查、查 metrics、查路由表、热重载日志级别。

## 1. 工程目录

本 day 复用 [Day2](../day02-userservice-demo) 的 `user.proto` 与 stub，只换 `server/main.go` 与 `trpc_go.yaml`，演示**同一份业务逻辑、不同配置**的玩法。

```
day03-config-and-admin/
├── README.md
├── go.mod
├── server/
│   └── main.go        # 复用 day02 service.go 的逻辑（重新贴一份，避免跨 module 引用）
├── server/service.go  # 拷贝自 day02
└── trpc_go.yaml       # 重点：admin / log / 多 service 配置
```

## 2. yaml 四段总览

```yaml
global:    {...}   # 全局环境信息
server:    {...}   # 服务端配置：监听哪些端口、用什么协议
client:    {...}   # 客户端默认配置：调下游时的超时/寻址默认值
plugins:   {...}   # 所有插件的配置入口（log/metrics/config/selector...）
```

> 框架启动顺序：先解析 yaml → init plugins → 注册 server.service → Serve。**所以如果某个插件配置错误，server 根本起不来**。

### 2.1 `global` 段

```yaml
global:
  namespace: Production       # 环境隔离（监控、配置、寻址都按 namespace 拆）
  env_name: prod-shenzhen     # 子环境（同 namespace 下的细分）
  container_name: user-svc-1  # 容器名（k8s 场景常见）
  enable_set: N               # set 化部署开关
  full_set_name: ""           # set 名（一般由部署平台注入）
```

经验：**namespace 决定一切的隔离边界**。开发机用 `Development`，CI 用 `Testing`，线上用 `Production`，绝不要混。

### 2.2 `server` 段

```yaml
server:
  app: study                     # 应用名
  server: user                   # 服务进程名
  bin_path: /usr/local/trpc/bin
  conf_path: /usr/local/trpc/conf
  data_path: /usr/local/trpc/data
  filter:                        # 全局 server filter（按顺序执行）
    - recovery                   # 推荐第一个：捕获 panic
  admin:                         # admin 端口
    ip: 127.0.0.1
    port: 11014
    read_timeout: 3000
    write_timeout: 60000
    enable_tls: false
  service:                       # 业务 service 列表（一个进程可以挂多个）
    - name: trpc.study.user.UserService
      ip: 127.0.0.1
      port: 8001
      network: tcp
      protocol: trpc
      timeout: 1000
      filter: [auth]             # 仅本 service 生效的 filter
```

要点：

- `filter` 出现在两层：`server.filter` 是全局，`server.service[].filter` 是单 service 的；它们**按数组顺序执行**；
- `admin` 是个**独立的 HTTP 监听**，专门暴露管理接口（健康检查、metrics、pprof、自定义命令）。生产环境务必绑内网 IP（如 `127.0.0.1` 或 sidecar 网卡），**不要暴露公网**。

### 2.3 `client` 段

```yaml
client:
  timeout: 1000                   # 全局默认超时
  namespace: Development          # 调下游时使用的 namespace
  filter: [degrade]               # 全局 client filter
  service:                        # 按 callee 服务名给出独立配置
    - callee: trpc.study.user.UserService
      target: ip://127.0.0.1:8001
      network: tcp
      protocol: trpc
      timeout: 800
      serialization: 0            # 0=PB 1=JSON 2=FlatBuffers
```

要点：业务代码里 `pb.NewUserServiceClientProxy()` 不需要传 `WithTarget`，框架会按 callee 名从 yaml 里查。**这是 day05 的核心套路**。

### 2.4 `plugins` 段

```yaml
plugins:
  log:                            # 日志插件
    default:                      # 默认 logger，业务 log.Infof 走这里
      - writer: console
        level: debug
      - writer: file
        level: info
        formatter: json
        writer_config:
          filename: ../log/trpc.log
          max_size: 100           # MB
          max_backups: 10
          max_age: 7              # 天
          compress: false
  metrics:                        # 监控插件，由 day07 接入 trpc-metrics-runtime
    prometheus: { ip: 0.0.0.0, port: 12017, path: /metrics }
  selector:                       # 名字服务插件，day05 详解
    polaris: { ... }
  config:                         # 配置中心插件
    rainbow: { ... }
```

## 3. 框架日志 API

❌ **错的写法**：

```go
fmt.Printf("create user %d\n", id)   // 不会进文件，不会上报远程日志中心
```

✅ **对的写法**：

```go
import "trpc.group/trpc-go/trpc-go/log"

log.Infof("create user id=%d", id)                              // 默认 logger
log.InfoContextf(ctx, "create user id=%d", id)                  // 带 context（推荐）
log.WithContextFields(ctx, "uid", "123").Infof("create user")   // 结构化字段
```

为什么：

1. 框架 `log` 接口背后挂了多个 `writer`（console / file / 远程上报），`fmt.Printf` 只走 stdout；
2. 带 `Context` 的版本会自动注入 trace_id / span_id，调用链才能串起来；
3. 业务代码越多，"统一日志规范"的边际收益越大 —— 排障效率指数级提升。

## 4. admin 端口

启用 admin 后，框架会暴露这些 HTTP 接口（默认）：

| 路径 | 作用 |
| --- | --- |
| `GET /cmds` | 列出所有已注册的命令 |
| `GET /cmds/loglevel?logger=default&level=debug` | **运行时改日志级别**（不重启进程） |
| `GET /cmds/config` | dump 当前生效的 yaml |
| `GET /metrics` | Prometheus 格式指标 |
| `GET /debug/pprof/...` | Go pprof 全套（heap、goroutine、profile） |
| `GET /healthz` | 健康检查（k8s liveness 用） |

可以在代码里挂自定义命令：

```go
import "trpc.group/trpc-go/trpc-go/admin"

admin.HandleFunc("/whoami", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "user-service v1.0")
})
```

## 5. 配置热加载

框架内置 `fsnotify` 监听 yaml：

- 改 `plugins.log.default[0].level` 从 `debug` → `info`，**保存**即时生效；
- 改 `server.service[].timeout` 通常需要重启（取决于具体 server 实现）；
- 插件级配置（如 metrics 上报地址）大多支持热生效。

> 别滥用热加载。生产环境推荐配合配置中心 + 灰度，不要直接 `vim` 线上 yaml。

## 6. 跑起来

```powershell
cd .\tRPC\day03-config-and-admin
trpc create -p ..\day02-userservice-demo\proto\user.proto -o . --rpconly
go mod tidy
go run .\server\
```

**验证 admin 端口**（另开窗口）：

```powershell
curl http://127.0.0.1:11014/cmds
curl http://127.0.0.1:11014/cmds/loglevel?logger=default
curl http://127.0.0.1:11014/whoami     # 自定义命令
```

**验证热加载**：

1. 服务运行中，编辑 `trpc_go.yaml`，把 `plugins.log.default[0].level` 改成 `error`，保存；
2. 重新 `go run .\client\`（用 day02 的 client），观察服务端日志没有新的 INFO 输出 → 热加载成功。

## 7. 验证标准

- [ ] `curl http://127.0.0.1:11014/cmds` 返回 JSON 命令列表；
- [ ] `curl http://127.0.0.1:11014/whoami` 输出 `user-service v1.0`；
- [ ] 热改 `level: debug → error` 后服务端日志被压制；
- [ ] 把代码里 `log.Infof` 换成 `fmt.Printf`，观察日志没有 trace_id（理解为什么不该这么写）。

## 8. 面试复盘

1. **`server.filter` 与 `server.service[].filter` 区别？** 前者作用于该 server 进程的所有 service；后者只作用于配置它的那一个 service。生产里 recovery / metrics / tracing 通常放全局。
2. **admin 端口为什么要单独开一个端口而不是和业务复用？** 隔离故障域和网络策略 —— 业务端口对外，admin 端口对内；业务端口被打满时 admin 仍可访问做诊断。
3. **`fmt.Printf` 和 `log.Infof` 在生产环境里的差别有多大？** 不可同日而语：丢失 trace_id（不能串调用链）、丢失级别控制、丢失结构化字段、不上报远程日志中心、不参与日志切割归档。
4. **热加载可以热改业务代码吗？** 不能。yaml 热加载只重新初始化插件配置，业务代码必须重启进程才能生效。生产里"零停机"靠的是滚动发布，不是热加载。
5. **如果 admin 端口意外暴露公网，最大的风险是什么？** pprof 接口可以被用来做 DoS（大量 profile 拉数据）、`/cmds/config` 会泄露内部配置（含密钥）、`/cmds/loglevel` 可以被攻击者改成 debug 来拖慢服务。
