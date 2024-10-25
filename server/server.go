package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"

	pb "DS-handin3/service/github.com/BirdyDK/DS-handin3"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedChittyChatServer
	participants map[string]chan string
	mu           sync.Mutex
	timestamp    int32
}

func NewServer() *server {
	return &server{
		participants: make(map[string]chan string),
	}
}

func (s *server) joinHandler(name string) string {
	s.mu.Lock()
	ch := make(chan string, 10)
	s.participants[name] = ch
	s.timestamp++
	s.mu.Unlock()
	message := "Participant " + name + " joined Chitty-Chat at Lamport time " + strconv.Itoa(int(s.timestamp))
	s.broadcast(message)
	return message
}

func (s *server) leaveHandler(name string) string {
	s.mu.Lock()
	delete(s.participants, name)
	s.timestamp++
	s.mu.Unlock()
	message := "Participant " + name + " left Chitty-Chat at Lamport time " + strconv.Itoa(int(s.timestamp))
	s.broadcast(message)
	return message
}

func (s *server) broadcast(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, ch := range s.participants {
		ch <- message
	}
	log.Printf("Broadcasted message:\r\n%s", message)
}

func (s *server) Join(ctx context.Context, in *pb.JoinRequest) (*pb.JoinResponse, error) {
	if in.GetTimestamp() > s.timestamp {
		s.timestamp = in.GetTimestamp()
	}
	message := s.joinHandler(in.Name)
	return &pb.JoinResponse{Message: message, Timestamp: s.timestamp}, nil
}

func (s *server) Leave(ctx context.Context, in *pb.LeaveRequest) (*pb.LeaveResponse, error) {
	if in.GetTimestamp() > s.timestamp {
		s.timestamp = in.GetTimestamp()
	}
	message := s.leaveHandler(in.Name)
	return &pb.LeaveResponse{Message: message, Timestamp: s.timestamp}, nil
}

func (s *server) Publish(ctx context.Context, in *pb.PublishRequest) (*pb.PublishResponse, error) {
	if in.GetTimestamp() > s.timestamp {
		s.timestamp = in.GetTimestamp()
	}
	if len(in.Message) > 128 {
		return nil, fmt.Errorf("message length exceeds 128 characters")
	}
	s.timestamp++
	message := "(Lamport " + strconv.Itoa(int(s.timestamp)) + ") " + in.Name + ": " + in.Message
	s.broadcast(message)
	return &pb.PublishResponse{Message: message, Timestamp: s.timestamp}, nil
}

func (s *server) Subscribe(in *pb.SubscribeRequest, stream pb.ChittyChat_SubscribeServer) error {
	if in.GetTimestamp() > s.timestamp {
		s.timestamp = in.GetTimestamp()
	}
	s.mu.Lock()
	ch := s.participants[in.Name]
	s.mu.Unlock()
	if ch == nil {
		return fmt.Errorf("participant not found")
	}
	for msg := range ch {
		if err := stream.Send(&pb.SubscribeResponse{Message: msg, Timestamp: s.timestamp}); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChittyChatServer(s, NewServer())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
