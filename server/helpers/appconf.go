package helpers

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type AppConf struct {
	Conf Application `yaml:"application"`
}

type Application struct {
	LogLevel       int8   `yaml:"logLevel"`
	OrgName        string `yaml:"orgName"`
	OrgAdmin       string `yaml:"orgAdmin"`
}

func LoadAppConf() (*AppConf, error) {
	confPath := GetConfigPath("app.yaml")
	yamlFile, err := ioutil.ReadFile(confPath)
	if err != nil {
		logger.Errorf("yamlFile.Get err: %v ", err)
		return nil, err
	}
	conf := new(AppConf)
	if err = yaml.Unmarshal(yamlFile, conf); err != nil {
		logger.Errorf("Unmarshal: %v", err)
		return nil, err
	}
	return conf, nil
}
