package zsync

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestContextWaitGroup(t *testing.T) {
	ctx := context.Background()
	wg := NewContextWaitGroup(ctx)

	wg.Go(func(ctx context.Context) error {
		return nil
	})

	if err := wg.Wait(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestContextWaitGroup_Error(t *testing.T) {
	ctx := context.Background()
	wg := NewContextWaitGroup(ctx)

	wg.Go(func(ctx context.Context) error {
		return errors.New("test error")
	})

	if err := wg.Wait(); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestContextWaitGroup_Cancel(t *testing.T) {
	ctx := context.Background()
	wg := NewContextWaitGroup(ctx)

	startTime := time.Now()
	wg.Go(func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(100 * time.Millisecond):
				return nil
			}
		}
	})

	go func() {
		time.Sleep(100 * time.Millisecond)
		wg.Cancel()
	}()

	if err := wg.Wait(); err != nil {
		t.Fatal("expected no error, got", err)
	}
	if time.Since(startTime) < 100*time.Millisecond {
		t.Fatal("expected to cancel after 100ms, took", time.Since(startTime))
	}
}

func TestContextWaitGroup_ErrorInRoutine(t *testing.T) {
	ctx := context.Background()
	wg := NewContextWaitGroup(ctx)

	var ErrTest = errors.New("test error")

	wg.Go(func(ctx context.Context) error {
		return ErrTest
	})

	wg.Go(func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(100 * time.Millisecond):
				return nil
			}
		}
	})

	if err := wg.Wait(); err != ErrTest {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestContextWaitGroup_CancelInRoutine(t *testing.T) {
	ctx := context.Background()
	wg := NewContextWaitGroup(ctx)

	wg.Go(func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(100 * time.Millisecond):
				return nil
			}
		}
	})

	wg.Go(func(ctx context.Context) error {
		wg.Cancel()
		return nil
	})

	if err := wg.Wait(); err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
