package api_test

import (
	"runtime"
	"testing"
)

func TestFoundationBaseline(t *testing.T) {
	if runtime.GOOS == "" {
		t.Fatal("expected non-empty GOOS")
	}
	t.Logf("PayBridge API Foundation Baseline initialized on %s/%s (Go %s)", runtime.GOOS, runtime.GOARCH, runtime.Version())
}
