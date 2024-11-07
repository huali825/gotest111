package context

import (
	"context"
	"testing"
)

func TestCtxCode(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key1", "value1")
	val := ctx.Value("key1")
	if val != "value1" {
		t.Error("value1 not found")
	}
	t.Log(val)

}

func TestCtxCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ch := ctx.Done()
		<-ch
	}()
}
