package mapper

import (
	"fmt"
	"lab4/entity"
	pb "lab4/generated/external_api"
	"net"
	"time"
	"utils"
)

func ExternalMapPlayersToMap[U any](input []*entity.Player, mapper func(player *entity.Player) U) map[string]U {
	output := make(map[string]U)
	for _, v := range input {
		output[v.Address] = mapper(v)
	}
	return output
}

func ExternalMapGamePlayersToPlayersMap(input []*pb.GamePlayer, mapper func(player *pb.GamePlayer) *entity.Player) map[string]*entity.Player {
	output := make(map[string]*entity.Player)
	for _, v := range input {
		output[fmt.Sprintf("%s:%d", *v.IpAddress, *v.Port)] = mapper(v)
	}
	return output
}

func ExternalMapProtoRole(r pb.NodeRole) entity.Role {
	switch r {
	case pb.NodeRole_NORMAL:
		return entity.Default
	case pb.NodeRole_MASTER:
		return entity.Master
	case pb.NodeRole_VIEWER:
		return entity.Viewer
	case pb.NodeRole_DEPUTY:
		return entity.Default
	}
	panic("can't map role")
}

func ExternalMapRoleToProto(r entity.Role) pb.NodeRole {
	switch r {
	case entity.Default:
		return pb.NodeRole_NORMAL
	case entity.Master:
		return pb.NodeRole_MASTER
	case entity.Viewer:
		return pb.NodeRole_VIEWER
	case entity.Deputy:
		return pb.NodeRole_DEPUTY
	default:
		panic("unhandled default case")
	}
}

func ExternalMapPlayerToProto(p *entity.Player) *pb.GamePlayer {
	role := ExternalMapRoleToProto(p.Role)
	status := pb.PlayerType_HUMAN
	score := p.Snake.GetScore()
	addr, _ := net.ResolveUDPAddr("udp4", p.Address)
	ip := addr.IP.String()
	port := int32(addr.Port)
	return &pb.GamePlayer{
		Name:      &p.Nickname,
		Id:        &p.Id,
		Role:      &role,
		Type:      &status,
		Score:     &score,
		IpAddress: &ip,
		Port:      &port,
	}
}

func ExternalMapProtoPlayer(p *pb.GamePlayer) *entity.Player {
	return &entity.Player{
		Nickname: *p.Name,
		Role:     ExternalMapProtoRole(*p.Role),
		Id:       *p.Id,
		Address:  fmt.Sprintf("%s:%d", *p.IpAddress, *p.Port),
	}
}

func ExternalMapSnakeToProtoSnake(snake *entity.Snake, playerId int32) *pb.GameState_Snake {
	state := pb.GameState_Snake_ALIVE
	dir := ExternalMapDirectionToProtoDirection(snake.Direction)
	points := make([]*pb.GameState_Coord, 0, snake.Body.Size())
	for i := 0; i < snake.Body.Size(); i++ {
		point, _ := snake.Body.Get(i)
		points = append(points, ExternalMapPointToProtoCoord(point))
	}
	return &pb.GameState_Snake{
		PlayerId:      &playerId,
		State:         &state,
		HeadDirection: &dir,
		Points:        points,
	}
}

func ExternalMapProtoGameStateToPlayers(state *pb.GameState, senderId int32, raddr string) map[string]*entity.Player {
	res := make(map[string]*entity.Player)
	players := ExternalMapGamePlayersToPlayersMap(state.GetPlayers().Players, ExternalMapProtoPlayer)
	for _, player := range players {
		for _, snake := range state.Snakes {
			if *snake.PlayerId == player.Id {
				player.Snake = ExternalMapProtoSnakeToSnake(snake)
			}
		}
		if player.Id == senderId {
			player.Address = raddr
		}
		res[player.Address] = player
	}
	return res
}

func ExternalMapProtoSnakeToSnake(snake *pb.GameState_Snake) *entity.Snake {
	body := utils.NewDeque[entity.Point]()
	for _, coord := range snake.GetPoints() {
		point := ExternalMapProtoCoordToPoint(coord)
		body.PushBack(*point)
	}
	return &entity.Snake{
		Direction: ExternalMapProtoDirectionToDirection(*snake.HeadDirection),
		Body:      body,
	}
}

func ExternalMapProtoCoordToPoint(coord *pb.GameState_Coord) *entity.Point {
	return &entity.Point{
		X: coord.GetX(),
		Y: coord.GetY(),
	}
}

func ExternalMapGameToGameState(game *entity.Game) *pb.GameState {
	order := game.GetStateOrder()
	food := MapList(game.GameField.Food.GetList(), ExternalMapPointToProtoCoord)
	players := MapMapToList(game.Players, ExternalMapPlayerToProto)
	snakes := make([]*pb.GameState_Snake, 0, len(game.Players))
	for _, player := range game.Players {
		snakes = append(snakes, ExternalMapSnakeToProtoSnake(player.Snake, player.Id))
	}
	return &pb.GameState{
		StateOrder: &order,
		Foods:      food,
		Players:    &pb.GamePlayers{Players: players},
		Snakes:     snakes,
	}
}

func ExternalMapPointToProtoCoord(point entity.Point) *pb.GameState_Coord {
	return &pb.GameState_Coord{
		X: &point.X,
		Y: &point.Y,
	}
}

func ExternalMapAnnouncement(g *pb.GameAnnouncement, address string, senderId int32) *entity.Game {
	result := &entity.Game{
		Name:    *g.GameName,
		Players: ExternalMapGamePlayersToPlayersMap(g.Players.Players, ExternalMapProtoPlayer),
		GameField: entity.NewGameField(
			*g.GetConfig().Height,
			*g.GetConfig().Width,
			*g.GetConfig().FoodStatic,
		),
		Delay:   time.Duration(*g.GetConfig().StateDelayMs),
		Address: address,
	}
	for _, player := range result.Players {
		if player.Id == senderId {
			player.Address = address
		}
		if player.Role == entity.Master {
			result.Master = player
		}
		if player.Role == entity.Deputy {
			result.Deputy = player
		}
	}
	return result
}

func ExternalMapGameToAnnouncement(currentGame *entity.Game) *pb.GameAnnouncement {
	delay := int32(currentGame.Delay.Milliseconds())
	return &pb.GameAnnouncement{
		GameName: &currentGame.Name,
		Config: &pb.GameConfig{
			Width:        &currentGame.GameField.Width,
			Height:       &currentGame.GameField.Height,
			FoodStatic:   &currentGame.GameField.FoodStatic,
			StateDelayMs: &delay,
		},
		Players: &pb.GamePlayers{
			Players: MapMapToList(
				FilterMap(currentGame.Players, func(a string, p *entity.Player) bool { return p.Snake != nil }),
				ExternalMapPlayerToProto,
			),
		},
	}
}

func ExternalMapProtoDirectionToDirection(dir pb.Direction) entity.Direction {
	switch dir {
	case pb.Direction_UP:
		return entity.Up
	case pb.Direction_DOWN:
		return entity.Down
	case pb.Direction_LEFT:
		return entity.Left
	case pb.Direction_RIGHT:
		return entity.Right
	default:
		panic("can't map direction")
	}
}

func ExternalMapDirectionToProtoDirection(dir entity.Direction) pb.Direction {
	switch dir {
	case entity.Up:
		return pb.Direction_UP
	case entity.Down:
		return pb.Direction_DOWN
	case entity.Left:
		return pb.Direction_LEFT
	case entity.Right:
		return pb.Direction_RIGHT
	default:
		panic("can't map direction")
	}
}
