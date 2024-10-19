package context

import (
	"lab4/entity"
	"lab4/grpc/client"
	"lab4/timer"
	"time"
)

type Context struct {
	CurrentGame   *entity.Game
	Self          *entity.Player
	Games         map[string]*entity.Game
	Timer         *timer.Timer
	GameSocket    string
	DesktopClient *client.S2CClient
}

func NewContext(gsocket string) *Context {
	return &Context{
		GameSocket: gsocket,
		Games:      make(map[string]*entity.Game),
		Timer:      timer.NewTimer(time.Second),
	}
}

func (ctx *Context) InitNewGame(h, w, food, delay int32, nickname, gameName string) {
	players := make(map[string]*entity.Player)
	ctx.CurrentGame = entity.NewGame(
		h,
		w,
		food,
		delay,
		gameName,
		players,
		nil,
	)
	ctx.CurrentGame.Address = ctx.GameSocket
	ctx.Games[ctx.GameSocket] = ctx.CurrentGame
	self := ctx.InitNewPlayer(nickname, ctx.GameSocket)
	self.Role = entity.Master
	ctx.Self = self
	ctx.CurrentGame.Master = self
	ctx.Timer.Delay = ctx.CurrentGame.Delay
}

func (ctx *Context) InitNewPlayer(nickname, address string) *entity.Player {
	player := entity.NewPlayer(nickname)
	player.SetAddress(address)
	head := ctx.CurrentGame.GameField.FindEmptyRandom()
	dir := entity.NewDirectionRandom()
	tail := ctx.CurrentGame.GameField.NextPos(head, entity.GetOpposite(dir))
	player.InitSnake(head, tail, dir)
	ctx.CurrentGame.Players[address] = player
	return player
}
