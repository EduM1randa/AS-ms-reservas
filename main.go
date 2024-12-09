package main

import (
	"log"
	"net"

	"ms-reservas/controllers"
	"ms-reservas/database"
	pb "ms-reservas/protos_pb/proto"
	"ms-reservas/server"

	"google.golang.org/grpc"
)

func main() {
	client := database.ConnectMongoDB()
	controllers.SetMongoClient(client)

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterReservationServiceServer(s, &server.Server{})
	pb.RegisterTableServiceServer(s, &server.Server{})

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
