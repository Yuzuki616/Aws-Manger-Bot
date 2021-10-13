package tgbot

import (
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (p *TgBot) AgaManger(bot *tb.Bot) {
	key := tb.ReplyMarkup{}
	createAga := key.Data("创建Aga", "create_aga")
	bot.Handle(&createAga, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请输入Aga的备注(不要重复): ")
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 15}
	})
	listAga := key.Data("列出Aga", "list_aga")
	delAga := key.Data("删除Aga", "del_aga")
	key.Inline(key.Row(createAga, listAga), key.Row(delAga))
	bot.Handle("/AgaManger", func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请选择要进行的操作", key)
		if err != nil {
			log.Error("Edit message error: ")
		}
	})
}
