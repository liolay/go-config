package main

import (
	"time"
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"
	"os"
)
type ConfigInfo struct {
	Name       string
	Content    []byte
	UpdateTime time.Time
}


func main() {

	for range time.Tick(5 * time.Second) {
		reader := strings.NewReader("")
		resp,_ := http.Post("http://localhost:6973/refresh","application/json;charset=utf-8",reader)
		defer resp.Body.Close()

		responseBody, _ := ioutil.ReadAll(resp.Body)

		var configInfos []ConfigInfo
		json.Unmarshal(responseBody, &configInfos)

		for _, config := range configInfos {
			ioutil.WriteFile(config.Name+"_new",config.Content,os.ModePerm)
		}
	}

}
