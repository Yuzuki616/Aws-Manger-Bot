package tgbot

import (
	"github.com/338317/Aws-Manger-Bot/conf"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"os"
	"time"
)

type TgBot struct {
	Config    *conf.Conf
	TypeKey   *tb.ReplyMarkup
	AmiKey    *tb.ReplyMarkup
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
/QuotaManger 配额相关操作`)
		if err != nil {
			log.Warning("Send message error: ", err)
		}
	})
	p.KeyManger(bot)
	p.QuotaManger(bot)
	p.Ec2Manger(bot)
	p.setRegionKey(bot)
	p.GlobalMess(bot)
	bot.Start()
}
