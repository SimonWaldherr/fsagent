package main

import (
	"os"
	"time"
)

// FileLastModified returns the Time a file was last modified.
func FileLastModified(filename string) (*time.Time, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	statInfo, _ := f.Stat()
	modTime := statInfo.ModTime()
	return &modTime, nil
}
