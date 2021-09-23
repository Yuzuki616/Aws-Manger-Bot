package tgbot

import (
	"github.com/338317/Aws-Manger-Bot/aws"
	"github.com/338317/Aws-Manger-Bot/conf"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
)

func (p *TgBot) GlobalMess(bot *tb.Bot) {
	bot.Handle(tb.OnText, func(m *tb.Message) {
		if _, ok := p.State[m.Sender.ID]; ok {
			switch p.State[m.Sender.ID].Parent {
			case 0:
				_, err := bot.Send(m.Sender, "请输入密钥ID: ")
				if err != nil {
					log.Println("Send Message error: ", err)
				}
				p.State[m.Sender.ID].Data = make(map[string]string)
				p.State[m.Sender.ID].Data["name"] = m.Text
				p.State[m.Sender.ID].Parent++
			case 1:
				_, err := bot.Send(m.Sender, "请输入密钥: ")
				if err != nil {
					log.Println("Send Message error: ", err)
				}
				p.State[m.Sender.ID].Data["id"] = m.Text
				p.State[m.Sender.ID].Parent++
			case 2:
				defer delete(p.State, m.Sender.ID)
				if _, ok := p.Config.UserInfo[m.Sender.ID]; !ok {
					p.Config.UserInfo[m.Sender.ID] = &conf.UserData{
						UserName:  m.Sender.Username,
						AwsSecret: map[string]*conf.AwsSecret{},
					}
				} else if _, ok := p.Config.UserInfo[m.Sender.ID].AwsSecret[p.State[m.Sender.ID].Data["name"]]; ok {
					_, err := bot.Send(m.Sender, "密钥已存在！！")
					if err != nil {
						log.Println("Send Message error: ", err)
					}
					return
				}
				p.Config.UserInfo[m.Sender.ID].AwsSecret[p.State[m.Sender.ID].Data["name"]] = &conf.AwsSecret{
					Id:     p.State[m.Sender.ID].Data["id"],
					Secret: m.Text,
				}
				saveErr := p.Config.SaveConfig()
				if saveErr != nil {
					_, sendErr := bot.Send(m.Sender, "添加密钥失败！")
					if sendErr != nil {
						log.Println(sendErr)
					}
					log.Println("Save config error: ", saveErr)
					return
				}
				_, sendErr := bot.Send(m.Sender, "已添加密钥！")
				if sendErr != nil {
					log.Println("Send message error: ", sendErr)
				}
			case 3:
				defer delete(p.State, m.Sender.ID)
				if _, ok := p.Config.UserInfo[m.Sender.ID].AwsSecret[m.Text]; ok {
					saveErr := p.Config.SaveConfig()
					if saveErr != nil {
						_, sendErr := bot.Send(m.Sender, "删除失败！")
						if sendErr != nil {
							log.Println(sendErr)
						}
						log.Println("Send message error: ", sendErr)
						return
					}
					_, sendErr := bot.Send(m.Sender, "删除成功！")
					if sendErr != nil {
						log.Println("Send message error: ", sendErr)
					}
				} else {
					_, err := bot.Send(m.Sender, "密钥不存在！")
					if err != nil {
						log.Println("Send message error: ", err)
					}
				}
			case 4:
				defer delete(p.State, m.Sender.ID)
				if _, ok := p.Config.UserInfo[m.Sender.ID].AwsSecret[m.Text]; ok {
					p.Config.UserInfo[m.Sender.ID].NowKey = m.Text
					saveErr := p.Config.SaveConfig()
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
			case 5:
				p.State[m.Sender.ID].Data = make(map[string]string)
				p.State[m.Sender.ID].Data["name"] = m.Text
				p.State[m.Sender.ID].Parent++
				_, err := bot.Send(m.Sender, "请选择地区: ", p.RegionKey)
				if err != nil {
					log.Println("Send message error: ", err)
				}
				p.State[m.Sender.ID].Parent = 999
			case 6:
				defer delete(p.State, m.Sender.ID)
				newRt, newErr := aws.New(p.State[m.Sender.ID].Data["region"],
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Id,
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Secret)
				if newErr != nil {
					_, sendErr := bot.Send(m.Sender, "删除失败！")
					if sendErr != nil {
						log.Println("Send message error: ", sendErr)
					}
					return
				}
				delErr := newRt.DeleteEc2(m.Text)
				if delErr != nil {
					_, sendErr := bot.Send(m.Sender, "删除失败！")
					if sendErr != nil {
						log.Println("Send message error: ", sendErr)
					}
					log.Println("delete instance error: ", delErr)
					return
				}
				_, sendErr2 := bot.Send(m.Sender, "删除成功！")
				if sendErr2 != nil {
					log.Println("Send message error: ", sendErr2)
				}
			case 7:
				newRt, newErr := aws.New(p.State[m.Sender.ID].Data["region"],
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Id,
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Secret)
				if newErr != nil {
					_, sendErr := bot.Send(m.Sender, "删除失败！")
					if sendErr != nil {
						log.Println("Send message error: ", sendErr)
					}
					return
				}
				ip, chErr := newRt.ChangeEc2Ip(m.Text)
				if chErr != nil {
					log.Println("Change ip error: ", chErr)
					_, sendErr := bot.Send(m.Sender, "删除失败！")
					if sendErr != nil {
						log.Println("Send message error: ", sendErr)
					}
					return
				}
				_, sendErr := bot.Send(m.Sender, "更换ip成功，新的IP为: ", *ip)
				if sendErr != nil {
					log.Println("Send message error: ", sendErr)
				}
			case 8:
				p.State[m.Sender.ID].Data["ami"] = m.Text
				defer delete(p.State, m.Sender.ID)
				p.createEc2(bot, &tb.Callback{Sender: m.Sender, Message: m})
			case 9:
				p.State[m.Sender.ID].Data["type"] = m.Text
				_, err := bot.Send(m.Sender, "请选择操作系统", p.AmiKey)
				if err != nil {
					log.Println("Send message error: ", err)
				}
			default:
				return
			}
		}
	})
}
