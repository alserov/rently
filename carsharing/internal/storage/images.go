package storage

import (
	"context"
	"fmt"
	"github.com/Shopify/go-storage"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
)

type ImageStorage interface {
	Save(ctx context.Context, bucket string, f io.Reader) (string, error)
	Get(ctx context.Context, path string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
}

const path = "./files/images"

func NewImageStorage() ImageStorage {
	return &imageStorage{
		s: storage.NewLocalFS(path),
	}
}

type imageStorage struct {
	s storage.FS
}

func (is imageStorage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	f, err := is.s.Open(ctx, fmt.Sprintf("/%s", path), &storage.ReaderOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %v", err)
	}

	return f, nil
}

func (is imageStorage) Save(ctx context.Context, bucket string, f io.Reader) (string, error) {
	id := uuid.New().String()

	w, err := is.s.Create(ctx, fmt.Sprintf("/%s/%s", bucket, id), &storage.WriterOptions{})
	defer w.Close()
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	if _, err = w.Write(data); err != nil {
		return "", fmt.Errorf("failed to save image: %v", err)
	}

	return id, nil
}

func (is imageStorage) Delete(ctx context.Context, key string) error {
	if err := is.s.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}
	return nil
}
