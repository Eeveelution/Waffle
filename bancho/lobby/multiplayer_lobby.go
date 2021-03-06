package lobby

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/osu/base_packet_structures"
	"sync"
)

type MultiplayerLobby struct {
	MultiChannel        *chat.Channel
	MatchInformation    base_packet_structures.MultiplayerMatch
	MatchHost           LobbyClient
	MultiClients        [8]LobbyClient
	PlayersLoaded       [8]bool
	PlayerSkipRequested [8]bool
	PlayerCompleted     [8]bool
	PlayerFailed        [8]bool
	LastScoreFrames     [8]base_packet_structures.ScoreFrame
	MatchInfoMutex      sync.Mutex
	InProgress          bool
}

// Join gets called when a client is attempting to join the lobby
func (multiLobby *MultiplayerLobby) Join(client LobbyClient, password string) bool {
	//TODO@(Furball): currently there's a bug where the lobby can only have 7 players instead of the max 8

	//if they input the wrong password, join failed
	if multiLobby.MatchInformation.GamePassword != password {
		return false
	}

	multiLobby.MatchInfoMutex.Lock()

	//Inform everyone of the client, just in case they don't know them yet
	for n := 0; n != 8; n++ {
		if multiLobby.MultiClients[n] != nil {
			multiLobby.MultiClients[n].BanchoOsuUpdate(client.GetRelevantUserStats(), client.GetUserStatus())
		}
	}

	//Search for an Empty spot
	for i := 0; i != 8; i++ {
		if multiLobby.MatchInformation.SlotStatus[i] == base_packet_structures.MultiplayerMatchSlotStatusOpen {
			//Set the slot to them as well as join #multiplayer
			multiLobby.SetSlot(int32(i), client)
			multiLobby.MultiChannel.Join(client)

			multiLobby.MatchInfoMutex.Unlock()

			//Update everyone
			multiLobby.UpdateMatch()

			//Join success
			return true
		}
	}

	multiLobby.MatchInfoMutex.Unlock()

	return false
}

// SetSlot is used to set a slot to a player
func (multiLobby *MultiplayerLobby) SetSlot(slot int32, client LobbyClient) {
	//Handle for if a player is passed here, it can also be null which just sets the slot to be empty
	if client != nil {
		//Set slot nformation
		multiLobby.MatchInformation.SlotUserId[slot] = client.GetUserId()
		multiLobby.MatchInformation.SlotStatus[slot] = base_packet_structures.MultiplayerMatchSlotStatusNotReady
		multiLobby.MultiClients[slot] = client

		//Set teams, if necessary
		if multiLobby.MatchInformation.MatchTeamType == base_packet_structures.MultiplayerMatchTypeTagTeamVs || multiLobby.MatchInformation.MatchTeamType == base_packet_structures.MultiplayerMatchTypeTeamVs {
			if slot%2 == 0 {
				multiLobby.MatchInformation.SlotTeam[slot] = base_packet_structures.MultiplayerSlotTeamRed
			} else {
				multiLobby.MatchInformation.SlotTeam[slot] = base_packet_structures.MultiplayerSlotTeamBlue
			}
		}
	} else {
		//Set the slot to empty
		multiLobby.MatchInformation.SlotUserId[slot] = -1

		//If it's not locked, make it open
		if multiLobby.MatchInformation.SlotStatus[slot] != base_packet_structures.MultiplayerMatchSlotStatusLocked {
			multiLobby.MatchInformation.SlotStatus[slot] = base_packet_structures.MultiplayerMatchSlotStatusOpen
		}

		//Set team to neutral and make there be no client in that spot
		multiLobby.MatchInformation.SlotTeam[slot] = base_packet_structures.MultiplayerSlotTeamNeutral
		multiLobby.MultiClients[slot] = nil
	}
}

// MoveSlot moves a player from one slot to the other
func (multiLobby *MultiplayerLobby) MoveSlot(oldSlot int, newSlot int) {
	if oldSlot == newSlot {
		return
	}

	currentStatus := multiLobby.MatchInformation.SlotStatus[oldSlot]

	multiLobby.SetSlot(int32(newSlot), multiLobby.MultiClients[oldSlot])
	multiLobby.SetSlot(int32(oldSlot), nil)

	multiLobby.MatchInformation.SlotStatus[newSlot] = currentStatus
}

