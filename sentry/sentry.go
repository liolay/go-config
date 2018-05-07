package sentry

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"net/http"
)

func Watch(path string) (*fsnotify.Watcher, chan fsnotify.Op) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	//defer watcher.Close()

	signal := make(chan fsnotify.Op)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)

				switch {
				case event.Op&fsnotify.Write == fsnotify.Write:
					signal <- fsnotify.Write
				case event.Op&fsnotify.Create == fsnotify.Create:
					signal <- fsnotify.Create
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					signal <- fsnotify.Remove
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	return watcher, signal
}

func main() {
	watcher, signal := Watch("/Users/liolay/config-repo")
	go func() {
		for {
			<-signal
			log.Println("=========")
		}
	}()
	log.Fatal(http.ListenAndServe(":8080", nil))
	watcher.Close()
}
