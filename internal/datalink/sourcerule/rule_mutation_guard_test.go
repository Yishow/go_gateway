package sourcerule

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRuleMutationGuardNestedCallsOtherRuleAndReleasedContext(t *testing.T) {
	service := &Service{}
	ctx, release := service.AcquireRuleMutation(context.Background(), "one")
	_, nestedRelease := service.AcquireRuleMutation(ctx, "one")
	nestedRelease()
	_, otherRelease := service.AcquireRuleMutation(context.Background(), "two")
	otherRelease()
	waiter := make(chan struct{})
	go func() { _, end := service.AcquireRuleMutation(context.Background(), "one"); end(); close(waiter) }()
	select {
	case <-waiter:
		release()
		t.Fatal("nested release unlocked outer operation")
	case <-time.After(20 * time.Millisecond):
	}
	release()
	release()
	select {
	case <-waiter:
	case <-time.After(time.Second):
		t.Fatal("waiting rule did not progress")
	}
	fresh, finish := service.AcquireRuleMutation(context.Background(), "one")
	_ = fresh
	reused := make(chan struct{})
	go func() { _, end := service.AcquireRuleMutation(ctx, "one"); end(); close(reused) }()
	select {
	case <-reused:
		finish()
		t.Fatal("released context bypassed an active rule guard")
	case <-time.After(20 * time.Millisecond):
	}
	finish()
	select {
	case <-reused:
	case <-time.After(time.Second):
		t.Fatal("released context deadlocked")
	}
	// Completion channels close only after each release removes its registry use.
	_, end := service.AcquireRuleMutation(context.Background(), "one")
	end()
	service.ruleGuardsMu.Lock()
	defer service.ruleGuardsMu.Unlock()
	require.Empty(t, service.ruleGuards)
}
