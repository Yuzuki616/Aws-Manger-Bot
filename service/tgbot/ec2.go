package tgbot

import (
	"github.com/338317/Aws-Manger-Bot/aws"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"os"
	"time"
)

const (
	//Aws Ami
	debian10    = "debian-10-amd64-20210329-591"
	ubuntu2004  = "ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-20210430"
	redhat8     = "RHEL_HA-8.4.0_HVM-20210504-x86_64-2-Hourly2-GP2"
	windows2019 = "Windows_Server-2019-English-Full-Base-2021.09.15"
	//Wavelength Zones
	tokyoWl  = "ap-northeast-1-wl1-nrt-wlz-1"
	seoulWl  = "ap-northeast-2-wl1-cjj-wlz-1"
	londonWL = "eu-west-2-wl1-lon-wlz-1"
	oregonWl = "us-west-2-wl1-phx-wlz-1"
)

func keySave(key string) string {
	tempName := time.Unix(time.Now().Unix(), 0).Format("./_2006-01-02_15:04:05.tmp")
	err := ioutil.WriteFile(tempName, []byte(key), 0644)
	if err != nil {
		log.Println("Save key file error:", err)
	}
	return tempName
} //缓存ssh密钥

func (p *TgBot) createEc2(bot *tb.Bot, c *tb.Callback) {
	if _, ok := p.Config.UserInfo[c.Sender.ID]; ok {
		_, err := bot.Edit(c.Message, "正在创建EC2...")
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		if p.Config.UserInfo[c.Sender.ID].NowKey == "" {
			_, err := bot.Edit(c.Message, "请先通过/KeyManger命令选择密钥")
			if err != nil {
				log.Error("Edit message error: ", err)
			}
		} else {
			awsO, newErr := aws.New(p.State[c.Sender.ID].Data["region"],
				p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Id,
				p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Secret)
			if newErr != nil {
				_, err := bot.Send(c.Sender, "创建失败!")
				if err != nil {
					log.Error("Send message error: ", err)
				}
				log.Error(newErr)
				return
			}
			var amiId string
			if _, ok := p.State[c.Sender.ID].Data["amiId"]; ok {
				amiId = p.State[c.Sender.ID].Data["amiId"]
			} else {
				amiTmp, amiErr := awsO.GetAmiId(p.State[c.Sender.ID].Data["ami"])
				if amiErr != nil {
					_, err := bot.Send(c.Sender, "创建失败!")
					if err != nil {
						log.Error("Send message error: ", err)
					}
					log.Error("Get ami ID error: ", amiErr)
					return
				}
				if amiTmp == "" {
					log.Error("Get ami ID error: Not found ami")
					return
				}
				amiId = amiTmp
			}
			var InstanceInfo *aws.Ec2Info
			if _, ok := p.State[c.Sender.ID].Data["zone"]; ok {
				sub, caErr := awsO.GetSubnetInfo()
				if caErr != nil {
					_, err := bot.Send(c.Sender, "创建失败!")
					if err != nil {
						log.Error("Send message error: ", err)
					}
					log.Error("Get gateway info error: ", caErr)
					return
				}
				var subId1 string
				if len(sub.Subnets) == 0 {
					subId, wlErr := awsO.CreateWl(p.State[c.Sender.ID].Data["zone"])
					if wlErr != nil {
						_, err := bot.Send(c.Sender, "创建失败!")
						if err != nil {
							log.Error("Send message error: ", err)
						}
						log.Error("Create wavelength error: ", wlErr)
						return
					}
					subId1 = subId
				} else {
					subId1 = *sub.Subnets[0].SubnetId
				}
				creRt, creErr := awsO.CreateEc2Wl(subId1, amiId, p.State[c.Sender.ID].Data["name"])
				if creErr != nil {
					_, err := bot.Send(c.Sender, "创建失败!")
					if err != nil {
						log.Error("Send message error: ", err)
					}
					log.Error("Create ec2wl error: ", creErr)
					return
				}
				InstanceInfo = creRt
			} else {
				creRt, creErr := awsO.CreateEc2(amiId,
					p.State[c.Sender.ID].Data["type"],
					p.State[c.Sender.ID].Data["name"])
				if creErr != nil {
					_, err := bot.Send(c.Sender, "创建失败!")
					if err != nil {
						log.Error("Send message error: ", err)
					}
					log.Error(creErr)
					return
				}
				InstanceInfo = creRt
			}
			_, err := bot.Send(c.Sender, "已添加到创建队列，正在等待创建...")
			if err != nil {
				log.Error("Send message error: ", err)
			}
			for true {
				getRt, getErr := awsO.GetEc2Info(*InstanceInfo.InstanceId)
				if getErr != nil {
					_, err := bot.Send(c.Sender, "获取实例信息失败！")
					if err != nil {
						log.Error("Send message error: ", err)
					}
					log.Error(getErr)
					return
				}
				if *getRt.Status == "running" {
					fileName := keySave(*InstanceInfo.Key)
					if getRt.Ip == nil {
						getRt.Ip = InstanceInfo.Ip
					}
					log.Info(getRt.Ip)
					_, err := bot.Send(c.Sender, "创建成功！\nUbuntu默认用户名ubuntu, Debian默认用户名admin\n\n实例信息: \n备注: "+*getRt.Name+
						"\n实例ID: "+*getRt.InstanceId+
						"\nIP: "+*getRt.Ip+"\nSSH密钥: ")
					if err != nil {
						log.Error("Send message error: ", err)
					}
					_, sendErr := bot.SendAlbum(c.Sender,
						tb.Album{&tb.Document{File: tb.FromDisk(fileName), FileName: *InstanceInfo.Name + "_key.pem"}})
					if sendErr != nil {
						log.Error("Send file error: ", sendErr)
					}
					removeErr := os.Remove(fileName)
					if removeErr != nil {
						log.Error("Remove temp file error: ", removeErr)
					}
					break
				}
				time.Sleep(time.Second * 3)
			}
		}
	} else {
		_, err := bot.Edit(c.Message, "请先通过/KeyManger命令添加密钥")
		if err != nil {
			log.Error("Edit message error: ", err)
		}
	}
} //创建Ec2

