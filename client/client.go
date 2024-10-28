package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"time"

	pb "DS-handin3/service/github.com/BirdyDK/DS-handin3"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewChittyChatClient(conn)

	var timestamp int32 = 0

	reader := bufio.NewReader(os.Stdin)
	log.Printf("Enter username:")
	name, _ := reader.ReadString('\n')
	name = name[:len(name)-2] // Trim the newline character

	join(client, name, timestamp)
	defer leave(client, name, timestamp)

	go listenBroadcasts(client, name, timestamp)

	for {
		msg, _ := reader.ReadString('\n')
		publish(client, name, msg[:len(msg)-1], timestamp) // Trim the newline character
	}
}

func join(client pb.ChittyChatClient, name string, timestamp int32) {
	timestamp++
	_, err := client.Join(context.Background(), &pb.JoinRequest{Name: name, Timestamp: timestamp})
	if err != nil {
		log.Fatalf("could not join: %v", err)
	}
	log.Printf("Joined chat as %s", name)
}

func leave(client pb.ChittyChatClient, name string, timestamp int32) {
	timestamp++
	_, err := client.Leave(context.Background(), &pb.LeaveRequest{Name: name, Timestamp: timestamp})
	if err != nil {
		log.Fatalf("could not leave: %v", err)
	}
	log.Printf("Left chat as %s", name)
}

func publish(client pb.ChittyChatClient, name, message string, timestamp int32) {
	timestamp++
	if len(message) > 128 {
		log.Println("Message length exceeds 128 characters")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.Publish(ctx, &pb.PublishRequest{Name: name, Message: message, Timestamp: timestamp})
	if err != nil {
		log.Fatalf("could not publish: %v", err)
	}
	log.Printf("Server Response:\r\n%s", r.GetMessage())
}

func listenBroadcasts(client pb.ChittyChatClient, name string, timestamp int32) {
	timestamp++
	stream, err := client.Subscribe(context.Background(), &pb.SubscribeRequest{Name: name, Timestamp: timestamp})
	if err != nil {
		log.Fatalf("could not subscribe: %v", err)
	}
	for {
		in, err := stream.Recv()
		if err != nil {
			log.Fatalf("failed to receive broadcast: %v", err)
		}

		if in.GetTimestamp() > timestamp {
			timestamp = in.GetTimestamp()
		}

		log.Printf("Broadcast:\r\n%s", in.GetMessage())
	}
}
