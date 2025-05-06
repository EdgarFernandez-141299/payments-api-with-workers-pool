package storage

import (
	"context"
	"fmt"
	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	"io"
	"os"
	"path/filepath"

	"gitlab.com/clubhub.ai1/gommon/logger"
)

const (
	StaticFilesDir = "static_files"

	errCreateStaticDir  = "failed to create static files directory: %w"
	errCreateDirStruct  = "failed to create directory structure: %w"
	errCreateFile       = "failed to create file: %w"
	errWriteFileContent = "failed to write file content: %w"
)

type DiskStorageAdapter struct {
	baseDir string
	logger  logger.LoggerInterface
}

func NewDiskStorageAdapter(logger logger.LoggerInterface) (*DiskStorageAdapter, error) {
	if err := os.MkdirAll(StaticFilesDir, 0755); err != nil {
		return nil, fmt.Errorf(errCreateStaticDir, err)
	}

	return &DiskStorageAdapter{
		baseDir: StaticFilesDir,
		logger:  logger,
	}, nil
}

func (d *DiskStorageAdapter) Store(oldCtx context.Context, path string, reader io.Reader) (string, error) {
	var fileURL string
	err := decorators.TraceDecoratorNoReturn(oldCtx, "DiskAdapter.Store", func(ctx context.Context, span decorators.Span) error {
		fullPath := filepath.Join(d.baseDir, path)

		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf(errCreateDirStruct, err)
		}

		file, err := os.Create(fullPath)
		if err != nil {
			return fmt.Errorf(errCreateFile, err)
		}
		defer file.Close()

		_, err = io.Copy(file, reader)
		if err != nil {
			return fmt.Errorf(errWriteFileContent, err)
		}

		absPath, err := filepath.Abs(fullPath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}
		fileURL = absPath

		return nil
	})
	return fileURL, err
}
