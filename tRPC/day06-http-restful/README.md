# Day 6 · HTTP & RESTful 多协议

**主题**：让同一份业务实现，**同时**通过 tRPC 私有协议和 HTTP/RESTful 对外提供服务。理解 tRPC-Go 的"同业务、多协议"工程范式。

**核心目标**：

- 能说清 tRPC-Go 三种 HTTP 服务的定位差异；
- 会在 proto 里写 `google.api.http` 注解；
- 会用一份 yaml 同时挂出 `protocol: trpc` 和 `protocol: restful` 两个 service；
- 知道 tRPC 错误码到 HTTP 状态码的映射表，能用 curl 直接调 tRPC 服务。

## 1. tRPC-Go 三种 HTTP 服务

| 类型 | 是否需要 proto | 路径风格 | 适用场景 |
| --- | --- | --- | --- |
| **泛 HTTP 标准服务** | ❌ 不需要 | 你自己定，像 `net/http` 一样 | 文件上传、静态资源、与已有 HTTP 系统集成 |
| **泛 HTTP RPC 服务** | ✅ 共享 RPC 桩代码 | `/{service}/{method}`，POST + JSON body | 给前端/移动端用 RPC 同款契约，但走 HTTP |
| **泛 HTTP RESTful 服务** | ✅ proto + `google.api.http` 注解 | RESTful 风格 `/v1/users/{id}` | 标准 REST API，前端最舒服 |

**核心 demo 用第三种 RESTful**，最贴近后端面试常考的"如何把 RPC 服务暴露成 REST API"。

## 2. proto 注解写法

```protobuf
syntax = "proto3";

package trpc.study.user;
option go_package = "day06-http-restful/stub/trpc/study/user";

import "google/api/annotations.proto";

service UserService {
  rpc CreateUser(CreateUserReq) returns (CreateUserRsp) {
    option (google.api.http) = {
      post: "/v1/users"
      body: "*"                     // 整个 req 作为 JSON body
    };
  }

  rpc GetUser(GetUserReq) returns (GetUserRsp) {
    option (google.api.http) = {
      get: "/v1/users/{id}"         // path 参数自动绑定到 req.id
    };
  }
}
```

注解语义：

- `post: "/v1/users"` → 这个 RPC 方法对外的 HTTP 是 `POST /v1/users`；
- `body: "*"` → req 整体作为 JSON body；可以写 `body: "user"` 表示只把 `req.user` 作为 body；
- `get: "/v1/users/{id}"` → `{id}` 占位符自动从 URL 解析并绑定到 `req.id` 字段。

## 3. 工程目录

```
day06-http-restful/
├── README.md
├── go.mod
├── proto/
│   ├── user.proto
│   └── google/api/annotations.proto    # google api 标准注解（trpc create 自动拉）
├── stub/
├── server/
│   ├── main.go
│   └── service.go
├── client/
│   ├── trpc_main.go      # 用 tRPC 协议调用
│   └── http_main.go      # 用 net/http 直接 curl 风格调用
└── trpc_go.yaml          # 同一个 service 同时暴露 trpc + restful
```

## 4. yaml：同业务多协议

```yaml
server:
  service:
    - name: trpc.study.user.UserService           # tRPC 协议入口
      ip: 127.0.0.1
      port: 8003
      protocol: trpc

    - name: trpc.study.user.UserService.RESTful   # RESTful 协议入口（同一份业务）
      ip: 127.0.0.1
      port: 8004
      protocol: restful
```

业务代码 `RegisterUserServiceService` 只调一次，但框架会在两个端口同时挂 —— **业务实现零改动，一份代码、两个协议**。

## 5. tRPC 错误码 → HTTP 状态码映射

