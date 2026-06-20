package compare

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var ErrUnsafeArchivePath = errors.New("compare: unsafe archive path")

func ExtractTar(reader io.Reader, dest string) error {
	tr := tar.NewReader(reader)
	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read tar header: %w", err)
		}
		if err := extractEntry(tr, dest, header); err != nil {
			return err
		}
	}
}

func extractEntry(reader io.Reader, dest string, header *tar.Header) error {
	target, err := safeArchivePath(dest, header.Name)
	if err != nil {
		return err
	}

	switch header.Typeflag {
	case tar.TypeDir:
		return os.MkdirAll(target, modePerm(header.Mode))
	case tar.TypeReg, tar.TypeRegA:
		if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
			return fmt.Errorf("create parent directory: %w", err)
		}
		return writeExtractedFile(reader, target, modePerm(header.Mode))
	default:
		return nil
	}
}

func safeArchivePath(dest string, name string) (string, error) {
	clean := filepath.Clean(filepath.FromSlash(name))
	if clean == "." || filepath.IsAbs(clean) || clean == ".." {
		return "", fmt.Errorf("%w: %s", ErrUnsafeArchivePath, name)
	}
	if strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("%w: %s", ErrUnsafeArchivePath, name)
	}
	return filepath.Join(dest, clean), nil
}

func writeExtractedFile(reader io.Reader, path string, mode os.FileMode) (err error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return fmt.Errorf("create extracted file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close extracted file: %w", closeErr))
		}
	}()

	if _, err := io.Copy(file, reader); err != nil {
		return fmt.Errorf("write extracted file: %w", err)
	}
	return nil
}

func modePerm(mode int64) os.FileMode {
	perm := os.FileMode(mode) & 0o777
	if perm == 0 {
		return 0o600
	}
	return perm
}
