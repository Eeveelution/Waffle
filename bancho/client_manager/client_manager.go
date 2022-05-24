package client_manager

import (
	"sync"
)

var clientList []OsuClient
var clientsById map[int32]OsuClient
var clientsByName map[string]OsuClient
var clientMutex sync.Mutex

// InitializeClientManager initializes the ClientManager
func InitializeClientManager() {
	clientsById = make(map[int32]OsuClient)
	clientsByName = make(map[string]OsuClient)
}

// LockClientList locks the client list, disallowing other threads from accessing until it's done
func LockClientList() {
	clientMutex.Lock()
}

// UnlockClientList unlocks the client list, allowing other threads to access it freely
func UnlockClientList() {
	clientMutex.Unlock()
}

// GetClientList returns a list of currently online and registered clients
func GetClientList() []OsuClient {
	return clientList
}

// GetClientById gets a client, assuming it's online, by their UserID
func GetClientById(id int32) OsuClient {
	value, exists := clientsById[id]

	if !exists {
		return nil
	}

	return value
}

// GetClientByName gets a client, assuming it's online, by their Username
func GetClientByName(username string) OsuClient {
	value, exists := clientsByName[username]

	if !exists {
		return nil
	}

	return value
}

// RegisterClient adds the Client to all the client lists it owns, it does NOT inform client's of its existence.
func RegisterClient(client OsuClient) {
	clientList = append(clientList, client)
	clientsById[client.GetUserId()] = client
	clientsByName[client.GetUserData().Username] = client
}

// UnregisterClient removes the Client from all the client lists it owns, it does NOT inform client's that it left
func UnregisterClient(client OsuClient) {
	LockClientList()

	for index, value := range clientList {
		if value == client {
			clientList = append(clientList[0:index], clientList[index+1:]...)
		}
	}

	delete(clientsById, client.GetUserId())
	delete(clientsByName, client.GetUserData().Username)

	UnlockClientList()
}

func GetClientCount() int {
	return len(clientList)
}
