package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type UserData struct {
	UserName  string                `yaml:"user_name"`
	NowKey    string                `yaml:"now_key"`
	AwsSecret map[string]*AwsSecret `yaml:"aws_secret"`
}

type AwsSecret struct {
	Name   string `yaml:"name"`
	Id     string `yaml:"id"`
	Secret string `yaml:"secret"`
}

type Conf struct {
	BotToken string
	UserInfo map[int]*UserData `yaml:"user_info"`
}

func New() *Conf {
	return &Conf{}
}

func (c *Conf) LoadConfig() error {
	r, readErr := ioutil.ReadFile("./config.yml")
	if readErr != nil {
		if os.IsNotExist(readErr) {
			var token string
			fmt.Println("Input telegram token: ")
			_, scanErr := fmt.Scan(&token)
			if scanErr != nil {
				log.Fatalln("Get telegram bot token error: ", scanErr)
			}
			c.BotToken = token
			c.UserInfo = map[int]*UserData{0: {}}
			writeErr := c.SaveConfig()
			return writeErr
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
