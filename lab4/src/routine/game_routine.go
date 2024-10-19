package routine

import (
	"lab4/context"
	"lab4/entity"
	"lab4/network"
	"log"
	"utils"
)

func StartGameRoutine(ctx *context.Context, nm *network.NodeMaster) {
	timerChan := ctx.Timer.Subscribe()
	for {
		<-timerChan
		go masterRoutine(ctx, nm)
		go checkDisconnections(ctx, nm)
		go sendGamesInfo(ctx)
		if ctx.DesktopClient != nil && ctx.CurrentGame != nil {
			go ctx.DesktopClient.FillSquares(ctx.CurrentGame.GameField.GetAllNotEmpty())
			go ctx.DesktopClient.SendGameInfoFull(ctx.CurrentGame)
		}
	}
}

func masterRoutine(ctx *context.Context, nm *network.NodeMaster) {
	if ctx.CurrentGame != nil && len(ctx.CurrentGame.Players) != 0 && ctx.Self.IsMaster() {
		if ctx.CurrentGame.GameField.IsMoving {
			log.Fatal("GameField is still moving. Can't start next move")
		}
		for _, player := range ctx.CurrentGame.Players {
			ctx.CurrentGame.GameField.MoveSnake(player.Snake)
		}
		for ctx.CurrentGame.GameField.Food.Size() < ctx.CurrentGame.GetFoodNumber() {
			ctx.CurrentGame.GameField.SpawnApple(ctx.CurrentGame.GameField.FindEmptyRandom())
		}
		for _, player := range ctx.CurrentGame.Players {
			if player.Id == ctx.Self.Id {
				continue
			}
			nm.SendGameState(player.Id, player.Address)
		}
	}
}

func checkDisconnections(ctx *context.Context, nm *network.NodeMaster) {
	if ctx.CurrentGame == nil {
		return
	}
	for _, player := range ctx.CurrentGame.Players {
		if ctx.CurrentGame.IsPlayerInactive(player) {
			log.Printf("Player '%s' (%v) is inactive! [%s]", player.Nickname, player.Role, player.Address)
			if ctx.Self.Role == entity.Default && player.IsMaster() {
				ctx.CurrentGame.Master = ctx.CurrentGame.Deputy
				ctx.CurrentGame.Deputy.Role = entity.Master
				ctx.CurrentGame.Deputy = nil
			} else if ctx.Self.IsDeputy() && player.IsMaster() {
				ctx.Self.Role = entity.Master
				deputy := ctx.CurrentGame.ChooseDeputy(player.Id)
				if deputy != nil {
					nm.SendRoleChange(deputy.Id, deputy.Address, entity.Deputy, entity.Master)
				}
				for _, p := range ctx.CurrentGame.Players {
					if deputy != nil && p.Id == deputy.Id {
						continue
					}
					nm.SendRoleChange(p.Id, player.Address, entity.Default, entity.Master)
				}
			} else if ctx.Self.IsMaster() && player.IsDeputy() {
				deputy := ctx.CurrentGame.ChooseDeputy(player.Id)
				if deputy != nil {
					nm.SendRoleChange(deputy.Id, deputy.Address, entity.Deputy, entity.Master)
				}
			}
		}
	}
}

func sendGamesInfo(ctx *context.Context) {
	if ctx.DesktopClient != nil && len(ctx.Games) != 0 {
		go ctx.DesktopClient.SendGamesInfo(utils.Values(ctx.Games))
	}
}
