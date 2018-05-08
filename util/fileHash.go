package util

import (
	"path/filepath"
	"os"
	"log"
	"crypto/md5"
	"io"
	"encoding/hex"
)

func HashFile(path string, filter func(os.FileInfo) bool, hashFunc func(string, string)) map[string]string {
	fileHash := make(map[string]string)
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}

		if filter == nil || !filter(info) {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		hash := md5.New()
		if _, err := io.Copy(hash, file); err != nil {
			log.Fatal(err)
		}

		sum := hash.Sum(nil)
		fileHash[path] = hex.EncodeToString(sum)
		hashFunc(path, fileHash[path])
		return nil
	})
	return fileHash
}
