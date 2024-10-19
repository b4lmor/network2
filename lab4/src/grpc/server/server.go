package server

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	gameContext "lab4/context"
	"lab4/entity"
	pb "lab4/generated/internal_api"
	"lab4/mapper"
	"lab4/network"
	"log"
	"net"
)

type Server struct {
	Socket string
	ctx    *gameContext.Context
	nm     *network.NodeMaster
	pb.UnimplementedClientToServerServer
}

func NewServer(socket string, ctx *gameContext.Context, nm *network.NodeMaster) *Server {
	return &Server{Socket: socket, ctx: ctx, nm: nm}
}

func (s *Server) CreateGame(ctx context.Context, req *pb.CreateGameRequest) (*pb.CreateGameResponse, error) {
	log.Printf("CreateGame called with: %+v\n", req)
	s.ctx.InitNewGame(req.Width, req.Height, req.Food, int32(1000), req.Nickname, req.GameName)
	s.ctx.DesktopClient.SendGameInfoFull(s.ctx.CurrentGame)
	//s.ctx.DesktopClient.FillSquares(s.ctx.CurrentGame.GameField.GetAll())
	return &pb.CreateGameResponse{Success: true}, nil
}

func (s *Server) JoinGame(ctx context.Context, req *pb.JoinGameRequest) (*pb.JoinGameResponse, error) {
	log.Printf("JoinGame called with: %+v\n", req)
	s.nm.SendJoin(req.GetAddress(), req.GetNickname(), "")
	return &pb.JoinGameResponse{Success: true}, nil
}

func (s *Server) LeaveGame(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	log.Println("LeaveGame called")
	if s.ctx.CurrentGame == nil {
		return &emptypb.Empty{}, nil
	}
	go s.ctx.DesktopClient.EndGame(s.ctx.Self.Snake.GetScore())
	delete(s.ctx.CurrentGame.Players, s.ctx.Self.Address)

	if s.ctx.Self.Role == entity.Master || len(s.ctx.CurrentGame.Players) == 0 {
		delete(s.ctx.Games, s.ctx.CurrentGame.Address)
	}

	s.ctx.Self = nil
	s.ctx.CurrentGame = nil
	return &emptypb.Empty{}, nil
}

func (s *Server) SendDirection(ctx context.Context, req *pb.SendDirectionRequest) (*emptypb.Empty, error) {
	log.Printf("SendDirection called with: %+v\n", req)
	if s.ctx.Self.Role != entity.Viewer && s.ctx.Self.Role != entity.None {
		s.nm.SendSteer(
			s.ctx.CurrentGame.Master.Id,
			s.ctx.CurrentGame.Master.Address,
			mapper.InternalMapDirection(req.Direction),
		)
		s.ctx.Self.Snake.ChangeDirection(mapper.InternalMapDirection(req.Direction))
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) Start() {
	lis, err := net.Listen("tcp", s.Socket)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterClientToServerServer(grpcServer, s)
	log.Printf("Server is listening on socket %s ...", s.Socket)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
