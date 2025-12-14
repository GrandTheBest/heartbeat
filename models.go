package heartbeat

import (
	"time"

	"github.com/google/uuid"
)

// Server is a common server, which using for register Server
type Server struct {
	PacketLifetime time.Duration
}

// Client is a common server, which using for register Client
type Client struct {
	ConnectedTo     string // Server host
	PoolingCooldown time.Duration
}

// Peer is a common server, which using for register Peer
type Peer struct {
	UUID          uuid.UUID
	Name          string
	ConnectedTo   string        // Server host
	PulseCooldown time.Duration // Interval of status pushing
}
