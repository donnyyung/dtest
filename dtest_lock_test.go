package dtest

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/datawire/dlib/dlog"
)

func TestDTestLock(t *testing.T) {
	ctx := dlog.NewTestContext(t, false)
	ch := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Second)
		WithNamedMachineLock(ctx, "locktest", func(ctx context.Context) {
			select {
			case <-ch:
			default:
				dlog.Error(ctx, "Lock acquired twice")
				t.Fail()
			}
		})
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		WithNamedMachineLock(ctx, "locktest", func(ctx context.Context) {
			time.Sleep(3 * time.Second)
			close(ch)
		})
	}()
	wg.Wait()
}
