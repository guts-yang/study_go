# Day 7: 综合实战，使用 net/http 编写并发 Web 服务

本章对应文件：

- `server.go`

## 学习目标

- 使用标准库 `net/http` 搭建 Web 服务。
- 编写健康检查、列表、创建、详情接口。
- 使用 JSON 编解码。
- 理解 HTTP handler 的并发模型。
- 使用锁保护服务内存状态。

## 运行方式

启动服务：

```powershell
go run ./day07-web-server
```

服务默认监听：

```text
http://localhost:8080
```

## 接口示例

健康检查：

```powershell
curl http://localhost:8080/health
```

查看用户列表：

```powershell
curl http://localhost:8080/users
```

创建用户：

```powershell
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d "{\"name\":\"Charlie\"}"
```

查看用户详情：

```powershell
curl http://localhost:8080/users/1
```

## 知识点讲解

`net/http` 是 Go 标准库中非常成熟的 Web 基础库。一个 handler 通常长这样：

```go
func(w http.ResponseWriter, r *http.Request) {
}
```

`http.ResponseWriter` 用于写响应，`*http.Request` 表示请求信息。

Go 1.22 的 `http.ServeMux` 支持更清晰的路由写法：

```go
mux.HandleFunc("GET /users/{id}", handler)
```

请求中的路径参数可以通过 `r.PathValue("id")` 获取。

`net/http` 的重要工程点是：每个请求通常会在独立 goroutine 中处理。因此 handler 访问共享变量时必须考虑并发安全。本章的 `UserStore` 使用 `sync.RWMutex` 保护 `map[int]User` 和 `nextID`。

`RWMutex` 允许多个读请求并发进入，但写请求需要独占。对于读多写少的服务状态，这是常见优化。

## 重点代码

- `User`：响应 JSON 的业务模型。
- `UserStore`：并发安全的内存存储。
- `Server.routes()`：注册路由。
- `handleCreateUser()`：解析 JSON 请求体。
- `writeJSON()`：统一 JSON 响应。
- `logRequest()`：简单请求日志中间件。

## 动手练习

1. 增加 `DELETE /users/{id}` 接口。
2. 给 `POST /users` 增加名称长度校验。
3. 把 `UserStore` 的列表结果按 ID 排序。
4. 为 handler 编写 `httptest` 单元测试。

## 复盘问题

- 为什么 handler 中访问 map 要加锁？
- `Mutex` 和 `RWMutex` 有什么区别？
- 为什么后端服务要设置 `ReadTimeout` 和 `WriteTimeout`？
- 标准库 `net/http` 与 Web 框架的关系是什么？
- 如何把这个内存版服务演进成连接数据库的真实服务？
