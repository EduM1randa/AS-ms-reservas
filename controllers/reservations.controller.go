package controllers

import (
	"context"
	"log"
	m "ms-reservas/models"
	pb "ms-reservas/protos_pb/proto"

	"go.mongodb.org/mongo-driver/mongo"
)

var mongoClient *mongo.Client

func SetMongoClient(client *mongo.Client) {
	mongoClient = client
}

func CreateRes(reservation m.Reservation) error {
	collection := mongoClient.Database("reservations").Collection("reservations")
	_, err := collection.InsertOne(context.TODO(), reservation)
	if err != nil {
		log.Printf("Failed to insert reservation: %v", err)
		return err
	}
	return nil
}

func CreateReservationHandler(req *pb.Reservation) (*pb.Response, error) {
	reservation := m.Reservation{
		ID:              req.Id,
		UserId:          req.UserId,
		TableId:         req.TableId,
		ReservationDate: req.ReservationDate,
		GuestCount:      int(req.GuestCount),
		Status:          req.Status,
	}
	err := CreateRes(reservation)
	if err != nil {
		return &pb.Response{Message: "Failed to create reservation", Success: false}, err
	}
	return &pb.Response{Message: "Reservation created successfully", Success: true}, nil
}
