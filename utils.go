package render

import "os"

var root, _ = os.Getwd()

func init() {
	if path := os.Getenv("WEB_ROOT"); path != "" {
		root = path
	}
}

func isExistingDir(pth string) bool {
	if fi, err := os.Stat(pth); err == nil {
		return fi.Mode().IsDir()
	}
	return false
}
