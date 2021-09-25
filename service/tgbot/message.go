package tgbot

import (
	"github.com/338317/Aws-Manger-Bot/aws"
	"github.com/338317/Aws-Manger-Bot/conf"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
)

func (p *TgBot) GlobalMess(bot *tb.Bot) {
	bot.Handle(tb.OnText, func(m *tb.Message) {
		if _, ok := p.State[m.Sender.ID]; ok {
			switch p.State[m.Sender.ID].Parent {
			case 0: //获取密钥ID
				_, err := bot.Send(m.Sender, "请输入密钥ID: ")
				if err != nil {
					log.Error("Send Message error: ", err)
				}
				p.State[m.Sender.ID].Data = make(map[string]string)
				p.State[m.Sender.ID].Data["name"] = m.Text
				p.State[m.Sender.ID].Parent++
			case 1: //获取密钥
				_, err := bot.Send(m.Sender, "请输入密钥: ")
				if err != nil {
					log.Error("Send Message error: ", err)
				}
				p.State[m.Sender.ID].Data["id"] = m.Text
				p.State[m.Sender.ID].Parent++
			case 2: //添加密钥
				defer delete(p.State, m.Sender.ID)
				if _, ok := p.Config.UserInfo[m.Sender.ID]; !ok {
					p.Config.UserInfo[m.Sender.ID] = &conf.UserData{
						UserName:  m.Sender.Username,
						AwsSecret: map[string]*conf.AwsSecret{},
					}
				} else if _, ok := p.Config.UserInfo[m.Sender.ID].AwsSecret[p.State[m.Sender.ID].Data["name"]]; ok {
					_, err := bot.Send(m.Sender, "密钥已存在！！")
					if err != nil {
						log.Error("Send Message error: ", err)
					}
					return
				}
				p.Config.UserInfo[m.Sender.ID].AwsSecret[p.State[m.Sender.ID].Data["name"]] = &conf.AwsSecret{
					Id:     p.State[m.Sender.ID].Data["id"],
					Secret: m.Text,
				}
				conf.Lock.Lock()
				saveErr := p.Config.SaveConfig()
				conf.Lock.Unlock()
				if saveErr != nil {
					_, sendErr := bot.Send(m.Sender, "添加密钥失败！")
					if sendErr != nil {
						log.Error(sendErr)
					}
					log.Error("Save config error: ", saveErr)
					return
				}
				_, sendErr := bot.Send(m.Sender, "已添加密钥！")
				if sendErr != nil {
					log.Error("Send message error: ", sendErr)
				}
			case 3: //删除密钥
				defer delete(p.State, m.Sender.ID)
				if _, ok := p.Config.UserInfo[m.Sender.ID].AwsSecret[m.Text]; ok {
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
			case 4: //选择密钥
				defer delete(p.State, m.Sender.ID)
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
			case 5: //选择AWS地区
				p.State[m.Sender.ID].Data = make(map[string]string)
				p.State[m.Sender.ID].Data["name"] = m.Text
				p.State[m.Sender.ID].Parent++
				_, err := bot.Send(m.Sender, "请选择地区: ", p.RegionKey)
				if err != nil {
					log.Println("Send message error: ", err)
				}
				p.State[m.Sender.ID].Parent = 999
			case 6: //删除Ec2
				defer delete(p.State, m.Sender.ID)
				newRt, newErr := aws.New(p.State[m.Sender.ID].Data["region"],
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Id,
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Secret)
				if newErr != nil {
					_, sendErr := bot.Send(m.Sender, "删除失败！")
					if sendErr != nil {
						log.Error("Send message error: ", sendErr)
					}
					return
				}
				delErr := newRt.DeleteEc2(m.Text)
				if delErr != nil {
					_, sendErr := bot.Send(m.Sender, "删除失败！")
					if sendErr != nil {
						log.Error("Send message error: ", sendErr)
					}
					log.Error("delete instance error: ", delErr)
					return
				}
				_, sendErr2 := bot.Send(m.Sender, "删除成功！")
				if sendErr2 != nil {
					log.Error("Send message error: ", sendErr2)
				}
			case 7: //更换ip
				defer delete(p.State, m.Sender.ID)
				newRt, newErr := aws.New(p.State[m.Sender.ID].Data["region"],
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Id,
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Secret)
				if newErr != nil {
					_, sendErr := bot.Send(m.Sender, "删除失败！")
					if sendErr != nil {
						log.Error("Send message error: ", sendErr)
					}
					return
				}
				ip, chErr := newRt.ChangeEc2Ip(m.Text)
				if chErr != nil {
					log.Error("Change ip error: ", chErr)
					_, sendErr := bot.Send(m.Sender, "删除失败！")
					if sendErr != nil {
						log.Error("Send message error: ", sendErr)
					}
					return
				}
				_, sendErr := bot.Send(m.Sender, "更换ip成功，新的IP为: ", *ip)
				if sendErr != nil {
					log.Error("Send message error: ", sendErr)
				}
			case 8: //获取AMI
				p.State[m.Sender.ID].Data["ami"] = m.Text
				defer delete(p.State, m.Sender.ID)
				p.createEc2(bot, &tb.Callback{Sender: m.Sender, Message: m})
			case 9: //选择操作系统
				p.State[m.Sender.ID].Data["type"] = m.Text
				_, err := bot.Send(m.Sender, "请选择操作系统", p.AmiKey)
				if err != nil {
					log.Error("Send message error: ", err)
				}
			case 10: //获取quota code
				defer delete(p.State, m.Sender.ID)
				code := strings.Split(m.Text, " ")
				newRt, newErr := aws.New(p.State[m.Sender.ID].Data["region"],
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Id,
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Secret)
				if newErr != nil {
					log.Error(newErr)
				}
				quota, quotaErr := newRt.GetQuota(code[0], code[1])
				if quotaErr != nil {
					log.Error("Get Quota error: ", quotaErr)
					return
				}
				_, sendErr := bot.Send(m.Sender, "配额"+*quota.QuotaName+"的值为"+
					strconv.FormatFloat(*quota.Value, 'f', -1, 64))
				if sendErr != nil {
					log.Error("Send message error: ", sendErr)
				}
			case 11: //获取quota code和要提升至的值
				defer delete(p.State, m.Sender.ID)
				code := strings.Split(m.Text, " ")
				newRt, newErr := aws.New(p.State[m.Sender.ID].Data["region"],
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Id,
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Secret)
				if newErr != nil {
					log.Error(newErr)
				}
				des, parErr := strconv.ParseFloat(code[2], 64)
				if parErr != nil {
					_, sendErr := bot.Send(m.Sender, "修改失败")
					if sendErr != nil {
						log.Error("Send message error: ", sendErr)
					}
					log.Error("String to Float64 error: ", parErr)
					return
				}
				changeErr := newRt.ChangeQuota(code[0], code[1], des)
				if changeErr != nil {
					_, sendErr := bot.Send(m.Sender, "修改失败")
					if sendErr != nil {
						log.Error("Send message error: ", sendErr)
					}
					log.Error("Change quota error: ", changeErr)
					return
				}
				_, sendErr := bot.Send(m.Sender, "修改成功")
				if sendErr != nil {
					log.Error("Send message error: ", sendErr)
				}
			default: //直接跳出
				return
			}
		}
	})
}
