// Package user 定义 User 服务的协议契约（消息结构体）。
//
// ⚠️ 本 demo 故意不使用 protoc 自动生成代码，而是手写 message struct，
//    目的是让读者看清 tRPC-Go 桩代码（stub）的全貌，而不是被 .pb.go 的
//    自动生成内容遮蔽。
//
// ⚠️ 由于不用 protobuf，这里也不实现 proto.Message 接口；
//    序列化方式将通过框架的 codec.SerializationTypeJSON 走 JSON 编解码。
//    （tRPC 协议层依旧是真实的 tRPC 私有协议：16 字节帧头 + PB 包头 + JSON 包体）
package user

// User 是业务领域对象。
//
// 字段：
//   - ID:   用户 ID，由服务端分配（自增）。
//   - Name: 用户名，由客户端创建时传入。
//
// 注：JSON tag 是必须的，因为我们走 JSON 序列化；字段首字母必须大写，
// 否则 encoding/json 看不到非导出字段。
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ===================== GetUser =====================

// GetUserReq 是 GetUser 接口的请求体。
type GetUserReq struct {
	ID int64 `json:"id"`
}

// GetUserRsp 是 GetUser 接口的响应体。
type GetUserRsp struct {
	User *User `json:"user"`
}

// ===================== CreateUser =====================

// CreateUserReq 是 CreateUser 接口的请求体。
type CreateUserReq struct {
	Name string `json:"name"`
}

// CreateUserRsp 是 CreateUser 接口的响应体。
type CreateUserRsp struct {
	User *User `json:"user"`
}

// ===================== ListUser =====================

// ListUserReq 是 ListUser 接口的请求体（无参数，但仍需要一个空结构体，
// 因为 tRPC 桩代码要求每个 RPC 都有明确的 req / rsp 类型）。
type ListUserReq struct{}

// ListUserRsp 是 ListUser 接口的响应体。
type ListUserRsp struct {
	Users []*User `json:"users"`
	Total int     `json:"total"`
}
