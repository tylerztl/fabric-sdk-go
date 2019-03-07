package helpers

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type AppConf struct {
	Conf Application `yaml:"application"`
}

type Application struct {
	LogPath  string `yaml:"logPath"`
	LogLevel int8   `yaml:"logLevel"`
	OrgName  string `yaml:"orgName"`
	OrgAdmin string `yaml:"orgAdmin"`
}

var appConfig = new(AppConf)

func init() {
	confPath := GetConfigPath("app.yaml")
	yamlFile, err := ioutil.ReadFile(confPath)
	if err != nil {
		panic(fmt.Errorf("yamlFile.Get err[%s]", err))
	}
	if err = yaml.Unmarshal(yamlFile, appConfig); err != nil {
		panic(fmt.Errorf("yamlFile.Unmarshal err[%s]", err))
	}
}

func GetAppConf() *AppConf {
	return appConfig
}
