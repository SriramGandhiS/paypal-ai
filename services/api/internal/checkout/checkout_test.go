package checkout_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/paybridge/paybridge-ai/services/api/internal/checkout"
	"github.com/paybridge/paybridge-ai/services/api/internal/store"
	"github.com/paybridge/paybridge-ai/services/api/internal/tenancy"
)

func TestIsTerminal(t *testing.T) {
	terminalStates := []string{"CANCELLED", "CAPTURE_DENIED", "APPLICATION_RECONCILED"}
	for _, st := range terminalStates {
		if !checkout.IsTerminal(st) {
			t.Errorf("expected %s to be terminal", st)
		}
	}

	nonTerminalStates := []string{"ORDER_CREATING", "ORDER_CREATED", "APPROVAL_PENDING", "APPROVED", "CAPTURE_REQUESTED"}
	for _, st := range nonTerminalStates {
		if checkout.IsTerminal(st) {
			t.Errorf("expected %s to be non-terminal", st)
		}
	}
}

func TestCheckoutRoleEnforcement(t *testing.T) {
	st := store.NewFromDB(nil)
	svc := checkout.NewService(st)

	t.Run("receiver cannot create order attempt", func(t *testing.T) {
		ctx := tenancy.WithTenant(context.Background(), "018f3a21-9b34-7a12-89cd-456789000000")
		ctx = tenancy.WithUser(ctx, "user-123", "receiver")

		_, err := svc.CreateAttempt(ctx, "case-123", "req-abc-999")
		if err == nil {
			t.Fatal("expected error when receiver attempts order creation, got nil")
		}
	})
}

func TestCaptureLock_MutualExclusion(t *testing.T) {
	st := store.NewFromDB(nil)
	svc := checkout.NewService(st)
	orderID := "ORDER-MUTEX-TEST-001"

	// 1. First caller acquires lock
	release1, acquired1 := svc.AcquireCaptureLock(orderID)
	if !acquired1 || release1 == nil {
		t.Fatalf("expected first caller to acquire lock, got acquired=%v", acquired1)
	}

	// 2. Second caller attempting the same order must be rejected
	release2, acquired2 := svc.AcquireCaptureLock(orderID)
	if acquired2 || release2 != nil {
		t.Fatalf("expected second caller to be rejected, got acquired=%v", acquired2)
	}

	// 3. Releasing lock allows subsequent acquisition
	release1()

	release3, acquired3 := svc.AcquireCaptureLock(orderID)
	if !acquired3 || release3 == nil {
		t.Fatalf("expected subsequent caller to acquire released lock, got acquired=%v", acquired3)
	}
	release3()
}

