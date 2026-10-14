package repository

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
)

func newTestRepo(t *testing.T) *Repository {
	t.Helper()

	path := filepath.Join(t.TempDir(), "test.db")
	repo, err := New(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { repo.Close() })

	return repo
}

func TestIssueAndReset(t *testing.T) {
	ctx := context.Background()
	repo := newTestRepo(t)

	first, err := repo.Issue(ctx)
	if err != nil || first.QueueNumber != "A0" {
		t.Fatalf("first ticket: %+v, err: %v", first, err)
	}

	second, err := repo.Issue(ctx)
	if err != nil || second.QueueNumber != "A1" {
		t.Fatalf("second ticket: %+v, err: %v", second, err)
	}

	if err := repo.Reset(ctx); err != nil {
		t.Fatal(err)
	}

	current, err := repo.Current(ctx)
	if err != nil || current != "00" {
		t.Fatalf("after reset: %s, %v", current, err)
	}

	next, err := repo.Issue(ctx)
	if err != nil || next.QueueNumber != "A0" {
		t.Fatalf("next cycle: %+v, %v", next, err)
	}

	if next.Cycle == first.Cycle || next.ID == first.ID {
		t.Fatal("reset must create a distinct cycle and ticket")
	}

	old, err := repo.GetTicket(ctx, first.ID)
	if err != nil || old.Cycle != first.Cycle {
		t.Fatal("old ticket must remain identifiable")
	}
}

func TestQueueCapacity(t *testing.T) {
	ctx := context.Background()
	repo := newTestRepo(t)

	for i := 0; i < 260; i++ {
		ticket, err := repo.Issue(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if i == 259 && ticket.QueueNumber != "Z9" {
			t.Fatalf("expected Z9, got %s", ticket.QueueNumber)
		}
	}

	_, err := repo.Issue(ctx)
	if !errors.Is(err, ErrQueueFull) {
		t.Fatalf("expected queue full, got %v", err)
	}
}

func TestConcurrentIssuance(t *testing.T) {
	ctx := context.Background()
	repo := newTestRepo(t)

	const total = 100

	var wg sync.WaitGroup
	results := make(chan string, total)
	errs := make(chan error, total)

	for i := 0; i < total; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			ticket, err := repo.Issue(ctx)
			if err != nil {
				errs <- err
				return
			}
			results <- ticket.QueueNumber
		}()
	}

	wg.Wait()
	close(results)
	close(errs)

	for err := range errs {
		t.Errorf("concurrent issue failed: %v", err)
	}

	seen := make(map[string]bool)

	for number := range results {
		if seen[number] {
			t.Errorf("duplicate ticket: %s", number)
		}
		seen[number] = true
	}

	if len(seen) != total {
		t.Fatalf("expected %d unique tickets, got %d",
			total, len(seen))
	}

	for i := 0; i < total; i++ {
		expected := fmt.Sprintf("%c%d",
			'A'+rune(i/10), i%10)
		if !seen[expected] {
			t.Errorf("missing ticket %s", expected)
		}
	}
}
