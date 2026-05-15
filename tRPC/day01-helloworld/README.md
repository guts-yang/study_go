# Day 1 · Hello, tRPC-Go

**主题**：环境验证 + 跑通官方风格的 Greeter 服务，建立"proto → 桩代码 → server + client"的整体心智。

**核心目标**：

- 能用一句话说清 tRPC-Go 与 `net/http` 的差异；
- 会用 `trpc-cmdline` 生成桩代码，并解释生成出来的目录结构；
- 在本机 PowerShell 双窗口里跑通 server / client 通信。

## 1. tRPC-Go 是什么

tRPC-Go 是一个**带服务治理能力的 RPC 框架**，对比 day07 的 `net/http`：

| 维度 | `net/http`（day07） | tRPC-Go |
| --- | --- | --- |
| 接口契约 | 自定义 JSON + 路由字符串 | proto3 IDL，强类型、跨语言 |
| 路由 | `mux.HandleFunc("GET /users/{id}", ...)` | 框架根据 `service.method` 自动路由 |
| 序列化 | 你自己 `json.Marshal` | 框架按 `protocol` 选 PB/JSON/FlatBuffers |
| 治理 | 你自己写中间件、超时、重试 | filter / selector / 配置 / 监控插件化 |
| 客户端 | `http.Client` + URL 拼接 | `pb.NewGreeterClientProxy()`，像调用本地函数 |
| 多协议 | 只有 HTTP | tRPC 私有协议 / HTTP / RESTful / gRPC |

> 一句话：**`net/http` 给你"通信原语"，tRPC-Go 给你"工程脚手架"。**

## 2. 工程目录

```
day01-helloworld/
├── README.md           # 本文档
├── go.mod              # 独立 module（学习用，不污染顶层 study_go）
├── proto/
│   └── helloworld.proto
├── stub/               # 由 `trpc create` 生成
│   └── ...             # helloworld.pb.go / helloworld.trpc.go / helloworld_mock.go
├── server/
│   └── main.go
├── client/
│   └── main.go
└── trpc_go.yaml        # 服务端配置
```

> 每个 day 目录都是一个**独立可运行的迷你工程**，第一次进入需要 `go mod tidy`。这样设计是为了让你专注当天主题、不被旁边目录的依赖干扰。

## 3. proto 文件解读

```protobuf
syntax = "proto3";

package trpc.helloworld;

option go_package = "day01-helloworld/stub/trpc/helloworld";

service Greeter {
  rpc SayHello(HelloRequest) returns (HelloReply);
}

message HelloRequest {
  string msg = 1;
}

message HelloReply {
  string msg = 1;
}
```

要点：

- `package trpc.helloworld` —— **proto 包名**，会成为 RPC 路径前缀（最终是 `/trpc.helloworld.Greeter/SayHello`）；
- `option go_package` —— Go 代码生成路径；
- 字段编号 `= 1` 一旦上线**不可修改、不可重用**，老代码会按这个编号反序列化；
- service 里的方法对应到生成的 Go interface 方法。

## 4. 生成桩代码

### 4.1 安装 trpc-cmdline（一次性）

