package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCountConcurrently(t *testing.T) {
	got := CountConcurrently(10, 1000)
	want := 10000

	if got != want {
		t.Fatalf("CountConcurrently() = %d, want %d", got, want)
	}
}

func TestFetchWithContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := FetchWithContext(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("FetchWithContext() error = %v, want context deadline exceeded", err)
	}
}
