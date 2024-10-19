package network

import (
	"lab4/entity"
	pb "lab4/generated/external_api"
	"lab4/mapper"
	"log"
	"net"
	"time"
)

func (nm *NodeMaster) consumeGameMessage(msg *pb.GameMessage, addr *net.UDPAddr) {
	log.Printf("Received %v Message\n", msg.GetType())
	switch msg.GetType().(type) {
	case *pb.GameMessage_Error:
		nm.consumeError(*msg.GetError().ErrorMessage, addr)
	case *pb.GameMessage_Announcement:
		nm.consumeAnnouncement(msg, addr)
	case *pb.GameMessage_Join:
		nm.consumeJoin(msg, addr)
	case *pb.GameMessage_Ack:
		nm.consumeAck(msg, addr)
	case *pb.GameMessage_Steer:
		nm.consumeSteer(msg.GetSteer(), addr)
	case *pb.GameMessage_State:
		nm.consumeState(msg, addr)
	case *pb.GameMessage_RoleChange:
		nm.consumeRoleChange(msg, addr)
	case *pb.GameMessage_Ping:
		nm.consumePing(msg, addr)
	case *pb.GameMessage_Discover:
		nm.consumeDiscover(msg, addr)
	}
	if nm.ctx.CurrentGame != nil && nm.ctx.CurrentGame.Players[addr.String()] != nil {
		nm.ctx.CurrentGame.Players[addr.String()].LastActiveAt = time.Now()
	}
}

func (nm *NodeMaster) consumeError(s string, addr *net.UDPAddr) {
	log.Printf("Error Message: %v from %v\n", s, addr.String())
}

func (nm *NodeMaster) consumeAnnouncement(msg *pb.GameMessage, addr *net.UDPAddr) {
	announcement := msg.GetAnnouncement().Games[0]
	game := mapper.ExternalMapAnnouncement(announcement, addr.String(), *msg.SenderId)
	nm.ctx.Games[addr.String()] = game
}

func (nm *NodeMaster) consumeJoin(msg *pb.GameMessage, addr *net.UDPAddr) {
	join := msg.GetJoin()
	if nm.ctx.Self.Role != entity.Master || nm.ctx.CurrentGame == nil {
		return
	}
	if nm.ctx.CurrentGame.Players[addr.String()] != nil {
		nm.SendError(-1, addr.String(), "Multiaccounting is forbidden")
		return
	}

	player := nm.ctx.InitNewPlayer(*join.PlayerName, addr.String())
	player.SetId(nm.ctx.CurrentGame.GetNextPlayerId())
	if nm.ctx.CurrentGame.Deputy == nil {
		player.Role = entity.Deputy
		nm.ctx.CurrentGame.Deputy = player
	} else {
		player.Role = entity.Default
	}

	nm.SendAck(*msg.MsgSeq, player.Id, player.Address)
}

func (nm *NodeMaster) consumeAck(ack *pb.GameMessage, addr *net.UDPAddr) {
	if nm.ctx.Self == nil || nm.ctx.Self.IsMaster() || nm.ctx.CurrentGame == nil {
		return
	} else if nm.ctx.CurrentGame.Master.Address == addr.String() {
		log.Printf("Set Self id: %d\n", *ack.ReceiverId)
		nm.ctx.Self.Id = *ack.ReceiverId
		nm.ctx.CurrentGame.Master.Id = *ack.SenderId
	}
}

func (nm *NodeMaster) consumeSteer(steer *pb.GameMessage_SteerMsg, addr *net.UDPAddr) {
	if !nm.ctx.Self.IsMaster() || nm.ctx.CurrentGame == nil || nm.ctx.CurrentGame.Players[addr.String()] == nil {
		return
	}
	nm.ctx.CurrentGame.Players[addr.String()].Snake.ChangeDirection(
		mapper.ExternalMapProtoDirectionToDirection(steer.GetDirection()),
	)
}

func (nm *NodeMaster) consumeState(msg *pb.GameMessage, addr *net.UDPAddr) {
	nm.ctx.CurrentGame = nm.ctx.Games[addr.String()]
	if nm.ctx.CurrentGame == nil {
		log.Println("Current game is nil!")
		return
	}
	state := msg.GetState().GetState()

	newPlayers := mapper.ExternalMapProtoGameStateToPlayers(state, msg.GetSenderId(), addr.String())
	newSnakeCache := mapper.MapPlayersToGameFieldCache(newPlayers)
	newFood := mapper.MapListToSet(mapper.MapList(state.Foods, mapper.ExternalMapProtoCoordToPoint))

	nm.ctx.CurrentGame.Players = newPlayers
	nm.ctx.CurrentGame.GameField.SnakeCache = newSnakeCache
	nm.ctx.CurrentGame.GameField.Food = newFood

	for _, player := range newPlayers {
		if player.Id == msg.GetReceiverId() {
			nm.ctx.Self = player
		}
		if player.IsMaster() {
			nm.ctx.CurrentGame.Master = player
		} else if player.IsDeputy() {
			nm.ctx.CurrentGame.Deputy = player
		}
	}

	response := newGameMessage(msg.GetMsgSeq(), nm.ctx.Self.Id, msg.GetSenderId())
	response.Type = &pb.GameMessage_Ack{}
	nm.sendGameMessage(addr.String(), response)
}

func (nm *NodeMaster) consumeRoleChange(msg *pb.GameMessage, addr *net.UDPAddr) {
	var sender, receiver *entity.Player
	for _, v := range nm.ctx.CurrentGame.Players {
		if v.Id == msg.GetSenderId() {
			sender = v
		} else if v.Id == msg.GetReceiverId() {
			receiver = v
		}
	}
	if sender == nil || receiver == nil || receiver != nm.ctx.Self {
		nm.SendError(msg.GetSenderId(), addr.String(), "Bad sender & receiver ID")
		return
	}
	if sender.IsMaster() || (sender.IsDeputy() && nm.ctx.CurrentGame.IsPlayerInactive(nm.ctx.CurrentGame.Master)) {
		rc := msg.GetRoleChange()
		sender.Role = mapper.ExternalMapProtoRole(*rc.SenderRole)
		receiver.Role = mapper.ExternalMapProtoRole(*rc.ReceiverRole)
	}
	response := newGameMessage(msg.GetMsgSeq(), nm.ctx.Self.Id, msg.GetSenderId())
	response.Type = &pb.GameMessage_Ack{}
	nm.sendGameMessage(addr.String(), response)
}

func (nm *NodeMaster) consumePing(msg *pb.GameMessage, addr *net.UDPAddr) {
	response := newGameMessage(msg.GetMsgSeq(), nm.ctx.Self.Id, msg.GetSenderId())
	response.Type = &pb.GameMessage_Ack{}
	nm.sendGameMessage(addr.String(), response)
}

func (nm *NodeMaster) consumeDiscover(msg *pb.GameMessage, addr *net.UDPAddr) {
	response := newGameMessage(msg.GetMsgSeq(), nm.ctx.Self.Id, msg.GetSenderId())
	response.Type = &pb.GameMessage_Ack{}
	nm.sendGameMessage(addr.String(), response)
}
