package main

import (
	"lab4/cli"
	"lab4/context"
	"lab4/entity"
	"lab4/grpc/client"
	"lab4/grpc/server"
	"lab4/network"
	"lab4/routine"
	"time"
)

func main() {
	args := cli.Parse()
	ctx := context.NewContext(args.GSocket)
	nodeMaster := network.NewNodeMaster(ctx)

	if !args.IsRobot {
		ctx.DesktopClient = client.NewS2CClient(args.CSocket)
		go server.NewServer(args.SSocket, ctx, nodeMaster).Start()
	} else if args.CreateGame {
		ctx.InitNewGame(30, 30, 10, 1000, "robot-master", "robot-game")
	} else {
		go func() {
			time.Sleep(5 * time.Second)
			nodeMaster.SendJoin("192.168.0.79:8082", "mr_robot", "")
			for {
				time.Sleep(3 * time.Second)
				nodeMaster.SendSteer(0, "192.168.0.79:8082", entity.NewDirectionRandom())
			}
		}()
	}

	go nodeMaster.StartSendMulticastAnnouncementMsg()
	go nodeMaster.StartListenMulticastAnnouncement()
	go nodeMaster.StartListenUnicastMessages()
	ctx.Timer.Start()
	routine.StartGameRoutine(ctx, nodeMaster)
}
