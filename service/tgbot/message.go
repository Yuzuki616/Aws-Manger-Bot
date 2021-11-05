package tgbot

import (
	"github.com/Yuzuki999/Aws-Manger-Bot/aws"
	"github.com/Yuzuki999/Aws-Manger-Bot/conf"
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
				if p.Config.UserInfo[m.Sender.ID].NowKey == "" {
					p.Config.UserInfo[m.Sender.ID].NowKey = p.State[m.Sender.ID].Data["name"]
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
			case 4: //切换密钥
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
			case 5: //获取ec2备注
				p.State[m.Sender.ID].Data = make(map[string]string)
				p.State[m.Sender.ID].Data["name"] = m.Text
				_, err := bot.Send(m.Sender, "请选择地区: ", p.RegionKey)
				if err != nil {
					log.Println("Send message error: ", err)
				}
				p.State[m.Sender.ID].Parent = 999
			case 6: //删除Ec2
				defer delete(p.State, m.Sender.ID)
				newRt, newErr := aws.New(p.State[m.Sender.ID].Data["region"],
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Id,
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Secret,
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Proxy)
				if newErr != nil {
					_, sendErr := bot.Send(m.Sender, "删除失败！")
					if sendErr != nil {
						log.Error("Send message error: ", sendErr)
					}
					log.Error("Init aws obj error: ", sendErr)
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
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Secret,
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Proxy)
				if newErr != nil {
					_, sendErr := bot.Send(m.Sender, "更换失败！")
					if sendErr != nil {
						log.Error("Send message error: ", sendErr)
					}
					log.Error("Init aws obj error: ", sendErr)
					return
				}
				ip, chErr := newRt.ChangeEc2Ip(m.Text)
				if chErr != nil {
					log.Error("Change ip error: ", chErr)
					_, sendErr := bot.Send(m.Sender, "更换失败！")
					if sendErr != nil {
						log.Error("Send message error: ", sendErr)
					}
					return
				}
				_, sendErr := bot.Send(m.Sender, "更换ip成功，新的IP为: "+*ip)
				if sendErr != nil {
					log.Error("Send message error: ", sendErr)
				}
			case 8: //获取AMI
				p.State[m.Sender.ID].Data["amiID"] = m.Text
				_, err := bot.Send(m.Sender, "请输入硬盘大小: ")
				if err != nil {
					log.Error("Edit message error: ", err)
				}
				p.State[m.Sender.ID].Parent = 19
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
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Secret,
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Proxy)
				if newErr != nil {
					_, sendErr := bot.Send(m.Sender, "查看失败")
					if sendErr != nil {
						log.Error("Send message error: ", sendErr)
					}
					log.Error("Init aws obj error: ", newErr)
					return
				}
				quota, quotaErr := newRt.GetQuota(code[0], code[1])
				if quotaErr != nil {
					_, sendErr := bot.Send(m.Sender, "查看失败")
					if sendErr != nil {
						log.Error("Send message error: ", sendErr)
					}
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
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Secret,
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Proxy)
				if newErr != nil {
					_, sendErr := bot.Send(m.Sender, "修改失败")
					if sendErr != nil {
						log.Error("Send message error: ", sendErr)
					}
					log.Error("Init aws obj error: ", newErr)
					return
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
			case 12:
				defer delete(p.State, m.Sender.ID)
				awsRt, awsErr := aws.New(p.State[m.Sender.ID].Data["region"],
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Id,
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Secret,
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Proxy)
				if awsErr != nil {
					_, sendErr := bot.Send(m.Sender, "获取失败")
					if sendErr != nil {
						log.Error("Send message error: ", sendErr)
					}
					log.Error("Init aws oub error: ", awsErr)
					return
				}
				pass, passErr := awsRt.GetWindowsPassword(m.Text)
				if passErr != nil {
					log.Error("Get windows Password error: ", passErr)
					_, sendErr := bot.Send(m.Sender, "获取失败")
					if sendErr != nil {
						log.Error("Send message error: ", sendErr)
					}
					return
				}
				_, sendErr := bot.Send(m.Sender, *pass.PasswordData+"\n\n\n以上为RSA加密后的密码，请自行使用ssh密钥解密")
				if sendErr != nil {
					log.Error("Send message error: ", sendErr)
				}
			case 13: //获取Ec2 Wavelength备注
				p.State[m.Sender.ID].Data["name"] = m.Text
				p.State[m.Sender.ID].Data["type"] = "t3.medium"
				_, err := bot.Send(m.Sender, "请选择Ami: ", p.AmiKey)
				if err != nil {
					log.Println("Send message error: ", err)
				}
				p.State[m.Sender.ID].Parent = 999
			case 14: //暂停Ec2
				defer delete(p.State, m.Sender.ID)
				awsRt, awsErr := aws.New(p.State[m.Sender.ID].Data["region"],
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Id,
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Secret,
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Proxy)
				if awsErr != nil {
					_, sendErr := bot.Send(m.Sender, "暂停失败")
					if sendErr != nil {
						log.Error("Send message error: ", sendErr)
					}
					log.Error("Init aws oub error: ", awsErr)
					return
				}
				stopErr := awsRt.StopEc2(m.Text)
				if stopErr != nil {
					_, sendErr := bot.Send(m.Sender, "暂停失败")
					if sendErr != nil {
						log.Error("Send message error: ", sendErr)
					}
					log.Error("Stop ec2 error: ", stopErr)
					return
				}
				_, sendErr := bot.Send(m.Sender, "暂停成功")
				if sendErr != nil {
					log.Error("Send message error: ", sendErr)
				}
			case 15: //获取aga备注
				p.State[m.Sender.ID].Data = make(map[string]string)
				p.State[m.Sender.ID].Data["agaName"] = m.Text
				_, err := bot.Send(m.Sender, "请选择地区: ", p.RegionKey)
				if err != nil {
					log.Println("Send message error: ", err)
				}
				p.State[m.Sender.ID].Parent = 109
			case 16: //创建Aga
				defer delete(p.State, m.Sender.ID)
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
				creRt, creErr := awsRt.CreateAga(p.State[m.Sender.ID].Data["agaName"],
					p.State[m.Sender.ID].Data["region"],
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
			case 17: //启动Ec2
				defer delete(p.State, m.Sender.ID)
				awsRt, awsErr := aws.New(p.State[m.Sender.ID].Data["region"],
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Id,
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Secret,
					p.Config.UserInfo[m.Sender.ID].AwsSecret[p.Config.UserInfo[m.Sender.ID].NowKey].Proxy)
				if awsErr != nil {
					_, sendErr := bot.Send(m.Sender, "启动失败")
					if sendErr != nil {
						log.Error("Send message error: ", sendErr)
					}
					log.Error("Init aws oub error: ", awsErr)
					return
				}
				startErr := awsRt.StartEc2(m.Text)
				if startErr != nil {
					_, sendErr := bot.Send(m.Sender, "启动失败")
					if sendErr != nil {
						log.Error("send message error: ", sendErr)
					}
					log.Error("Start ec2 error: ", startErr)
					return
				}
				_, sendErr := bot.Send(m.Sender, "启动成功")
				if sendErr != nil {
					log.Error("send message error: ", sendErr)
				}
			case 18: //删除Aga
				defer delete(p.State, m.Sender.ID)
				awsRt, awsErr := aws.New(p.State[m.Sender.ID].Data["region"],
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
			case 19: //获取硬盘大小
				p.State[m.Sender.ID].Data["disk"] = m.Text
				defer delete(p.State, m.Sender.ID)
				p.createEc2(bot, &tb.Callback{Message: m})
			case 20: //获取代理地址
				p.State[m.Sender.ID].Data["name"] = m.Text
				_, err := bot.Send(m.Sender, "请输入代理地址: ")
				if err != nil {
					log.Error("Send message error: ", err)
				}
				p.State[m.Sender.ID].Parent = 21
			case 21: //添加代理地址
				defer delete(p.State, m.Sender.ID)
				p.Config.UserInfo[m.Sender.ID].AwsSecret[p.State[m.Sender.ID].Data["name"]].Proxy = m.Text
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
			case 22: //删除代理地址
				defer delete(p.State, m.Sender.ID)
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
			default: //直接跳出
				return
			}
		}
	})
}
