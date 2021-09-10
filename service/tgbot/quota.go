package tgbot

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
)

func (p *TgBot) getQuota(bot *tb.Bot, c *tb.Callback) {
	_, err := bot.Edit(c.Message, "请输入ServiceCode和QuotaCode(用空格隔开): ")
	if err != nil {
		log.Println("Edit message error: ", err)
	}
	p.State[c.Sender.ID] = &State{Parent: 6}
}

func (p TgBot) changeQuota(bot *tb.Bot, c *tb.Callback) {
	_, err := bot.Edit(c.Message, "请输入ServiceCode和QuotaCode和要提升至的数量(用空格隔开): ")
	if err != nil {
		log.Println("Edit message error: ", err)
	}
	p.State[c.Sender.ID] = &State{Parent: 7}
}

func (p *TgBot) QuotaManger(bot *tb.Bot) {
	key := tb.ReplyMarkup{}
	getQuota := key.Data("获取配额", "get_quota")
	bot.Handle(&getQuota, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请选择地区", p.RegionKey)
		if err != nil {
			log.Println("Edit message error:", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 102}
	})
	updateQuota := key.Data("更新配额", "update_quota")
	bot.Handle(&updateQuota, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请选择地区", p.RegionKey)
		if err != nil {
			log.Println("Edit message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 103}
	})
	key.Inline(key.Row(getQuota, updateQuota))
	bot.Handle("/QuotaManger", func(m *tb.Message) {
		_, err := bot.Send(m.Sender, "请选择要进行的操作: ", key)
		if err != nil {
			log.Println("Send message error: ", err)
		}
	})
}
