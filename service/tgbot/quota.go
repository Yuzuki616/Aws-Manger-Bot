package tgbot

import (
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	serviceCodeAcd = "ec2"
	quotaCodeAcd   = "L-1216C47A"
)

func (p *TgBot) QuotaManger(bot *tb.Bot) {
	quotaKey := &tb.ReplyMarkup{}
	acd := quotaKey.Data("查看标准EC2配额", "get_def")
	bot.Handle(&acd, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请选择地区", p.RegionKey)
		if err != nil {
			log.Println("Edit message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 106, Data: map[string]string{}}
	})
	other := quotaKey.Data("查看自定义配额", "get_other")
	bot.Handle(&other, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请选择地区", p.RegionKey)
		if err != nil {
			log.Println("Edit message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 104, Data: map[string]string{}}
	})
	quotaKey.Inline(quotaKey.Row(acd), quotaKey.Row(other))
	key := &tb.ReplyMarkup{}
	getQuota := key.Data("查看配额", "get_quota")
	bot.Handle(&getQuota, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action:  Get Quota")
		_, editErr := bot.Edit(c.Message, "请选择配额", quotaKey)
		if editErr != nil {
			log.Error("Edit message error: ", editErr)
		}
	})
	updateQuota := key.Data("更新配额", "update_quota")
	bot.Handle(&updateQuota, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action:  Update quota")
		_, err := bot.Edit(c.Message, "请选择地区", p.RegionKey)
		if err != nil {
			log.Println("Edit message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 105, Data: map[string]string{}}
	})
	key.Inline(key.Row(getQuota, updateQuota))
	bot.Handle("/QuotaManger", func(m *tb.Message) {
		_, err := bot.Send(m.Sender, "请选择要进行的操作: ", key)
		if err != nil {
			log.Println("Send message error: ", err)
		}
	})
}
