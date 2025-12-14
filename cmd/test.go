package main

import (
	"context"
	"log"
	"time"

	"github.com/GrandTheBest/heartbeat"
	"github.com/GrandTheBest/heartbeat/types"
	"github.com/google/uuid"
)

func main() {
	server := heartbeat.NewServer(5 * time.Second)
	peer := heartbeat.NewPeer("bot", "http://127.0.0.1:7000", 10*time.Second)
	client := heartbeat.NewClient("http://127.0.0.1:7000", 2*time.Second)

	vault := types.Vault{
		V: map[uuid.UUID]*types.StatusPacket{},
	}

	go func() {
		err := server.Start(7000, &vault)
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
		log.Println("Server started successfully")
	}()

	time.Sleep(1 * time.Second)

	packets := make(chan *types.StatusResponse)

	go peer.StartPushing(context.Background())
	go client.StartPooling(context.Background(), packets)

	for {
		packet := <-packets
		if len(packet.Packets) == 0 {
			log.Println("Skipping!")
			continue
		}
		log.Printf("UUID: %s\nLast Check: %s\nStatus: %d\n\n", packet.Packets[0].LastPulse.PeerUUID, time.Unix(packet.Packets[0].LastPulse.PulseTime, 0).Format("2006.01.02 15:04:05"), packet.Packets[0].Status)
	}
}
