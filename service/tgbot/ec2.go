package tgbot

import (
	"github.com/Yuzuki999/Aws-Manger-Bot/aws"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
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
			awsO, newErr := aws.New(p.Data[c.Sender.ID].Data["region"],
				p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Id,
				p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Secret,
				p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Proxy)
			if newErr != nil {
				_, err := bot.Send(c.Sender, "创建失败!")
				if err != nil {
					log.Error("Send message error: ", err)
				}
				log.Error(newErr)
				return
			}
			diskSize, parErr := strconv.ParseInt(p.Data[c.Sender.ID].Data["disk"], 10, 64)
			if parErr != nil {
				log.Error("String to int64 error: ", parErr)
			}
			var amiId string
			if _, ok := p.Data[c.Sender.ID].Data["amiId"]; ok {
				amiId = p.Data[c.Sender.ID].Data["amiId"]
			} else {
				amiTmp, amiErr := awsO.GetAmiId(p.Data[c.Sender.ID].Data["ami"])
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
			if _, ok := p.Data[c.Sender.ID].Data["zone"]; ok {
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
					subId, wlErr := awsO.CreateWl(p.Data[c.Sender.ID].Data["zone"])
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
				creRt, creErr := awsO.CreateEc2Wl(subId1, amiId, p.Data[c.Sender.ID].Data["region"], diskSize)
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
					p.Data[c.Sender.ID].Data["type"],
					p.Data[c.Sender.ID].Data["name"],
					diskSize)
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
					log.Error("Get ec2 info error:", getErr)
					return
				}
				if *getRt.Status == "running" {
					if getRt.Ip == nil {
						getRt.Ip = InstanceInfo.Ip
					}
					_, err := bot.Send(c.Sender, "创建成功！\nUbuntu默认用户名ubuntu, Debian默认用户名admin\n\n实例信息: \n备注: "+*getRt.Name+
						"\n实例ID: "+*getRt.InstanceId+
						"\nIP: "+*getRt.Ip+"\nSSH密钥: \n\n"+*InstanceInfo.Key)
					if err != nil {
						log.Error("Send message error: ", err)
					}
					break
				}
				time.Sleep(time.Second * 3)
			}
		}
	} else {
		_, err := bot.Send(c.Sender, "请先通过/KeyManger命令添加密钥")
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
			aws0, newErr := aws.New(p.Data[c.Sender.ID].Data["region"],
				p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Id,
				p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Secret,
				p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Proxy)
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
			delete(p.Data, c.Sender.ID)
		}
	} else {
		_, err := bot.Edit(c.Message, "请先通过/KeyManger命令添加密钥")
		if err != nil {
			log.Println("Edit message error: ", err)
		}
	}
}

