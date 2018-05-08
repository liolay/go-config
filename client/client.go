package main

import (
	"time"
	"strings"
	"os"
	"flag"
	"log"
	"net/url"
	"github.com/gorilla/websocket"
	"go-config/struct"
	"os/signal"
	"io"
)

var (
	homePaths  []string
	server     string
	tick       int
	connection *websocket.Conn
)

func init() {
	commandHomePath := flag.String("home", os.Getenv("HOME")+"/haha", "dirs to store config files from server,multi dirs use comma split")
	commandServer := flag.String("server", "localhost:5337", "server to fetch config files")
	commandTick := flag.Int("tick", 15, "interval second while reconnect to server ")
	flag.Parse()
	homePaths = strings.Split(*commandHomePath, ",")
	server = *commandServer
	tick = *commandTick
}

func createConnection() *websocket.Conn {
	serverUrl := url.URL{Scheme: "ws", Host: server, Path: "/sync"}
	log.Printf("connecting to %s", serverUrl.String())
	connection, _, err := websocket.DefaultDialer.Dial(serverUrl.String(), nil)
	if err != nil {
		log.Println("can't dial", err)
		log.Printf("retry after %d seconds...", tick)
		time.Sleep(time.Duration(tick) * time.Second)
		return createConnection()
	}
	return connection
}

func readMessage() {
	connection = createConnection()
	go func() {
		defer readMessage()
		for {
			var syncFd []_struct.SyncFileDescribe
			err := connection.ReadJSON(&syncFd)
			if err != nil && err != io.ErrUnexpectedEOF {
				log.Println("lost connection:", err)
				log.Printf("reconnect after %d seconds...", tick)
				time.Sleep(time.Duration(tick) * time.Second)
				return
			}
			sync(syncFd)
		}
	}()

	for _, root := range homePaths {
		connection.WriteJSON(_struct.NewFileDescribe(root))
	}
}
func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)

	readMessage()

	for {
		select {
		case <-interrupt:
			_ = connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			time.Sleep(time.Second)
			return
		}
	}

}

func sync(syncFileDescribes []_struct.SyncFileDescribe) {
	for _, fd := range syncFileDescribes {
		if fd.Root != "" {
			create(fd.Root+fd.Name, fd.Content, fd.UpdateTime)
			continue
		}

		for _, root := range homePaths {
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
	parentPath := string([]rune(filePath)[:strings.LastIndex(filePath, "/")])
	if _, err := os.Stat(parentPath); os.IsNotExist(err) {
		err = os.MkdirAll(parentPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}
