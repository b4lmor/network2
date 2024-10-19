package network

import (
	"lab4/entity"
	pb "lab4/generated/external_api"
	"lab4/mapper"
	"log"
)

func (nm *NodeMaster) SendError(receiverId int32, receiverAddr string, description string) {
	log.Println("Sending error to ", receiverAddr)
	msg := newGameMessage(nm.ctx.CurrentGame.GetNextMsgSeq(), nm.ctx.Self.Id, receiverId)
	msg.Type = &pb.GameMessage_Error{Error: &pb.GameMessage_ErrorMsg{ErrorMessage: &description}}
	nm.sendGameMessage(receiverAddr, msg)
}

func (nm *NodeMaster) SendAck(seq int64, receiverId int32, receiverAddr string) {
	log.Println("Sending Ack to ", receiverAddr)
	msg := newGameMessage(seq, nm.ctx.Self.Id, receiverId)
	msg.Type = &pb.GameMessage_Ack{}
	nm.sendGameMessage(receiverAddr, msg)
}

func (nm *NodeMaster) SendJoin(receiverAddr string, nickname, gameName string) {
	log.Println("Sending Join Message to ", receiverAddr)
	if nm.ctx.Games[receiverAddr] == nil {
		log.Println("No such game found in game list")
		return
	}
	nm.ctx.Self = entity.NewPlayer(nickname)
	nm.ctx.Self.Role = entity.Default
	nm.ctx.CurrentGame = nm.ctx.Games[receiverAddr]
	msg := newGameMessage(int64(0), -1, -1)
	role := pb.NodeRole_NORMAL
	msg.Type = &pb.GameMessage_Join{Join: &pb.GameMessage_JoinMsg{PlayerName: &nickname, GameName: &gameName, RequestedRole: &role}}
	nm.sendGameMessage(receiverAddr, msg)
}

func (nm *NodeMaster) SendSteer(receiverId int32, receiverAddr string, dir entity.Direction) {
	log.Println("Sending steer message to ", receiverAddr)
	if nm.ctx.CurrentGame == nil {
		return
	}
	msg := newGameMessage(nm.ctx.CurrentGame.GetNextMsgSeq(), nm.ctx.Self.Id, receiverId)
	dirMapped := mapper.ExternalMapDirectionToProtoDirection(dir)
	msg.Type = &pb.GameMessage_Steer{Steer: &pb.GameMessage_SteerMsg{Direction: &dirMapped}}
	nm.sendGameMessage(receiverAddr, msg)
}

func (nm *NodeMaster) SendGameState(receiverId int32, receiverAddr string) {
	log.Println("Sending game state to ", receiverAddr)
	state := mapper.ExternalMapGameToGameState(nm.ctx.CurrentGame)
	msg := newGameMessage(nm.ctx.CurrentGame.GetNextMsgSeq(), nm.ctx.Self.Id, receiverId)
	msg.Type = &pb.GameMessage_State{State: &pb.GameMessage_StateMsg{State: state}}
	nm.sendGameMessage(receiverAddr, msg)
}

func (nm *NodeMaster) SendRoleChange(receiverId int32, receiverAddress string, senderRole, receiverRole entity.Role) {
	sr := mapper.ExternalMapRoleToProto(senderRole)
	rr := mapper.ExternalMapRoleToProto(receiverRole)
	msg := newGameMessage(nm.ctx.CurrentGame.GetNextMsgSeq(), nm.ctx.Self.Id, receiverId)
	msg.Type = &pb.GameMessage_RoleChange{RoleChange: &pb.GameMessage_RoleChangeMsg{SenderRole: &sr, ReceiverRole: &rr}}
	nm.sendGameMessage(receiverAddress, msg)
}

func (nm *NodeMaster) sendAnnouncement(receiverId int32, receiverAddress string) {
	if nm.ctx.Self == nil || !nm.ctx.Self.IsMaster() {
		return
	}
	msg := newGameMessage(nm.ctx.CurrentGame.GetNextMsgSeq(), nm.ctx.Self.Id, receiverId)
	msg.Type = &pb.GameMessage_Announcement{
		Announcement: &pb.GameMessage_AnnouncementMsg{
			Games: []*pb.GameAnnouncement{mapper.ExternalMapGameToAnnouncement(nm.ctx.CurrentGame)},
		},
	}
	nm.sendGameMessage(receiverAddress, msg)
}
