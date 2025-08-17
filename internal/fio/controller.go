package fio

import (
	"github.com/Henry-GongZY/TAAS/internal/fio/storage"
)

const (
	// B represents 1 byte
	B = 1
	// KB represents 1 kilobyte (1024 bytes)
	KB uint64 = 1 << (10 * iota)
	// MB represents 1 megabyte (1024 KB)
	MB
	// GB represents 1 gigabyte (1024 MB)
	GB
	// TB represents 1 terabyte (1024 GB)
	TB
)

type FileConfig struct {
	MaxChunkSize int64
}

// FileController Manage upload and download.
type FileController struct {
	config        *FileConfig
	remoteStorage *storage.RemoteStorage
	fileDB        map[string]*FileInfo
}

func NewFileController(size int64, storage *storage.RemoteStorage) *FileController {
	return &FileController{
		config: &FileConfig{
			MaxChunkSize: size,
		},
		remoteStorage: storage,
		fileDB:        make(map[string]*FileInfo),
	}
}

// UploadFile Upload a whole file
func (c *FileController) UploadFile(filePath string) error {
	return nil
}

// DownloadFile Download a whole file
func (c *FileController) DownloadFile(fileName, savePath string) error {
	return nil
}
