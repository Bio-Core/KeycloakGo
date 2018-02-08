package keycloak

import (
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
func getInstance() *logger {
	once.Do(func() {
		logs = createLogger("./log/UserLogs.log")
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
		getInstance()
	}
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
