package common

import (
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
}

func NewSyncFileDescribe(rootPath string, repoPath string, clientFileDescribe FileDescribe) []SyncFileDescribe {
	var syncFileDescribe []SyncFileDescribe
	filepath.Walk(rootPath+repoPath, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() || strings.HasPrefix(info.Name(), ".") || strings.HasSuffix(info.Name(), ".md5") {
			return nil
		}

		fileRelativePath := strings.Replace(path, rootPath, "", 1)
		fileMd5Path := path + ".md5"
		fileMd5RelativePath := fileRelativePath + ".md5"
		clientHashBytes, _ := clientFileDescribe.Describe[fileRelativePath]
		hashBytes, _ := ioutil.ReadFile(fileMd5Path)
		if hashBytes != nil && bytes.Equal(hashBytes, clientHashBytes) {
			return nil
		}

		fileContent, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println("readFile:", err)
			return err
		}

		syncFileDescribe = append(syncFileDescribe, SyncFileDescribe{
			Root:       clientFileDescribe.Root,
			Name:       fileRelativePath,
			Content:    fileContent,
		})

		fileMd5Content, err := ioutil.ReadFile(fileMd5Path)
		if err != nil {
			log.Println("readFileMd5:", err)
			return err
		}
		syncFileDescribe = append(syncFileDescribe, SyncFileDescribe{
			Root:       clientFileDescribe.Root,
			Name:       fileMd5RelativePath,
			Content:    fileMd5Content,
		})

		return nil
	})
	return syncFileDescribe
}
