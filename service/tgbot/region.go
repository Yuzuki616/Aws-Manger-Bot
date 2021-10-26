package tgbot

import (
	"github.com/338317/Aws-Manger-Bot/aws"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
)

func (p *TgBot) regionHandle(bot *tb.Bot, c *tb.Callback) {
	if p.State[c.Sender.ID].Parent == 101 {
		p.listEc2(bot, c)
		return
	}
	if p.State[c.Sender.ID].Parent == 102 {
		_, err := bot.Edit(c.Message, "请输入要删除的实例ID: ")
		if err != nil {
			log.Error("Send message error: ", err)
		}
		p.State[c.Sender.ID].Parent = 6
		return
	}
	if p.State[c.Sender.ID].Parent == 103 {
		_, err := bot.Edit(c.Message, "请输入要更换IP的实例ID: ")
		if err != nil {
			log.Error("Send message error: ", err)
		}
		p.State[c.Sender.ID].Parent = 7
		return
	}
	if p.State[c.Sender.ID].Parent == 104 {
		_, err := bot.Edit(c.Message, "请输入ServiceCode和QuotaCode(用空格隔开): ")
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.State[c.Sender.ID].Parent = 10
		return
	}
	if p.State[c.Sender.ID].Parent == 105 {
		_, err := bot.Edit(c.Message, "请输入ServiceCode和QuotaCode和要提升至的数量(用空格隔开): ")
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.State[c.Sender.ID].Parent = 11
		return
	}
	if p.State[c.Sender.ID].Parent == 106 {
		defer delete(p.State, c.Sender.ID)
		newRt, newErr := aws.New(p.State[c.Sender.ID].Data["region"],
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
		return
	}
	if p.State[c.Sender.ID].Parent == 107 {
		_, err := bot.Edit(c.Message, "请输入实例ID: ")
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.State[c.Sender.ID].Parent = 12
		return
	}
	if p.State[c.Sender.ID].Parent == 108 {
		_, err := bot.Edit(c.Message, "请输入实例ID: ")
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.State[c.Sender.ID].Parent = 14
		return
	}
	if p.State[c.Sender.ID].Parent == 109 {
		_, err := bot.Edit(c.Message, "请输入要关联的Ec2实例ID: ")
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.State[c.Sender.ID].Parent = 16
		return
	}
	if p.State[c.Sender.ID].Parent == 110 {
		_, err := bot.Edit(c.Message, "请输入要实例ID: ")
		if err != nil {
			log.Error("Edit message error: ", err)
		}
		p.State[c.Sender.ID].Parent = 17
		return
	}
	_, err := bot.Edit(c.Message, "请选择EC2类型", p.TypeKey)
	if err != nil {
		log.Error("Edit message error: ", err)
	}

}

func (p *TgBot) setRegionKey(bot *tb.Bot) {
	regionKey := &tb.ReplyMarkup{}
	p.RegionKey = regionKey
	ohio := regionKey.Data("俄亥俄", "us-east-2")
	bot.Handle(&ohio, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["region"] = "us-east-2"
		p.regionHandle(bot, c)
	})
	virginia := regionKey.Data("弗吉尼亚", "us-east-1")
	bot.Handle(&virginia, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["region"] = "us-east-1"
		p.regionHandle(bot, c)
	})
	california := regionKey.Data("加利福尼亚", "us-west-1")
	bot.Handle(&california, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["region"] = "us-west-1"
		p.regionHandle(bot, c)
	})
	oregon := regionKey.Data("俄勒冈", "us-west-2")
	bot.Handle(&oregon, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["region"] = "us-west-2"
		p.regionHandle(bot, c)
	})
	hongKong := regionKey.Data("香港", "ap-east-1")
	bot.Handle(&hongKong, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["region"] = "ap-east-1"
		p.regionHandle(bot, c)
	})
	mumbai := regionKey.Data("孟买", "ap-south-1")
	bot.Handle(&mumbai, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["region"] = "ap-south-1"
		p.regionHandle(bot, c)
	})
	tokyo := regionKey.Data("东京", "ap-northeast-1")
	bot.Handle(&tokyo, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["region"] = "ap-northeast-1"
		p.regionHandle(bot, c)
	})
	osaka := regionKey.Data("大阪", "ap-northeast-3")
	bot.Handle(&osaka, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["region"] = "ap-northeast-3"
		p.regionHandle(bot, c)
	})
	seoul := regionKey.Data("首尔", "ap-northeast-2")
	bot.Handle(&seoul, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["region"] = "ap-northeast-2"
		p.regionHandle(bot, c)
	})
	singapore := regionKey.Data("新加坡", "ap-southeast-1")
	bot.Handle(&singapore, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["region"] = "ap-southeast-1"
		p.regionHandle(bot, c)
	})
	sydney := regionKey.Data("雪梨", "ap-southeast-2")
	bot.Handle(&sydney, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["region"] = "ap-southeast-2"
		p.regionHandle(bot, c)
	})
	caCentral := regionKey.Data("加拿大西部", "ca-central-1")
	bot.Handle(&caCentral, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["region"] = "ca-central-1"
		p.regionHandle(bot, c)
	})
	frankfurt := regionKey.Data("法兰克福", "eu-central-1")
	bot.Handle(&frankfurt, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["region"] = "eu-central-1"
		p.regionHandle(bot, c)
	})
	ireland := regionKey.Data("冰岛", "eu-west-1")
	bot.Handle(&ireland, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["region"] = "eu-west-1"
		p.regionHandle(bot, c)
	})
	london := regionKey.Data("伦敦", "eu-west-2")
	bot.Handle(&london, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["region"] = "eu-west-2"
		p.regionHandle(bot, c)
	})
	milan := regionKey.Data("米兰", "eu-south-1")
	bot.Handle(&milan, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["region"] = "eu-south-1"
		p.regionHandle(bot, c)
	})
	paris := regionKey.Data("巴黎", "eu-west-3")
	bot.Handle(&paris, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["region"] = "eu-west-3"
		p.regionHandle(bot, c)
	})
	stockholm := regionKey.Data("斯德哥尔摩", "eu-north-1")
	bot.Handle(&stockholm, func(c *tb.Callback) {
		p.State[c.Sender.ID].Data["region"] = "eu-north-1"
		p.regionHandle(bot, c)
	})
	regionKey.Inline(
		regionKey.Row(ohio, virginia, california),
		regionKey.Row(oregon, hongKong, mumbai),
		regionKey.Row(tokyo, osaka, seoul),
		regionKey.Row(singapore, sydney, caCentral),
		regionKey.Row(frankfurt, ireland, london),
		regionKey.Row(milan, paris, stockholm))
}
