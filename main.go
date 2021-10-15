package main

import (
	"github.com/338317/Aws-Manger-Bot/conf"
	"github.com/338317/Aws-Manger-Bot/service/tgbot"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"os"
)

const (
	version = "1.0.0"
)

func printInfo() {
	log.Info("Aws Manger Bot")
	log.Info("Version: ", version)
	log.Info("Github: https://github.com/yuzuki999/Aws-Manger-Bot")
}

func main() {
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "01-02 15:04:05",
		LogFormat:       "Aws-Manger-Bot | %time% | %lvl% >> %msg% \n",
	})
	printInfo()
	config := conf.New()
	log.Info("Loading config...")
	err := config.LoadConfig()
	if err != nil {
		log.Error("Load config error: ", err)
		os.Exit(1)
	}
	log.Info("Done")
	switch config.LogLevel {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	}
	bot := tgbot.New(config)
	bot.Start()
}
