package tgbot

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
)

func (p *TgBot) LightSailManger(bot *tb.Bot) {
	key := &tb.ReplyMarkup{}
	create := key.Data("创建LightSail", "create_ls")
	bot.Handle(&create, func(c *tb.Callback) {

	})
	list := key.Data("列出LightSail", "list_ls")
	change := key.Data("更换IP", "change_ip")
	del := key.Data("删除LightSail", "del_ls")
	key.Inline(key.Row(create), key.Row(list), key.Row(change, del))
	bot.Handle("/LightSailManger", func(m *tb.Message) {
		_, err := bot.Send(m.Sender, "请选择你要进行的操作", key)
		if err != nil {
			log.Println("Send message error: ", err)
		}
	})
}
