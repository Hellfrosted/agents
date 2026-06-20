package lockfile

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestAcquireLock_createsLockFileWithMetadata_whenLockAvailable(t *testing.T) {
	// Given
	lockPath := filepath.Join(t.TempDir(), ".skill-lock.json.lock")
	metadata := LockMetadata{PID: 1234, Host: "test-host", CreatedAt: "2026-06-20T00:00:00Z"}

	// When
	lock, err := AcquireLock(context.Background(), lockPath, metadata)
	if err != nil {
		t.Fatalf("AcquireLock returned error: %v", err)
	}
	defer lock.Release()

	// Then
	raw, err := os.ReadFile(lockPath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if !containsString(string(raw), `"pid":1234`) {
		t.Fatalf("lock metadata missing pid: %s", raw)
	}
}

func TestAcquireLock_returnsLockHeld_whenLockFileExists(t *testing.T) {
	// Given
	lockPath := filepath.Join(t.TempDir(), ".skill-lock.json.lock")
	if err := os.WriteFile(lockPath, []byte(`{"pid":1}`), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	// When
	_, err := AcquireLock(context.Background(), lockPath, LockMetadata{})

	// Then
	if !errors.Is(err, ErrLockHeld) {
		t.Fatalf("AcquireLock error = %v, want ErrLockHeld", err)
	}
}

func TestLockRelease_removesLockFile_whenHeld(t *testing.T) {
	// Given
	lockPath := filepath.Join(t.TempDir(), ".skill-lock.json.lock")
	lock, err := AcquireLock(context.Background(), lockPath, LockMetadata{})
	if err != nil {
		t.Fatalf("AcquireLock returned error: %v", err)
	}

	// When
	err = lock.Release()

	// Then
	if err != nil {
		t.Fatalf("Release returned error: %v", err)
	}
	if _, err := os.Stat(lockPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("lock file still exists or stat failed differently: %v", err)
	}
}

func TestAcquireLockWait_acquiresAfterRetry_whenExistingLockIsReleased(t *testing.T) {
	// Given
	lockPath := filepath.Join(t.TempDir(), ".skill-lock.json.lock")
	first, err := AcquireLock(context.Background(), lockPath, LockMetadata{PID: 1})
	if err != nil {
		t.Fatalf("AcquireLock returned error: %v", err)
	}
	retry := make(chan struct{}, 1)

	// When
	if err := first.Release(); err != nil {
		t.Fatalf("Release returned error: %v", err)
	}
	retry <- struct{}{}
	got, err := AcquireLockWait(context.Background(), LockWaitInput{
		Path:     lockPath,
		Metadata: LockMetadata{PID: 2},
		Retry:    retry,
	})

	// Then
	if err != nil {
		t.Fatalf("AcquireLockWait returned error: %v", err)
	}
	defer got.Release()
	raw, err := os.ReadFile(lockPath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if !containsString(string(raw), `"pid":2`) {
		t.Fatalf("lock metadata missing new pid: %s", raw)
	}
}

func TestAcquireLockWait_returnsTimeout_whenContextEndsBeforeRetry(t *testing.T) {
	// Given
	lockPath := filepath.Join(t.TempDir(), ".skill-lock.json.lock")
	first, err := AcquireLock(context.Background(), lockPath, LockMetadata{PID: 1})
	if err != nil {
		t.Fatalf("AcquireLock returned error: %v", err)
	}
	defer first.Release()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// When
	_, err = AcquireLockWait(ctx, LockWaitInput{Path: lockPath, Metadata: LockMetadata{PID: 2}})

	// Then
	if !errors.Is(err, ErrLockTimeout) {
		t.Fatalf("AcquireLockWait error = %v, want ErrLockTimeout", err)
	}
}

func containsString(text string, part string) bool {
	return len(part) == 0 || findString(text, part) >= 0
}

func findString(text string, part string) int {
	for i := 0; i+len(part) <= len(text); i++ {
		if text[i:i+len(part)] == part {
			return i
		}
	}
	return -1
}
