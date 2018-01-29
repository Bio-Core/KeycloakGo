package keycloak

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type logger struct {
	filename string
	*log.Logger
}

var userLog *logger
var logs *logger
var once sync.Once
var file *os.File
var fileStat os.FileInfo

//GetInstance returns a new logger to a file
func GetInstance() *logger {
	once.Do(func() {
		logs = createLogger("./logs/UserLogs.log")
	})
	return logs
}

func createLogger(fname string) *logger {
	file, _ = os.OpenFile(fname, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	fileStat, _ = file.Stat()
	return &logger{
		filename: fname,
		Logger:   log.New(file, "", log.Ldate|log.Ltime),
	}
}

func logAction(username string, a action, additional string) {
	info, _ := file.Stat()
	if info.Name() == fileStat.Name() {
		event := getAction(a)
		userLog.Println(username+": ", event, " ", additional)
	} else {
		GetInstance()
	}
	a1 := fileStat.Name()
	a2 := fileStat.Sys()
	a3 := fileStat.Size()
	b1 := info.Name()
	b2 := info.Sys()
	b3 := info.Size()
	same := os.SameFile(fileStat, info)
	fmt.Printf(a1, a2, a3, b1, b2, b3, same)
}

func getAction(a action) Action {
	switch a {
	case actionLogin:
		return ActionLogin

	case actionLogout:
		return ActionLogout

	case actionPageAccess:
		return ActionPageAccess
	}
	return ""
}