// TestTwoIndependentServiceInstances verifies multi-instance locking across separate service replicas.
func TestTwoIndependentServiceInstances(t *testing.T) {
	sharedLocker := checkout.NewSharedMemoryLocker()

	// Instance 1 and Instance 2 are two completely distinct service instances
	svc1 := checkout.NewServiceWithLocker(nil, sharedLocker)
	svc2 := checkout.NewServiceWithLocker(nil, sharedLocker)

	ctx := context.Background()
	tenantID := "018f3a21-9b34-7a12-89cd-456789000000"
	actorID1 := "usr-alice"
	actorID2 := "usr-bob"
	orderID := "ORDER-SHARED-REPLICA-999"

	t.Run("concurrent capture across independent instances allows exactly one winner", func(t *testing.T) {
		startGate := make(chan struct{})
		var wg sync.WaitGroup
		results := make(chan bool, 2)

		wg.Add(2)
		go func() {
			defer wg.Done()
			<-startGate
			rel, acq, err := svc1.AcquireDistributedCaptureLock(ctx, tenantID, actorID1, orderID)
			if err == nil && acq {
				results <- true
				time.Sleep(30 * time.Millisecond)
				rel()
			} else {
				results <- false
			}
		}()

		go func() {
			defer wg.Done()
			<-startGate
			rel, acq, err := svc2.AcquireDistributedCaptureLock(ctx, tenantID, actorID2, orderID)
			if err == nil && acq {
				results <- true
				time.Sleep(30 * time.Millisecond)
				rel()
			} else {
				results <- false
			}
		}()

		close(startGate)
		wg.Wait()
		close(results)

		var successCount, failureCount int
		for r := range results {
			if r {
				successCount++
			} else {
				failureCount++
			}
		}

		if successCount != 1 || failureCount != 1 {
			t.Fatalf("expected exactly 1 winner and 1 rejected caller, got success=%d, failure=%d", successCount, failureCount)
		}
	})

	t.Run("release by instance 1 allows instance 2 to acquire immediately", func(t *testing.T) {
		orderID := "ORDER-RELEASE-HANDOVER-001"
		rel1, acq1, err1 := svc1.AcquireDistributedCaptureLock(ctx, tenantID, actorID1, orderID)
		if !acq1 || err1 != nil {
			t.Fatalf("expected svc1 to acquire lock, got acq=%v, err=%v", acq1, err1)
		}

		// svc2 must fail while svc1 holds lock
		_, acq2Before, _ := svc2.AcquireDistributedCaptureLock(ctx, tenantID, actorID2, orderID)
		if acq2Before {
			t.Fatal("expected svc2 to be rejected while svc1 holds lock")
		}

		// svc1 releases
		rel1()

		// svc2 can now acquire
		rel2, acq2After, err2 := svc2.AcquireDistributedCaptureLock(ctx, tenantID, actorID2, orderID)
		if !acq2After || err2 != nil {
			t.Fatalf("expected svc2 to acquire lock after svc1 release, got acq=%v, err=%v", acq2After, err2)
		}
		rel2()
	})

	t.Run("lease expiration allows recovery by instance 2 even without explicit release", func(t *testing.T) {
		shortLocker := checkout.NewSharedMemoryLocker()
		shortSvc2 := checkout.NewServiceWithLocker(nil, shortLocker)

		orderID := "ORDER-TTL-EXPIRY-002"
		lockKey := "capture:" + tenantID + ":" + orderID

		// Acquire with 40ms TTL directly on locker to test expiry recovery
		_, acq1, err := shortLocker.TryAcquire(ctx, lockKey, 40*time.Millisecond)
		if !acq1 || err != nil {
			t.Fatalf("expected initial acquire, got acq=%v, err=%v", acq1, err)
		}

		// Immediate attempt from svc2 must be rejected
		_, acq2Immediate, _ := shortSvc2.AcquireDistributedCaptureLock(ctx, tenantID, actorID2, orderID)
		if acq2Immediate {
			t.Fatal("expected immediate acquire by svc2 to be rejected")
		}

		// Wait 60ms for lease to expire
		time.Sleep(60 * time.Millisecond)

		// svc2 should now successfully recover and acquire expired lease
		rel2, acq2Recovered, err2 := shortSvc2.AcquireDistributedCaptureLock(ctx, tenantID, actorID2, orderID)
		if !acq2Recovered || err2 != nil {
			t.Fatalf("expected svc2 to recover expired lease, got acq=%v, err=%v", acq2Recovered, err2)
		}
		rel2()
	})
}

// faultyLocker simulates database disconnection / failure
type faultyLocker struct {
	err error
}

func (f *faultyLocker) TryAcquire(ctx context.Context, key string, ttl time.Duration) (func(), bool, error) {
	return nil, false, f.err
}

func TestCaptureLock_DatabaseFailureHandling(t *testing.T) {
	dbErr := errors.New("connection to PostgreSQL pool timed out")
	fl := &faultyLocker{err: dbErr}
	svc := checkout.NewServiceWithLocker(nil, fl)

	ctx := context.Background()
	tenantID := "018f3a21-9b34-7a12-89cd-456789000000"
	actorID := "usr-actor"
	orderID := "ORDER-DB-FAIL-001"

	rel, acq, err := svc.AcquireDistributedCaptureLock(ctx, tenantID, actorID, orderID)
	if acq {
		t.Fatal("expected lock acquisition to fail during DB error")
	}
	if rel != nil {
		t.Fatal("expected release func to be nil on DB failure")
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected error %v, got %v", dbErr, err)
	}

	// Verify local in-process lock was safely cleaned up so no dangling locks remain
	relRetry, acqRetry := svc.AcquireCaptureLock(orderID)
	if !acqRetry || relRetry == nil {
		t.Fatal("expected in-process lock to be clean after DB failure, but was held")
	}
	relRetry()
}
