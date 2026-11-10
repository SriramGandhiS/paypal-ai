package inbox_test

import (
	"context"
	"testing"

	"github.com/paybridge/paybridge-ai/services/api/internal/inbox"
)

func TestProcessor_NilStoreSafe(t *testing.T) {
	p := inbox.NewProcessor(nil, nil)
	ctx := context.Background()

	n, err := p.ProcessBatch(ctx, "worker-1", 10)
	if err != nil {
		t.Fatalf("unexpected error with nil store: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 processed with nil store, got %d", n)
	}
}
