package main

import (
	"time"
	"encoding/json"
	"path/filepath"
	"os"
	"log"
	"strings"
)

type FDWrapper struct {
	Home string
	Fds  map[string]time.Time
}

func (fd *FDWrapper) Json() ([]byte, error) {
	return json.Marshal(fd)
}
func newFDWrapper(rootPath string) *FDWrapper {
	fds := make(map[string]time.Time)
	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}

		if err != nil {
			log.Println(err)
			return err
		}

		fds[strings.Replace(path, rootPath, "", 1)] = info.ModTime()
		return nil
	})
	return &FDWrapper{
		Home: rootPath,
		Fds:  fds,
	}
}
