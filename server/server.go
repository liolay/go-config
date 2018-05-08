package main

import (
	"go-config/util"
	"flag"
	"os"
	"strings"
	"github.com/gorilla/websocket"
	"strconv"
	"io/ioutil"
	"net/http"
	"log"
	"go-config/struct"
	"io"
)

var (
	homePath string
	repo     string
	port     string
	upgrader = websocket.Upgrader{} // use default options

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

func hashFile(path string, beforeWrite func(string, string)) {
	util.HashFile(path, func(fileInfo os.FileInfo) bool {
		return !(strings.HasSuffix(fileInfo.Name(), ".md5") || strings.HasPrefix(fileInfo.Name(), "."))
	}, func(file string, hash string) {
		if beforeWrite != nil {
			beforeWrite(file, hash)
		}
		ioutil.WriteFile(file+".md5", []byte(hash), os.ModePerm)
	})
}

func main() {
	fileChangeSignal := make(chan []_struct.SyncFileDescribe, 10)

	hashFile(homePath+repo, nil)

	watcher, changeSignal := util.WatchFile(homePath + repo)
	defer watcher.Close()
	go func() {
		for {
			event := <-changeSignal
			hashFile(event.Name, func(file string, newHash string) {
				hashFileName := file + ".md5"
				oldHash, e := ioutil.ReadFile(hashFileName)
				if e == nil && string(oldHash) != newHash {
					fileContent, err := ioutil.ReadFile(file)
					if err != nil {
						log.Println("file cant be sync to client", err)
					}
					fileChangeSignal <- []_struct.SyncFileDescribe{
						{
							Name:    strings.Replace(file, homePath, "", 1),
							Content: fileContent,
						},
						{
							Name:    strings.Replace(hashFileName, homePath, "", 1),
							Content: []byte(newHash),
						},
					}
				}
			})
		}
	}()

	http.HandleFunc("/sync", func(writer http.ResponseWriter, request *http.Request) {
		connection, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			log.Print("upgrade error:", err)
			return
		}
		defer connection.Close()

		done := make(chan struct{})
		go func() {
			defer close(done)
			for {
				var fd _struct.FileDescribe
				err := connection.ReadJSON(&fd)
				if err!=nil && err != io.ErrUnexpectedEOF {
					log.Println("client down", err)
					return
				}

				newSyncFileDescribe := _struct.NewSyncFileDescribe(homePath, repo, fd)
				err = connection.WriteJSON(newSyncFileDescribe)
			}
		}()

		for {
			select {
			case <-done:
				return
			case file := <-fileChangeSignal:
				err = connection.WriteJSON(file)
			}

		}
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))

}
