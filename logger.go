package keycloak

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

type logger struct {
	filename string
	*log.Logger
}

var userLog *logger
var appLog *logger

var logs *logger
var once sync.Once
var file *os.File
var fileStat os.FileInfo
var files []os.FileInfo

//GetInstance returns a new logger to a file
func getInstance(name string) *logger {
	once.Do(func() {
		logs = createLogger("./log/" + name)
	})
	return logs
}

func createLogger(fname string) *logger {
	file, _ = os.OpenFile(fname, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	fileStat, _ = file.Stat()
	files = append(files, fileStat)
	return &logger{
		filename: fname,
		Logger:   log.New(file, "", log.Ldate|log.Ltime),
	}
}

func logAction(username string, a Action, additional string) {
	info, _ := file.Stat()

	if info.Name() == files[0].Name() {
		userLog.Println(username+": ", a, " ", additional)
	} else {
		userLog = getInstance("UserLogs.log")
	}
}

//LogAction is an external call for logging actions into the file log
func LogAction(a Action, additional string) {
	username := GetUsername()
	info, _ := file.Stat()
	if info.Name() == files[1].Name() {
		appLog.Println(username+": ", a, " ", additional)
	} else {
		appLog = getInstance("AppLogs.log")
	}
}

//GetUsername gets the current users username
func GetUsername() string {
	client := &http.Client{}
	url := keycloakserver + "/auth/realms/" + realm + "/protocol/openid-connect/userinfo"
	req, _ := http.NewRequest("GET", url, nil)
	if token == nil {
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	//Check if token is still valid
	response, err := client.Do(req)
	if err != nil || response.Status != "200 OK" {
		return ""
	}
	body, _ := ioutil.ReadAll(response.Body)
	var f interface{}
	json.Unmarshal(body, &f)
	m := f.(map[string]interface{})
	username := m["preferred_username"].(string)
	return username
}

//GetEmail gets the current users username
func GetEmail() string {
	client := &http.Client{}
	url := keycloakserver + "/auth/realms/" + realm + "/protocol/openid-connect/userinfo"
	req, _ := http.NewRequest("GET", url, nil)
	if token == nil {
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	//Check if token is still valid
	response, err := client.Do(req)
	if err != nil || response.Status != "200 OK" {
		return ""
	}
	body, _ := ioutil.ReadAll(response.Body)
	var f interface{}
	json.Unmarshal(body, &f)
	m := f.(map[string]interface{})
	username := m["email"].(string)
	return username
}