| tRPC 错误码 | 含义 | HTTP 状态码 |
| --- | --- | --- |
| `RetServerDecodeFail (1)` | 请求解码失败 | 400 |
| `RetServerEncodeFail (2)` | 响应编码失败 | 500 |
| `RetServerNoService (101)` | 服务不存在 | 404 |
| `RetServerNoFunc (102)` | 方法不存在 | 404 |
| `RetServerTimeout (101)` | 处理超时 | 504 |
| `RetServerOverload (123)` | 过载 | 429 |
| `RetServerSystemErr (999)` | 系统错误 | 500 |
| `RetServerAuthFail (51)` | 鉴权失败 | 401 |
| `RetServerValidateFail (52)` | 参数校验失败 | 400 |
| 其它业务错误码 | 业务自定义 | 默认 500，可用 `restful.WithStatusCode` 覆写 |

业务侧自定义状态码：

```go
return nil, restful.WithStatusCode{
    StatusCode: http.StatusCreated,
    Err:        errs.New(200, "ok"),
}
```

## 6. 跑起来

```powershell
cd .\tRPC\day06-http-restful
trpc create -p .\proto\user.proto -o . --rpconly
go mod tidy

# 窗口 1：启服务（同时监听 8003/trpc 和 8004/restful）
go run .\server\

# 窗口 2：tRPC 客户端
go run .\client\ -mode trpc

# 窗口 3：HTTP 客户端
curl -X POST http://127.0.0.1:8004/v1/users `
  -H "Content-Type: application/json" `
  -d '{"name":"Alice"}'

curl http://127.0.0.1:8004/v1/users/1
```

期望响应：

```json
// POST /v1/users
{"user":{"id":1,"name":"Alice","createdAt":"1715500800"}}

// GET /v1/users/1
{"user":{"id":1,"name":"Alice","createdAt":"1715500800"}}

// GET /v1/users/9999
HTTP/1.1 500 Internal Server Error
{"code":404,"message":"user not found"}
```

> 注意：第三个请求返回的是 HTTP 500（业务错误码默认映射），如果要返回 HTTP 404，handler 里要用 `restful.WithStatusCode` 显式指定。

## 7. 验证标准

- [ ] tRPC 客户端走 8003 能正常返回；
- [ ] curl 走 8004 同一份业务也能返回；
- [ ] curl `GET /v1/users/9999` 拿到 JSON 错误体 `{"code":404, ...}`；
- [ ] 关掉 8003 端口（删 yaml 里 trpc service），HTTP 仍可工作 —— 验证多协议是真正解耦的；
- [ ] 把 proto 中 `body: "*"` 改成 `body: "user"`，重新 `trpc create`，观察 HTTP body 结构变化。

## 8. 面试复盘

1. **泛 HTTP RPC 与 RESTful 的核心区别？** 前者路径固定为 `/{service}/{method}` POST，body=JSON(req)，是"用 HTTP 模拟 RPC"；后者用 proto 注解定义 RESTful URL，是"真正的 REST API"。
2. **同一个业务实现挂多个协议会有性能损耗吗？** 几乎没有。每个协议各自有独立的端口监听 + 编解码栈，业务 handler 只被调用一次。多出来的成本是端口占用和少量 goroutine。
3. **RESTful 路由用前缀树而不是哈希表的原因？** 因为路径含参数（`/v1/users/{id}`），哈希表无法做带通配符的精确查找；前缀树可以同时支持静态段、参数段、通配段，且支持优先级（静态 > 参数 > 通配）。
4. **`body: "*"` 与 `body: "user"` 的差别在哪？** 前者整个 req 序列化进 body；后者只把 req 的 `user` 字段进 body，其它字段（如 `request_id`）走 query/header。后者更适合"既要 RESTful 风格、又要传辅助参数"。
5. **tRPC 错误码到 HTTP 状态码的映射，能不能让业务码也参与？** 能。`restful.WithStatusCode` 可以让 handler 显式指定 HTTP 状态码，框架 Skip 内置映射。但要克制 —— 同一个业务码在不同接口出现不同 HTTP 码会让客户端心智很乱。
