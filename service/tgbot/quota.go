package tgbot

import (
	"github.com/Yuzuki999/Aws-Manger-Bot/aws"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
	"time"
)

const (
	serviceCodeAcd = "ec2"
	quotaCodeAcd   = "L-1216C47A"
)

func (p *TgBot) QuotaManger(bot *tb.Bot) {
	quotaKey := &tb.ReplyMarkup{}
	acd := quotaKey.Data("查看标准EC2配额", "get_def")
	bot.Handle(&acd, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请选择地区", p.RegionKey)
		if err != nil {
			log.Println("Edit message error: ", err)
		}
		p.Data[c.Sender.ID].RegionChan = make(chan int)
		select {
		case <-p.Data[c.Sender.ID].RegionChan:
			defer delete(p.Data, c.Sender.ID)
			newRt, newErr := aws.New(p.Data[c.Sender.ID].Data["region"],
				p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Id,
				p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Secret,
				p.Config.UserInfo[c.Sender.ID].AwsSecret[p.Config.UserInfo[c.Sender.ID].NowKey].Proxy)
			if newErr != nil {
				log.Error(newErr)
			}
			quota, quotaErr := newRt.GetQuota(serviceCodeAcd, quotaCodeAcd)
			if quotaErr != nil {
				log.Error("Get quota error: ", quotaErr)
				_, editErr := bot.Edit(c.Message, "查看失败")
				if editErr != nil {
					log.Error("Edit message error: ", editErr)
				}
				return
			}
			_, editErr := bot.Edit(c.Message, "标准实例的配额为"+strconv.FormatFloat(*quota.Value, 'f', -1, 64))
			if editErr != nil {
				log.Error("Edit message error: ", editErr)
			}
		case <-time.After(30 * time.Second):
			_, editErr := bot.Edit(c.Message, "操作超时")
			if editErr != nil {
				log.Error("Edit message error: ", editErr)
			}
		}
	})
	other := quotaKey.Data("查看自定义配额", "get_other")
	bot.Handle(&other, func(c *tb.Callback) {
		_, err := bot.Edit(c.Message, "请选择地区", p.RegionKey)
		if err != nil {
			log.Println("Edit message error: ", err)
		}
		p.Data[c.Sender.ID].RegionChan = make(chan int)
		select {
		case <-p.Data[c.Sender.ID].RegionChan:
			_, err := bot.Edit(c.Message, "请输入ServiceCode和QuotaCode(用空格隔开): ")
			if err != nil {
				log.Error("Edit message error: ", err)
			}
			p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
				defer delete(p.Data, m.Sender.ID)
				code := strings.Split(m.Text, " ")
				newRt, newErr := aws.New(p.Data[m.Sender.ID].Data["region"],
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
	quotaKey.Inline(quotaKey.Row(acd), quotaKey.Row(other))
	key := &tb.ReplyMarkup{}
	getQuota := key.Data("查看配额", "get_quota")
	bot.Handle(&getQuota, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action:  Get Quota")
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
		_, editErr := bot.Edit(c.Message, "请选择配额", quotaKey)
		if editErr != nil {
			log.Error("Edit message error: ", editErr)
		}
	})
	updateQuota := key.Data("更新配额", "update_quota")
	bot.Handle(&updateQuota, func(c *tb.Callback) {
		log.Info("User: ", c.Sender.FirstName, " ",
			c.Sender.LastName, " ID: ", c.Sender.ID, " Action:  Update quota")
		p.Data[c.Sender.ID] = &Data{Data: map[string]string{}}
		_, err := bot.Edit(c.Message, "请选择地区", p.RegionKey)
		if err != nil {
			log.Println("Edit message error: ", err)
		}
		p.Data[c.Sender.ID].RegionChan = make(chan int)
		select {
		case <-p.Data[c.Sender.ID].RegionChan:
			_, editErr := bot.Edit(c.Message, "请输入ServiceCode和QuotaCode和要提升至的数量(用空格隔开): ")
			if editErr != nil {
				log.Error("Edit message error: ", editErr)
			}
			p.Session.SessionAdd(c.Sender.ID, func(m *tb.Message) {
				defer delete(p.Data, m.Sender.ID)
				code := strings.Split(m.Text, " ")
				newRt, newErr := aws.New(p.Data[m.Sender.ID].Data["region"],
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
	key.Inline(key.Row(getQuota, updateQuota))
	bot.Handle("/QuotaManger", func(m *tb.Message) {
		_, err := bot.Send(m.Sender, "请选择要进行的操作: ", key)
		if err != nil {
			log.Println("Send message error: ", err)
		}
	})
}
