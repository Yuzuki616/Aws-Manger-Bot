package tgbot

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
)

func (p *TgBot) KeyManger(bot *tb.Bot) {
	key := &tb.ReplyMarkup{}
	addKey := key.Data("添加密钥", "addKey")
	listKey := key.Data("查看密钥", "listKey")
	delKey := key.Data("删除密钥", "delKey")
	swKey := key.Data("切换密钥", "swKey")
	key.Inline(key.Row(addKey, delKey), key.Row(listKey, swKey))
	bot.Handle("/KeyManger", func(m *tb.Message) {
		if m.Private() {
			_, err := bot.Send(m.Sender, "请选择你要进行的操作", key)
			if err != nil {
				log.Println("Send Message error: ", err)
			}
		} else {
			_, err := bot.Reply(m, "请私聊Bot使用")
			if err != nil {
				log.Println("Reply message error: ")
			}
		}
	})
	bot.Handle(&addKey, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请输入密钥备注: ")
		if err != nil {
			log.Println("Edit Message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 0}
	})
	bot.Handle(&listKey, func(c *tb.Callback) {
		tmp := "当前使用的密钥: " + p.Config.UserInfo[c.Sender.ID].NowKey + "\n\n已添加的密钥: "
		for key, val := range p.Config.UserInfo[c.Sender.ID].AwsSecret {
			tmp += "\n\n备注: " + key + "\nID: " + val.Id + "\n密钥: " + val.Secret
		}
		_, err := bot.Edit(c.Message, tmp)
		if err != nil {
			log.Println("Edit Message error: ", err)
		}
	})
	bot.Handle(&delKey, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请输入要删除的密钥备注：")
		if err != nil {
			log.Println("Edit Message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 3}
	})
	bot.Handle(&swKey, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请输入要使用的密钥备注：")
		if err != nil {
			log.Println("Edit Message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 4}
	})
}
