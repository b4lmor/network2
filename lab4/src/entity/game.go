package entity

import (
	"time"
)

type Game struct {
	Name            string
	Players         map[string]*Player
	Master          *Player
	Deputy          *Player
	GameField       *GameField
	Delay           time.Duration
	Address         string
	messageSeq      int64
	stateOrder      int32
	currentPlayerId int32
}

func NewGame(h, w, food, delay int32, name string, players map[string]*Player, master *Player) *Game {
	return &Game{
		Name:      name,
		Players:   players,
		Master:    master,
		GameField: NewGameField(h, w, food),
		Delay:     time.Duration(delay) * time.Millisecond,
	}
}

func (g *Game) GetFoodNumber() int32 {
	return g.GameField.FoodStatic + int32(len(g.Players))
}

func (g *Game) GetNextPlayerId() int32 {
	g.currentPlayerId++
	return g.currentPlayerId
}

func (g *Game) GetNextMsgSeq() int64 {
	g.messageSeq++
	return g.messageSeq
}

func (g *Game) GetStateOrder() int32 {
	g.stateOrder++
	return g.stateOrder
}

func (g *Game) ChooseDeputy(exclude int32) *Player {
	if len(g.Players) < 2 {
		return nil
	}
	for _, player := range g.Players {
		if player.Role == Default && player.Id != exclude {
			player.Role = Deputy
			g.Deputy = player
			return player
		}
	}
	return nil
}

func (g *Game) IsPlayerInactive(player *Player) bool {
	if player == nil {
		return false
	}
	return time.Now().Sub(player.LastActiveAt) > (g.Delay*time.Millisecond)*4/5
}
