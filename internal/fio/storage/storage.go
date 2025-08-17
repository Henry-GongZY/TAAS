package storage

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"sync"
)

// RemoteStorage 是远程存储服务的抽象接口
type RemoteStorage interface {
	UploadChunk(chunkName string, chunkData io.Reader) (string, error)
	DownloadChunk(chunkName string) (io.Reader, error)
}

// MockRemoteStorage 模拟远程存储
type MockRemoteStorage struct {
	chunks map[string][]byte
	mu     sync.Mutex
}

func NewMockRemoteStorage() *MockRemoteStorage {
	return &MockRemoteStorage{
		chunks: make(map[string][]byte),
	}
}

func (s *MockRemoteStorage) UploadChunk(chunkName string, chunkData io.Reader) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := io.ReadAll(chunkData)
	if err != nil {
		return "", err
	}
	s.chunks[chunkName] = data
	// 模拟计算上传块的MD5
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (s *MockRemoteStorage) DownloadChunk(chunkName string) (io.Reader, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, ok := s.chunks[chunkName]
	if !ok {
		return nil, fmt.Errorf("chunk not found: %s", chunkName)
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}
