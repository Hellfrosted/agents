package lockfile

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrLockHeld    = errors.New("lockfile: lock already held")
	ErrLockTimeout = errors.New("lockfile: lock timeout")
)

type LockMetadata struct {
	PID       int    `json:"pid"`
	Host      string `json:"host"`
	CreatedAt string `json:"createdAt"`
}

type Lock struct {
	path string
}

type LockWaitInput struct {
	Path     string
	Metadata LockMetadata
	Retry    <-chan struct{}
}

func AcquireLock(ctx context.Context, path string, metadata LockMetadata) (Lock, error) {
	if err := ctx.Err(); err != nil {
		return Lock{}, fmt.Errorf("acquire lock context: %w", err)
	}

	raw, err := json.Marshal(metadata)
	if err != nil {
		return Lock{}, fmt.Errorf("marshal lock metadata: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return Lock{}, fmt.Errorf("create lock directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return Lock{}, ErrLockHeld
		}
		return Lock{}, fmt.Errorf("create lock file: %w", err)
	}

	if _, err := file.Write(append(raw, '\n')); err != nil {
		return Lock{}, errors.Join(
			fmt.Errorf("write lock metadata: %w", err),
			file.Close(),
		)
	}
	if err := file.Close(); err != nil {
		return Lock{}, fmt.Errorf("close lock file: %w", err)
	}
	return Lock{path: path}, nil
}

func AcquireLockWait(ctx context.Context, input LockWaitInput) (Lock, error) {
	for {
		lock, err := AcquireLock(ctx, input.Path, input.Metadata)
		if err == nil {
			return lock, nil
		}
		if !errors.Is(err, ErrLockHeld) {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return Lock{}, fmt.Errorf("%w: %w", ErrLockTimeout, err)
			}
			return Lock{}, err
		}

		select {
		case <-ctx.Done():
			return Lock{}, fmt.Errorf("%w: %w", ErrLockTimeout, ctx.Err())
		case <-input.Retry:
		}
	}
}

func (l Lock) Release() error {
	if l.path == "" {
		return nil
	}
	if err := os.Remove(l.path); err != nil {
		return fmt.Errorf("remove lock file: %w", err)
	}
	return nil
}