func (p *TgBot) Ec2Manger(bot *tb.Bot) {
	diskKey := &tb.ReplyMarkup{}
	G8 := diskKey.Data("8GB", "8")
	bot.Handle(&G8, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["disk"] = "8"
		p.createEc2(bot, c)
	})
	G32 := diskKey.Data("32GB", "32")
	bot.Handle(&G32, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["disk"] = "32"
		p.createEc2(bot, c)
	})
	G64 := diskKey.Data("64GB", "64")
	bot.Handle(&G64, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["disk"] = "64"
		p.createEc2(bot, c)
	})
	otherG := diskKey.Data("自定义大小", "Other")
	bot.Handle(&otherG, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请输入硬盘大小: ")
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
			p.Data[c.Sender.ID].Data["disk"] = m.Text
			p.Session[c.Sender.ID].Channel <- true
		})
		select {
		case tmp := <-p.Session[c.Sender.ID].Channel:
			if tmp != true {
				return
			}
			p.createEc2(bot, c)
			p.Session.SessionDel(c.Sender.ID)
		case <-time.After(30 * time.Second):
			_, editErr := bot.Edit(c.Message, "操作超时")
			if editErr != nil {
				log.Error("Edit message error: ", editErr)
			}
			p.Session.SessionDel(c.Sender.ID)
		}
	})
	diskKey.Inline(diskKey.Row(G8, G32, G64), diskKey.Row(otherG))
	amiKey := &tb.ReplyMarkup{}
	debian := amiKey.Data("Debian10", "debian10")
	p.AmiKey = amiKey
	bot.Handle(&debian, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["ami"] = debian10
		_, err := bot.Edit(c.Message, "请选择硬盘大小", diskKey)
		if err != nil {
			log.Error("Edit message error: ", err)
		}
	})
	ubuntu := amiKey.Data("Ubuntu20.04", "ubuntu2004")
	bot.Handle(&ubuntu, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["ami"] = ubuntu2004
		_, err := bot.Edit(c.Message, "请选择硬盘大小", diskKey)
		if err != nil {
			log.Error("Edit message error: ", err)
		}
	})
	redhat := amiKey.Data("Redhat8", "redhat8")
	bot.Handle(&redhat, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["ami"] = redhat8
		_, err := bot.Edit(c.Message, "请选择硬盘大小", diskKey)
		if err != nil {
			log.Error("Edit message error: ", err)
		}
	})
	windows := amiKey.Data("Windows2019", "windows2019")
	bot.Handle(&windows, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["ami"] = windows2019
		_, err := bot.Edit(c.Message, "请选择硬盘大小", diskKey)
		if err != nil {
			log.Error("Edit message error: ", err)
		}
	})
	otherAmi := amiKey.Data("其他系统", "other")
	bot.Handle(&otherAmi, func(c *tb.Callback) {
		_, editErr := bot.Edit(c.Message, "请输入Ami ID: ")
		if editErr != nil {
			log.Println(editErr)
		}
		p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
			p.Data[c.Sender.ID].Data["ami"] = m.Text
			p.Session[c.Sender.ID].Channel <- true
		})
		select {
		case tmp := <-p.Session[c.Sender.ID].Channel:
			if tmp != true {
				return
			}
			_, err := bot.Send(c.Sender, "请选择硬盘大小:", diskKey)
			if err != nil {
				log.Error("Edit message error: ", err)
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
	amiKey.Inline(amiKey.Row(debian, ubuntu, redhat), amiKey.Row(windows), amiKey.Row(otherAmi))
	typeKey := &tb.ReplyMarkup{}
	p.TypeKey = typeKey
	t2 := typeKey.Data("t2.micro", "t2micro")
	bot.Handle(&t2, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["type"] = "t2.micro"
		_, err := bot.Edit(c.Message, "请选择操作系统", amiKey)
		if err != nil {
			log.Println("Edit message error: ", err)
		}
	})
	t3 := typeKey.Data("t3.micro", "t3micro")
	bot.Handle(&t3, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["type"] = "t3.micro"
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
		p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
			p.Data[c.Sender.ID].Data["type"] = m.Text
			p.Session[c.Sender.ID].Channel <- true
		})
		select {
		case tmp := <-p.Session[c.Sender.ID].Channel:
			if tmp != true {
				return
			}
			_, err := bot.Send(c.Sender, "请选择操作系统", amiKey)
			if err != nil {
				log.Error("Send message error: ", err)
				p.Session.SessionDel(c.Sender.ID)
			} else {
				return
			}
		case <-time.After(30 * time.Second):
			_, sendErr := bot.Edit(c.Message, "操作超时")
			if sendErr != nil {
				log.Error("Edit message error: ", sendErr)
			}
			p.Session.SessionDel(c.Sender.ID)
		}
	})
	typeKey.Inline(typeKey.Row(t2, t3), typeKey.Row(otherType))
	regionWl := &tb.ReplyMarkup{}
	tokyo := regionWl.Data("东京", "tokyo_wl")
	bot.Handle(&tokyo, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data = map[string]string{
			"region": "ap-northeast-1",
			"zone":   tokyoWl,
			"type":   "t3.medium",
		}
		_, sendErr := bot.Send(c.Sender, "请选择Ami: ", amiKey)
		if sendErr != nil {
			log.Println("Send message error: ", sendErr)
		}
	})
	seoul := regionWl.Data("首尔", "seoul_wl")
	bot.Handle(&seoul, func(c *tb.Callback) {
		p.Data[c.Sender.ID] = &Data{
			Data: map[string]string{
				"region": "ap-northeast-2",
				"zone":   seoulWl,
				"type":   "t3.medium",
			}}
		_, sendErr := bot.Send(c.Sender, "请选择Ami: ", amiKey)
		if sendErr != nil {
			log.Println("Send message error: ", sendErr)
		}
	})
	london := regionWl.Data("伦敦", "london_wl")
	bot.Handle(&london, func(c *tb.Callback) {
		p.Data[c.Sender.ID] = &Data{
			Data: map[string]string{
				"region": "eu-west-2",
				"zone":   londonWL,
				"type":   "t3.medium",
			}}
		_, sendErr := bot.Send(c.Sender, "请选择Ami: ", amiKey)
		if sendErr != nil {
			log.Println("Send message error: ", sendErr)
		}
	})
	oregon := regionWl.Data("俄勒冈", "oregon_wl")
	bot.Handle(&oregon, func(c *tb.Callback) {
		p.Data[c.Sender.ID] = &Data{
			Data: map[string]string{
				"region": "us-west-2",
				"zone":   oregonWl,
				"type":   "t3.medium",
			}}
		_, sendErr := bot.Send(c.Sender, "请选择Ami: ", amiKey)
		if sendErr != nil {
			log.Println("Send message error: ", sendErr)
		}
	})
	regionWl.Inline(regionWl.Row(tokyo, seoul), regionWl.Row(london, oregon))
	key := &tb.ReplyMarkup{}
	newEc2 := key.Data("创建EC2", "createEc2")
	bot.Handle(&newEc2, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action:  Create ec2")
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
		_, editErr2 := bot.Edit(c.Message, "请选择地区: ", p.RegionKey)
		if editErr2 != nil {
			log.Println("Edit message error: ", editErr2)
		}
		p.Data[c.Sender.ID].RegionChan = make(chan int)
		select {
		case <-p.Data[c.Sender.ID].RegionChan:
			_, err := bot.Edit(c.Message, "请输入将要创建的ec2的备注: ")
			if err != nil {
				log.Println("Edit message error: ", err)
			}
			p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
				p.Data[c.Sender.ID].Data["name"] = m.Text
				p.Session[c.Sender.ID].Channel <- true
			})
			select {
			case tmp := <-p.Session[c.Sender.ID].Channel:
				if tmp != true {
					return
				}
				p.Session.SessionDel(c.Sender.ID)
				_, editErr := bot.Edit(c.Message, "请选择EC2类型", p.TypeKey)
				if editErr != nil {
					log.Error("Edit message error: ", editErr)
				}
			case <-time.After(30 * time.Second):
				_, sendErr := bot.Edit(c.Message, "操作超时")
				if sendErr != nil {
					log.Error("Edit message error: ", sendErr)
				}
				p.Session.SessionDel(c.Sender.ID)
			}

		case <-time.After(30 * time.Second):
			_, editErr := bot.Edit(c.Message, "操作超时")
			if editErr != nil {
				log.Error("Edit message error: ", editErr)
			}
		}
	})
	newEc2Wl := key.Data("创建Ec2Wl", "createEc2Wl")
	bot.Handle(&newEc2Wl, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action:  Create ec2 wavelength")
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
		_, editErr := bot.Edit(c.Message, "请输入将要创建的ec2的备注: ")
		if editErr != nil {
			log.Error("Edit message error: ", editErr)
		}
		p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
			p.Data[c.Sender.ID].Data["name"] = m.Text
			p.Session[c.Sender.ID].Channel <- true
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
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
		_, editErr2 := bot.Edit(c.Message, "请选择Wavelength地区: ", regionWl)
		if editErr2 != nil {
			log.Error("Edit message error: ", editErr2)
		}
	})
	listEc2 := key.Data("列出EC2", "listEc2")
	bot.Handle(&listEc2, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action:  List ec2")
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
		_, err := bot.Edit(c.Message, "请选择AWS区域: ", p.RegionKey)
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.Data[c.Sender.ID].RegionChan = make(chan int)
		select {
		case <-p.Data[c.Sender.ID].RegionChan:
			p.listEc2(bot, c)
		case <-time.After(30 * time.Second):
			_, editErr := bot.Edit(c.Message, "操作超时")
			if editErr != nil {
				log.Error("Edit message error: ", editErr)
			}
		}
	})
	stopEc2 := key.Data("暂停Ec2", "stopEc2")
	bot.Handle(&stopEc2, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action:  Stop ec2")
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
		_, err := bot.Edit(c.Message, "请选择地区： ", p.RegionKey)
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.Data[c.Sender.ID].RegionChan = make(chan int)
		select {
		case <-p.Data[c.Sender.ID].RegionChan:
			_, editErr := bot.Edit(c.Message, "请输入实例ID: ")
			if editErr != nil {
				log.Error("Edit message error: ", editErr)
			}
			p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
				defer delete(p.Data, m.Sender.ID)
				awsRt, awsErr := aws.New(p.Data[m.Sender.ID].Data["region"],
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
		case <-time.After(30 * time.Second):
			_, editErr := bot.Edit(c.Message, "操作超时")
			if editErr != nil {
				log.Error("Edit message error: ", editErr)
			}
		}
	})
	startEc2 := key.Data("启动Ec2", "startEc2")
	bot.Handle(&startEc2, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action: Start ec2")
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
		_, err := bot.Edit(c.Message, "请选择地区： ", p.RegionKey)
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.Data[c.Sender.ID].RegionChan = make(chan int)
		select {
		case <-p.Data[c.Sender.ID].RegionChan:
			_, editErr := bot.Edit(c.Message, "请输入要启动的实例ID: ")
			if editErr != nil {
				log.Error("Edit message error: ", editErr)
			}
			p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
				defer delete(p.Data, m.Sender.ID)
				awsRt, awsErr := aws.New(p.Data[m.Sender.ID].Data["region"],
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
				p.Session[m.Sender.ID].Channel <- true
			})
			select {
			case tmp := <-p.Session[c.Sender.ID].Channel:
				p.Session.SessionDel(c.Sender.ID)
				if tmp != true {
					return
				}
			case <-time.After(30 * time.Second):
				_, sendErr := bot.Edit(c.Message, "操作超时")
				if sendErr != nil {
					log.Error("Edit message error: ", sendErr)
				}
				p.Session.SessionDel(c.Sender.ID)
			}
		case <-time.After(30 * time.Second):
			_, editErr := bot.Edit(c.Message, "操作超时")
			if editErr != nil {
				log.Error("Edit message error: ", editErr)
			}
		}
	})
	delEc2 := key.Data("删除Ec2", "delEc2")
	bot.Handle(&delEc2, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action: Delete ec2")
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
		_, err := bot.Edit(c.Message, "请选择地区: ", p.RegionKey)
		if err != nil {
			log.Println("Edit message error: ", err)
		}
		p.Data[c.Sender.ID].RegionChan = make(chan int)
		select {
		case <-p.Data[c.Sender.ID].RegionChan:
			_, editErr := bot.Edit(c.Message, "请输入要删除的实例ID: ")
			if editErr != nil {
				log.Error("Send message error: ", editErr)
			}
			p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
				defer delete(p.Data, m.Sender.ID)
				newRt, newErr := aws.New(p.Data[m.Sender.ID].Data["region"],
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
		case <-time.After(30 * time.Second):
			_, editErr := bot.Edit(c.Message, "操作超时")
			if editErr != nil {
				log.Error("Edit message error: ", editErr)
			}
		}
	})
	chIp := key.Data("更换IP", "changeIp")
	bot.Handle(&chIp, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action: Change Ip")
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
		_, editErr := bot.Edit(c.Message, "请选择地区: ", p.RegionKey)
		if editErr != nil {
			log.Error("Edit message error: ", editErr)
		}
		p.Data[c.Sender.ID].RegionChan = make(chan int)
		select {
		case <-p.Data[c.Sender.ID].RegionChan:
			_, editErr := bot.Edit(c.Message, "请输入要更换IP的实例ID: ")
			if editErr != nil {
				log.Error("Send message error: ", editErr)
			}
			p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
				defer delete(p.Data, m.Sender.ID)
				newRt, newErr := aws.New(p.Data[m.Sender.ID].Data["region"],
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
		case <-time.After(30 * time.Second):
			_, editErr := bot.Edit(c.Message, "操作超时")
			if editErr != nil {
				log.Error("Edit message error: ", editErr)
			}
		}
	})
	getPassword := key.Data("提取Windows密码", "get_password")
	bot.Handle(&getPassword, func(c *tb.Callback) {
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action: Get windows password")
		_, editErr := bot.Edit(c.Message, "请选择地区: ", p.RegionKey)
		if editErr != nil {
			log.Error("Edit message error: ", editErr)
		}
		p.Data[c.Sender.ID].RegionChan = make(chan int)
		select {
		case <-p.Data[c.Sender.ID].RegionChan:
			_, err := bot.Edit(c.Message, "请输入实例ID: ")
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
		case <-time.After(30 * time.Second):
			_, editErr := bot.Edit(c.Message, "操作超时")
			if editErr != nil {
				log.Error("Edit message error: ", editErr)
			}
		}
	})
	key.Inline(key.Row(newEc2, listEc2), key.Row(newEc2Wl), key.Row(getPassword),
		key.Row(stopEc2, startEc2), key.Row(delEc2, chIp))
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
