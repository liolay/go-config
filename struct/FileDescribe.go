package _struct

import (
	"encoding/json"
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

func (fd *FileDescribe) Json() ([]byte, error) {
	return json.Marshal(fd)
}

func ParseFileDescribe(jsonByte []byte) (FileDescribe, error) {
	var fd FileDescribe
	err := json.Unmarshal(jsonByte, &fd)
	return fd, err
}

func NewFileDescribe(rootPath string) *FileDescribe {
	describe := make(map[string][]byte)
	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}

		if err != nil {
			log.Println(err)
			return err
		}

		hashBytes, err := ioutil.ReadFile(path + ".md5")
		if err == nil {
			describe[strings.Replace(path, rootPath, "", 1)] = hashBytes
		}
		return nil
	})
	return &FileDescribe{
		Root:     rootPath,
		Describe: describe,
	}
}



