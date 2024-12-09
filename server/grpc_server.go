package server

import (
	"context"
	"ms-reservas/controllers"
	pb "ms-reservas/protos_pb/proto"
)

type Server struct {
	pb.UnimplementedReservationServiceServer
	pb.UnimplementedTableServiceServer
}

func (s *Server) CreateReservation(ctx context.Context, req *pb.CreateReservationRequest) (*pb.Response, error) {
	return controllers.CreateReservationHandler(req)
}

func (s *Server) GetReservationByID(ctx context.Context, req *pb.GetReservationByIDRequest) (*pb.Reservation, error) {
	return controllers.GetByIdHandler(req)
}

func (s *Server) GetReservationsByUserID(ctx context.Context, req *pb.GetReservationsByUserIDRequest) (*pb.Reservations, error) {
	return controllers.GetReservationsByUserIDHandler(req)
}

func (s *Server) GetReservationsByDate(ctx context.Context, req *pb.GetReservationsByDateRequest) (*pb.Reservations, error) {
	return controllers.GetReservationsByDateHandler(req)
}

func (s *Server) UpdateReservation(ctx context.Context, req *pb.UpdateReservationRequest) (*pb.Response, error) {
	return controllers.UpdateReservationHandler(req)
}

func (s *Server) DeleteReservation(ctx context.Context, req *pb.DeleteReservationRequest) (*pb.Response, error) {
	return controllers.DeleteReservationHandler(req)
}

// Implementación de los métodos del servicio de mesas
func (s *Server) CreateTable(ctx context.Context, req *pb.CreateTableRequest) (*pb.Response, error) {
	return controllers.CreateTableHandler(req)
}

func (s *Server) GetTables(ctx context.Context, req *pb.Empty) (*pb.Tables, error) {
	return controllers.GetTablesHandler(req)
}

func (s *Server) UpdateTable(ctx context.Context, req *pb.UpdateTableRequest) (*pb.Response, error) {
	return controllers.UpdateTableHandler(req)
}

func (s *Server) GetAvailableTables(ctx context.Context, req *pb.GetAvailableTablesRequest) (*pb.Tables, error) {
	return controllers.GetAvailableTablesHandler(req)
}
