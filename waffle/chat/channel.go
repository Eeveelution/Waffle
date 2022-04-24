package chat

import (
	"sync"
)

type Channel struct {
	Name            string
	Description     string
	ReadPrivileges  int32
	WritePrivileges int32
	Autojoin        bool
	Clients         []ChatClient
	ClientMutex     sync.Mutex
}

func (channel *Channel) Join(client ChatClient) bool {
	if (channel.ReadPrivileges & client.GetUserPrivileges()) <= 0 {
		return false
	}

	channel.ClientMutex.Lock()

	for _, chatUser := range channel.Clients {
		//Check for duplicate client
		if chatUser.GetUserId() == client.GetUserId() {
			channel.ClientMutex.Unlock()
			return true
		}
	}

	channel.Clients = append(channel.Clients, client)
	channel.ClientMutex.Unlock()

	return true
}

func (channel *Channel) Leave(client ChatClient) {
	channel.ClientMutex.Lock()

	for index, value := range channel.Clients {
		if value == client {
			channel.Clients = append(channel.Clients[0:index], channel.Clients[index+1:]...)
		}
	}

	channel.ClientMutex.Unlock()
}

func (channel *Channel) SendMessage(sendingClient ChatClient, message string, target string) {
	if (channel.WritePrivileges & sendingClient.GetUserPrivileges()) <= 0 {
		sendingClient.SendChatMessage("WaffleBot", "You're not allowed to post in this channel! Your message has been discarded.", target)
		return
	}

	channel.ClientMutex.Lock()

	for _, client := range channel.Clients {
		if client != sendingClient {
			client.SendChatMessage(sendingClient.GetUsername(), message, target)
		}
	}

	channel.ClientMutex.Unlock()
}
