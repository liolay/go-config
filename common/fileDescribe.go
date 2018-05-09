package common

import (
	"path/filepath"
	"os"
	"log"
	"strings"
	"io/ioutil"
)

type FileDescribe struct {
	Root     string
	Describe map[string][]byte
}

func NewFileDescribe(rootPath string) *FileDescribe {
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		os.MkdirAll(rootPath, os.ModePerm)
	}

	describe := make(map[string][]byte)
	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}

		if err != nil {
			log.Println(err)
			return err
		}

		if hashBytes, err := ioutil.ReadFile(path + ".md5"); err == nil {
			describe[strings.Replace(path, rootPath, "", 1)] = hashBytes
		}
		return nil
	})
	return &FileDescribe{
		Root:     rootPath,
		Describe: describe,
	}
}
