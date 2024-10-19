package entity

import "time"

type Player struct {
	Snake        *Snake
	Role         Role
	Nickname     string
	Address      string
	Id           int32
	LastActiveAt time.Time
}

func NewPlayer(nickname string) *Player {
	return &Player{Role: None, Nickname: nickname, LastActiveAt: time.Now()}
}

func (p *Player) SetAddress(address string) {
	p.Address = address
}

func (p *Player) SetId(id int32) {
	p.Id = id
}

func (p *Player) InitSnake(head, tail Point, direction Direction) {
	p.Snake = NewSnake(head, tail, direction)
}

func (p *Player) IsMaster() bool {
	return p.Role == Master
}

func (p *Player) IsDeputy() bool {
	return p.Role == Deputy
}
