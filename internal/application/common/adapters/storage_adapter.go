package adapters

import (
	"context"
	"io"
)

type StorageAdapter interface {
	Store(ctx context.Context, path string, file io.Reader) (string, error)
}
