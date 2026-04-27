package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Day 7 目标：
// 1. 使用标准库 net/http 编写极简 Web 服务。
// 2. 理解每个请求会在独立 goroutine 中处理，因此 handler 内共享状态必须并发安全。
// 3. 使用 JSON、路由、超时配置搭出一个接近真实后端服务的雏形。

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type UserStore struct {
	mu     sync.RWMutex
	nextID int
	users  map[int]User
}

func NewUserStore() *UserStore {
	return &UserStore{
		nextID: 1,
		users:  make(map[int]User),
	}
}

func (s *UserStore) Create(name string) User {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := User{
		ID:        s.nextID,
		Name:      name,
		CreatedAt: time.Now(),
	}
	s.users[user.ID] = user
	s.nextID++
	return user
}

func (s *UserStore) List() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

func (s *UserStore) Get(id int) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	return user, ok
}

type Server struct {
	store *UserStore
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /users", s.handleListUsers)
	mux.HandleFunc("POST /users", s.handleCreateUser)
	mux.HandleFunc("GET /users/{id}", s.handleGetUser)
	return logRequest(mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleListUsers(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.store.List())
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}
	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
		return
	}

	user := s.store.Create(req.Name)
	writeJSON(w, http.StatusCreated, user)
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
		return
	}

	user, ok := s.store.Get(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("encode response:", err)
	}
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s cost=%s", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	store := NewUserStore()
	store.Create("Alice")
	store.Create("Bob")

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      (&Server{store: store}).routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("server listening on http://localhost:8080")
	log.Println("try: curl http://localhost:8080/health")

	// net/http 的并发模型：
	// - Server 接收到连接后，会为请求分配 goroutine 执行 handler。
	// - 因此多个请求可能同时访问 UserStore。
	// - UserStore 使用 RWMutex：读请求可并发，写请求独占，保护 map 和 nextID。
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
