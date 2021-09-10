package tgbot

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"main/conf"
	"time"
)

type TgBot struct {
	Config    *conf.Conf
	TypeKey   *tb.ReplyMarkup
	State     map[int]*State
	RegionKey *tb.ReplyMarkup
}

type State struct {
	Data   map[string]string
	Parent int
}

func New(Config *conf.Conf) *TgBot {
	return &TgBot{Config: Config}
}

func (p *TgBot) Start() {
	p.State = make(map[int]*State)
	bot, _ := tb.NewBot(tb.Settings{
		Token:  p.Config.BotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	p.KeyManger(bot)
	p.Ec2Manger(bot)
	p.setRegionKey(bot)
	p.GlobalMess(bot)
	bot.Start()
}
