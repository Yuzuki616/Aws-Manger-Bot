package tgbot

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func (p *TgBot) setRegionKey(bot *tb.Bot) {
	regionKey := &tb.ReplyMarkup{}
	p.RegionKey = regionKey
	ohio := regionKey.Data("俄亥俄", "us-east-2")
	bot.Handle(&ohio, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["region"] = "us-east-2"
		p.Data[c.Sender.ID].RegionChan <- 0
	})
	virginia := regionKey.Data("弗吉尼亚", "us-east-1")
	bot.Handle(&virginia, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["region"] = "us-east-1"
		p.Data[c.Sender.ID].RegionChan <- 0
	})
	california := regionKey.Data("加利福尼亚", "us-west-1")
	bot.Handle(&california, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["region"] = "us-west-1"
		p.Data[c.Sender.ID].RegionChan <- 0
	})
	oregon := regionKey.Data("俄勒冈", "us-west-2")
	bot.Handle(&oregon, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["region"] = "us-west-2"
		p.Data[c.Sender.ID].RegionChan <- 0
	})
	hongKong := regionKey.Data("香港", "ap-east-1")
	bot.Handle(&hongKong, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["region"] = "ap-east-1"
		p.Data[c.Sender.ID].RegionChan <- 0
	})
	mumbai := regionKey.Data("孟买", "ap-south-1")
	bot.Handle(&mumbai, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["region"] = "ap-south-1"
		p.Data[c.Sender.ID].RegionChan <- 0
	})
	tokyo := regionKey.Data("东京", "ap-northeast-1")
	bot.Handle(&tokyo, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["region"] = "ap-northeast-1"
		p.Data[c.Sender.ID].RegionChan <- 0
	})
	osaka := regionKey.Data("大阪", "ap-northeast-3")
	bot.Handle(&osaka, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["region"] = "ap-northeast-3"
		p.Data[c.Sender.ID].RegionChan <- 0
	})
	seoul := regionKey.Data("首尔", "ap-northeast-2")
	bot.Handle(&seoul, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["region"] = "ap-northeast-2"
		p.Data[c.Sender.ID].RegionChan <- 0
	})
	singapore := regionKey.Data("新加坡", "ap-southeast-1")
	bot.Handle(&singapore, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["region"] = "ap-southeast-1"
		p.Data[c.Sender.ID].RegionChan <- 0
	})
	sydney := regionKey.Data("雪梨", "ap-southeast-2")
	bot.Handle(&sydney, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["region"] = "ap-southeast-2"
		p.Data[c.Sender.ID].RegionChan <- 0
	})
	caCentral := regionKey.Data("加拿大西部", "ca-central-1")
	bot.Handle(&caCentral, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["region"] = "ca-central-1"
		p.Data[c.Sender.ID].RegionChan <- 0
	})
	frankfurt := regionKey.Data("法兰克福", "eu-central-1")
	bot.Handle(&frankfurt, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["region"] = "eu-central-1"
		p.Data[c.Sender.ID].RegionChan <- 0
	})
	ireland := regionKey.Data("冰岛", "eu-west-1")
	bot.Handle(&ireland, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["region"] = "eu-west-1"
		p.Data[c.Sender.ID].RegionChan <- 0
	})
	london := regionKey.Data("伦敦", "eu-west-2")
	bot.Handle(&london, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["region"] = "eu-west-2"
		p.Data[c.Sender.ID].RegionChan <- 0
	})
	milan := regionKey.Data("米兰", "eu-south-1")
	bot.Handle(&milan, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["region"] = "eu-south-1"
		p.Data[c.Sender.ID].RegionChan <- 0
	})
	paris := regionKey.Data("巴黎", "eu-west-3")
	bot.Handle(&paris, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["region"] = "eu-west-3"
		p.Data[c.Sender.ID].RegionChan <- 0
	})
	stockholm := regionKey.Data("斯德哥尔摩", "eu-north-1")
	bot.Handle(&stockholm, func(c *tb.Callback) {
		p.Data[c.Sender.ID].Data["region"] = "eu-north-1"
		p.Data[c.Sender.ID].RegionChan <- 0
	})
	regionKey.Inline(
		regionKey.Row(ohio, virginia, california),
		regionKey.Row(oregon, hongKong, mumbai),
		regionKey.Row(tokyo, osaka, seoul),
		regionKey.Row(singapore, sydney, caCentral),
		regionKey.Row(frankfurt, ireland, london),
		regionKey.Row(milan, paris, stockholm))
}
