package common

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

func NewSyncFileDescribe(rootPath string, repoPath string, clientFileDescribe FileDescribe) []SyncFileDescribe {
	var syncFileDescribe []SyncFileDescribe
	filepath.Walk(rootPath+repoPath, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}

		fileRelativePath := strings.Replace(path, rootPath, "", 1)
		hashBytes, _ := ioutil.ReadFile(path + ".md5")

		clientHashBytes, _ := clientFileDescribe.Describe[fileRelativePath]

		if hashBytes != nil && bytes.Equal(hashBytes, clientHashBytes) {
			return nil
		}
		println(path)
		println("hashBytes:",hashBytes)
		println("clientHashBytes:",clientHashBytes)
		fileContent, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println("readFile:", err)
			return err
		}

		syncFileDescribe = append(syncFileDescribe, SyncFileDescribe{
			Root:       clientFileDescribe.Root,
			Name:       fileRelativePath,
			Content:    fileContent,
			UpdateTime: info.ModTime(),
		})

		return nil
	})
	return syncFileDescribe
}

