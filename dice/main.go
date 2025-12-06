package main

import (
	"context"
	"log"
	"math/rand/v2"
	"net"

	pb "github.com/unwelcome/devops4/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type DiceService struct {
	pb.UnimplementedDiceServiceServer
}

func (s *DiceService) Roll(ctx context.Context, req *pb.RollRequest) (*pb.RollResponse, error) {
	num := rand.IntN(6) + 1 // Генерирует 1..6
	log.Printf("OperationID: %s Generated: %d", req.GetOperationId(), num)
	return &pb.RollResponse{Number: int32(num)}, nil
}

func main() {
	creds, err := credentials.NewServerTLSFromFile("/certs/server.crt", "/certs/server.key")
	if err != nil {
		log.Fatalf("Failed to load TLS keys: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterDiceServiceServer(s, &DiceService{})

	log.Println("Dice Service started on :50051 with TLS")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
