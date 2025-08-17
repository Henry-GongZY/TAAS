package fio

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// FileInfo represents a file or directory's metadata.
type FileInfo struct {
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	IsDir    bool       `json:"isDir"`
	Size     int64      `json:"size,omitempty"`
	ModTime  time.Time  `json:"modTime,omitempty"`
	Children []FileInfo `json:"children,omitempty"`
}

// DirTree represents a directory tree and provides methods for saving and loading.
type DirTree struct {
	Root FileInfo `json:"root"`
}

// buildDirTree is a recursive helper function to build the directory tree.
func buildDirTree(path string) (FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return FileInfo{}, err
	}

	fileInfo := FileInfo{
		Name:    info.Name(),
		Path:    path,
		IsDir:   info.IsDir(),
		Size:    info.Size(),
		ModTime: info.ModTime(),
	}

	if !info.IsDir() {
		return fileInfo, nil // Return if it's a file
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return fileInfo, err
	}

	fileInfo.Children = make([]FileInfo, 0, len(entries))
	for _, entry := range entries {
		childPath := filepath.Join(path, entry.Name())
		childInfo, err := buildDirTree(childPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to read %s - %v\n", childPath, err)
			continue
		}
		fileInfo.Children = append(fileInfo.Children, childInfo)
	}

	return fileInfo, nil
}

// NewDirTreeFromPath creates a new DirTree from local filesystem(path).
func NewDirTreeFromPath(rootPath string) (*DirTree, error) {
	rootInfo, err := buildDirTree(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build directory tree: %w", err)
	}
	return &DirTree{Root: rootInfo}, nil
}

//// NewDirTree creates an empty DirTree.
//func NewDirTree() (*DirTree, error) {
//
//}

// SaveToJson serializes the directory tree to a JSON file.
func (dt *DirTree) SaveToJson(outputPath string) error {
	jsonData, err := json.MarshalIndent(dt.Root, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	fmt.Printf("Directory structure saved to %s\n", outputPath)
	return nil
}

// LoadFromJson reads a JSON file and returns a new DirTree instance.
func LoadFromJson(filePath string) (*DirTree, error) {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	var rootInfo FileInfo
	err = json.Unmarshal(jsonData, &rootInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	fmt.Printf("Directory structure loaded from %s\n", filePath)
	return &DirTree{Root: rootInfo}, nil
}
