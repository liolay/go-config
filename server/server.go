package main

import (
	"time"
	"net/http"
	"encoding/json"
	"os"
	"log"
	"path/filepath"
	"io/ioutil"
)

type Param struct {
	folder string
	stat   map[string]time.Time
}

type ConfigInfo struct {
	Name       string
	Content    []byte
	UpdateTime time.Time
}

func main() {
	http.HandleFunc("/refresh", func(writer http.ResponseWriter, request *http.Request) {
		//body := request.Body
		//if body != nil {
		//	if jsonString, err := ioutil.ReadAll(body); err == nil {
		//		var param []Param
		//		refreshPo := json.Unmarshal(jsonString, &param)
		//		fmt.Printf("%v", refreshPo)
		//	}
		//}

		writer.Header().Add("Content-Type", "application/json;charset=utf-8");
		writer.WriteHeader(200)

		configInfos := []ConfigInfo{}
		serverConfigFolder := os.Getenv("HOME") + "/config/"
		filepath.Walk(serverConfigFolder, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			bytes, _ := ioutil.ReadFile(path)
			configInfos = append(configInfos, ConfigInfo{
				Name:       path,
				Content:    bytes,
				UpdateTime: info.ModTime(),
			})

			return nil
		})

		bytes, _ := json.Marshal(configInfos)
		os.Stdout.Write(bytes)
		writer.Write(bytes)
	})

	log.Fatal(http.ListenAndServe(":6973", nil))

}
