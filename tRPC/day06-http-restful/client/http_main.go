// Package main 用标准 net/http 直调 RESTful 端口 8004。
//
// 命令：go run .\client\ -mode http
// 演示 tRPC RESTful 端口对外就是普通 HTTP，与任何语言的 HTTP 客户端都互通。
//go:build http

package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func main() {
	// 1) POST /v1/users
	body := bytes.NewBufferString(`{"name":"Alice"}`)
	rsp, err := http.Post("http://127.0.0.1:8004/v1/users", "application/json", body)
	must(err)
	defer rsp.Body.Close()
	dump("POST /v1/users", rsp)

	// 2) GET /v1/users/1
	rsp2, err := http.Get("http://127.0.0.1:8004/v1/users/1")
	must(err)
	defer rsp2.Body.Close()
	dump("GET /v1/users/1", rsp2)

	// 3) GET /v1/users/9999 → 404 业务码
	rsp3, err := http.Get("http://127.0.0.1:8004/v1/users/9999")
	must(err)
	defer rsp3.Body.Close()
	dump("GET /v1/users/9999", rsp3)
}

func dump(label string, rsp *http.Response) {
	body, _ := io.ReadAll(rsp.Body)
	fmt.Printf("[%s] http=%d body=%s\n", label, rsp.StatusCode, string(body))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
