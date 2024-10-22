package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"time"

	pb "../service/github.com/BirdyDK/DS-handin3"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewChittyChatClient(conn)

	name := "participant_name"
	join(client, name)
	defer leave(client, name)

	go listenBroadcasts(client)

	reader := bufio.NewReader(os.Stdin)
	for {
		msg, _ := reader.ReadString('\n')
		publish(client, name, msg)
	}
}

func join(client pb.ChittyChatClient, name string) {
	_, err := client.Join(context.Background(), &pb.JoinRequest{Name: name})
	if err != nil {
		log.Fatalf("could not join: %v", err)
	}
}

func leave(client pb.ChittyChatClient, name string) {
	_, err := client.Leave(context.Background(), &pb.LeaveRequest{Name: name})
	if err != nil {
		log.Fatalf("could not leave: %v", err)
	}
}

func publish(client pb.ChittyChatClient, name, message string) {
	if len(message) > 128 {
		log.Println("Message length exceeds 128 characters")
		return
	}
	_, err := client.Publish(context.Background(), &pb.PublishRequest{Name: name, Message: message})
	if err != nil {
		log.Fatalf("could not publish: %v", err)
	}
}

func listenBroadcasts(client pb.ChittyChatClient) {
	for {
		// Your logic to listen for broadcast messages
		// This could involve a streaming RPC method
		// For example purposes, it will just loop indefinitely
		time.Sleep(time.Second) // Simulate waiting for messages
	}
}
