package tgbot

import (
	"github.com/Yuzuki999/Aws-Manger-Bot/conf"
	"github.com/Yuzuki999/Aws-Manger-Bot/session"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"os"
	"time"
)

type TgBot struct {
	Config    *conf.Conf
	TypeKey   *tb.ReplyMarkup
	AmiKey    *tb.ReplyMarkup
	Data      map[int]*Data
	RegionKey *tb.ReplyMarkup
	Session   session.Session
}

type Data struct {
	RegionChan chan int
	Data       map[string]string
}

func New(Config *conf.Conf) *TgBot {
	return &TgBot{Config: Config}
}

func (p *TgBot) Start() {
	p.Data = map[int]*Data{}
	p.Session = session.Session{}
	bot, err := tb.NewBot(tb.Settings{
		Token:  p.Config.BotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Error("Bot start error: ", err)
		os.Exit(1)
	}
	bot.Handle("/start", func(m *tb.Message) {
		_, err := bot.Send(m.Sender,
			`欢迎使用Aws Manger Bot
请使用 /KeyManger 指令添加并选择密钥，然后使用相应指令管理相应资源
指令列表:

/Ec2Manger Ec2相关操作
/AgaManger Aga相关操作
/QuotaManger 配额相关操作`)
		if err != nil {
			log.Error("Send message error: ", err)
		}
	})
	p.KeyManger(bot)
	p.QuotaManger(bot)
	p.Ec2Manger(bot)
	p.setRegionKey(bot)
	p.AgaManger(bot)
	bot.Handle(tb.OnText, func(m *tb.Message) {
		if p.Session.SessionCheck(m.Sender.ID) {
			p.Session.SessionHandle(m.Sender.ID, m)
		}
	})
	bot.Start()
}
