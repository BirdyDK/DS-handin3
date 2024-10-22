package main

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"

	pb "DS-handin3/service/github.com/BirdyDK/DS-handin3"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedChittyChatServer
}

var timestamp int32 = 0

func UpdateTimestamp(newTime int32) int32 {
	if newTime > timestamp {
		timestamp = newTime
	}
	timestamp++
	return timestamp
}

func (s *server) Broadcast(ctx context.Context, in *pb.BroadcastRequest) (*pb.BroadcastResponse, error) {
	return &pb.BroadcastResponse{Message: in.Message + strconv.Itoa(int(in.Timestamp))}, nil
}

func (s *server) Join(ctx context.Context, in *pb.JoinRequest) (*pb.JoinResponse, error) {
	return &pb.JoinResponse{Message: "Participant " + in.Name + " joined Chitty-Chat"}, nil
}

func (s *server) Leave(ctx context.Context, in *pb.LeaveRequest) (*pb.LeaveResponse, error) {
	return &pb.LeaveResponse{Message: "Participant " + in.Name + " left Chitty-Chat"}, nil
}

func (s *server) Publish(ctx context.Context, in *pb.PublishRequest) (*pb.PublishResponse, error) {
	// Example logic to handle publishing a message
	message := "Published message from " + in.Name + ": " + in.Message
	return &pb.PublishResponse{Message: message}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChittyChatServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	for {

		time.Sleep(time.Second)
	}
}
