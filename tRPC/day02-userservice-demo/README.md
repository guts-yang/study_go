# Day 2 · UserService 端到端 Demo

**主题**：用 tRPC-Go 重写 [`day07-web-server`](../../day07-web-server) 的用户服务（CreateUser / GetUser），跑通 server + client 全链路。这是整个 7 天里**最核心的一天**，后面所有进阶专题都会回到这个 demo 上叠加。

**核心目标**：

- 能独立设计一个 proto 文件并解释每个字段的含义；
- 能用 `trpc create` 生成桩代码并读懂 `*.trpc.go` 的接口；
- 能从 day07 的 `net/http` handler 平滑迁移业务逻辑（map + RWMutex 这部分原样保留）；
- 跑通 client → server → client 的双向调用，并能识别框架打印的关键日志。

## 1. 需求与 day07 对照

| 维度 | day07 net/http | day02 tRPC-Go |
| --- | --- | --- |
| 接口契约 | `POST /users` `{"name":"Alice"}` | `service UserService { rpc CreateUser(CreateUserReq) returns (CreateUserRsp); }` |
| 返回结构 | `User{id,name,created_at}` 直接 JSON | `message User { uint64 id=1; string name=2; int64 created_at=3; }` |
| 路由 | `mux.HandleFunc` 字符串拼路径 | proto 自动派生 `/trpc.study.user.UserService/CreateUser` |
| 错误 | `writeJSON(w, 404, ...)` | `return nil, errs.New(404, "user not found")` |
| 并发安全 | `sync.RWMutex + map[int]User` | **完全相同**：`sync.RWMutex + map[uint64]*pb.User` |
| 启动 | `http.Server.ListenAndServe()` | `trpc.NewServer().Serve()` |

> 🔑 **核心心智**：tRPC 没有改变并发安全的责任划分。"handler 在多 goroutine 中被并发调用"这条 Go 后端铁律仍然成立，存储层的锁照样要写。

## 2. 工程目录

```
day02-userservice-demo/
├── README.md
├── go.mod
├── proto/
│   └── user.proto
├── stub/                 # 由 trpc create 生成
├── server/
│   ├── main.go           # 注册入口
│   └── service.go        # UserService 业务实现（map + RWMutex）
├── client/
│   └── main.go           # 顺序调用 CreateUser → GetUser
└── trpc_go.yaml
```

## 3. proto 设计要点

```protobuf
syntax = "proto3";

package trpc.study.user;

option go_package = "day02-userservice-demo/stub/trpc/study/user";

service UserService {
  rpc CreateUser(CreateUserReq) returns (CreateUserRsp);
  rpc GetUser   (GetUserReq)    returns (GetUserRsp);
}

message User {
  uint64 id         = 1;
  string name       = 2;
  int64  created_at = 3;  // unix 秒；proto 没有 time.Time 类型
}

message CreateUserReq { string name = 1; }
message CreateUserRsp { User   user = 1; }

message GetUserReq    { uint64 id   = 1; }
message GetUserRsp    { User   user = 1; }
```

设计要点：

- **包名分层** `trpc.study.user`：`trpc.<app>.<module>` 是社区惯例，方便监控按层级聚合；
- **字段编号永不重用**：哪怕字段被删除，编号也要 `reserved`；
- **没有 time.Time**：proto3 内置时间需要 `google.protobuf.Timestamp`，简单场景用 `int64` unix 秒更直观；
- **Req/Rsp 包一层而不是裸 User**：以后加分页、加 trace_id、加扩展字段都不会破坏接口；
- 方法返回 `(Rsp, error)`：成功时 `error == nil`；失败时**只返回 error**（rsp 设为 nil），由 `errs` 统一表达。

## 4. 生成桩代码

```powershell
cd .\tRPC\day02-userservice-demo
trpc create -p .\proto\user.proto -o . --rpconly
go mod tidy
```

生成产物（关注接口签名）：

```go
// stub/trpc/study/user/user.trpc.go 中会有类似定义：

type UserServiceService interface {
    CreateUser(ctx context.Context, req *CreateUserReq) (*CreateUserRsp, error)
    GetUser   (ctx context.Context, req *GetUserReq)    (*GetUserRsp,    error)
}

func RegisterUserServiceService(s server.Service, svr UserServiceService) { ... }

type UserServiceClientProxy interface { /* 客户端方法签名相同 */ }
func NewUserServiceClientProxy(opts ...client.Option) UserServiceClientProxy { ... }
```

**只需要做两件事**：让你的 struct 实现 `UserServiceService` 接口；用 `RegisterUserServiceService` 注册。

## 5. 服务端实现

`server/service.go` 把 day07 的存储层"原样"搬过来，只是 map 的 value 换成了 `*pb.User`：

