package chat

import (
	"sync"
)

type Channel struct {
	Name         string
	Description  string
	AdminChannel bool
	Clients      []*AdminPrivilegable
	ClientMutex  sync.Mutex
}

func (channel Channel) Join(client AdminPrivilegable) bool {
	if channel.AdminChannel && client.IsOfAdminPrivileges() == false {
		return false
	}

	channel.ClientMutex.Lock()
	channel.Clients = append(channel.Clients, &client)
	channel.ClientMutex.Unlock()

	return true
}
