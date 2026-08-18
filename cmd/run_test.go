package cmd

import (
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWaitForShutdown(t *testing.T) {
	t.Run("returns on SIGTERM", func(t *testing.T) {
		assertUnblocksOn(t, syscall.SIGTERM)
	})

	t.Run("returns on SIGINT", func(t *testing.T) {
		assertUnblocksOn(t, syscall.SIGINT)
	})
}

// assertUnblocksOn sends sig to the test process and asserts that
// waitForShutdown stops blocking.
func assertUnblocksOn(t *testing.T, sig syscall.Signal) {
	t.Helper()

	done := make(chan struct{})
	go func() {
		waitForShutdown()
		close(done)
	}()

	// Give waitForShutdown time to register its handler, otherwise the signal
	// reaches the default handler and terminates the test process.
	time.Sleep(100 * time.Millisecond)

	assert.NoError(t, syscall.Kill(syscall.Getpid(), sig))

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("waitForShutdown did not return on %s", sig)
	}
}