```go
type userImpl struct {
    mu    sync.RWMutex
    seq   uint64
    store map[uint64]*pb.User
}

func newUserImpl() *userImpl {
    return &userImpl{store: make(map[uint64]*pb.User)}
}

func (u *userImpl) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.CreateUserRsp, error) {
    if strings.TrimSpace(req.GetName()) == "" {
        // errs.New 比 errors.New 多一个错误码字段；client 端通过 errs.Code(err) 拿到
        return nil, errs.New(400, "name is required")
    }

    u.mu.Lock()
    defer u.mu.Unlock()

    u.seq++
    user := &pb.User{
        Id:        u.seq,
        Name:      req.GetName(),
        CreatedAt: time.Now().Unix(),
    }
    u.store[user.Id] = user

    log.Infof("CreateUser ok id=%d name=%s", user.Id, user.Name)
    return &pb.CreateUserRsp{User: user}, nil
}

func (u *userImpl) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserRsp, error) {
    u.mu.RLock()
    defer u.mu.RUnlock()

    user, ok := u.store[req.GetId()]
    if !ok {
        return nil, errs.New(404, "user not found")
    }
    return &pb.GetUserRsp{User: user}, nil
}
```

`server/main.go` 只做注册和启动：

```go
func main() {
    s := trpc.NewServer()
    pb.RegisterUserServiceService(s, newUserImpl())
    if err := s.Serve(); err != nil {
        log.Fatal(err)
    }
}
```

## 6. 客户端调用

`client/main.go` 顺序演示一次"建用户 → 取用户"，并故意访问一个不存在的 id 来观察错误码：

```go
func main() {
    proxy := pb.NewUserServiceClientProxy(
        client.WithTarget("ip://127.0.0.1:8001"),
        client.WithTimeout(time.Second),
    )

    ctx := context.Background()

    crsp, err := proxy.CreateUser(ctx, &pb.CreateUserReq{Name: "Alice"})
    if err != nil {
        panic(err)
    }
    fmt.Printf("CreateUser ok: id=%d name=%s\n", crsp.GetUser().GetId(), crsp.GetUser().GetName())

    grsp, err := proxy.GetUser(ctx, &pb.GetUserReq{Id: crsp.GetUser().GetId()})
    if err != nil {
        panic(err)
    }
    fmt.Printf("GetUser ok:    id=%d name=%s\n", grsp.GetUser().GetId(), grsp.GetUser().GetName())

    // 故意访问不存在的 id，演示 errs 错误码透传
    _, err = proxy.GetUser(ctx, &pb.GetUserReq{Id: 9999})
    fmt.Printf("GetUser 9999 -> code=%d msg=%s\n", errs.Code(err), errs.Msg(err))
}
```

要点：

- `errs.Code(err)` / `errs.Msg(err)` 是从 server 端 `errs.New(code, msg)` 反序列化回来的，**跨进程保留了语义**；
- 注意 day01 用的是 8000 端口、day02 起改用 8001（避免端口冲突）。

## 7. 配置文件

```yaml
global:
  namespace: Development

server:
  app: study
  server: user
  service:
    - name: trpc.study.user.UserService
      ip: 127.0.0.1
      port: 8001
      network: tcp
      protocol: trpc
      timeout: 1000

plugins:
  log:
    default:
      - writer: console
        level: debug
```

## 8. 跑起来

**窗口 1**：

```powershell
cd .\tRPC\day02-userservice-demo
trpc create -p .\proto\user.proto -o . --rpconly
go mod tidy
go run .\server\
```

**窗口 2**：

```powershell
cd .\tRPC\day02-userservice-demo
go run .\client\
```

期望输出：

```
CreateUser ok: id=1 name=Alice
GetUser ok:    id=1 name=Alice
GetUser 9999 -> code=404 msg=user not found
```

## 9. 验证标准

- [ ] `stub/trpc/study/user/` 下能看到 `user.pb.go`、`user.trpc.go`、`user_mock.go`；
- [ ] server 启动日志含 `service:trpc.study.user.UserService launch success`；
- [ ] client 三行输出全部命中预期；
- [ ] 把 `req.Name` 改成空字符串重跑，观察到 `code=400 msg=name is required`；
- [ ] 改完 service 实现后**不需要重启 client**（因为 client 是临时进程），但**需要重启 server**。

## 10. 面试复盘

1. **proto 字段编号为什么不能改？** 字段编号写进二进制 wire format，反序列化按编号匹配。改编号 = 字段语义错位，会无声破坏所有线上消费者。
2. **Proxy 为什么是并发安全的？** Proxy 内部维护连接池、filter 链是只读的，每次 RPC 用 `context` 携带请求级状态，不存在共享可变状态。
3. **从 day07 迁移过来的最大代价是什么？** 不是代码量，而是**契约方式的改变**：HTTP 时代靠"约定 + 文档"，RPC 时代靠"proto 文件 + 桩代码"。一旦上线，proto 就成为消费者依赖的二进制协议，演进必须严格守"只加不删不改编号"。
4. **`errs.New(404, "user not found")` 中的 404 是 HTTP 状态码吗？** **不是**。tRPC 的 errs code 是框架自定义的整数空间。HTTP 协议下确实有"业务错误码 → HTTP 状态码"的映射表（day06 会展开），但本质上 errs code 是跨协议的语义错误码。
5. **如果要支持"按 name 查询用户"，应该新增一个 RPC 还是给 GetUserReq 加字段？** 优先**新增 RPC**（如 `QueryUserByName`），原因：(1) 单一职责；(2) 监控分桶清晰；(3) 后续加缓存、限流策略可独立配。给现有 Req 加可选字段是一种"协议腐化"，要克制。
