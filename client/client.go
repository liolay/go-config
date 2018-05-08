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
	homePaths []string
	server    string
	tick      int
)

func init() {
	commandHomePath := flag.String("home", os.Getenv("HOME")+"/haha", "dirs to store config files from server,multi dirs use comma split")
	commandServer := flag.String("server", "localhost:5337", "server to fetch config files")
	commandTick := flag.Int("tick", 15, "interval second between two fetch request")
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
		log.Println("can't dial,retry after 3 seconds:", err)
		time.Sleep(3 * time.Second)
		return createConnection()
	}
	return connection
}

var con *websocket.Conn
func readMessage() *websocket.Conn {
	connection := createConnection()
	con = connection
	go func() {
		defer readMessage()
		for {
			var syncFd []_struct.SyncFileDescribe
			err := connection.ReadJSON(&syncFd)
			if err != nil && err != io.ErrUnexpectedEOF {
				log.Println("server is down,reconnect after 2 seconds:", err)
				time.Sleep(2 * time.Second)
				return
			}
			sync(syncFd)
		}
	}()

	for _, root := range homePaths {
		connection.WriteJSON(_struct.NewFileDescribe(root))
	}

	return connection
}
func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	//done := make(chan struct{})
	connection := readMessage()
	//connection := createConnection()
	//defer connection.Close()
	//
	//done := make(chan struct{})
	//go func() {
	//	defer main()
	//
	//	defer close(done)
	//	for {
	//		var syncFd []_struct.SyncFileDescribe
	//		err := connection.ReadJSON(&syncFd)
	//		if err != nil && err != io.ErrUnexpectedEOF {
	//			log.Println("server is down,reconnect after 2 seconds:", err)
	//			time.Sleep(2 * time.Second)
	//			return
	//		}
	//		sync(syncFd)
	//	}
	//
	//}()
	//
	//for _, root := range homePaths {
	//	connection.WriteJSON(_struct.NewFileDescribe(root))
	//}

	for {

		select {
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			println(connection == con)
			err := connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("connection.WriteMessage(websocket.CloseMessage:", err)
			}
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
