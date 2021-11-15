package session

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

type Session map[int]*Data

type Data struct {
	handle  func(m *tb.Message)
	Channel chan bool
}

func (p *Session) SessionAdd(id int, handle func(m *tb.Message)) {
	if _, ok := (*p)[id]; ok {
		close((*p)[id].Channel)
	}
	(*p)[id] = &Data{
		Channel: make(chan bool),
		handle:  handle,
	}
}

func (p *Session) SessionCheck(id int) bool {
	_, ok := (*p)[id]
	return ok
}

func (p *Session) SessionHandle(id int, m *tb.Message) {
	(*p)[id].handle(m)
}

func (p *Session) SessionDel(id int) {
	delete(*p, id)
}
