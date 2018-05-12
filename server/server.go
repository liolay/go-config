package main

import (
	"go-config/util"
	"flag"
	"github.com/gorilla/websocket"
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
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"path/filepath"
	"github.com/julienschmidt/httprouter"
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

func patternMatch(pattern string, strs ...string) bool {
	if len(strings.Trim(pattern, "*")) != 0 {
		for _, str := range strs {
			if strings.HasPrefix(pattern, "*") {
				if strings.HasSuffix(str, strings.Trim(pattern, "*")) {
					return true
				}
				continue
			}

			if strings.HasPrefix(str, strings.Trim(pattern, "*")) {
				return true
			}
		}
		return false
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

			if !patternMatch(app, appNode.Name) {
				continue
			}
			if !patternMatch(profile, strings.Split(appNode.Profile, ",")...) {
				continue
			}
			route.Repo = buildRepoUrl(route.Repo, appNode.Name, appNode.Profile)
			route.Model = repoModel(route.Repo)
			return &route
		}
	}
	return &util.RRoute{Repo: buildRepoUrl(config.DefaultRepo, appNode.Name, appNode.Profile), Model: repoModel(config.DefaultRepo)}
}

func buildLocalRepoPath(repoUrl string) string {
	localRepoPath := strings.TrimSuffix(repoUrl, ".git")
	return config.HomePath + localRepoPath[strings.LastIndex(localRepoPath, "/"):]
}

func getRepo(repoUrl string, localRepoPath string) *git.Repository {
	if _, err := os.Stat(localRepoPath); os.IsNotExist(err) {
		lock, _ := locks.LoadOrStore(localRepoPath, new(sync.Mutex))
		mutex := lock.(*sync.Mutex)
		mutex.Lock()
		defer mutex.Unlock()

		if _, err := os.Stat(localRepoPath); err == nil {
			return util.OpenLocalRepo(localRepoPath)
		}

		return util.Clone(localRepoPath, []byte(config.SshKey), repoUrl)
	}
	return util.OpenLocalRepo(localRepoPath)
}

func readProfileFile(file *object.File, profile string) *common.ServerPushedFile {
	for _, prof := range strings.Split(profile, ",") {
		if strings.Index(file.Name, "-") != -1 && prof != "" && !strings.HasSuffix(file.Name, "-"+prof+filepath.Ext(file.Name)) {
			continue
		}
		reader, err := file.Reader()
		if err != nil {
			log.Println(err)
			return nil
		}

		bytes, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Println(err)
			return nil
		}

		return &common.ServerPushedFile{
			Name:    file.Name,
			Content: bytes,
		}
	}
	return nil
}

func findConfigFiles(repo *git.Repository, model util.RepoModel, app string, profile string, label string) []common.ServerPushedFile {
	var files []common.ServerPushedFile

	iterator := util.FileIterator(repo, label)
	if iterator == nil {
		return files
	}

	switch {
	case model == util.OnlyOne:
		iterator.ForEach(func(file *object.File) error {
			if strings.Contains(file.Name, "/") {

				if !strings.HasPrefix(file.Name, app+"/") {
					return nil
				}

				if profileFile := readProfileFile(file, profile); profileFile != nil {

					profileFile.App = app
					files = append(files, *profileFile)
				}

				return nil
			}

			if profileFile := readProfileFile(file, profile); profileFile != nil {
				profileFile.App = app
				files = append(files, *profileFile)
			}
			return nil
		})
	case model == util.AppOne:
		iterator.ForEach(func(file *object.File) error {
			if profileFile := readProfileFile(file, profile); profileFile != nil {
				profileFile.App = app
				files = append(files, *profileFile)
			}
			return nil
		})
	case model == util.AppProfileOne:
		iterator.ForEach(func(file *object.File) error {
			if profileFile := readProfileFile(file, ""); profileFile != nil {
				profileFile.App = app
				files = append(files, *profileFile)
			}
			return nil
		})
	default:
		log.Printf("unsupported model '%d'", model)
	}
	return files
}

func syncSingleAppFiles(connection *websocket.Conn, clientApp util.AppNode) {
	route := findRoute(clientApp)
	log.Printf("find route[app:%s,profile:%s,label:%s]:%s", clientApp.Name, clientApp.Profile, clientApp.Label, route.Repo)
	repo := getRepo(route.Repo, buildLocalRepoPath(route.Repo))
	if repo == nil {
		log.Println("cant find repo from disk,check you repostory url")
		return
	}

	files := findConfigFiles(repo, route.Model, clientApp.Name, clientApp.Profile, clientApp.Label)
	if files != nil {
		bytes, err := json.Marshal(files)
		if err != nil {
			log.Println(err)
			return
		}

		connection.WriteJSON(common.NewServerPushFileMessage(bytes))
	} else {
		log.Printf("no files for app:%s,profile:%s,label:%s", clientApp.Name, clientApp.Profile, clientApp.Label)
	}
}

var writeChannels = new(sync.Map)

func syncFile(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	connection, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	writeChannel := make(chan util.AppNode)
	writeChannels.Store(writeChannel, writeChannel)

	defer func() {
		writeChannels.Delete(writeChannel)
		connection.Close()
	}()

	done := make(chan struct{})
	appNodeConfig := make([]util.AppNode, 5)
	go func() {
		defer close(done)
		for {
			var message common.Message
			if err := connection.ReadJSON(&message); err != nil && err != io.ErrUnexpectedEOF {
				log.Println("client down", err)
				return
			}

			if common.ClientConnect == message.MessageType {
				if err := json.Unmarshal(message.Data, &appNodeConfig); err != nil {
					connection.WriteJSON(common.NewClientConnectReplyMessage([]byte(err.Error())))
					return
				}

				for _, clientApp := range appNodeConfig {
					syncSingleAppFiles(connection, clientApp)
				}
			} else {
				log.Println("unsupported message type")
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case app := <-writeChannel:
			for _, clientApp := range appNodeConfig {
				if app.Name != "" && app.Name != clientApp.Name {
					continue
				}
				if app.Profile != "" && app.Profile != clientApp.Profile {
					continue
				}
				if app.Label != "" && app.Label != clientApp.Label {
					continue
				}
				syncSingleAppFiles(connection, clientApp)
			}
		}
	}
}

func refresh(_ http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	writeChannels.Range(func(key, value interface{}) bool {
		value.(chan util.AppNode) <- util.AppNode{Name: ps.ByName("app"), Label: ps.ByName("label"), Profile: ps.ByName("profile")}
		return true
	})
}

func main() {
	router := httprouter.New()
	router.GET("/sync", syncFile)
	router.GET("/refresh", refresh)
	router.GET("/refresh/:app", refresh)
	router.GET("/refresh/:app/:profile", refresh)
	router.GET("/refresh/:app/:profile/:label", refresh)
	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}
