package keycloak

import (
	"encoding/json"
	"fmt"
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
var files map[string]os.FileInfo

//GetInstance returns a new logger to a file
func getInstance(name string) *logger {
	logs = createLogger("./log/" + name)
	return logs
}

func createLogger(fname string) *logger {
	if files == nil {
		files = make(map[string]os.FileInfo, 2)
	}
	file, _ = os.OpenFile(fname, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	fileStat, _ = file.Stat()
	files[fname] = fileStat
	return &logger{
		filename: fname,
		Logger:   log.New(file, "", log.Ldate|log.Ltime),
	}
}

func logAction(username string, a Action, additional string) {
	if files["./log/UserLogs.log"] == nil {
		userLog = getInstance("UserLogs.log")
	}
	userLog.Println(username+": ", a, " ", additional)
}

//LogAction is an external call for logging actions into the file log
func LogAction(a Action, additional, tokenString string) {
	username := GetUsername(tokenString)
	email := GetEmail(tokenString)
	fmt.Printf("got username and gmail")
	if files["./log/AppLogs.log"] == nil {
		fmt.Printf("Found file not made")
		appLog = getInstance("AppLogs.log")
		fmt.Printf("Made file")
	}
	appLog.Println(username+";"+email+": ", a, " ", additional)
}

//GetUsername gets the current users username
func GetUsername(tokenString string) string {
	client := &http.Client{}
	url := keycloakserver + "/auth/realms/" + realm + "/protocol/openid-connect/userinfo"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
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
func GetEmail(tokenString string) string {
	client := &http.Client{}
	url := keycloakserver + "/auth/realms/" + realm + "/protocol/openid-connect/userinfo"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
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
