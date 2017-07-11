package assetfs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// AssetFileSystem AssetFS based on FileSystem
type AssetFileSystem struct {
	Paths []string
}

// RegisterPath register view paths
func (fs *AssetFileSystem) RegisterPath(pth string) error {
	if _, err := os.Stat(pth); !os.IsNotExist(err) {
		var existing bool
		for _, p := range fs.Paths {
			if p == pth {
				existing = true
				break
			}
		}
		if !existing {
			fs.Paths = append(fs.Paths, pth)
		}
		return nil
	}
	return errors.New("not found")
}

// PrependPath prepend path to view paths
func (fs *AssetFileSystem) PrependPath(pth string) error {
	if _, err := os.Stat(pth); !os.IsNotExist(err) {
		var existing bool
		for _, p := range fs.Paths {
			if p == pth {
				existing = true
				break
			}
		}
		if !existing {
			fs.Paths = append([]string{pth}, fs.Paths...)
		}
		return nil
	}
	return errors.New("not found")
}

// Asset get content with name from assetfs
func (fs *AssetFileSystem) Asset(name string) ([]byte, error) {
	for _, pth := range fs.Paths {
		if _, err := os.Stat(filepath.Join(pth, name)); err == nil {
			return ioutil.ReadFile(filepath.Join(pth, name))
		}
	}
	return []byte{}, fmt.Errorf("%v not found", name)
}

// Glob list matched files from assetfs
func (fs *AssetFileSystem) Glob(pattern string) (matches []string, err error) {
	for _, pth := range fs.Paths {
		if results, err := filepath.Glob(filepath.Join(pth, pattern)); err == nil {
			for _, result := range results {
				matches = append(matches, strings.TrimPrefix(result, pth))
			}
		}
	}
	return
}

// Compile compile assetfs
func (fs *AssetFileSystem) Compile() error {
	return nil
}
