package tgbot

import (
	"github.com/338317/Aws-Manger-Bot/aws"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	debian10   = "ami-0c7ea5497c02abcaf"
	ubuntu2004 = "ami-0df99b3a8349462c6"
	centos8    = "ami-000eaef4896b0e4dc"
)

func keySave(key string) string {
	tmepName := time.Unix(time.Now().Unix(), 0).Format("./_2006-01-02_15:04:05.tmp")
	err := ioutil.WriteFile(tmepName, []byte(key), 0644)
	if err != nil {
		log.Println("Save key file error:", err)
	}
	return tmepName
}

func (p *TgBot) createEc2(bot *tb.Bot, c *tb.Callback) {
	if _, ok := p.Config.UserInfo[c.Sender.ID]; ok {
		_, err := bot.Edit(c.Message, "正在创建EC2...")
		if err != nil {
			log.Println("Edit message error: ", err)
		}
		if p.Config.UserInfo[c.Sender.ID].NowKey == "" {
			_, err := bot.Edit(c.Message, "请先通过/KeyManger命令选择密钥")
			if err != nil {
				log.Println("Edit message error: ", err)
			}
		} else {
			awsO, newErr := aws.New(p.State[c.Sender.ID].Data["region"],
				p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Id,
				p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Secret)
			if newErr != nil {
				_, err := bot.Send(c.Sender, "创建失败!")
				if err != nil {
					log.Println("Send message error: ", err)
				}
				log.Println(newErr)
				return
			}
			creRt, creErr := awsO.CreateEc2(p.State[c.Sender.ID].Data["ami"],
				p.State[c.Sender.ID].Data["type"],
				p.State[c.Sender.ID].Data["name"])
			if creErr != nil {
				_, err := bot.Send(c.Sender, "创建失败!")
				if err != nil {
					log.Println("Send message error: ", err)
				}
				log.Println(creErr)
				return
			}
			_, err := bot.Send(c.Sender, "已添加到创建队列，正在等待创建...")
			if err != nil {
				log.Println("Send message error: ", err)
			}
			for true {
				getRt, getErr := awsO.GetEc2Info(*creRt.InstanceId)
				if getErr != nil {
					_, err := bot.Send(c.Sender, "获取实例信息失败！")
					if err != nil {
						log.Println("Send message error: ", err)
					}
					log.Println(getErr)
					return
				}
				if *getRt.Status == "running" {
					fileName := keySave(*creRt.Key)
					_, err := bot.Send(c.Sender, "创建成功！\n实例信息: \n备注: "+*getRt.Name+
						"\n实例ID: "+*getRt.InstanceId+
						"\nIP: "+*getRt.Ip+"\nSSH密钥: ")
					if err != nil {
						log.Println("Send message error: ", err)
					}
					_, sendErr := bot.SendAlbum(c.Sender,
						tb.Album{&tb.Document{File: tb.FromDisk(fileName), FileName: *creRt.Name + "_key.pem"}})
					if sendErr != nil {
						log.Println("Send file error: ", sendErr)
					}
					removeErr := os.Remove(fileName)
					if removeErr != nil {
						log.Println("Remove temp file error: ", removeErr)
					}
					break
				}
				time.Sleep(time.Second * 3)
			}
		}
	} else {
		_, err := bot.Edit(c.Message, "请先通过/KeyManger命令添加密钥")
		if err != nil {
			log.Println("Edit message error: ", err)
		}
	}
}

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
	debian := amiKey.Data("Debian10", debian10)
	p.AmiKey = amiKey
	bot.Handle(&debian, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["ami"] = debian10
		defer delete(p.State, c.Sender.ID)
		p.createEc2(bot, c)
	})
	ubuntu := amiKey.Data("Ubuntu20.04", ubuntu2004)
	bot.Handle(&ubuntu, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["ami"] = ubuntu2004
		defer delete(p.State, c.Sender.ID)
		p.createEc2(bot, c)
	})
	centos := amiKey.Data("Centos8", centos8)
	bot.Handle(&centos, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["ami"] = centos8
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
	amiKey.Inline(amiKey.Row(debian, ubuntu, centos), amiKey.Row(otherAmi))
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
	key := &tb.ReplyMarkup{}
	newEc2 := key.Data("创建EC2", "createEc2")
	bot.Handle(&newEc2, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请输入将要创建的ec2的备注: ")
		if err != nil {
			log.Println("Edit message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 5}
	})
	listEc2 := key.Data("列出EC2", "listEc2")
	bot.Handle(&listEc2, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请选择AWS区域: ", p.RegionKey)
		if err != nil {
			log.Println("Edit message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 101, Data: map[string]string{}}
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
			log.Println("Edit message error: ", err)
		}
		p.State[c.Sender.ID] = &State{Parent: 103, Data: map[string]string{}}
	})
	key.Inline(key.Row(newEc2, listEc2), key.Row(delEc2, chIp))
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
