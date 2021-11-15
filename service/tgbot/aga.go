package tgbot

import (
	"github.com/Yuzuki999/Aws-Manger-Bot/aws"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

func (p *TgBot) AgaManger(bot *tb.Bot) {
	key := &tb.ReplyMarkup{}
	createAga := key.Data("创建Aga", "create_aga")
	bot.Handle(&createAga, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action: Create aga")
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
		_, editErr := bot.Edit(c.Message, "请输入Aga的备注(不要重复): ")
		if editErr != nil {
			log.Error("Edit message error: ", editErr)
		}
		p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
			p.Data[m.Sender.ID].Data["agaName"] = m.Text
			p.Session[m.Sender.ID].Channel <- true
		})
		select {
		case tmp := <-p.Session[c.Sender.ID].Channel:
			if tmp != true {
				_, sendErr := bot.Send(c.Sender, "请选择地区: ", p.RegionKey)
				if sendErr != nil {
					log.Println("Send message error: ", sendErr)
				}
				return
			}
			p.Data[c.Sender.ID].RegionChan = make(chan int)
			select {
			case <-p.Data[c.Sender.ID].RegionChan:
				_, err := bot.Edit(c.Message, "请输入要关联的Ec2实例ID: ")
				if err != nil {
					log.Error("Edit message error: ", err)
				}
				p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
					defer delete(p.Data, m.Sender.ID)
					awsRt, awsErr := aws.New("us-west-2",
						p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Id,
						p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Secret,
						p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Proxy)
					if awsErr != nil {
						_, sendErr := bot.Send(m.Sender, "创建失败")
						if sendErr != nil {
							log.Error("Send message error: ", sendErr)
						}
						log.Error("Init aws obj error: ", awsErr)
						return
					}
					creRt, creErr := awsRt.CreateAga(p.Data[m.Sender.ID].Data["agaName"],
						p.Data[m.Sender.ID].Data["region"],
						m.Text)
					if creErr != nil {
						_, sendErr := bot.Send(m.Sender, "创建失败")
						if sendErr != nil {
							log.Error("Send message error: ", sendErr)
						}
						log.Error("Create aga error: ", creErr)
					}
					var ip string
					for _, v := range creRt.Ip {
						for _, v2 := range v.IpAddresses {
							ip += *v2 + "\n"
						}
					}
					_, sendErr := bot.Send(m.Sender, "创建成功! "+"\n\n备注: "+
						*creRt.Name+"\nIP: \n"+ip+"\n状态: "+*creRt.Status+"\nArn: \n"+creRt.Arn)
					if sendErr != nil {
						log.Error("Send message error: ", sendErr)
					}
					p.Session[m.Sender.ID].Channel <- true
				})
			case <-time.After(30 * time.Second):
				_, editErr := bot.Edit(c.Message, "操作超时")
				if editErr != nil {
					log.Error("Edit message error: ", editErr)
				}
			}
			select {
			case tmp := <-p.Session[c.Sender.ID].Channel:
				if tmp != true {
					return
				}
			case <-time.After(30 * time.Second):
				_, sendErr := bot.Edit(c.Message, "操作超时")
				if sendErr != nil {
					log.Error("Edit message error: ", sendErr)
				}
			}
			p.Session.SessionDel(c.Sender.ID)
		case <-time.After(30 * time.Second):
			_, sendErr := bot.Edit(c.Message, "操作超时")
			if sendErr != nil {
				log.Error("Edit message error: ", sendErr)
			}
			p.Session.SessionDel(c.Sender.ID)
		}
	})
	listAga := key.Data("列出Aga", "list_aga")
	bot.Handle(&listAga, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action: List aga")
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
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
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
		_, err := bot.Edit(c.Message, "请输入Arn: ")
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
			defer delete(p.Data, m.Sender.ID)
			awsRt, awsErr := aws.New(p.Data[m.Sender.ID].Data["region"],
				p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Id,
				p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Secret,
				p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Proxy)
			if awsErr != nil {
				_, sendErr := bot.Send(m.Sender, "删除失败")
				if sendErr != nil {
					log.Error("Send message error: ", sendErr)
				}
				log.Error("Init aws oub error: ", awsErr)
				return
			}
			delErr := awsRt.DeleteAga(m.Text)
			if delErr != nil {
				_, sendErr := bot.Send(m.Sender, "删除失败")
				if sendErr != nil {
					log.Error("Send message error: ", sendErr)
				}
				log.Error("Delete aga error: ", awsErr)
				return
			}
			p.Session[m.Sender.ID].Channel <- true
		})
		select {
		case tmp := <-p.Session[c.Sender.ID].Channel:
			if !tmp {
				return
			}
			p.Session.SessionDel(c.Sender.ID)
		case <-time.After(30 * time.Second):
			_, sendErr := bot.Edit(c.Message, "操作超时")
			if sendErr != nil {
				log.Error("Edit message error: ", sendErr)
			}
			p.Session.SessionDel(c.Sender.ID)
		}
	})
	key.Inline(key.Row(createAga, listAga), key.Row(delAga))
	bot.Handle("/AgaManger", func(m *tb.Message) {
		_, err := bot.Send(m.Sender, "请选择要进行的操作", key)
		if err != nil {
			log.Error("Edit message error: ", err)
		}
	})
}
