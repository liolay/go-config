package main

import (
	"time"
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"
	"os"
	"fmt"
	"path/filepath"
	"bytes"
	"flag"
)

type ConfigInfo struct {
	Name       string
	Content    []byte
	UpdateTime time.Time
}

var (
	homePaths []string
	server    string
	tick      int
	force     bool
)

func init() {
	commandHomePath := flag.String("home", os.Getenv("HOME")+"/.diamond-client", "dirs to store config files from server,multi dirs use comma split")
	commandServer := flag.String("server", "http://localhost:5337", "server to fetch config files")
	commandTick := flag.Int("tick", 15, "interval second between two fetch request")
	commandForce := flag.Bool("force", false, "whether or not force pull all config files regardless changes")
	flag.Parse()

	homePaths = strings.Split(*commandHomePath, ",")
	server = *commandServer
	tick = *commandTick
	force = *commandForce
}

func main() {
	for range time.Tick(time.Duration(tick) * time.Second) {
		for _, homePath := range homePaths {
			configInfos, err := fetch(homePath)
			if err != nil {
				fmt.Println(err)
				continue
			}

			createFile(configInfos, homePath)
		}
	}

}

func createFile(configInfos []ConfigInfo, configRoot string) {
	for _, config := range configInfos {
		filePath := configRoot + config.Name
		createParentDir(filePath)

		configFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer configFile.Close()
		configFile.Write(config.Content)
		os.Chtimes(filePath, config.UpdateTime, config.UpdateTime)
	}
}

func fetch(configRoot string) ([]ConfigInfo, error) {
	buffer := bytes.NewBufferString("")
	if !force {
		fileStateBytes, e := json.Marshal(fds(configRoot))
		if e != nil {
			fmt.Println(e)
		}
		buffer = bytes.NewBuffer(fileStateBytes)
	}
	resp, err := http.Post(server+"/sync", "application/json;charset=utf-8", buffer)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var configInfos []ConfigInfo
	err = json.Unmarshal(responseBody, &configInfos)
	return configInfos, err
}

func fds(rootPath string) map[string]time.Time {
	fds := make(map[string]time.Time)

	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}

		if err != nil {
			fmt.Print(err)
			return err
		}

		fds[strings.Replace(path, rootPath, "", 1)] = info.ModTime()

		return nil
	})

	return fds
}

func createParentDir(filePath string) {
	parentPath := string([]rune(filePath)[:strings.LastIndex(filePath, "/")])
	if _, err := os.Stat(parentPath); os.IsNotExist(err) {
		err = os.MkdirAll(parentPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}
