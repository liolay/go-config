package main

import (
	"time"
	"os"
	"flag"
	"log"
	"net/url"
	"github.com/gorilla/websocket"
	"go-config/common"
	"io"
	"go-config/util"
	"io/ioutil"
	"os/signal"
	"encoding/json"
)

const defaultConfig = "configClient.yml"

var (
	config     *util.ClientConfig
	connection *websocket.Conn
	ticker     *time.Ticker
)

func init() {
	flag.Parse()
	configFile := flag.Arg(0)
	if configFile == "" {
		log.Printf("no config file point,use default %s locate...", defaultConfig)
		configFile = defaultConfig
	}
	fileContent, err := ioutil.ReadFile(configFile);
	if err != nil {
		panic(err)
	}
	config, err = util.ParseClientConfig(fileContent)
	if err != nil {
		panic(err)
	}
}

func createConnection() *websocket.Conn {
	serverUrl := url.URL{Scheme: "ws", Host: config.Server, Path: "/sync"}
	log.Printf("connecting to %s", serverUrl.String())
	connection, _, err := websocket.DefaultDialer.Dial(serverUrl.String(), nil)
	if err != nil {
		log.Println("can't dial", err)
		log.Printf("retry after %d seconds...", config.Tick)
		ticker = time.NewTicker(time.Duration(config.Tick) * time.Second)
		<-ticker.C
		return createConnection()
	}
	return connection
}

func readMessage() {
	connection = createConnection()
	go func() {
		defer readMessage()
		for {
			var message common.Message
			if err := connection.ReadJSON(&message); err != nil && err != io.ErrUnexpectedEOF {
				log.Println("lost connection:", err)
				log.Printf("reconnect after %d seconds...", config.Tick)
				ticker = time.NewTicker(time.Duration(config.Tick) * time.Second)
				<-ticker.C
				return
			}

			if message.MessageType == common.ClientConnectReply {
				err := message.Data
				if err != nil {
					log.Fatalln("client config file cant be parsed by server")
					panic(string(err))
				}
				continue
			}

			if message.MessageType == common.ServerPushFile {
				files := make([]common.ServerPushedFile, 5)
				json.Unmarshal(message.Data, files)

				for _, file := range files {
					for _, app := range config.App {
						if file.App == app.Name {
							for _, home := range app.HomePath {
								fileFolder := home + "/" + file.App
								if _, err := os.Stat(fileFolder); os.IsNotExist(err) {
									os.MkdirAll(fileFolder, os.ModePerm)
								}
							}
							break
						}
					}
				}
				continue
			}

			//sync(syncFd)
		}
	}()

	configBytes, err := json.Marshal(config.App)

	if err != nil {
		panic(err)
	}
	connection.WriteJSON(common.NewClientConnectMessage(configBytes))
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	done := make(chan struct{})
	go func() {
		<-interrupt
		if connection != nil {
			_ = connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			time.Sleep(time.Second)
		}
		panic("process be killed")
	}()

	readMessage()

	<-done
}

//
//func sync(syncFileDescribes []common.SyncFileDescribe) {
//	for _, fd := range syncFileDescribes {
//		if fd.Root != "" {
//			create(fd.Root+fd.Name, fd.Content)
//			continue
//		}
//
//		for _, root := range config.HomePath {
//			create(root+fd.Name, fd.Content)
//		}
//	}
//}
//
//func create(filePath string, content []byte) {
//	createParentDir(filePath)
//	ioutil.WriteFile(filePath, content, os.ModePerm)
//}
//
//func createParentDir(filePath string) {
//	parentPath := string(filePath[:strings.LastIndex(filePath, "/")])
//	if _, err := os.Stat(parentPath); os.IsNotExist(err) {
//		if err = os.MkdirAll(parentPath, os.ModePerm); err != nil {
//			panic(err)
//		}
//	}
//}
