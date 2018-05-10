package main

import (
	"go-config/util"
	"flag"
	"github.com/gorilla/websocket"
	"strconv"
	"io/ioutil"
	"net/http"
	"log"
	"go-config/common"
	"io"
	"sync"
	"encoding/json"
	"strings"
	"os"
	"gopkg.in/src-d/go-git.v4"
)

const (
	defaultConfig      = "configServer.yml"
	appPlaceholder     = "${app}"
	profilePlaceholder = "${profile}"
)

var (
	config   *util.ServerConfig
	upgrader = websocket.Upgrader{} // use default options

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
	config, err = util.ParseServerConfig(fileContent)
	if err != nil {
		panic(err)
	}
}

var locks = new(sync.Map)

func buildRepoUrl(repo string, app string, profile string) string {
	repo = strings.Replace(repo, appPlaceholder, app, -1)
	return strings.Replace(repo, profilePlaceholder, profile, -1)
}

func repoModel(repo string) util.RepoModel {
	var repoModel util.RepoModel = 0
	if strings.Contains(repo, appPlaceholder) {
		repoModel++
	}
	if strings.Contains(repo, profilePlaceholder) {
		repoModel++
	}
	return repoModel
}

func patterMatch(pattern string, str string) bool {
	if len(strings.Trim(pattern, "*")) != 0 {
		if strings.HasPrefix(pattern, "*") && strings.HasSuffix(str, strings.Trim(pattern, "*")) {
			return false
		}
		if strings.HasSuffix(pattern, "*") && strings.HasPrefix(str, strings.Trim(pattern, "*")) {
			return false
		}
	}
	return true
}

func findRoute(appNode util.AppNode) *util.RRoute {
	for _, route := range config.Route {
		for _, p := range route.Pattern {
			if !strings.Contains(p, "/") {
				p = p + "/*"
			}

			pattern := strings.Split(p, "/")
			app := pattern[0]
			profile := pattern[1]

			if !patterMatch(app, appNode.Name) {
				continue
			}
			if !patterMatch(profile, appNode.Profile) {
				continue
			}
			route.Repo = buildRepoUrl(route.Repo, appNode.Name, appNode.Profile)
			route.Model = repoModel(route.Repo)
			return &route
		}
	}
	return &util.RRoute{Repo: buildRepoUrl(config.DefaultRepo, appNode.Name, appNode.Profile), Model: repoModel(config.DefaultRepo)}
}

func cloneRepo(repoUrl string) *git.Repository {
	localRepoPath := strings.TrimSuffix(repoUrl, ".git")
	localRepoPath = config.HomePath + localRepoPath[strings.LastIndex(localRepoPath, "/"):]

	if _, err := os.Stat(localRepoPath); os.IsNotExist(err) {
		lock, _ := locks.LoadOrStore(localRepoPath, new(sync.Mutex))
		mutex := lock.(*sync.Mutex)
		mutex.Lock()
		defer mutex.Unlock()
		if _, err := os.Stat(localRepoPath); err == nil || !os.IsNotExist(err) {
			return openLocalRepo(localRepoPath)
		}

		return clone(config.HomePath, []byte(config.SshKey), repoUrl)
	}
	return openLocalRepo(localRepoPath)
}

func syncFile(writer http.ResponseWriter, request *http.Request) {
	connection, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer connection.Close()

	done := make(chan struct{})
	go func() {
		for {
			var message common.Message
			if err := connection.ReadJSON(&message); err != nil && err != io.ErrUnexpectedEOF {
				log.Println("client down", err)
				return
			}

			if common.ClientConnect == message.MessageType {
				appNodeConfig := make([]util.AppNode, 5)
				if err := json.Unmarshal(message.Data, appNodeConfig); err != nil {
					connection.WriteJSON(common.NewClientConnectReplyMessage([]byte(err.Error())))
				}

				for _, clientApp := range appNodeConfig {
					route := findRoute(clientApp)
					repo := cloneRepo(route.Repo)
					if repo == nil {
						log.Fatalln("cant find repo from disk,check you repostory url")
						continue
					}
					//todo 开始判断仓库应用模式，推送文件

					checkOut(repo, clientApp.Label)
					switch {
					case route.Model == util.OnlyOne:
					case route.Model == util.AppOne:
					case route.Model == util.AppProfileOne:
					}
				}
			} else {
				panic("unsupported message type")
			}
		}
	}()

	<-done
	//for {
	//	select {
	//	case <-done:
	//		return
	//	case file := <-fileChangeSignal:
	//		err = connection.WriteJSON(file)
	//	}
	//
	//}
}

func refresh(writer http.ResponseWriter, request *http.Request) {

}

func main() {
	//watcher, changeSignal := util.WatchFile(homePath + repo)
	//defer watcher.Close()
	//go func() {
	//	for {
	//		event := <-changeSignal
	//		hashFile(event.Name, func(file string, newHash string) {
	//			hashFileName := file + ".md5"
	//			if oldHash, _ := ioutil.ReadFile(hashFileName); string(oldHash) != newHash {
	//
	//				fileContent, err := ioutil.ReadFile(file)
	//				if err != nil {
	//					log.Println("file cant be sync to client", err)
	//					return
	//				}
	//
	//				fileChangeSignal <- []common.SyncFileDescribe{
	//					{
	//						Name:    strings.Replace(file, homePath, "", 1),
	//						Content: fileContent,
	//					},
	//					{
	//						Name:    strings.Replace(hashFileName, homePath, "", 1),
	//						Content: []byte(newHash),
	//					},
	//				}
	//			}
	//		})
	//	}
	//}()

	http.HandleFunc("/sync", syncFile)
	http.HandleFunc("/refresh", refresh)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))
}