// UpdateMatch tells everyone inside the match and the lobby about the new happenings of the match
func (multiLobby *MultiplayerLobby) UpdateMatch() {
	for i := 0; i != 8; i++ {
		if multiLobby.MultiClients[i] != nil {
			multiLobby.MultiClients[i].BanchoMatchUpdate(multiLobby.MatchInformation)
		}
	}

	//Distribute in multiLobby as well
	BroadcastToLobby(func(client LobbyClient) {
		client.BanchoMatchUpdate(multiLobby.MatchInformation)
	})
}

// Part handles a player leaving the match
func (multiLobby *MultiplayerLobby) Part(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	slot := multiLobby.GetSlotFromUserId(client.GetUserId())

	//If they somehow don't exist, ignore
	if slot == -1 {
		return
	}

	//Reset their slot
	multiLobby.SetSlot(int32(slot), nil)

	//If they were the host, handle that separately, as we need to pass on the host
	if multiLobby.MatchHost == client {
		multiLobby.HandleHostLeave(slot)
	}

	//Make them leave #multiplayer
	client.BanchoChannelRevoked("#multiplayer")

	multiLobby.MultiChannel.Leave(client)

	//If there's nobody inside, disband
	if multiLobby.GetUsedUpSlots() == 0 {
		multiLobby.Disband()
	}

	//Tell everyone about it
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

// Disband is called when everyone leaves the match
func (multiLobby *MultiplayerLobby) Disband() {
	RemoveMultiMatch(multiLobby.MatchInformation.MatchId)
}

// HandleHostLeave handles the host leaving, as we need to pass on the host
func (multiLobby *MultiplayerLobby) HandleHostLeave(slot int) {
	//If nobody's there anymore, disband
	if multiLobby.GetUsedUpSlots() == 0 {
		multiLobby.Disband()
	}

	//Search for a new host
	for i := 0; i != 8; i++ {
		if multiLobby.MultiClients[i] != nil {
			//If a client is found, set them to be the new host
			multiLobby.MatchHost = multiLobby.MultiClients[i]

			//We can move them freely if the match isn't in progress, as the slot IDs don't have to be preserved, unlike during gameplay
			if !multiLobby.InProgress {
				multiLobby.MoveSlot(i, slot)
			}

			//Tell the new client they're host now
			multiLobby.MatchHost.BanchoMatchTransferHost()

			multiLobby.MatchInformation.HostId = multiLobby.MatchHost.GetUserId()
		}
	}

	multiLobby.UpdateMatch()
}

// TryChangeSlot gets called when a player tries to change slot
func (multiLobby *MultiplayerLobby) TryChangeSlot(client LobbyClient, slotId int) {
	multiLobby.MatchInfoMutex.Lock()

	//Refuse if the slot is occupied or locked
	if multiLobby.MatchInformation.SlotStatus[slotId] == base_packet_structures.MultiplayerMatchSlotStatusLocked || (multiLobby.MatchInformation.SlotStatus[slotId]&base_packet_structures.MultiplayerMatchSlotStatusHasPlayer) > 0 {
		return
	}

	//Move them to that slot and tell everyone
	multiLobby.MoveSlot(multiLobby.GetSlotFromUserId(client.GetUserId()), slotId)
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

// ChangeTeam gets called when a player is trying to change their team
func (multiLobby *MultiplayerLobby) ChangeTeam(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	clientSlot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if clientSlot == -1 {
		return
	}

	//Flip colors
	if multiLobby.MatchInformation.SlotTeam[clientSlot] == base_packet_structures.MultiplayerSlotTeamRed {
		multiLobby.MatchInformation.SlotTeam[clientSlot] = base_packet_structures.MultiplayerSlotTeamBlue
	} else if multiLobby.MatchInformation.SlotTeam[clientSlot] == base_packet_structures.MultiplayerSlotTeamBlue {
		multiLobby.MatchInformation.SlotTeam[clientSlot] = base_packet_structures.MultiplayerSlotTeamRed
	} else {
		multiLobby.MatchInformation.SlotTeam[clientSlot] = base_packet_structures.MultiplayerSlotTeamRed
	}

	//Tell everyone
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

// TransferHost gets called when the host willingly gives up their host
func (multiLobby *MultiplayerLobby) TransferHost(client LobbyClient, slotId int) {
	multiLobby.MatchInfoMutex.Lock()

	if multiLobby.MatchHost != client {
		return
	}

	//set the new host
	multiLobby.MatchHost = multiLobby.MultiClients[slotId]
	multiLobby.MatchInformation.HostId = multiLobby.MatchHost.GetUserId()

	//Tell them about it
	multiLobby.MatchHost.BanchoMatchTransferHost()

	//Tell everyone
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

// ReadyUp gets called when a player has clicked the Ready button
func (multiLobby *MultiplayerLobby) ReadyUp(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	clientSlot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if clientSlot == -1 {
		return
	}

	//Set them to be ready and tell everyone they're ready
	multiLobby.MatchInformation.SlotStatus[clientSlot] = base_packet_structures.MultiplayerMatchSlotStatusReady
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

// Unready gets called when a player has changed their mind about being ready and pressed the not ready button
func (multiLobby *MultiplayerLobby) Unready(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	clientSlot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if clientSlot == -1 {
		return
	}

	//Set them to be not ready and tell everyone
	multiLobby.MatchInformation.SlotStatus[clientSlot] = base_packet_structures.MultiplayerMatchSlotStatusNotReady
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

// ChangeSettings gets called when the host of the lobby changes some settings
func (multiLobby *MultiplayerLobby) ChangeSettings(client LobbyClient, matchSettings base_packet_structures.MultiplayerMatch) {
	multiLobby.MatchInfoMutex.Lock()

	if multiLobby.MatchHost != client {
		return
	}

	//Update the settings and tell everyone
	multiLobby.MatchInformation = matchSettings
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

// ChangeMods gets called when the host of the lobby changes which mods are going to get played
func (multiLobby *MultiplayerLobby) ChangeMods(client LobbyClient, newMods int32) {
	multiLobby.MatchInfoMutex.Lock()

	if multiLobby.MatchHost != client {
		return
	}

	//Set new mods and tell everyone
	multiLobby.MatchInformation.ActiveMods = uint16(newMods)
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

// LockSlot gets called when the host attempts to lock/unlock a slot
func (multiLobby *MultiplayerLobby) LockSlot(client LobbyClient, slotId int) {
	multiLobby.MatchInfoMutex.Lock()

	if multiLobby.MatchHost != client {
		return
	}

	//don't allow the host to kick themselves by locking their slot
	if multiLobby.MultiClients[slotId] == multiLobby.MatchHost {
		return
	}

	//If we lock a slot with a player inside, we kick them
	if (multiLobby.MatchInformation.SlotStatus[slotId] & base_packet_structures.MultiplayerMatchSlotStatusHasPlayer) > 0 {
		droppedClient := multiLobby.MultiClients[slotId]

		multiLobby.MatchInfoMutex.Unlock()

		droppedClient.LeaveCurrentMatch()
		droppedClient.BanchoMatchUpdate(multiLobby.MatchInformation)

		multiLobby.MatchInfoMutex.Lock()
	}

	//If it's locked already, make it open
	if multiLobby.MatchInformation.SlotStatus[slotId] == base_packet_structures.MultiplayerMatchSlotStatusLocked {
		multiLobby.MatchInformation.SlotStatus[slotId] = base_packet_structures.MultiplayerMatchSlotStatusOpen

		multiLobby.UpdateMatch()
		multiLobby.MatchInfoMutex.Unlock()

		return
	}

	//Don't allow all slots to be locked
	if multiLobby.GetOpenSlotCount() > 2 && multiLobby.MatchInformation.SlotStatus[slotId] == base_packet_structures.MultiplayerMatchSlotStatusOpen {
		multiLobby.MatchInformation.SlotStatus[slotId] = base_packet_structures.MultiplayerMatchSlotStatusLocked
	}

	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

// InformNoBeatmap gets called when a player happens to be missing the map thats about to be played
func (multiLobby *MultiplayerLobby) InformNoBeatmap(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	slot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if slot == -1 {
		return
	}

	//Mark them as missing the map and tell everyone
	multiLobby.MatchInformation.SlotStatus[slot] = base_packet_structures.MultiplayerMatchSlotStatusMissingMap
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

// InformGotBeatmap gets called whenever a player has now gotten the beatmap that they were missing earlier
func (multiLobby *MultiplayerLobby) InformGotBeatmap(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	slot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if slot == -1 {
		return
	}

	//Set them to be not ready and tell everyone
	multiLobby.MatchInformation.SlotStatus[slot] = base_packet_structures.MultiplayerMatchSlotStatusNotReady
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

// InformLoadComplete gets called when a player has loaded into the game
func (multiLobby *MultiplayerLobby) InformLoadComplete(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	slot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if slot == -1 {
		return
	}

	//Set their slot to be fully loaded
	multiLobby.PlayersLoaded[slot] = true

	//Check if everyone has loaded in, if yes then tell everyone that everyone's ready and begin!
	if multiLobby.HaveAllPlayersLoaded() {
		for i := 0; i != 8; i++ {
			if multiLobby.MultiClients[i] != nil {
				multiLobby.MultiClients[i].BanchoMatchAllPlayersLoaded()
			}
		}
	}

	multiLobby.MatchInfoMutex.Unlock()
}

// InformScoreUpdate this happens every time a player hits a circle or gets a slidertick or whatever
func (multiLobby *MultiplayerLobby) InformScoreUpdate(client LobbyClient, scoreFrame base_packet_structures.ScoreFrame) {
	multiLobby.MatchInfoMutex.Lock()

	slot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if slot == -1 {
		return
	}

	//Set their slot id
	scoreFrame.Id = uint8(slot)
	//Currently unused, but could be useful to display statistics after the match had ended and stuff
	multiLobby.LastScoreFrames[slot] = scoreFrame

	//Tell everyone about their new score
	for i := 0; i != 8; i++ {
		if multiLobby.MultiClients[i] != nil {
			multiLobby.MultiClients[i].BanchoMatchScoreUpdate(scoreFrame)
		}
	}

	multiLobby.MatchInfoMutex.Unlock()
}

// InformCompletion gets called whenever a client has finished playing a map
func (multiLobby *MultiplayerLobby) InformCompletion(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	slot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if slot == -1 {
		return
	}

	//Set them to be completed
	multiLobby.PlayerCompleted[slot] = true

	//Check if everyone completed
	if multiLobby.HaveAllPlayersCompleted() {
		//Set the match to no longer be in progress
		multiLobby.InProgress = false

		for i := 0; i != 8; i++ {
			//Reset all states
			multiLobby.PlayerCompleted[i] = false
			multiLobby.PlayerSkipRequested[i] = false
			multiLobby.PlayersLoaded[i] = false
			multiLobby.PlayerFailed[i] = false

			if multiLobby.MultiClients[i] != nil {
				multiLobby.MatchInformation.SlotStatus[i] = base_packet_structures.MultiplayerMatchSlotStatusNotReady

				multiLobby.MultiClients[i].BanchoMatchComplete()
			}
		}
	}

	//Tell everyone
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

// InformPressedSkip gets called when a player pressed skip in multi
func (multiLobby *MultiplayerLobby) InformPressedSkip(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	slot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if slot == -1 {
		return
	}

	//Set their slot to be skipped
	multiLobby.PlayerSkipRequested[slot] = true

	//Tell everyone that they skipped
	for i := 0; i != 8; i++ {
		if multiLobby.MultiClients[i] != nil {
			multiLobby.MultiClients[i].BanchoMatchPlayerSkipped(int32(slot))
		}
	}

	//If everyone skipped, tell everyone that it's okay to skip
	if multiLobby.HaveAllPlayersSkipped() {
		for i := 0; i != 8; i++ {
			if multiLobby.MultiClients[i] != nil {
				multiLobby.MultiClients[i].BanchoMatchSkip()
			}
		}
	}

	multiLobby.MatchInfoMutex.Unlock()
}

// InformFailed gets called whenever a client fails
func (multiLobby *MultiplayerLobby) InformFailed(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	slot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if slot == -1 {
		return
	}

	//Set them as failed
	multiLobby.PlayerFailed[slot] = true

	//Tell everyone they failed
	for i := 0; i != 8; i++ {
		if multiLobby.MultiClients[i] != nil {
			multiLobby.MultiClients[i].BanchoMatchPlayerFailed(int32(slot))
		}
	}

	multiLobby.MatchInfoMutex.Unlock()
}

// StartGame gets called whenever the host starts the game
func (multiLobby *MultiplayerLobby) StartGame(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	if multiLobby.MatchHost != client {
		return
	}

	//Sets the game to be in progress
	multiLobby.InProgress = true

	//Tell everyone to start
	for i := 0; i != 8; i++ {
		if multiLobby.MultiClients[i] != nil {
			multiLobby.MatchInformation.SlotStatus[i] = base_packet_structures.MultiplayerMatchSlotStatusPlaying

			multiLobby.MultiClients[i].BanchoMatchStart(multiLobby.MatchInformation)
		}
	}

	//Tell everyone, in lobby aswell
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}
