package main

import "lab4/network"

func main() {
	master := &network.NodeMaster{}
	slave := &network.NodeSlave{}

	go slave.ListenMulticastAnnouncementMsg()
	master.SendMulticastAnnouncementMsg()
}
