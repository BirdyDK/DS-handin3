package main

import (
	"context"
	"log"
	"net"

	pb "DS-handin3/service/github.com/BirdyDK/DS-handin3"

	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Broadcast(ctx context.Context, in *pb.BroadcastRequest) (*pb.BroadcastResponse, error) {
	return &pb.BroadcastResponse{Message: "Hello " + in.Message + in.Timestamp}, nil
}

func (s *server) Publish(ctx context.Context, in *pb.PublishRequest) (*pb.PublishResponse, error) {
	return &pb.PublishResponse{Message: "Hello " + in.Name + in.Message}, nil
}

func (s *server) Join(ctx context.Context, in *pb.JoinRequest) (*pb.JoinResponse, error) {
	// Broadcast join message to all participants
	// Implement logical clock logic if needed
	return &pb.JoinResponse{Message: "Participant " + in.Name + " joined Chitty-Chat"}, nil
}

func (s *server) Leave(ctx context.Context, in *pb.LeaveRequest) (*pb.LeaveResponse, error) {
	// Broadcast leave message to all participants
	// Implement logical clock logic if needed
	return &pb.LeaveResponse{Message: "Participant " + in.Name + " left Chitty-Chat"}, nil
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
}
