package main

import (
	"fmt"
	"github.com/338317/Aws-Manger-Bot/conf"
	"github.com/338317/Aws-Manger-Bot/service/tgbot"
	"log"
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
