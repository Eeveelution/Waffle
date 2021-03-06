package irc_clients

import (
	"Waffle/bancho/irc/irc_messages"
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/database"
	"Waffle/helpers/serialization"
	"time"
)

func (client *IrcClient) GetRelevantUserStats() database.UserStats {
	return database.UserStats{
		UserID:         0,
		Mode:           0,
		Rank:           0,
		RankedScore:    0,
		TotalScore:     0,
		Level:          0,
		Accuracy:       0,
		Playcount:      0,
		CountSSH:       0,
		CountSS:        0,
		CountSH:        0,
		CountS:         0,
		CountA:         0,
		CountB:         0,
		CountC:         0,
		CountD:         0,
		Hit300:         0,
		Hit100:         0,
		Hit50:          0,
		HitMiss:        0,
		HitGeki:        0,
		HitKatu:        0,
		ReplaysWatched: 0,
	}
}

func (client *IrcClient) GetUserStatus() base_packet_structures.StatusUpdate {
	return base_packet_structures.StatusUpdate{
		Status:          serialization.OsuStatusIdle,
		StatusText:      "on IRC",
		BeatmapChecksum: "No!",
		CurrentMods:     0,
		Playmode:        0,
		BeatmapId:       -1,
	}
}

func (client *IrcClient) GetUserData() database.User {
	return client.UserData
}

func (client *IrcClient) GetClientTimezone() int32 {
	return 0
}

func (client *IrcClient) BanchoHandleOsuQuit(userId int32) {
	//TODO: do this
}

func (client *IrcClient) BanchoSpectatorJoined(userId int32) {
	//We don't do anything here cuz no spectator over IRC
}

func (client *IrcClient) BanchoSpectatorLeft(userId int32) {
	//We don't do anything here cuz no spectator over IRC
}

func (client *IrcClient) BanchoFellowSpectatorJoined(userId int32) {
	//We don't do anything here cuz no spectator over IRC
}

func (client *IrcClient) BanchoFellowSpectatorLeft(userId int32) {
	//We don't do anything here cuz no spectator over IRC
}

func (client *IrcClient) BanchoSpectatorCantSpectate(userId int32) {
	//We don't do anything here cuz no spectator over IRC
}

func (client *IrcClient) BanchoSpectateFrames(frameBundle base_packet_structures.SpectatorFrameBundle) {
	//We don't do anything here cuz no spectator over IRC
}

func (client *IrcClient) BanchoIrcMessage(message base_packet_structures.Message) {
	client.packetQueue <- irc_messages.IrcSendPrivMsg(message.Sender, message.Target, message.Message)
}

func (client *IrcClient) BanchoOsuUpdate(stats database.UserStats, update base_packet_structures.StatusUpdate) {

}

func (client *IrcClient) BanchoPresence(user database.User, stats database.UserStats, timezone int32) {

}

func (client *IrcClient) GetIdleTimes() (lastRecieve time.Time, logon time.Time) {
	return client.lastReceive, client.logonTime
}

func (client *IrcClient) GetFormattedJoinedChannels() string {
	channelString := ""

	for _, value := range client.joinedChannels {
		if value.ReadPrivileges == 0 {
			channelString += value.Name + " "
		}
	}

	return channelString
}
