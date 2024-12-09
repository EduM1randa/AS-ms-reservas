package main

import (
	"context"
	"log"
	"net"

	"ms-reservas/controllers"
	"ms-reservas/database"
	pb "ms-reservas/protos_pb/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedReservationServiceServer
	pb.UnimplementedTableServiceServer
}

func (s *server) CreateReservation(ctx context.Context, req *pb.CreateReservationRequest) (*pb.Response, error) {
	return controllers.CreateReservationHandler(req)
}

func (s *server) GetReservationByID(ctx context.Context, req *pb.GetReservationByIDRequest) (*pb.Reservation, error) {
	return controllers.GetByIdHandler(req)
}

func (s *server) GetReservationsByUserID(ctx context.Context, req *pb.GetReservationsByUserIDRequest) (*pb.Reservations, error) {
	return controllers.GetReservationsByUserIDHandler(req)
}

func (s *server) GetReservationsByDate(ctx context.Context, req *pb.GetReservationsByDateRequest) (*pb.Reservations, error) {
	return controllers.GetReservationsByDateHandler(req)
}

func (s *server) UpdateReservation(ctx context.Context, req *pb.UpdateReservationRequest) (*pb.Response, error) {
	return controllers.UpdateReservationHandler(req)
}

func (s *server) DeleteReservation(ctx context.Context, req *pb.DeleteReservationRequest) (*pb.Response, error) {
	return controllers.DeleteReservationHandler(req)
}

func (s *server) CreateTable(ctx context.Context, req *pb.CreateTableRequest) (*pb.Response, error) {
	return controllers.CreateTableHandler(req)
}

func (s *server) GetTables(ctx context.Context, req *pb.Empty) (*pb.Tables, error) {
	return controllers.GetTablesHandler(req)
}

func (s *server) UpdateTable(ctx context.Context, req *pb.UpdateTableRequest) (*pb.Response, error) {
	return controllers.UpdateTableHandler(req)
}

func (s *server) GetAvailableTables(ctx context.Context, req *pb.GetAvailableTablesRequest) (*pb.Tables, error) {
	return controllers.GetAvailableTablesHandler(req)
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
	pb.RegisterTableServiceServer(s, &server{})

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