按 [`tRPC/README.md`](../README.md#3-trpc-cmdline脚手架) 安装，验证：

```powershell
trpc version
```

### 4.2 生成

```powershell
# 切到本目录
cd .\tRPC\day01-helloworld

# --rpconly 表示只生成桩代码，不生成 main.go / yaml 模板（我们已经手写了）
trpc create -p .\proto\helloworld.proto -o . --rpconly
```

生成产物（关注三类文件）：

| 文件 | 作用 |
| --- | --- |
| `helloworld.pb.go` | proto 消息体的 Go 结构（`HelloRequest` / `HelloReply`） |
| `helloworld.trpc.go` | RPC 接口与代理（`GreeterService` 接口、`RegisterGreeterService`、`NewGreeterClientProxy`） |
| `helloworld_mock.go` | gomock 自动生成的 mock，用于单测（Day 7 会用） |

### 4.3 拉依赖

```powershell
go mod tidy
```

## 5. 服务端代码（`server/main.go`）

服务端做三件事：实现接口、注册到 server、启动服务。

```go
package main

import (
    "context"

    pb "day01-helloworld/stub/trpc/helloworld"

    trpc "trpc.group/trpc-go/trpc-go"
    "trpc.group/trpc-go/trpc-go/log"
)

// greeterImpl 实现 pb.GreeterService 接口
type greeterImpl struct{}

func (g *greeterImpl) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
    log.Infof("收到请求: %s", req.GetMsg())
    return &pb.HelloReply{Msg: "Hello, " + req.GetMsg()}, nil
}

func main() {
    s := trpc.NewServer()                       // 读取 trpc_go.yaml
    pb.RegisterGreeterService(s, &greeterImpl{})
    if err := s.Serve(); err != nil {           // 阻塞，监听信号优雅退出
        log.Fatal(err)
    }
}
```

> ⚠️ 上面 `import` 中的 `day01-helloworld/stub/...` 路径取决于你本目录 `go.mod` 的 module 名以及 proto 中 `option go_package` 的写法。生成完桩代码后，看 `helloworld.pb.go` 顶部的 `package` 与导入路径，对齐即可。

## 6. 客户端代码（`client/main.go`）

```go
package main

import (
    "context"
    "fmt"
    "time"

    pb "day01-helloworld/stub/trpc/helloworld"

    _ "trpc.group/trpc-go/trpc-go" // 触发框架 init
    "trpc.group/trpc-go/trpc-go/client"
)

func main() {
    proxy := pb.NewGreeterClientProxy(
        client.WithTarget("ip://127.0.0.1:8000"), // 直连寻址
        client.WithTimeout(time.Second),
    )

    rsp, err := proxy.SayHello(context.Background(), &pb.HelloRequest{Msg: "world"})
    if err != nil {
        panic(err)
    }
    fmt.Println(rsp.GetMsg()) // 期望：Hello, world
}
```

要点：

- `client.WithTarget("ip://127.0.0.1:8000")`：协议头 `ip://` 表示直连；后面会接触 `polaris://`、`dns://` 等服务发现协议；
- 客户端**不需要 yaml**，所有配置都可以通过 `client.With*` Options 给出；
- Proxy 是并发安全的，业务代码里通常作为 package 级单例。

## 7. 服务端配置（`trpc_go.yaml`）

```yaml
global:
  namespace: Development

server:
  app: study
  server: helloworld
  service:
    - name: trpc.helloworld.Greeter   # 必须与 proto package + service 一致
      ip: 127.0.0.1
      port: 8000
      network: tcp
      protocol: trpc                  # 默认私有协议；可换 http
      timeout: 1000                   # 单次请求最大处理时间，毫秒

plugins:
  log:
    default:
      - writer: console
        level: debug
```

要点：

- `service[].name` 是 RPC 路由依据，**必须**和 proto 的 `package + service` 拼起来一致；
- `protocol: trpc` 是 tRPC 私有协议（二进制、高性能），后续 day06 会换成 http；
- `plugins` 段是后续所有插件的入口（log、metrics、config、selector...）。

## 8. 跑起来（PowerShell 双窗口）

**窗口 1：启动服务端**

```powershell
cd .\tRPC\day01-helloworld
go mod tidy
go run .\server\
```

期望日志：

```
... service:trpc.helloworld.Greeter launch success, tcp 127.0.0.1:8000 ...
```

**窗口 2：跑客户端**

```powershell
cd .\tRPC\day01-helloworld
go run .\client\
```

期望输出：

```
Hello, world
```

服务端窗口同时会打印 `收到请求: world`。

## 9. 验证标准

- [ ] `trpc version` 能正常输出版本号；
- [ ] `trpc create` 生成的 stub 目录有 `*.pb.go`、`*.trpc.go`、`*_mock.go` 三类文件；
- [ ] server 启动日志包含 `launch success`；
- [ ] client 控制台输出 `Hello, world`；
- [ ] 把 client 的 `Msg` 改成中文 `你好`，重跑能得到 `Hello, 你好`。

## 10. 面试复盘

1. **proto 字段编号能不能改？** 不能。线上消费者按编号反序列化，改了等于换字段。可以新加字段（用新编号）、可以标记 `reserved`，绝不能改。
2. **`pb.NewGreeterClientProxy()` 返回的 Proxy 是否并发安全？** 是的，内部连接池 + filter 链都是无锁/读多写少结构，实践上常作为包级单例。
3. **`trpc.NewServer()` 没显式传配置文件路径，它怎么知道用 `trpc_go.yaml`？** 框架默认从启动目录寻找 `trpc_go.yaml`，可用 `-conf` 命令行参数指定路径。
4. **`service.name` 为什么必须和 proto 包名拼一致？** 因为 tRPC 私有协议帧里把"完整方法名" `/trpc.helloworld.Greeter/SayHello` 作为路由 key，server 端依据它分发到注册好的 service。
5. **`ip://127.0.0.1:8000` 这种 target 在生产里能用吗？** 一般不能。生产里要用名字服务（北极星、DNS、k8s service），day05 会展开。
