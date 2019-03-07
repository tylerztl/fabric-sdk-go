package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/logs"
)

var logger *logs.BeeLogger

func init() {
	appConfig := GetAppConf()
	if appConfig == nil {
		panic("AppConf is nil")
	}
	logger = logs.NewLogger()
	config := make(map[string]interface{})
	config["filename"] = appConfig.Conf.LogPath
	config["level"] = appConfig.Conf.LogLevel
	configStr, err := json.Marshal(config)
	if err != nil {
		panic(fmt.Errorf("logger marshal err[%s]", err))
	}
	err = logger.SetLogger(logs.AdapterConsole, string(configStr))
	err = logger.SetLogger(logs.AdapterFile, string(configStr))
	if err != nil {
		panic(fmt.Errorf("logger SetLogger err[%s]", err))
	}
	logger.EnableFuncCallDepth(true)
	logger.SetLogFuncCallDepth(4)
}

func GetLogger() *logs.BeeLogger {
	return logger
}
