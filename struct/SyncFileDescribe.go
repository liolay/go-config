package _struct

import (
	"time"
	"path/filepath"
	"os"
	"strings"
	"io/ioutil"
	"log"
	"bytes"
)

type SyncFileDescribe struct {
	Root       string
	Name       string
	Content    []byte
	UpdateTime time.Time
}

func NewSyncFileDescribe(rootPath string, repoPath string, clientFileDescribe FileDescribe) *[]SyncFileDescribe {
	var syncFileDescribe []SyncFileDescribe
	filepath.Walk(rootPath+repoPath, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}

		fileRelativePath := strings.Replace(path, rootPath, "", 1)
		hashBytes, _ := ioutil.ReadFile(fileRelativePath + ".md5")

		clientHashBytes, _ := clientFileDescribe.Describe[fileRelativePath]
		if hashBytes != nil || bytes.Equal(hashBytes, clientHashBytes) {
			return nil
		}

		fileContent, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println("readFile:", err)
			return err
		}

		syncFileDescribe = append(syncFileDescribe, SyncFileDescribe{
			Root:       rootPath,
			Name:       fileRelativePath,
			Content:    fileContent,
			UpdateTime: info.ModTime(),
		})

		return nil
	})
	return &syncFileDescribe
}

