package tgbot

import (
	"time"

	"github.com/Yuzuki999/Aws-Manger-Bot/aws"
	log "github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
)

const (
	ubuntu2004Ls = "ubuntu_20_04"
	debian10Ls   = "debian_10"
)

const (
	nano    = "nano_"
	micro   = "micro_"
	small   = "small_"
	medium  = "medium_"
	large   = "large_"
	xlarge  = "xlarge_"
	dxlarge = "2xlarge_"
)

var zone = map[string]string{
	"ap-northeast-1": "2_0",
	"ap-northeast-2": "2_0",
	"ap-south-1":     "2_1",
	"ap-southeast-1": "2_0",
	"ap-southeast-2": "2_2",
	"ca-central-1":   "2_0",
	"eu-central-1":   "2_0",
	"eu-west-1":      "2_0",
	"eu-west-2":      "2_0",
	"eu-west-3":      "2_0",
	"us-east-1":      "2_0",
	"us-east-2":      "2_0",
	"us-west-2":      "2_0",
}
var typeKey *telebot.ReplyMarkup

func (p *TgBot) lightSailInit(bot *telebot.Bot) *telebot.ReplyMarkup {
	key := &telebot.ReplyMarkup{}
	var r []telebot.Row
	var t telebot.Row
	i := 0
	for k := range zone {
		tmp := key.Data(k, k)
		bot.Handle(&tmp, func(c *telebot.Callback) {
			p.Data[c.Sender.ID].Data["zone"] = k
			p.Data[c.Sender.ID].RegionChan <- 0
		})
		i++
		t = append(t, tmp)
		if i == 3 {
			r = append(r, t)
			t = telebot.Row{}
			i = 0
		}
	}
	key.Inline(r...)
	return key
}

func (p *TgBot) createLightsail(bot *telebot.Bot, c *telebot.Callback) {
	awsO, newErr := aws.New(p.Data[c.Sender.ID].Data["zone"],
		p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Id,
		p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Secret,
		p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Proxy)
	if newErr != nil {
		log.Error("Init aws sdk error: ", newErr)
		_, err := bot.Send(c.Sender, "创建失败!")
		if err != nil {
			log.Error("Send message error: ", err)
		}
		return
	}
	lsRt, lsErr := awsO.CreateLs(p.Data[c.Sender.ID].Data["name"], p.Data[c.Sender.ID].Data["zone"]+"a", p.Data[c.Sender.ID].Data["blueprint"],
		p.Data[c.Sender.ID].Data["type"]+zone[p.Data[c.Sender.ID].Data["zone"]])
	if lsErr != nil {
		log.Error("Create lightsail error: ", lsErr)
		_, err := bot.Send(c.Sender, "创建失败!")
		if err != nil {
			log.Error("Send message error: ", err)
		}
		return
	}
	lsInfo, infoErr := awsO.GetLsInfo(*lsRt.Name)
	if infoErr != nil {
		log.Error("Get lightsail info error: ", infoErr)
	}
	_, sendErr := bot.Send(c.Sender, "创建成功!\n备注: "+
		*lsInfo.Name+"\n状态: "+
		*lsInfo.Status+"IP\n"+
		*lsInfo.Ip)
	if sendErr != nil {
		log.Error("Send message error: ", sendErr)
	}
}

