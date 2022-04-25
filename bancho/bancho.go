package bancho

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/clients"
	"Waffle/bancho/database"
	"Waffle/bancho/lobby"
	"fmt"
	"net"
)

type Bancho struct {
	Server net.Listener
}

func CreateBancho() *Bancho {
	bancho := new(Bancho)

	chat.InitializeChannels()                //Initializes Chat channels
	database.Initialize()                    //Initialized Database related endeavors
	client_manager.InitializeClientManager() //Initializes the client manager
	lobby.InitializeLobby()                  //Initializes the multi lobby
	clients.CreateWaffleBot()                //Creates WaffleBot

	//Creates the TCP server under which Waffle runs
	listener, err := net.Listen("tcp", "127.0.0.1:13381")

	if err != nil {
		fmt.Printf("Failed to Create TCP Server on 127.0.0.1:13381")
	}

	bancho.Server = listener

	return bancho
}

func (bancho *Bancho) RunBancho() {
	fmt.Printf("Running Bancho on 127.0.0.1:13381\n")

	for {
		//Accept connections
		conn, err := bancho.Server.Accept()
		fmt.Printf("Connection Accepted!\n")

		if err != nil {
			continue
		}

		//Handle new connection
		go clients.HandleNewClient(conn)
	}
}
