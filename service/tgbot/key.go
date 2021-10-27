package tgbot

import (
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (p *TgBot) KeyManger(bot *tb.Bot) {
	key := &tb.ReplyMarkup{}
	addKey := key.Data("添加密钥", "addKey")
	listKey := key.Data("查看密钥", "listKey")
	delKey := key.Data("删除密钥", "delKey")
	swKey := key.Data("切换密钥", "swKey")
	addProxy := key.Data("添加代理", "addProxy")
	delProxy := key.Data("删除代理", "delProxy")
	key.Inline(key.Row(addKey, delKey), key.Row(listKey, swKey), key.Row(addProxy, delProxy))
	bot.Handle("/KeyManger", func(m *tb.Message) {
		if m.Private() {
			_, err := bot.Send(m.Sender, "请选择你要进行的操作", key)
			if err != nil {
				log.Error("Send Message error: ", err)
			}
		} else {
			_, err := bot.Reply(m, "请私聊Bot使用")
			if err != nil {
				log.Error("Reply message error: ")
			}
		}
	})
	bot.Handle(&addKey, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action:  Add key")
		_, err := bot.Edit(c.Message, "请输入密钥备注: ")
		if err != nil {
			log.Error("Edit Message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 0}
	})
	bot.Handle(&listKey, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action:  List key")
		tmp := "当前使用的密钥: " + p.Config.UserInfo[c.Sender.ID].NowKey + "\n\n已添加的密钥: "
		if p.Config.UserInfo[c.Sender.ID].AwsSecret == nil {
			_, err := bot.Edit(c.Message, "当前未添加任何密钥")
			if err != nil {
				log.Error("Edit Message error: ", err)
			}
			return
		}
		for key, val := range p.Config.UserInfo[c.Sender.ID].AwsSecret {
			tmp += "\n\n备注: " + key + "\nID: " + val.Id + "\n密钥: " + val.Secret
		}
		_, err := bot.Edit(c.Message, tmp)
		if err != nil {
			log.Error("Edit Message error: ", err)
		}
	})
	bot.Handle(&delKey, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action:  Delete Key")
		_, err := bot.Edit(c.Message, "请输入要删除的密钥备注：")
		if err != nil {
			log.Error("Edit Message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 3}
	})
	bot.Handle(&swKey, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action:  Switch key")
		_, err := bot.Edit(c.Message, "请输入要使用的密钥备注：")
		if err != nil {
			log.Error("Edit Message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 4}
	})
	bot.Handle(&addProxy, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action:  Add proxy for key")
		_, err := bot.Edit(c.Message, "请输入要添加代理的密钥备注: ")
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 20, Data: map[string]string{}}
	})
	bot.Handle(&delKey, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action:  Del proxy for key")
		_, err := bot.Edit(c.Message, "请输入要删除代理的密钥备注: ")
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 22, Data: map[string]string{}}
	})
}
