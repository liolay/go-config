package main

import (
	"time"
	"strings"
	"encoding/json"
	"os"
	"path/filepath"
	"flag"
	"log"
	"net/url"
	"github.com/gorilla/websocket"
	"os/signal"
)

type ConfigInfo struct {
	Name       string
	Content    []byte
	UpdateTime time.Time
}

var (
	homePaths []string
	server    string
	bootstrap string
	tick      int
	force     bool
)

func init() {
	commandHomePath := flag.String("home", os.Getenv("HOME")+"/.diamond-client", "dirs to store config files from server,multi dirs use comma split")
	commandServer := flag.String("server", "http://localhost:5337", "server to fetch config files")
	commandBootstrap := flag.String("server", "http://localhost:5337", "server to fetch config files")
	commandTick := flag.Int("tick", 15, "interval second between two fetch request")
	commandForce := flag.Bool("force", false, "whether or not force pull all config files regardless changes")
	flag.Parse()

	homePaths = strings.Split(*commandHomePath, ",")
	server = *commandServer
	bootstrap = *commandBootstrap
	tick = *commandTick
	force = *commandForce
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)

	serverUrl := url.URL{Scheme: "ws", Host: bootstrap, Path: "/echo"}
	log.Printf("connecting to %s", serverUrl.String())

	connection, _, err := websocket.DefaultDialer.Dial(serverUrl.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer connection.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := connection.ReadMessage()
			if err != nil {
				log.Println(err)
				continue
			}

			var configInfos []ConfigInfo
			err = json.Unmarshal(message, &configInfos)
			if err != nil {
				log.Println(err)
				continue
			}

			for _, homePath := range homePaths {
				createFile(configInfos, homePath)
			}
		}
	}()

	ticker := time.NewTicker(time.Duration(tick) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			for _, homePath := range homePaths {
				fileStateBytes, err := newFDWrapper(homePath).Json()
				if err != nil {
					log.Println(err)
				}

				err = connection.WriteMessage(websocket.TextMessage, fileStateBytes)
				if err != nil {
					log.Println("write:", err)
					continue
				}
			}

		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println(err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}

	}
}


func createFile(configInfos []ConfigInfo, configRoot string) {
	for _, config := range configInfos {
		filePath := configRoot + config.Name
		createParentDir(filePath)

		configFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			log.Println(err)
			continue
		}
		defer configFile.Close()
		configFile.Write(config.Content)
		os.Chtimes(filePath, config.UpdateTime, config.UpdateTime)
	}
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
