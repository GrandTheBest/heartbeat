package types

import (
	"time"

	"github.com/google/uuid"
)

// Status is a peer status. See constants.go
type Status int

// PulsePacket is a basis structure, which describes a peer status.
// It used for pushing status on server from peer(aka Agent)
type PulsePacket struct {
	PeerUUID  uuid.UUID     `json:"peer_uuid"`
	PeerName  string        `json:"peer_name"`  // Agent name
	PulseTime int64         `json:"pulse_time"` // Unix time
	Cooldown  time.Duration `json:"cooldown"`   // Also that «PulseCooldown»
}

// StatusPacket is a basic structure, which using for tell client about status of some peer.
type StatusPacket struct {
	LastPulse PulsePacket `json:"last_pulse"` // Packet, which sent by peer
	Status    Status      `json:"status"`     // Status of peer
}

// ResponsePacket is a basic response packet.
// It's using for response on any request
type ResponsePacket struct {
	Ok      bool   `json:"ok"` // Status
	Message string `json:"message"`
}

// StatusResponse isn't a basic response packet :3
// It's using for response on «Get» request.
// «Get» is a method, NOT HTTP Method!
type StatusResponse struct {
	Ok      bool           `json:"ok"`      // Status
	Packets []StatusPacket `json:"packets"` // Status packets
}

// PacketVault is a basic type, which using for save pull packets
type PacketVault map[uuid.UUID]*StatusPacket
