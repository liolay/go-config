package main

import (
	"time"
	"strings"
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
			var syncFd []common.SyncFileDescribe
			if err := connection.ReadJSON(&syncFd); err != nil && err != io.ErrUnexpectedEOF {
				log.Println("lost connection:", err)
				log.Printf("reconnect after %d seconds...", config.Tick)
				ticker = time.NewTicker(time.Duration(config.Tick) * time.Second)
				<-ticker.C
				return
			}
			sync(syncFd)
		}
	}()

	for _, root := range config.HomePath {
		connection.WriteJSON(common.NewFileDescribe(root))
	}
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

func sync(syncFileDescribes []common.SyncFileDescribe) {
	for _, fd := range syncFileDescribes {
		log.Printf("%v",fd)
		if fd.Root != "" {
			create(fd.Root+fd.Name, fd.Content, fd.UpdateTime)
			continue
		}

		for _, root := range config.HomePath {
			create(root+fd.Name, fd.Content, fd.UpdateTime)
		}
	}
}

func create(filePath string, content []byte, time time.Time) {
	createParentDir(filePath)

	configFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Println(err)
		return
	}
	defer configFile.Close()
	configFile.Write(content)
	os.Chtimes(filePath, time, time)
}

func createParentDir(filePath string) {
	parentPath := string(filePath[:strings.LastIndex(filePath, "/")])
	if _, err := os.Stat(parentPath); os.IsNotExist(err) {
		if err = os.MkdirAll(parentPath, os.ModePerm); err != nil {
			panic(err)
		}
	}
}
