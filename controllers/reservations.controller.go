package controllers

import (
	"context"
	"fmt"
	"log"
	"time"

	m "ms-reservas/models"
	pb "ms-reservas/protos_pb/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var mongoClient *mongo.Client

func SetMongoClient(client *mongo.Client) {
	mongoClient = client
}

// CREATE
func CreateRes(reservation m.Reservation) error {
	if reservation.UserId == "" {
		return fmt.Errorf("userID is required")
	}
	if reservation.TableId == "" {
		return fmt.Errorf("tableID is required")
	}
	if reservation.ReservationDate == "" {
		return fmt.Errorf("reservationDate is required")
	}
	if reservation.GuestCount == 0 {
		return fmt.Errorf("guestCount is required")
	}
	if reservation.Status == "" {
		return fmt.Errorf("status is required")
	}
	collection := mongoClient.Database("reservations").Collection("reservations")
	_, err := collection.InsertOne(context.TODO(), reservation)
	if err != nil {
		log.Printf("failed to insert reservation: %v", err)
		return err
	}
	return nil
}

func CreateReservationHandler(req *pb.CreateReservationRequest) (*pb.Response, error) {
	reservation := m.Reservation{
		UserId:          req.UserId,
		TableId:         req.TableId,
		ReservationDate: req.ReservationDate,
		GuestCount:      int(req.GuestCount),
		Status:          req.Status,
		CreateAt:        time.Now(),
	}
	err := CreateRes(reservation)
	if err != nil {
		return &pb.Response{Message: "Failed to create reservation", Success: false}, err
	}
	return &pb.Response{Message: "Reservation created successfully", Success: true}, nil
}

// GET BY ID
func GetByIdHandler(req *pb.GetReservationByIDRequest) (*pb.Reservation, error) {
	id := req.Id
	reservation, err := GetReservationByID(id)
	if err != nil {
		return nil, err
	}
	return &pb.Reservation{
		Id:              reservation.ID,
		UserId:          reservation.UserId,
		TableId:         reservation.TableId,
		ReservationDate: reservation.ReservationDate,
		GuestCount:      int32(reservation.GuestCount),
		Status:          reservation.Status,
	}, nil
}

func GetReservationByID(id string) (*m.Reservation, error) {
	collection := mongoClient.Database("reservations").Collection("reservations")
	var reservation m.Reservation
	err := collection.FindOne(context.TODO(), m.Reservation{ID: id}).Decode(&reservation)
	if err != nil {
		log.Printf("Failed to find reservation: %v", err)
		return nil, err
	}
	return &reservation, nil
}

// GET BY USER ID
func GetReservationsByUserIDHandler(req *pb.GetReservationsByUserIDRequest) (*pb.Reservations, error) {
	userID := req.UserId
	reservations, err := GetReservationsByUserID(userID)
	if err != nil {
		return nil, err
	}
	var pbReservations []*pb.Reservation
	for _, reservation := range reservations {
		pbReservations = append(pbReservations, &pb.Reservation{
			Id:              reservation.ID,
			UserId:          reservation.UserId,
			TableId:         reservation.TableId,
			ReservationDate: reservation.ReservationDate,
			GuestCount:      int32(reservation.GuestCount),
			Status:          reservation.Status,
		})
	}
	return &pb.Reservations{Reservations: pbReservations}, nil
}

func GetReservationsByUserID(userID string) ([]m.Reservation, error) {
	collection := mongoClient.Database("reservations").Collection("reservations")
	cursor, err := collection.Find(context.TODO(), bson.M{"user_id": userID})
	if err != nil {
		log.Printf("Failed to find reservations: %v", err)
		return nil, err
	}
	var reservations []m.Reservation
	if err = cursor.All(context.TODO(), &reservations); err != nil {
		log.Printf("Failed to decode reservations: %v", err)
		return nil, err
	}
	return reservations, nil
}

// GET BY DATE
func GetReservationsByDateHandler(req *pb.GetReservationsByDateRequest) (*pb.Reservations, error) {
	date := req.ReservationDate
	reservations, err := GetReservationsByDate(date)
	if err != nil {
		return nil, err
	}
	var pbReservations []*pb.Reservation
	for _, reservation := range reservations {
		pbReservations = append(pbReservations, &pb.Reservation{
			Id:              reservation.ID,
			UserId:          reservation.UserId,
			TableId:         reservation.TableId,
			ReservationDate: reservation.ReservationDate,
			GuestCount:      int32(reservation.GuestCount),
			Status:          reservation.Status,
			CreateAt:        reservation.CreateAt.Format(time.RFC3339),
			UpdateAt:        reservation.UpdateAt.Format(time.RFC3339),
		})
	}
	return &pb.Reservations{Reservations: pbReservations}, nil
}

func GetReservationsByDate(date string) ([]m.Reservation, error) {
	collection := mongoClient.Database("reservations").Collection("reservations")
	cursor, err := collection.Find(context.TODO(), bson.M{"reservation_date": date})
	if err != nil {
		log.Printf("failed to find reservations: %v", err)
		return nil, err
	}
	var reservations []m.Reservation
	if err = cursor.All(context.TODO(), &reservations); err != nil {
		log.Printf("failed to decode reservations: %v", err)
		return nil, err
	}
	return reservations, nil
}

// UPDATE
func UpdateReservationHandler(req *pb.UpdateReservationRequest) (*pb.Response, error) {
	id := req.Id
	update := bson.M{}
	if req.TableId != "" {
		update["table_id"] = req.TableId
	}
	if req.ReservationDate != "" {
		update["reservation_date"] = req.ReservationDate
	}
	if req.GuestCount != 0 {
		update["guest_count"] = req.GuestCount
	}
	if req.Status != "" {
		update["status"] = req.Status
	}
	update["update_at"] = time.Now()

	err := UpdateReservation(id, update)
	if err != nil {
		return &pb.Response{Message: "Failed to update reservation", Success: false}, err
	}
	return &pb.Response{Message: "Reservation updated successfully", Success: true}, nil
}

func UpdateReservation(id string, update bson.M) error {
	collection := mongoClient.Database("reservations").Collection("reservations")
	_, err := collection.UpdateOne(context.TODO(), bson.M{"id": id}, bson.M{"$set": update})
	if err != nil {
		log.Printf("Failed to update reservation: %v", err)
		return err
	}
	return nil
}

// DELETE
func DeleteReservationHandler(req *pb.DeleteReservationRequest) (*pb.Response, error) {
	id := req.Id
	err := DeleteReservation(id)
	if err != nil {
		return &pb.Response{Message: "Failed to delete reservation", Success: false}, err
	}
	return &pb.Response{Message: "Reservation deleted successfully", Success: true}, nil
}

func DeleteReservation(id string) error {
	collection := mongoClient.Database("reservations").Collection("reservations")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil {
		log.Printf("Failed to delete reservation: %v", err)
		return err
	}
	return nil
}
