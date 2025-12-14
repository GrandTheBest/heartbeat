package heartbeat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/GrandTheBest/heartbeat/internal"
	"github.com/GrandTheBest/heartbeat/types"
)

// NewServer returns a new server instance
func NewServer(pl time.Duration) *Server {
	return &Server{
		PacketLifetime: pl,
	}
}

// NewClient returns a new client instance.
// «host» is a server host, to which client must be connected
// «cooldown» is a cooldown for «StartPooling» method
func NewClient(host string, cooldown time.Duration) *Client {
	return &Client{
		ConnectedTo:     host,
		PoolingCooldown: cooldown,
	}
}

// NewPeer returns a new peer instance.
// «name» is a name of peer. It should be unique
// «host» is a server host, to which peer must be connected.
// «cooldown» is a cooldown, by which the peer will send Pulse
func NewPeer(name string, host string, cooldown time.Duration) *Peer {
	return &Peer{
		Name:          name,
		ConnectedTo:   host,
		PulseCooldown: cooldown,
	}
}

// Start — starts a server
func (s *Server) Start(port int, v *types.Vault) error {
	http.HandleFunc("/pulse", func(w http.ResponseWriter, r *http.Request) {
		resp := types.ResponsePacket{
			Ok:      true,
			Message: "ok",
		}

		bodyRaw, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Failed to read body: %v", err)
			resp = types.ResponsePacket{
				Ok:      false,
				Message: "internal",
			}
		}
		defer r.Body.Close()

		var body types.PulsePacket
		err = json.Unmarshal(bodyRaw, &body)
		if err != nil {
			fmt.Printf("Failed to unmarshal body: %v", err)
			resp = types.ResponsePacket{
				Ok:      false,
				Message: "internal",
			}
		}

		s.Receive(v, body)
		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			fmt.Printf("Failed to encode: %v", err)
			return
		}
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		packetsRaw := s.Give(v)

		packets := make([]types.StatusPacket, 0, len(packetsRaw))
		for _, v := range packetsRaw {
			packets = append(packets, v)
		}

		err := json.NewEncoder(w).Encode(types.StatusResponse{
			Ok:      true,
			Packets: packets,
		})
		if err != nil {
			fmt.Printf("Failed to encode: %v", err)
		}
	})

	go func() {
		for {
			v.Lock()
			for i, p := range v.V {
				if time.Now().Unix()-p.LastPulse.PulseTime > int64(p.LastPulse.Cooldown.Seconds()) {
					l := p
					l.Status = FailureStatus
					v.V[i] = l
				}
			}
			v.Unlock()
			time.Sleep(s.PacketLifetime)
		}
	}()

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// Send is a method, which using for send PulsePacket from peer
func (p *Peer) Send() (bool, error) {
	resp, err := internal.CallP(
		fmt.Sprintf("%s/pulse", p.ConnectedTo),
		types.PulsePacket{
			PeerName:  p.Name,
			PulseTime: time.Now().Unix(),
		},
	)

	if err != nil {
		return false, err
	}

	return resp.Ok, nil
}

// StartPushing is a simple method, that provides simple pooling to send pulse packets
func (p *Peer) StartPushing(ctx context.Context) {
	ticker := time.NewTicker(p.PulseCooldown)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if _, err := p.Send(); err != nil {
				log.Printf("Failed to send pulse: %v", err)
			}
		}
	}
}

// Receive is a method, which allows server to receive PulsePacket
func (s *Server) Receive(v *types.Vault, p types.PulsePacket) {
	v.Lock()
	v.V[p.PeerName] = &types.StatusPacket{
		Status:    ActiveStatus,
		LastPulse: p,
	}
	v.Unlock()
}

// Give is a method, which gives all Status packets to client
func (s *Server) Give(v *types.Vault) []types.StatusPacket {
	v.Lock()
	packets := make([]types.StatusPacket, 0, len(v.V))
	for _, val := range v.V {
		packets = append(packets, *val)
	}
	v.Unlock()
	return packets
}

// Get is method, which using for get all status packets from server
func (c *Client) Get() (*types.StatusResponse, error) {
	resp, err := internal.CallG(
		fmt.Sprintf("%s/get", c.ConnectedTo),
	)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// StartPooling is a simple method, that provides simple pooling to get pulse packets
func (c *Client) StartPooling(ctx context.Context, ch chan *types.StatusResponse) {
	ticker := time.NewTicker(c.PoolingCooldown)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			resp, err := c.Get()
			if err != nil {
				log.Printf("Failed to get pulse packets: %v", err)
			}
			select {
			case ch <- resp:
			default:
			}
		}
	}
}
