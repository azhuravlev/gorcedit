package main

import (
	"errors"
	"os"
)

// FileInfo returns os information about the file and its existence flag
func FileInfo(name string) (os.FileInfo, bool, error) {
	fi, err := os.Stat(name)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	return fi, err == nil, err
}
