// Package main 是 day06-http-restful 的客户端，演示 tRPC 和 HTTP 两种调用方式。
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	pb "day06-http-restful/stub/trpc/study/user"

	_ "git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/client"
)

func main() {
	// --- 方式 1：tRPC 协议调用（port 8001）---
	fmt.Println("=== tRPC 调用 ===")
	proxy := pb.NewUserServiceClientProxy(
		client.WithTarget("ip://127.0.0.1:8001"),
		client.WithTimeout(time.Second),
	)
	crsp, err := proxy.CreateUser(context.Background(), &pb.CreateUserReq{Name: "Eve"})
	if err != nil {
		fmt.Printf("tRPC CreateUser error: %v\n", err)
	} else {
		fmt.Printf("tRPC CreateUser ok: id=%d\n", crsp.GetUser().GetId())
	}

	// --- 方式 2：HTTP JSON 调用（port 8080）---
	fmt.Println("\n=== HTTP/JSON 调用 ===")
	httpCreateUser("Frank")
}

func httpCreateUser(name string) {
	body, _ := json.Marshal(map[string]string{"name": name})
	resp, err := http.Post(
		"http://127.0.0.1:8080/trpc.study.user.UserService/CreateUser",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		fmt.Printf("HTTP CreateUser error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	fmt.Printf("HTTP CreateUser response: %s\n", string(data))
}
