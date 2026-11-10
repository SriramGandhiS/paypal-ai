package outbox_test

import (
	"context"
	"testing"

	"github.com/paybridge/paybridge-ai/services/api/internal/outbox"
)

type dummyHandler struct {
	called bool
}

func (h *dummyHandler) HandleJob(ctx context.Context, job outbox.Job) error {
	h.called = true
	return nil
}

func TestDispatcher_NilStoreSafe(t *testing.T) {
	d := outbox.NewDispatcher(nil, nil)
	handler := &dummyHandler{}
	d.RegisterHandler("DISPATCH_WEBHOOK", handler)

	ctx := context.Background()
	n, err := d.ProcessNextBatch(ctx, "worker-1", 10)
	if err != nil {
		t.Fatalf("unexpected error with nil store: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 processed with nil store, got %d", n)
	}
}
