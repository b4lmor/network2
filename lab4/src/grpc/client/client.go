package client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"lab4/entity"
	pb "lab4/generated/internal_api"
	"lab4/mapper"
	"log"
	"time"
)

type S2CClient struct {
	client  pb.ServerToClientClient
	timeout time.Duration
}

func NewS2CClient(socket string) *S2CClient {
	conn, err := grpc.NewClient(
		socket,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	return &S2CClient{client: pb.NewServerToClientClient(conn), timeout: 2 * time.Second}
}

func (cl *S2CClient) FillSquares(changed map[entity.Point]entity.Cell) {
	ctx, cancel := context.WithTimeout(context.Background(), cl.timeout)
	defer cancel()
	request := mapper.InternalMapToFillSquaresRequest(changed)
	_, err := cl.client.FillSquares(ctx, request)
	if err != nil {
		log.Fatalf("can't fill squares: %v", err)
	}
}

func (cl *S2CClient) SendGameInfoFull(game *entity.Game) {
	ctx, cancel := context.WithTimeout(context.Background(), cl.timeout)
	defer cancel()
	request := &pb.SendGameInfoFullRequest{
		GameName:       game.Name,
		Address:        game.Master.Address,
		MasterNickname: game.Master.Nickname,
		Height:         game.GameField.Height,
		Width:          game.GameField.Width,
		MaxFood:        game.GetFoodNumber(),
		Players: mapper.MapMapToList(
			mapper.FilterMap(game.Players, func(a string, p *entity.Player) bool { return p.Snake != nil }),
			mapper.InternalMapPlayer,
		),
	}
	_, err := cl.client.SendGameInfoFull(ctx, request)
	if err != nil {
		log.Fatalf("can't send game info: %v", err)
	}
}

func (cl *S2CClient) EndGame(score int32) {
	ctx, cancel := context.WithTimeout(context.Background(), cl.timeout)
	defer cancel()
	request := &pb.EndGameRequest{
		Score: score,
	}
	_, err := cl.client.EndGame(ctx, request)
	if err != nil {
		log.Fatalf("can't end game: %v", err)
	}
}

func (cl *S2CClient) SendGamesInfo(games []*entity.Game) {
	ctx, cancel := context.WithTimeout(context.Background(), cl.timeout)
	defer cancel()
	gameInfos := mapper.MapList(games, mapper.InternalMapGame)
	request := &pb.SendGamesInfoRequest{GameInfo: gameInfos}
	_, err := cl.client.SendGamesInfo(ctx, request)
	if err != nil {
		log.Printf("[Warning] :: Can't send games info: %v", err)
	}
}
