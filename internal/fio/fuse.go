package fio

import (
	"context"
	"fmt"
	go_fuse "github.com/hanwen/go-fuse/v2/fuse"
	"github.com/hanwen/go-fuse/v2/fuse/nodefs"
	"github.com/hanwen/go-fuse/v2/fuse/pathfs"
	"path/filepath"
	"sync"
	"syscall"
)

// MyFS is our FUSE filesystem implementation.
// It wraps a DirTree structure and exposes it to the kernel.
type MyFS struct {
	pathfs.FileSystem
	mu sync.Mutex // Mutex for concurrent access
	dt *DirTree
}

// GetAttr is used by the kernel to get file/directory attributes.
func (fs *MyFS) GetAttr(ctx context.Context, name string, out *go_fuse.Attr) syscall.Errno {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	node := fs.findNode(name)
	if node == nil {
		return syscall.ENOENT
	}

	out.Mode = syscall.S_IFREG
	if node.IsDir {
		out.Mode = syscall.S_IFDIR
	}
	out.Size = uint64(node.Size)
	out.Mtime = uint64(node.ModTime.Unix())
	out.Ctime = uint64(node.ModTime.Unix())
	out.Nlink = 1

	return syscall.Errno(0)
}

// OpenDir is called by the kernel when a directory is opened.
func (fs *MyFS) OpenDir(ctx context.Context, name string) ([]go_fuse.DirEntry, syscall.Errno) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	node := fs.findNode(name)
	if node == nil {
		return nil, syscall.ENOENT
	}

	if !node.IsDir {
		return nil, syscall.ENOTDIR
	}

	entries := make([]go_fuse.DirEntry, 0, len(node.Children))
	for _, child := range node.Children {
		mode := uint32(syscall.S_IFREG)
		if child.IsDir {
			mode = syscall.S_IFDIR
		}
		entries = append(entries, go_fuse.DirEntry{
			Name: child.Name,
			Mode: mode,
		})
	}

	return entries, syscall.Errno(0)
}

// MyFile implements the nodefs.File interface for our FUSE filesystem.
type MyFile struct {
	data []byte
	nodefs.File
}

// Read implements the Read interface for our file.
func (f *MyFile) Read(dest []byte, off int64) (go_fuse.ReadResult, go_fuse.Status) {
	end := int64(len(f.data))
	if off >= end {
		return nil, go_fuse.OK
	}

	count := copy(dest, f.data[off:])
	return go_fuse.ReadResultData(dest[:count]), go_fuse.OK
}

// Open is called by the kernel when a file is opened.
func (fs *MyFS) Open(ctx context.Context, name string, flags uint32) (nodefs.File, syscall.Errno) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	node := fs.findNode(name)
	if node == nil {
		return nil, syscall.ENOENT
	}

	if node.IsDir {
		return nil, syscall.EISDIR
	}

	// This is a read-only file system, so we can return a simple file with its content.
	// For this example, we will just return a placeholder content for files.
	// You could load content from an external source or provide it in the FileInfo struct.
	content := []byte(fmt.Sprintf("This is the content of %s\n", node.Name))

	// Use our custom MyFile type instead of nodefs.NewFile
	file := &MyFile{
		data: content,
		File: nodefs.NewDefaultFile(),
	}

	return file, syscall.Errno(0)
}

// findNode recursively searches for a node by its path.
func (fs *MyFS) findNode(path string) *FileInfo {
	if path == "" || path == "." {
		return &fs.dt.Root
	}
	parts := filepath.SplitList(path)
	currentNode := &fs.dt.Root

	for _, part := range parts {
		found := false
		for _, child := range currentNode.Children {
			if child.Name == part {
				currentNode = &child
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}
	return currentNode
}
