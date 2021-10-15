package conf

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"sync"
)

type UserData struct {
	UserName  string                `yaml:"User_Name"`
	NowKey    string                `yaml:"Now_Key"`
	AwsSecret map[string]*AwsSecret `yaml:"Aws_Secret"`
}

type AwsSecret struct {
	Id     string `yaml:"id"`
	Secret string `yaml:"secret"`
}

type Conf struct {
	LogLevel string            `yaml:"Log_Level"`
	BotToken string            `yaml:"Bot_Token"`
	UserInfo map[int]*UserData `yaml:"User_Info"`
}

func New() *Conf {
	return &Conf{}
}

var Lock = sync.RWMutex{}

func (c *Conf) LoadConfig() error {
	r, readErr := ioutil.ReadFile("./config.yml")
	if readErr != nil {
		if os.IsNotExist(readErr) {
			log.Error("Config file not found")
			log.Error("Write default config file")
			c.LogLevel = "error"
			c.UserInfo = map[int]*UserData{0: {}}
			c.BotToken = "Tg Bot Token"
			writeErr := c.SaveConfig()
			if writeErr != nil {
				log.Error("Write file error: ", writeErr)
				os.Exit(1)
			}
			log.Error("已将默认配置文件写出，请填写bot token后重新启动")
			log.Error("Exit")
			os.Exit(1)
		}
		return readErr
	}
	unmErr := yaml.Unmarshal(r, c)
	if unmErr != nil {
		return unmErr
	}
	return nil
}

func (c *Conf) SaveConfig() error {
	rt, marErr := yaml.Marshal(c)
	if marErr != nil {
		return marErr
	}
	writeErr := ioutil.WriteFile("./config.yml", rt, 0644)
	if writeErr != nil {
		return writeErr
	}
	return nil
}
