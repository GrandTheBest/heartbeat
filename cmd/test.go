package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/GrandTheBest/heartbeat/pkg"
	"github.com/GrandTheBest/heartbeat/types"
)

func main() {
	server := pkg.NewServer()
	peer := pkg.NewPeer("bot", "http://127.0.0.1:7000", 10*time.Second)
	client := pkg.NewClient("http://127.0.0.1:7000", 30*time.Second)

	vault := types.PacketVault{}

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
		fmt.Printf("UUID: %s\nLast Check: %s\n\n", packet.Packets[0].LastPulse.PeerUUID, time.Unix(packet.Packets[0].LastPulse.PulseTime, 0).Format("2006.01.02 15:04:05"))
	}
}
