package main

import (
	"context"
	"log"
	"net"

	"ms-reservas/controllers"
	"ms-reservas/database"
	pb "ms-reservas/protos_pb/proto" // Importa el paquete generado

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedReservationServiceServer
}

// CreateReservation implementa pb.ReservationServiceServer
func (s *server) CreateReservation(ctx context.Context, req *pb.Reservation) (*pb.Response, error) {
	return controllers.CreateReservationHandler(req)
}

// SayHello implementa pb.ReservationServiceServer
func (s *server) SayHello(ctx context.Context, req *pb.Message) (*pb.Message, error) {
	return &pb.Message{Body: "Hello, " + req.Body}, nil
}

func main() {
	client := database.ConnectMongoDB()
	controllers.SetMongoClient(client)

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterReservationServiceServer(s, &server{})

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
