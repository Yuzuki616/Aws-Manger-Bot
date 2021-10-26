package tgbot

import (
	"github.com/338317/Aws-Manger-Bot/aws"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (p *TgBot) AgaManger(bot *tb.Bot) {
	key := &tb.ReplyMarkup{}
	createAga := key.Data("创建Aga", "create_aga")
	bot.Handle(&createAga, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action: Create aga")
		_, err := bot.Edit(c.Message, "请输入Aga的备注(不要重复): ")
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 15, Data: map[string]string{}}
	})
	listAga := key.Data("列出Aga", "list_aga")
	bot.Handle(&listAga, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action: List aga")
		awsRt, awsErr := aws.New("us-west-2",
			p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Id,
			p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Secret,
			p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Proxy)
		if awsErr != nil {
			_, sendErr := bot.Send(c.Sender, "列出失败")
			if sendErr != nil {
				log.Error("Send message error: ", sendErr)
			}
			log.Error("Init aws obj error: ", awsErr)
			return
		}
		aga, agaErr := awsRt.ListAga()
		if agaErr != nil {
			_, sendErr := bot.Send(c.Sender, "列出失败")
			if sendErr != nil {
				log.Error("Send message error: ", sendErr)
			}
			log.Error("List aga error: ", agaErr)
			return
		}
		agaList := "已创建的Aga: \n\n"
		for _, v := range aga {
			var ip string
			for _, i := range v.IpSets[0].IpAddresses {
				ip += *i + "\n"
			}
			agaList += "\n备注: " + *v.Name + "\n状态: " + *v.Status +
				"\nIP: \n" + ip + "\nArn: " + *v.AcceleratorArn
		}
		_, sendErr := bot.Edit(c.Message, agaList)
		if sendErr != nil {
			log.Error("Send message error: ", sendErr)
		}
	})
	delAga := key.Data("删除Aga", "del_aga")
	bot.Handle(&delAga, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action: Delete aga")
		_, err := bot.Edit(c.Message, "请输入Arn: ")
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 18, Data: map[string]string{}}
	})
	key.Inline(key.Row(createAga, listAga), key.Row(delAga))
	bot.Handle("/AgaManger", func(m *tb.Message) {
		_, err := bot.Send(m.Sender, "请选择要进行的操作", key)
		if err != nil {
			log.Error("Edit message error: ", err)
		}
	})
}
