package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Hellfrosted/agents/internal/skup/config"
	"github.com/Hellfrosted/agents/internal/skup/lockfile"
)

type lockSnapshot struct {
	exists bool
	raw    []byte
}

func withLockfilePreservation(ctx context.Context, resolved config.Resolved, request Request, action func() error) error {
	lockPath := filepath.Join(resolved.AgentsHome, ".skill-lock.json")
	lock, err := lockfile.AcquireLock(ctx, lockPath+".lock", lockMetadata(request))
	if err != nil {
		if errors.Is(err, lockfile.ErrLockHeld) {
			return fmt.Errorf("%w: %s", lockfile.ErrLockTimeout, lockPath+".lock")
		}
		return err
	}
	defer func() {
		if err := lock.Release(); err != nil {
			writeError(request.Stderr, err)
		}
	}()

	snapshot, err := readLockSnapshot(lockPath)
	if err != nil {
		return err
	}
	actionErr := action()
	if actionErr != nil {
		if restoreErr := restoreLockSnapshot(lockPath, snapshot); restoreErr != nil {
			return errors.Join(actionErr, restoreErr)
		}
		return actionErr
	}
	if err := preserveLockSnapshot(lockPath, snapshot); err != nil {
		return err
	}
	return nil
}

func lockMetadata(request Request) lockfile.LockMetadata {
	host, _ := os.Hostname()
	return lockfile.LockMetadata{
		PID:       os.Getpid(),
		Host:      host,
		CreatedAt: nowUTC(request).Format(time.RFC3339Nano),
	}
}

func readLockSnapshot(path string) (lockSnapshot, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return lockSnapshot{}, nil
		}
		return lockSnapshot{}, fmt.Errorf("read lockfile snapshot: %w", err)
	}
	return lockSnapshot{exists: true, raw: append([]byte(nil), raw...)}, nil
}

func restoreLockSnapshot(path string, snapshot lockSnapshot) error {
	if !snapshot.exists {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove lockfile after failed command: %w", err)
		}
		return nil
	}
	return writeRawFile(path, snapshot.raw)
}

func preserveLockSnapshot(path string, snapshot lockSnapshot) error {
	if !snapshot.exists {
		return nil
	}
	after, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return writeRawFile(path, snapshot.raw)
		}
		return fmt.Errorf("read lockfile after command: %w", err)
	}
	merged, err := mergeLockJSON(snapshot.raw, after)
	if err != nil {
		return writeRawFile(path, snapshot.raw)
	}
	return writeRawFile(path, merged)
}

func mergeLockJSON(before []byte, after []byte) ([]byte, error) {
	var beforeFields map[string]json.RawMessage
	if err := json.Unmarshal(before, &beforeFields); err != nil {
		return nil, err
	}
	var afterFields map[string]json.RawMessage
	if err := json.Unmarshal(after, &afterFields); err != nil {
		return nil, err
	}
	mergeTopLevelFields(afterFields, beforeFields)
	mergeSkills(afterFields, beforeFields)
	raw, err := json.MarshalIndent(afterFields, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal preserved lockfile: %w", err)
	}
	return append(raw, '\n'), nil
}

func mergeTopLevelFields(after map[string]json.RawMessage, before map[string]json.RawMessage) {
	for name, raw := range before {
		if name == "version" || name == "skills" {
			continue
		}
		after[name] = mergeObjects(after[name], raw)
	}
}

func mergeSkills(after map[string]json.RawMessage, before map[string]json.RawMessage) {
	var beforeSkills map[string]json.RawMessage
	if err := json.Unmarshal(before["skills"], &beforeSkills); err != nil {
		return
	}
	var afterSkills map[string]json.RawMessage
	if err := json.Unmarshal(after["skills"], &afterSkills); err != nil || afterSkills == nil {
		afterSkills = map[string]json.RawMessage{}
	}
	for name, raw := range beforeSkills {
		afterSkills[name] = mergeObjects(afterSkills[name], raw)
	}
	raw, err := json.Marshal(afterSkills)
	if err == nil {
		after["skills"] = raw
	}
}

func mergeObjects(after json.RawMessage, before json.RawMessage) json.RawMessage {
	if len(after) == 0 {
		return append(json.RawMessage(nil), before...)
	}
	var beforeObject map[string]json.RawMessage
	var afterObject map[string]json.RawMessage
	if json.Unmarshal(before, &beforeObject) != nil || json.Unmarshal(after, &afterObject) != nil {
		return append(json.RawMessage(nil), after...)
	}
	for name, raw := range beforeObject {
		afterObject[name] = mergeObjects(afterObject[name], raw)
	}
	raw, err := json.Marshal(afterObject)
	if err != nil {
		return append(json.RawMessage(nil), after...)
	}
	return raw
}

func writeRawFile(path string, raw []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create lockfile directory: %w", err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".skill-lock.json.tmp-*")
	if err != nil {
		return fmt.Errorf("create lockfile temp file: %w", err)
	}
	tmpPath := tmp.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tmpPath)
		}
	}()
	if _, err := tmp.Write(raw); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write lockfile temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close lockfile temp file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("replace lockfile: %w", err)
	}
	cleanup = false
	return nil
}
