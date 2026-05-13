package storage

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mob/ddd-template/internal/app/port"
	"os"
	"path/filepath"
	"strings"
)

type LocalStorage struct {
	storagePath string
}

func NewLocalStorage(storagePath string) port.FileStorage {
	return &LocalStorage{storagePath: filepath.Clean(storagePath)}
}

func (s *LocalStorage) Put(key string, content io.Reader) error {
	path, err := s.pathForKey(key)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, content)
	return err
}

func (s *LocalStorage) Get(key string) (io.ReadCloser, error) {
	path, err := s.pathForKey(key)
	if err != nil {
		return nil, err
	}

	return os.Open(path)
}

func (s *LocalStorage) Delete(key string) error {
	path, err := s.pathForKey(key)
	if err != nil {
		return err
	}

	err = os.Remove(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	return err
}

func (s *LocalStorage) pathForKey(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("storage key cannot be empty")
	}

	cleanKey := filepath.Clean(key)
	if filepath.IsAbs(cleanKey) || cleanKey == "." || strings.HasPrefix(cleanKey, ".."+string(os.PathSeparator)) || cleanKey == ".." {
		return "", fmt.Errorf("invalid storage key: %q", key)
	}

	path := filepath.Join(s.storagePath, cleanKey)
	root, err := filepath.Abs(s.storagePath)
	if err != nil {
		return "", err
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	if absPath != root && !strings.HasPrefix(absPath, root+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid storage key: %q", key)
	}

	return path, nil
}