func (p *TgBot) listEc2(bot *tb.Bot, c *tb.Callback) {
	if _, ok := p.Config.UserInfo[c.Sender.ID]; ok {
		if p.Config.UserInfo[c.Sender.ID].NowKey == "" {
			_, err := bot.Edit(c.Message, "请先通过/KeyManger命令选择密钥")
			if err != nil {
				log.Println("Edit message error: ", err)
			}
		} else {
			delErr := bot.Delete(c.Message)
			if delErr != nil {
				log.Println("Delete message error: ", delErr)
			}
			aws0, newErr := aws.New(p.State[c.Sender.ID].Data["region"],
				p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Id,
				p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Secret)
			if newErr != nil {
				log.Println(newErr)
			}
			listRt, listErr := aws0.ListEc2()
			if listErr != nil {
				log.Println(listErr)
			}
			tmp := "已创建的Ec2实例: \n\n"
			for _, ec2Info := range listRt {
				if len(ec2Info.Instances[0].Tags) != 0 {
					tmp += "\n备注: " + *ec2Info.Instances[0].Tags[0].Value
				}
				tmp += "\n状态: " + *ec2Info.Instances[0].State.Name +
					"\n实例ID: " + *ec2Info.Instances[0].InstanceId
				if ec2Info.Instances[0].PublicIpAddress != nil {
					tmp += "\nIP: " + *ec2Info.Instances[0].PublicIpAddress
				}
			}
			_, err := bot.Send(c.Sender, tmp)
			if err != nil {
				log.Println("Send message error: ", err)
			}
			delete(p.State, c.Sender.ID)
		}
	} else {
		_, err := bot.Edit(c.Message, "请先通过/KeyManger命令添加密钥")
		if err != nil {
			log.Println("Edit message error: ", err)
		}
	}
}

