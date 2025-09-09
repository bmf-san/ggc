package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MemoryFileSystem implements FileSystem in memory for testing
type MemoryFileSystem struct {
	files map[string][]byte
	dirs  map[string]bool
}

// NewMemoryFileSystem creates a new memory-based filesystem for testing
func NewMemoryFileSystem() *MemoryFileSystem {
	mfs := &MemoryFileSystem{
		files: make(map[string][]byte),
		dirs:  make(map[string]bool),
	}
	// Root directory always exists
	mfs.dirs["/"] = true
	mfs.dirs["."] = true
	return mfs
}

func (mfs *MemoryFileSystem) ReadFile(filename string) ([]byte, error) {
	data, exists := mfs.files[filename]
	if !exists {
		return nil, &os.PathError{Op: "open", Path: filename, Err: os.ErrNotExist}
	}
	// Return a copy to prevent external modification
	result := make([]byte, len(data))
	copy(result, data)
	return result, nil
}

func (mfs *MemoryFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	// Ensure parent directory exists
	dir := filepath.Dir(filename)
	if !mfs.dirs[dir] && dir != "." && dir != "/" {
		return &os.PathError{Op: "open", Path: filename, Err: os.ErrNotExist}
	}
	
	// Store a copy to prevent external modification
	mfs.files[filename] = make([]byte, len(data))
	copy(mfs.files[filename], data)
	return nil
}

func (mfs *MemoryFileSystem) MkdirAll(path string, perm os.FileMode) error {
	// Normalize path
	path = filepath.Clean(path)
	
	// Mark all parent directories as existing
	parts := strings.Split(path, string(filepath.Separator))
	currentPath := ""
	
	for i, part := range parts {
		if i == 0 && part == "" {
			currentPath = "/"
		} else if currentPath == "/" {
			currentPath = "/" + part
		} else if currentPath == "" {
			currentPath = part
		} else {
			currentPath = filepath.Join(currentPath, part)
		}
		mfs.dirs[currentPath] = true
	}
	return nil
}

func (mfs *MemoryFileSystem) Remove(name string) error {
	if _, exists := mfs.files[name]; exists {
		delete(mfs.files, name)
		return nil
	}
	if mfs.dirs[name] {
		delete(mfs.dirs, name)
		return nil
	}
	return &os.PathError{Op: "remove", Path: name, Err: os.ErrNotExist}
}

func (mfs *MemoryFileSystem) Rename(oldpath, newpath string) error {
	data, exists := mfs.files[oldpath]
	if !exists {
		return &os.PathError{Op: "rename", Path: oldpath, Err: os.ErrNotExist}
	}
	
	// Ensure target directory exists
	dir := filepath.Dir(newpath)
	if !mfs.dirs[dir] && dir != "." && dir != "/" {
		return &os.PathError{Op: "rename", Path: newpath, Err: os.ErrNotExist}
	}
	
	mfs.files[newpath] = data
	delete(mfs.files, oldpath)
	return nil
}

func (mfs *MemoryFileSystem) Stat(name string) (os.FileInfo, error) {
	if data, exists := mfs.files[name]; exists {
		return &memoryFileInfo{
			name: filepath.Base(name),
			size: int64(len(data)),
		}, nil
	}
	if mfs.dirs[name] {
		return &memoryFileInfo{
			name:  filepath.Base(name),
			isDir: true,
		}, nil
	}
	return nil, &os.PathError{Op: "stat", Path: name, Err: os.ErrNotExist}
}

func (mfs *MemoryFileSystem) CreateTemp(dir, pattern string) (File, error) {
	// Ensure directory exists
	if !mfs.dirs[dir] && dir != "." && dir != "/" {
		return nil, &os.PathError{Op: "createtemp", Path: dir, Err: os.ErrNotExist}
	}
	
	// Generate a unique temporary filename
	name := filepath.Join(dir, fmt.Sprintf("temp_%s_%d", pattern, len(mfs.files)))
	return &memoryFile{name: name, fs: mfs}, nil
}

func (mfs *MemoryFileSystem) Chmod(name string, mode os.FileMode) error {
	// Memory filesystem doesn't need permission handling for testing
	// Just verify the file exists
	if _, exists := mfs.files[name]; !exists && !mfs.dirs[name] {
		return &os.PathError{Op: "chmod", Path: name, Err: os.ErrNotExist}
	}
	return nil
}

// memoryFile implements File interface
type memoryFile struct {
	name   string
	buffer bytes.Buffer
	fs     *MemoryFileSystem
	closed bool
}

func (mf *memoryFile) Write(p []byte) (n int, err error) {
	if mf.closed {
		return 0, os.ErrClosed
	}
	return mf.buffer.Write(p)
}

func (mf *memoryFile) Close() error {
	if mf.closed {
		return nil
	}
	mf.fs.files[mf.name] = mf.buffer.Bytes()
	mf.closed = true
	return nil
}

func (mf *memoryFile) Name() string {
	return mf.name
}

// memoryFileInfo implements os.FileInfo
type memoryFileInfo struct {
	name  string
	size  int64
	isDir bool
}

func (mfi *memoryFileInfo) Name() string       { return mfi.name }
func (mfi *memoryFileInfo) Size() int64        { return mfi.size }
func (mfi *memoryFileInfo) Mode() os.FileMode  { 
	if mfi.isDir {
		return os.ModeDir | 0755
	}
	return 0644 
}
func (mfi *memoryFileInfo) ModTime() time.Time { return time.Now() }
func (mfi *memoryFileInfo) IsDir() bool        { return mfi.isDir }
func (mfi *memoryFileInfo) Sys() interface{}   { return nil }
