package tgbot

import (
	"github.com/Yuzuki999/Aws-Manger-Bot/conf"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
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
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
		_, err := bot.Edit(c.Message, "请输入密钥备注: ")
		if err != nil {
			log.Error("Edit Message error: ", err)
		}
		p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
			p.Data[m.Sender.ID] = &Data{Data: map[string]string{"name": m.Text}}
			p.Session[m.Sender.ID].Channel <- true
		})
		select {
		case tmp := <-p.Session[c.Sender.ID].Channel:
			if tmp != true {
				return
			}
			_, editErr := bot.Send(c.Sender, "请输入密钥ID: ")
			if editErr != nil {
				log.Error("Send Message error: ", editErr)
			}
			p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
				p.Data[m.Sender.ID].Data["id"] = m.Text
				_, sendErr := bot.Send(m.Sender, "请输入密钥: ")
				if sendErr != nil {
					log.Error("Send Message error: ", sendErr)
				}
				p.Session[m.Sender.ID].Channel <- true
			})
			select {
			case tmp := <-p.Session[c.Sender.ID].Channel:
				if tmp != true {
					return
				}
				p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
					defer delete(p.Data, c.Sender.ID)
					if _, ok := p.Config.UserInfo[c.Sender.ID]; !ok {
						p.Config.UserInfo[c.Sender.ID] = &conf.UserData{
							UserName:  c.Sender.Username,
							AwsSecret: map[string]*conf.AwsSecret{},
						}
					} else if _, ok := p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Data[c.Sender.ID].Data["name"]]; ok {
						_, sendErr := bot.Send(c.Sender, "密钥已存在！！")
						if sendErr != nil {
							log.Error("Send Message error: ", sendErr)
						}
						return
					}
					p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Data[c.Sender.ID].Data["name"]] = &conf.AwsSecret{
						Id:     p.Data[c.Sender.ID].Data["id"],
						Secret: m.Text,
					}
					if p.Config.UserInfo[c.Sender.ID].NowKey == "" {
						p.Config.UserInfo[c.Sender.ID].NowKey = p.Data[c.Sender.ID].Data["name"]
					}
					conf.Lock.Lock()
					saveErr := p.Config.SaveConfig()
					conf.Lock.Unlock()
					if saveErr != nil {
						_, sendErr := bot.Send(c.Sender, "添加密钥失败！")
						if sendErr != nil {
							log.Error(sendErr)
						}
						log.Error("Save config error: ", saveErr)
						return
					}
					_, sendErr := bot.Send(c.Sender, "已添加密钥！")
					if sendErr != nil {
						log.Error("Send message error: ", sendErr)
					}
					p.Session[m.Sender.ID].Channel <- true
				})
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
	bot.Handle(&listKey, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action:  List key")
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
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
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
		_, err := bot.Edit(c.Message, "请输入要删除的密钥备注：")
		if err != nil {
			log.Error("Edit Message error: ", err)
		}
		p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
			defer delete(p.Data, m.Sender.ID)
			if _, ok := p.Config.UserInfo[m.Sender.ID].AwsSecret[m.Text]; ok {
				if p.Config.UserInfo[m.Sender.ID].NowKey == m.Text {
					p.Config.UserInfo[m.Sender.ID].NowKey = ""
				}
				delete(p.Config.UserInfo[m.Sender.ID].AwsSecret, m.Text)
				conf.Lock.Lock()
				saveErr := p.Config.SaveConfig()
				conf.Lock.Unlock()
				if saveErr != nil {
					_, sendErr := bot.Send(m.Sender, "删除失败！")
					if sendErr != nil {
						log.Error(sendErr)
					}
					log.Error("Send message error: ", sendErr)
					return
				}
				_, sendErr := bot.Send(m.Sender, "删除成功！")
				if sendErr != nil {
					log.Error("Send message error: ", sendErr)
				}
			} else {
				_, err := bot.Send(m.Sender, "密钥不存在！")
				if err != nil {
					log.Error("Send message error: ", err)
				}
			}
			p.Session[m.Sender.ID].Channel <- true
		})
		select {
		case tmp := <-p.Session[c.Sender.ID].Channel:
			if tmp != true {
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
	bot.Handle(&swKey, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action:  Switch key")
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
		_, err := bot.Edit(c.Message, "请输入要使用的密钥备注：")
		if err != nil {
			log.Error("Edit Message error: ", err)
		}
		p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
			defer delete(p.Data, m.Sender.ID)
			if _, ok := p.Config.UserInfo[m.Sender.ID].AwsSecret[m.Text]; ok {
				p.Config.UserInfo[m.Sender.ID].NowKey = m.Text
				conf.Lock.Lock()
				saveErr := p.Config.SaveConfig()
				conf.Lock.Unlock()
				if saveErr != nil {
					log.Println("Save config error: ", saveErr)
					_, sendErr := bot.Send(m.Sender, "切换密钥失败！")
					if sendErr != nil {
						log.Println("Send message error: ", sendErr)
					}
					return
				}
				_, sendErr := bot.Send(m.Sender, "切换密钥成功！")
				if sendErr != nil {
					log.Println("Send message error: ", sendErr)
				}
			} else {
				_, err := bot.Send(m.Sender, "切换密钥失败！该密钥不存在！")
				if err != nil {
					log.Println("Send message error: ", err)
				}
			}
			p.Session[m.Sender.ID].Channel <- true
		})
		select {
		case tmp := <-p.Session[c.Sender.ID].Channel:
			if tmp != true {
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
	bot.Handle(&addProxy, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action:  Add proxy for key")
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
		_, err := bot.Edit(c.Message, "请输入要添加代理的密钥备注: ")
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
			p.Data[m.Sender.ID].Data["name"] = m.Text
			p.Session[m.Sender.ID].Channel <- true
		})
		select {
		case tmp := <-p.Session[c.Sender.ID].Channel:
			if tmp != true {
				return
			}
			_, sendErr := bot.Send(c.Sender, "请输入代理地址: ")
			if sendErr != nil {
				log.Error("Send message error: ", sendErr)
			}
			p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
				defer delete(p.Data, m.Sender.ID)
				p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Data[m.Sender.ID].Data["name"]].Proxy = m.Text
				conf.Lock.Lock()
				saveErr := p.Config.SaveConfig()
				conf.Lock.Unlock()
				if saveErr != nil {
					_, sendErr := bot.Send(m.Sender, "添加失败")
					if sendErr != nil {
						log.Error("Send message error: ", sendErr)
					}
					return
				}
				_, sendErr := bot.Send(m.Sender, "添加成功")
				if sendErr != nil {
					log.Error("Send message error: ", sendErr)
				}
				p.Session[m.Sender.ID].Channel <- true
			})
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
	bot.Handle(&delKey, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action:  Del proxy for key")
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
		_, err := bot.Edit(c.Message, "请输入要删除代理的密钥备注: ")
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
			defer delete(p.Data, m.Sender.ID)
			p.Config.UserInfo[m.Sender.ID].AwsSecret[m.Text].Proxy = ""
			conf.Lock.Lock()
			saveErr := p.Config.SaveConfig()
			conf.Lock.Unlock()
			if saveErr != nil {
				_, sendErr := bot.Send(m.Sender, "删除失败")
				if sendErr != nil {
					log.Error("Send message error: ", sendErr)
				}
				return
			}
			_, sendErr := bot.Send(m.Sender, "删除成功")
			if sendErr != nil {
				log.Error("Send message error: ", sendErr)
			}
			p.Session[m.Sender.ID].Channel <- true
		})
		select {
		case tmp := <-p.Session[c.Sender.ID].Channel:
			if tmp != true {
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
}