func (p *TgBot) Ec2Manger(bot *tb.Bot) {
	amiKey := &tb.ReplyMarkup{}
	debian := amiKey.Data("Debian10", "debian10")
	p.AmiKey = amiKey
	bot.Handle(&debian, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["ami"] = debian10
		defer delete(p.State, c.Sender.ID)
		p.createEc2(bot, c)
	})
	ubuntu := amiKey.Data("Ubuntu20.04", "ubuntu2004")
	bot.Handle(&ubuntu, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["ami"] = ubuntu2004
		defer delete(p.State, c.Sender.ID)
		p.createEc2(bot, c)
	})
	redhat := amiKey.Data("Redhat8", "redhat8")
	bot.Handle(&redhat, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["ami"] = redhat8
		defer delete(p.State, c.Sender.ID)
		p.createEc2(bot, c)
	})
	windows := amiKey.Data("Windows2019", "windows2019")
	bot.Handle(&windows, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["ami"] = windows2019
		defer delete(p.State, c.Sender.ID)
		p.createEc2(bot, c)
	})
	otherAmi := amiKey.Data("其他系统", "other")
	bot.Handle(&otherAmi, func(c *tb.Callback) {
		_, editErr := bot.Edit(c.Message, "请输入Ami ID: ")
		if editErr != nil {
			log.Println(editErr)
		}
		p.State[c.Sender.ID].Parent = 8
	})
	amiKey.Inline(amiKey.Row(debian, ubuntu, redhat), amiKey.Row(windows), amiKey.Row(otherAmi))
	typeKey := &tb.ReplyMarkup{}
	p.TypeKey = typeKey
	t2 := typeKey.Data("t2.micro", "t2micro")
	bot.Handle(&t2, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["type"] = "t2.micro"
		_, err := bot.Edit(c.Message, "请选择操作系统", amiKey)
		if err != nil {
			log.Println("Edit message error: ", err)
		}
	})
	t3 := typeKey.Data("t3.micro", "t3micro")
	bot.Handle(&t3, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["type"] = "t3.micro"
		_, err := bot.Edit(c.Message, "请选择操作系统", amiKey)
		if err != nil {
			log.Println("Edit message error: ", err)
		}
	})
	otherType := typeKey.Data("其他类型", "otherType")
	bot.Handle(&otherType, func(c *tb.Callback) {
		_, editErr := bot.Edit(c.Message, "请输入ec2类型: ")
		if editErr != nil {
			log.Println(editErr)
		}
		p.State[c.Sender.ID].Parent = 9
	})
	typeKey.Inline(typeKey.Row(t2, t3), typeKey.Row(otherType))
	regionWl := &tb.ReplyMarkup{}
	tokyo := regionWl.Data("东京", "tokyo_wl")
	bot.Handle(&tokyo, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请输入备注:  ")
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.State[c.Sender.ID] = &State{
			Parent: 13,
			Data: map[string]string{
				"region": "ap-northeast-1",
				"zone":   tokyoWl,
			}}
	})
	seoul := regionWl.Data("首尔", "seoul_wl")
	bot.Handle(&seoul, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请选择类型: ", typeKey)
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.State[c.Sender.ID] = &State{
			Parent: 13,
			Data: map[string]string{
				"region": "ap-northeast-2",
				"zone":   seoulWl,
			}}
	})
	london := regionWl.Data("伦敦", "london_wl")
	bot.Handle(&london, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请选择类型: ", typeKey)
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.State[c.Sender.ID] = &State{
			Parent: 13,
			Data: map[string]string{
				"region": "eu-west-2",
				"zone":   londonWL,
			}}
	})
	oregon := regionWl.Data("俄勒冈", "oregon_wl")
	bot.Handle(&oregon, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请选择类型: ", typeKey)
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.State[c.Sender.ID] = &State{
			Parent: 13,
			Data: map[string]string{
				"region": "us-west-2",
				"zone":   oregonWl,
			}}
	})
	regionWl.Inline(regionWl.Row(tokyo, seoul), regionWl.Row(london, oregon))
	key := &tb.ReplyMarkup{}
	newEc2 := key.Data("创建EC2", "createEc2")
	bot.Handle(&newEc2, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请输入将要创建的ec2的备注: ")
		if err != nil {
			log.Println("Edit message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 5}
	})
	newEc2Wl := key.Data("创建Ec2Wl", "createEc2Wl")
	bot.Handle(&newEc2Wl, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请选择Wavelength地区: ", regionWl)
		if err != nil {
			log.Println("Edit message error: ", err)
		}
	})
	listEc2 := key.Data("列出EC2", "listEc2")
	bot.Handle(&listEc2, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请选择AWS区域: ", p.RegionKey)
		if err != nil {
			log.Println("Edit message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 101, Data: map[string]string{}}
	})
	stopEc2 := key.Data("暂停Ec2", "stopEc2")
	bot.Handle(&stopEc2, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请选择地区： ", p.RegionKey)
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 108, Data: map[string]string{}}
	})
	delEc2 := key.Data("删除Ec2", "delEc2")
	bot.Handle(&delEc2, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请选择地区: ", p.RegionKey)
		if err != nil {
			log.Println("Edit message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 102, Data: map[string]string{}}
	})
	chIp := key.Data("更换IP", "changeIp")
	bot.Handle(&chIp, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请选择地区: ", p.RegionKey)
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 103, Data: map[string]string{}}
	})
	getPassword := key.Data("提取Windows密码", "get_password")
	bot.Handle(&getPassword, func(c *tb.Callback) {
		_, editErr := bot.Edit(c.Message, "请选择地区: ", p.RegionKey)
		if editErr != nil {
			log.Error("Edit message error: ", editErr)
		}
		p.State[c.Sender.ID] = &State{Parent: 107, Data: map[string]string{}}
	})
	key.Inline(key.Row(newEc2, listEc2), key.Row(newEc2Wl), key.Row(getPassword), key.Row(delEc2, chIp))
	bot.Handle("/Ec2Manger", func(m *tb.Message) {
		if m.Private() {
			_, err := bot.Send(m.Sender, "请选择你要进行的操作", key)
			if err != nil {
				log.Println("Send message error: ", err)
			}
		} else {
			_, err := bot.Send(m.Sender, "请私聊Bot使用")
			if err != nil {
				log.Println("Send message error: ", err)
			}
		}
	})
}
