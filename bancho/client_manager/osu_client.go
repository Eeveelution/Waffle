package client_manager

import (
	"Waffle/bancho/packets"
	"Waffle/database"
)

// OsuClient defines an Interface of what we need from client.Client to be able to manage this client in the ClientLists
type OsuClient interface {
	GetUserId() int32
	GetPacketQueue() chan packets.BanchoPacket
	GetRelevantUserStats() database.UserStats
	GetUserStatus() packets.StatusUpdate
	GetUserData() database.User
	GetClientTimezone() int32
	InformSpectatorJoin(client OsuClient)
	InformSpectatorLeft(client OsuClient)
	InformSpectatorCantSpectate(client OsuClient)
	CleanupClient(reason string)
	Cut()
	GetAwayMessage() string
	HandleBeatmapInfoRequest(infoRequest packets.BeatmapInfoRequest)
}
