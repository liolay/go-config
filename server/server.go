package main

import (
	"time"
	"net/http"
	"encoding/json"
	"os"
	"log"
	"path/filepath"
	"io/ioutil"
	"strings"
	"flag"
	"bufio"
	"strconv"
)

type ConfigInfo struct {
	Name       string
	Content    []byte
	UpdateTime time.Time
}

var (
	homePath string
	repo     string
	port     string
)

func init() {
	commandHomePath := flag.String("home", os.Getenv("HOME"), "config server home path")
	commandRepo := flag.String("repo", "config-repo", "config server repo")
	commandPort := flag.Int("port", 5337, "config server running port")
	flag.Parse()

	homePath = strings.TrimSuffix(*commandHomePath, "/")
	repo = "/" + strings.Trim(*commandRepo, "/")
	port = strconv.Itoa(*commandPort)
}

func main() {
	http.HandleFunc("/sync", func(writer http.ResponseWriter, request *http.Request) {
		clientFDS := parseClientFDS(request)
		var configInfos []ConfigInfo

		filepath.Walk(homePath+repo, func(path string, info os.FileInfo, err error) error {
			if info == nil || info.IsDir() {
				return nil
			}

			filePath := strings.Replace(path, homePath, "", 1)
			if !changed(clientFDS, filePath, info.ModTime()) {
				return nil
			}

			fileContent, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			configInfos = append(configInfos, ConfigInfo{
				Name:       filePath,
				Content:    fileContent,
				UpdateTime: info.ModTime(),
			})

			return nil
		})

		repoBytes, _ := json.Marshal(configInfos)

		writer.Header().Add("Content-Type", "application/json;charset=utf-8")
		writer.WriteHeader(200)
		writer.Write(repoBytes)
	})

	console := bufio.NewWriter(os.Stdout)
	console.WriteString("server stated at port:" + port + ",home:" + homePath + ",repo:" + repo)
	console.Flush()

	log.Fatal(http.ListenAndServe(":"+port, nil))

}

func parseClientFDS(request *http.Request) map[string]time.Time {
	var requestFDS map[string]time.Time
	if request.Body != nil {
		if jsonString, err := ioutil.ReadAll(request.Body); err == nil {
			json.Unmarshal(jsonString, &requestFDS)
		}
	}
	return requestFDS
}

func changed(clientFDS map[string]time.Time, filepath string, modTime time.Time) bool {
	clientModTime, present := clientFDS[filepath]
	if present {
		return !clientModTime.Equal(modTime)
	}
	return true
}