func (p *TgBot) LightSailManger(bot *telebot.Bot) {
	zoneKey := p.lightSailInit(bot)
	blueprintKey := &telebot.ReplyMarkup{}
	debian10 := blueprintKey.Data("Debian10", debian10Ls)
	bot.Handle(&debian10, func(c *telebot.Callback) {
		p.Data[c.Sender.ID].Data["blueprint"] = debian10Ls
		p.createLightsail(bot, c)
	})
	ubuntu2004 := blueprintKey.Data("Ubuntu2004", ubuntu2004Ls)
	bot.Handle(&ubuntu2004, func(c *telebot.Callback) {
		p.Data[c.Sender.ID].Data["blueprint"] = debian10Ls
		p.createLightsail(bot, c)
	})
	blueprintKey.Inline(blueprintKey.Row(debian10, ubuntu2004))
	typeKey = &telebot.ReplyMarkup{}
	nanoKey := typeKey.Data("1c0.5g", nano)
	bot.Handle(&nanoKey, func(c *telebot.Callback) {
		p.Data[c.Sender.ID].Data["type"] = nano
		_, editErr := bot.Edit(c.Message, "请选择操作系统", blueprintKey)
		if editErr != nil {
			log.Error("Edit message error", editErr)
		}
	})
	microKey := typeKey.Data("1c1g", micro)
	bot.Handle(&microKey, func(c *telebot.Callback) {
		p.Data[c.Sender.ID].Data["type"] = micro
		_, editErr := bot.Edit(c.Message, "请选择操作系统", blueprintKey)
		if editErr != nil {
			log.Error("Edit message error", editErr)
		}
	})
	smallKey := typeKey.Data("1c2g", small)
	bot.Handle(&smallKey, func(c *telebot.Callback) {
		p.Data[c.Sender.ID].Data["type"] = small
		_, editErr := bot.Edit(c.Message, "请选择操作系统", blueprintKey)
		if editErr != nil {
			log.Error("Edit message error", editErr)
		}
	})
	mediumKey := typeKey.Data("2c4g", medium)
	bot.Handle(&mediumKey, func(c *telebot.Callback) {
		p.Data[c.Sender.ID].Data["type"] = medium
		_, editErr := bot.Edit(c.Message, "请选择操作系统", blueprintKey)
		if editErr != nil {
			log.Error("Edit message error", editErr)
		}
	})
	largeKey := typeKey.Data("2c8g", large)
	bot.Handle(&largeKey, func(c *telebot.Callback) {
		p.Data[c.Sender.ID].Data["type"] = large
		_, editErr := bot.Edit(c.Message, "请选择操作系统", blueprintKey)
		if editErr != nil {
			log.Error("Edit message error", editErr)
		}
	})
	xLargeKey := typeKey.Data("4c16g", xlarge)
	bot.Handle(&xLargeKey, func(c *telebot.Callback) {
		p.Data[c.Sender.ID].Data["type"] = xlarge
		_, editErr := bot.Edit(c.Message, "请选择操作系统", blueprintKey)
		if editErr != nil {
			log.Error("Edit message error", editErr)
		}
	})
	dXLargeKey := typeKey.Data("8c32g", dxlarge)
	bot.Handle(&dXLargeKey, func(c *telebot.Callback) {
		p.Data[c.Sender.ID].Data["type"] = dxlarge
		_, editErr := bot.Edit(c.Message, "请选择操作系统", blueprintKey)
		if editErr != nil {
			log.Error("Edit message error", editErr)
		}
	})
	typeKey.Inline(typeKey.Row(nanoKey, microKey, smallKey), typeKey.Row(mediumKey, largeKey, xLargeKey), typeKey.Row(dXLargeKey))
	key := &telebot.ReplyMarkup{}
	createLs := key.Data("创建Lightsail", "createLs")
	bot.Handle(&createLs, func(c *telebot.Callback) {
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
		p.Data[c.Sender.ID].RegionChan = make(chan int)
		_, editErr := bot.Edit(c.Message, "请选择地区", zoneKey)
		if editErr != nil {
			log.Error("Send message error", editErr)
		}
		select {
		case <-p.Data[c.Sender.ID].RegionChan:
			_, editErr := bot.Edit(c.Message, "请输入Lightsail实例的备注: ")
			if editErr != nil {
				log.Error("Send message error: ", editErr)
			}
			p.Session.SessionAdd(c.Sender.ID, func(m *telebot.Message) {
				p.Data[m.Sender.ID].Data["name"] = m.Text
				p.Session[m.Sender.ID].Channel <- true
			})
			select {
			case tmp := <-p.Session[c.Sender.ID].Channel:
				if tmp != true {
					return
				}
				p.Session.SessionDel(c.Sender.ID)
				_, sendErr := bot.Send(c.Sender, "请选择类型: ", typeKey)
				if sendErr != nil {
					log.Error("Send message error", sendErr)
				}
			case <-time.After(30 * time.Second):
				p.Session.SessionDel(c.Sender.ID)
				_, editErr := bot.Edit(c.Message, "操作超时")
				if editErr != nil {
					log.Error("Edit message error: ", editErr)
				}
			case <-time.After(30 * time.Second):
				_, editErr := bot.Edit(c.Message, "操作超时")
				if editErr != nil {
					log.Error("Edit message error: ", editErr)
				}
			}
		}
	})
	delLs := key.Data("删除Lightsail", "deleteLs")
	bot.Handle(&delLs, func(c *telebot.Callback) {
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
		p.Data[c.Sender.ID].RegionChan = make(chan int)
		_, editErr := bot.Edit(c.Message, "请选择地区", zoneKey)
		if editErr != nil {
			log.Error("Edit message error: ", editErr)
		}
		select {
		case <-p.Data[c.Sender.ID].RegionChan:
			_, editErr := bot.Edit(c.Message, "请输入备注: ")
			if editErr != nil {
				log.Error("Edit message error: ", editErr)
			}
			p.Session.SessionAdd(c.Sender.ID, func(m *telebot.Message) {
				p.Data[c.Sender.ID].Data["name"] = m.Text
				p.Session[c.Sender.ID].Channel <- true
			})
			select {
			case <-p.Session[c.Sender.ID].Channel:
				awsO, newErr := aws.New(p.Data[c.Sender.ID].Data["zone"],
					p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Id,
					p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Secret,
					p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Proxy)
				if newErr != nil {
					log.Error("Init aws sdk error: ", newErr)
					_, err := bot.Send(c.Sender, "删除失败!")
					if err != nil {
						log.Error("Send message error: ", err)
					}
					return
				}
				delErr := awsO.DeleteLs(p.Data[c.Sender.ID].Data["name"])
				if delErr != nil {
					log.Error("Delete lightsail error: ", delErr)
					_, err := bot.Send(c.Sender, "删除失败!")
					if err != nil {
						log.Error("Send message error: ", err)
					}
					return
				}
				_, err := bot.Send(c.Sender, "删除成功!")
				if err != nil {
					log.Error("Send message error: ", err)
				}
			case <-time.After(30 * time.Second):
				_, sendErr := bot.Send(c.Sender, "操作超时")
				if sendErr != nil {
					log.Error("Edit message error: ", sendErr)
				}
			}
		case <-time.After(30 * time.Second):
			_, editErr := bot.Edit(c.Message, "操作超时")
			if editErr != nil {
				log.Error("Edit message error: ", editErr)
			}
		}
	})
	stopLs := key.Data("停止Lightsail", "stopLs")
	bot.Handle(&stopLs, func(c *telebot.Callback) {
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}, RegionChan: make(chan int)}
		_, editErr := bot.Edit(c.Message, "请选择地区", zoneKey)
		if editErr != nil {
			log.Error("Edit message error: ", editErr)
		}
		select {
		case <-p.Data[c.Sender.ID].RegionChan:
			_, editErr := bot.Edit(c.Message, "请输入备注: ")
			if editErr != nil {
				log.Error("Edit message error: ", editErr)
			}
			p.Session.SessionAdd(c.Sender.ID, func(m *telebot.Message) {
				p.Data[c.Sender.ID].Data["name"] = m.Text
				p.Session[c.Sender.ID].Channel <- true
			})
			select {
			case <-p.Session[c.Sender.ID].Channel:
				awsO, newErr := aws.New(p.Data[c.Sender.ID].Data["zone"],
					p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Id,
					p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Secret,
					p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Proxy)
				if newErr != nil {
					log.Error("Init aws sdk error: ", newErr)
					_, err := bot.Send(c.Sender, "停止失败!")
					if err != nil {
						log.Error("Send message error: ", err)
					}
					return
				}
				stopErr := awsO.StopLs(p.Data[c.Sender.ID].Data["name"])
				if stopErr != nil {
					log.Error("Stop lightsail error: ", stopErr)
					_, err := bot.Send(c.Sender, "停止失败!")
					if err != nil {
						log.Error("Send message error: ", err)
					}
					return
				}
				_, err := bot.Send(c.Sender, "停止成功!")
				if err != nil {
					log.Error("Send message error: ", err)
				}
			case <-time.After(30 * time.Second):
				_, sendErr := bot.Send(c.Sender, "操作超时")
				if sendErr != nil {
					log.Error("Send message error: ", sendErr)
				}
			}
		case <-time.After(30 * time.Second):
			_, editErr := bot.Edit(c.Message, "操作超时")
			if editErr != nil {
				log.Error("Edit message error: ", editErr)
			}
		}
	})
	startLs := key.Data("启动Lightsail", "startLs")
	bot.Handle(&startLs, func(c *telebot.Callback) {
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}, RegionChan: make(chan int)}
		_, editErr := bot.Edit(c.Message, "请选择地区", zoneKey)
		if editErr != nil {
			log.Error("Edit message error: ", editErr)
		}
		select {
		case <-p.Data[c.Sender.ID].RegionChan:
			_, editErr := bot.Edit(c.Message, "请输入备注: ")
			if editErr != nil {
				log.Error("Edit message error: ", editErr)
			}
			p.Session.SessionAdd(c.Sender.ID, func(m *telebot.Message) {
				p.Data[c.Sender.ID].Data["name"] = m.Text
				p.Session[c.Sender.ID].Channel <- true
			})
			select {
			case <-p.Session[c.Sender.ID].Channel:
				awsO, newErr := aws.New(p.Data[c.Sender.ID].Data["zone"],
					p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Id,
					p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Secret,
					p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Proxy)
				if newErr != nil {
					log.Error("Init aws sdk error: ", newErr)
					_, err := bot.Send(c.Sender, "启动失败!")
					if err != nil {
						log.Error("Send message error: ", err)
					}
					return
				}
				startErr := awsO.StartLs(p.Data[c.Sender.ID].Data["name"])
				if startErr != nil {
					log.Error("Start lightsail error: ", startErr)
					_, err := bot.Send(c.Sender, "启动失败!")
					if err != nil {
						log.Error("Send message error: ", err)
					}
					return
				}
				_, err := bot.Send(c.Sender, "启动成功!")
				if err != nil {
					log.Error("Send message error: ", err)
				}
			case <-time.After(30 * time.Second):
				_, sendErr := bot.Send(c.Sender, "操作超时")
				if sendErr != nil {
					log.Error("Send message error: ", sendErr)
				}
			}
		case <-time.After(30 * time.Second):
			_, editErr := bot.Edit(c.Message, "操作超时")
			if editErr != nil {
				log.Error("Edit message error: ", editErr)
			}
		}
	})
	changeIp := key.Data("更换IP", "changeIp")
	bot.Handle(&changeIp, func(c *telebot.Callback) {
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}, RegionChan: make(chan int)}
		_, editErr := bot.Edit(c.Message, "请选择地区", zoneKey)
		if editErr != nil {
			log.Error("Edit message error: ", editErr)
		}
		select {
		case <-p.Data[c.Sender.ID].RegionChan:
			_, editErr := bot.Edit(c.Message, "请输入备注: ")
			if editErr != nil {
				log.Error("Edit message error: ", editErr)
			}
			p.Session.SessionAdd(c.Sender.ID, func(m *telebot.Message) {
				p.Data[c.Sender.ID].Data["name"] = m.Text
				p.Session[c.Sender.ID].Channel <- true
			})
			select {
			case <-p.Session[c.Sender.ID].Channel:
				awsO, newErr := aws.New(p.Data[c.Sender.ID].Data["zone"],
					p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Id,
					p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Secret,
					p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Proxy)
				if newErr != nil {
					log.Error("Init aws sdk error: ", newErr)
					_, err := bot.Send(c.Sender, "更换失败!")
					if err != nil {
						log.Error("Send message error: ", err)
					}
					return
				}
				changeErr := awsO.ChangeLsIp(p.Data[c.Sender.ID].Data["name"])
				if changeErr != nil {
					log.Error("Change lightsail ip error: ", changeErr)
					_, err := bot.Send(c.Sender, "更换失败!")
					if err != nil {
						log.Error("Send message error: ", err)
					}
					return
				}
				ls, getErr := awsO.GetLsInfo(p.Data[c.Sender.ID].Data["name"])
				if getErr != nil {
					log.Error("Get lightsail info error: ", getErr)
				}
				_, err := bot.Send(c.Sender, "更换成功!新的IP为: \n"+*ls.Ip)
				if err != nil {
					log.Error("Send message error: ", err)
				}
			case <-time.After(30 * time.Second):
				_, sendErr := bot.Send(c.Sender, "操作超时")
				if sendErr != nil {
					log.Error("Send message error: ", sendErr)
				}
			}
		case <-time.After(30 * time.Second):
			_, editErr := bot.Edit(c.Message, "操作超时")
			if editErr != nil {
				log.Error("Edit message error: ", editErr)
			}
		}
	})
	key.Inline(key.Row(createLs, delLs), key.Row(stopLs, startLs), key.Row(changeIp))
	bot.Handle("/LightSailManger", func(m *telebot.Message) {
		if m.Private() {
			mess := p.CheckKey(m.Sender.ID)
			if mess != "" {
				_, err := bot.Send(m.Sender, mess)
				if err != nil {
					log.Println("Send message error: ", err)
				}
				return
			}
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
