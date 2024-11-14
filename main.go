package main

import (
	"log"
	"net"

	"ms-reservas/controllers"
	"ms-reservas/database"

	"google.golang.org/grpc"
)

func main() {
	client := database.ConnectMongoDB()

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	controllers.SetMongoClient(client)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC over port: %v", err)
	}
}
