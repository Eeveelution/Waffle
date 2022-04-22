package packets

import (
	"bytes"
	"encoding/binary"
)

func BanchoSendSpectatorLeft(packetQueue chan BanchoPacket, userId int32) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, userId)

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          BanchoSpectatorLeft,
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue <- packet
}
