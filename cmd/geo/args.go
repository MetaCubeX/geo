package main

import (
	"errors"
	"os"
	"path"
	"strings"
)

var (
	dbType string
	dbPath string
)

func find() (paths []string, err error) {
	files, err := os.ReadDir(workingDir)
	if err != nil {
		return nil, os.ErrInvalid
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		} else {
			name := f.Name()
			if strings.HasSuffix(name, ".mmdb") ||
				strings.HasSuffix(name, ".db") ||
				strings.HasSuffix(name, ".dat") ||
				strings.HasSuffix(name, ".metadb") {
				paths = append(paths, path.Join(workingDir, name))
			}
		}
	}
	if len(paths) == 0 {
		err = errors.New("failed to find a geo database, please specific one through argument")
	}
	return
}
