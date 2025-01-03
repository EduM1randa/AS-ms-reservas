package controllers

import (
	"context"
	"fmt"
	"log"
	"time"

	m "ms-reservas/models"
	pb "ms-reservas/protos_pb/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	if reservation.ReservationTime == "" {
		return fmt.Errorf("reservationTime is required")
	}
	if reservation.GuestCount == 0 {
		return fmt.Errorf("guestCount is required")
	}
	if reservation.Status == "" {
		return fmt.Errorf("status is required")
	}

	validStatuses := map[string]bool{
		"confirmada": true,
		"cancelada":  true,
		"completada": true,
	}
	if !validStatuses[reservation.Status] {
		return fmt.Errorf("invalid status, expected one of: confirmada, cancelada, completada")
	}

	const dateFormat = "02-01-2006"
	_, err := time.Parse(dateFormat, reservation.ReservationDate)
	if err != nil {
		return fmt.Errorf("invalid date format, expected dd-mm-yyyy")
	}

	const timeFormat = "15:04"
	reservationTime, err := time.Parse(timeFormat, reservation.ReservationTime)
	if err != nil {
		return fmt.Errorf("invalid time format, expected HH:MM")
	}

	fmt.Println(reservationTime.Minute())

	if reservationTime.Minute() != 0 {
		return fmt.Errorf("reservation time must be end in 00")
	}

	collection := mongoClient.Database("reservations-db").Collection("reservations")
	_, err = collection.InsertOne(context.TODO(), reservation)
	if err != nil {
		log.Printf("failed to insert reservation: %v", err)
		return err
	}
	return nil
}

func CreateReservationHandler(req *pb.CreateReservationRequest) (*pb.Response, error) {
	exists, err := ReservationExists(req.TableId, req.ReservationDate, req.ReservationTime)
	if err != nil {
		return &pb.Response{Message: "Failed to check existing reservations", Success: false}, err
	}
	if exists {
		return &pb.Response{Message: "Reservation already exists for this table and date", Success: false}, nil
	}

	reservation := m.Reservation{
		UserId:          req.UserId,
		TableId:         req.TableId,
		ReservationDate: req.ReservationDate,
		ReservationTime: req.ReservationTime,
		GuestCount:      int(req.GuestCount),
		Status:          req.Status,
		CreateAt:        time.Now(),
	}
	err = CreateRes(reservation)
	if err != nil {
		return &pb.Response{Message: "Failed to create reservation", Success: false}, err
	}

	err = UpdateTableIsReserved(req.TableId, true)
	if err != nil {
		return &pb.Response{Message: "Failed to update table status", Success: false}, err
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
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("invalid id format: %v", err)
		return nil, err
	}

	collection := mongoClient.Database("reservations-db").Collection("reservations")
	var reservation m.Reservation
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&reservation)
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
	collection := mongoClient.Database("reservations-db").Collection("reservations")
	cursor, err := collection.Find(context.TODO(), bson.M{"userid": userID})
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
	collection := mongoClient.Database("reservations-db").Collection("reservations")
	cursor, err := collection.Find(context.TODO(), bson.M{"reservationdate": date})
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
		update["tableid"] = req.TableId
	}
	if req.ReservationDate != "" {
		update["reservationdate"] = req.ReservationDate
	}
	if req.GuestCount != 0 {
		update["guestcount"] = req.GuestCount
	}
	if req.Status != "" {
		update["status"] = req.Status
		validStatuses := map[string]bool{
			"confirmada": true,
			"cancelada":  true,
			"completada": true,
		}
		if !validStatuses[req.Status] {
			return nil, fmt.Errorf("invalid status, expected one of: confirmada, cancelada, completada")
		}
	}
	update["updateat"] = time.Now()

	err := UpdateReservation(id, update)
	if err != nil {
		return &pb.Response{Message: "Failed to update reservation", Success: false}, err
	}
	return &pb.Response{Message: "Reservation updated successfully", Success: true}, nil
}

func UpdateReservation(id string, update bson.M) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("invalid id format: %v", err)
		return nil
	}

	collection := mongoClient.Database("reservations-db").Collection("reservations")
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, bson.M{"$set": update})
	if err != nil {
		log.Printf("Failed to update reservation: %v", err)
		return err
	}

	if status, ok := update["status"].(string); ok && (status == "completada" || status == "cancelada") {
		var reservation m.Reservation
		err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&reservation)
		if err != nil {
			log.Printf("Failed to find updated reservation: %v", err)
			return err
		}

		err = UpdateTableIsReserved(reservation.TableId, false)
		if err != nil {
			log.Printf("Failed to update table status: %v", err)
			return err
		}
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
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("invalid id format: %v", err)
		return nil
	}

	collection := mongoClient.Database("reservations-db").Collection("reservations")
	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		log.Printf("Failed to delete reservation: %v", err)
		return err
	}
	return nil
}

func ReservationExists(tableId, reservationDate, reservationTime string) (bool, error) {
	collection := mongoClient.Database("reservations-db").Collection("reservations")

	// Convertir la hora de la reserva a un objeto time.Time
	const timeFormat = "15:04"
	_, err := time.Parse(timeFormat, reservationTime)
	if err != nil {
		return false, fmt.Errorf("invalid time format, expected HH:MM")
	}

	// Buscar reservas que se superpongan con la nueva reserva
	filter := bson.M{
		"tableid":         tableId,
		"reservationdate": reservationDate,
		"reservationtime": reservationTime,
	}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Printf("failed to check existing reservations: %v", err)
		return false, err
	}
	return count > 0, nil
}
