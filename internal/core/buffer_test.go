package core

import (
	"bytes"
	"testing"
)

func TestPooledBufferLifecycle(t *testing.T) {
	t.Parallel()

	b := NewPooledBuffer()
	if _, err := b.WriteString("payload"); err != nil {
		t.Fatal(err)
	}
	if b.String() != "payload" {
		t.Errorf("got %q, want %q", b.String(), "payload")
	}
	if err := b.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	// A buffer obtained after Close must start empty regardless of which
	// pooled instance is returned.
	next := NewPooledBuffer()
	if next.Len() != 0 {
		t.Errorf("reused buffer must be reset, got %q", next.String())
	}
	if err := next.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
}

func TestPooledBufferReadFrom(t *testing.T) {
	t.Parallel()

	b := NewPooledBuffer()
	defer func() {
		_ = b.Close()
	}()

	if _, err := b.ReadFrom(bytes.NewReader([]byte("hello"))); err != nil {
		t.Fatal(err)
	}
	if b.String() != "hello" {
		t.Errorf("got %q, want %q", b.String(), "hello")
	}
}
