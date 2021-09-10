package main

import (
	"fmt"
	"log"
	"main/conf"
	"main/service/tgbot"
)
import "flag"

const (
	version = "0.0.1"
	name    = "Aws-Manger-Bot"
)

func main() {
	flag.Parse()
	fmt.Printf("%s %s\n", name, version)
	config := conf.New()
	err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}
	bot := tgbot.New(config)
	bot.Start()
}
