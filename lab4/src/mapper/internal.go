package mapper

import (
	"lab4/entity"
	pb "lab4/generated/internal_api"
	"log"
)

func InternalMapToFillSquaresRequest(pointsMap map[entity.Point]entity.Cell) *pb.FillSquaresRequest {
	var request pb.FillSquaresRequest
	request.Points = make([]*pb.PointColor, 0)
	for p, color := range pointsMap {
		pointColor := &pb.PointColor{
			Point: &pb.Point{X: int32(p.X), Y: int32(p.Y)},
			Color: int32(color),
		}
		request.Points = append(request.Points, pointColor)
	}
	return &request
}

func InternalMapDirection(direction pb.Directions) entity.Direction {
	switch direction {
	case pb.Directions_UP:
		return entity.Up
	case pb.Directions_DOWN:
		return entity.Down
	case pb.Directions_LEFT:
		return entity.Left
	case pb.Directions_RIGHT:
		return entity.Right
	default:
		log.Fatal("can't map direction")
	}
	return entity.Up
}

func InternalMapRole(role entity.Role) pb.NodeRole {
	switch role {
	case entity.Master:
		return pb.NodeRole_MASTER
	case entity.Deputy:
		return pb.NodeRole_DEPUTY
	case entity.Default:
		return pb.NodeRole_NORMAL
	case entity.Viewer:
		return pb.NodeRole_VIEWER
	default:
		panic("can't map role to node role")
	}
}

func InternalMapPlayer(player *entity.Player) *pb.Player {
	return &pb.Player{
		Nickname: player.Nickname,
		Score:    player.Snake.GetScore(),
		Role:     InternalMapRole(player.Role),
	}
}

func InternalMapGame(g *entity.Game) *pb.GameInfoShort {
	return &pb.GameInfoShort{
		GameName:    g.Name,
		Address:     g.Address,
		PlayerCount: int32(len(g.Players)),
	}
}
