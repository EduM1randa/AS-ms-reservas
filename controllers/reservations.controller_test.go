package controllers_test

import (
	"context"
	"testing"

	"ms-reservas/controllers"
	"ms-reservas/database"
	"ms-reservas/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateRes(t *testing.T) {
	// Conectar a MongoDB Atlas
	client := database.ConnectMongoDB()
	defer client.Disconnect(context.TODO())

	controllers.SetMongoClient(client)

	// Definir los casos de prueba
	tests := []struct {
		name          string
		reservation   models.Reservation
		expectedError string
	}{
		{
			name: "Missing UserID",
			reservation: models.Reservation{
				TableId:         "table1",
				ReservationDate: "01-01-2023",
				ReservationTime: "12:30",
				GuestCount:      2,
				Status:          "confirmed",
			},
			expectedError: "userID is required",
		},
		{
			name: "Invalid Date Format",
			reservation: models.Reservation{
				UserId:          "user1",
				TableId:         "table1",
				ReservationDate: "2023-01-01",
				ReservationTime: "12:30",
				GuestCount:      2,
				Status:          "confirmed",
			},
			expectedError: "invalid date format, expected dd-mm-yyyy",
		},
		// Agrega más casos de prueba según sea necesario
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := controllers.CreateRes(tt.reservation)
			if err == nil || err.Error() != tt.expectedError {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestGetReservationByID(t *testing.T) {
	// Conectar a MongoDB Atlas
	client := database.ConnectMongoDB()
	defer client.Disconnect(context.TODO())

	controllers.SetMongoClient(client)

	// Definir el ID de la reserva
	id := "60d5ec49f1d2c2a1d4e8b0c8"
	objectID, _ := primitive.ObjectIDFromHex(id)

	// Insertar un documento de prueba
	collection := client.Database("reservations-db").Collection("reservations")
	_, err := collection.InsertOne(context.TODO(), bson.M{
		"_id":              objectID,
		"userId":           "user1",
		"tableId":          "table1",
		"reservationDate":  "01-01-2023",
		"reservationTime":  "12:30",
		"guestCount":       2,
		"status":           "confirmed",
	})
	if err != nil {
		t.Fatalf("Failed to insert test document: %v", err)
	}

	// Llamar a la función GetReservationByID
	reservation, err := controllers.GetReservationByID(id)
	if err != nil {
		t.Fatalf("Failed to get reservation by ID: %v", err)
	}

	// Verificar los resultados
	if reservation.ID != id {
		t.Errorf("expected ID %v, got %v", id, reservation.ID)
	}
	if reservation.UserId != "user1" {
		t.Errorf("expected UserId user1, got %v", reservation.UserId)
	}
	if reservation.TableId != "table1" {
		t.Errorf("expected TableId table1, got %v", reservation.TableId)
	}
	if reservation.ReservationDate != "01-01-2023" {
		t.Errorf("expected ReservationDate 01-01-2023, got %v", reservation.ReservationDate)
	}
	if reservation.ReservationTime != "12:30" {
		t.Errorf("expected ReservationTime 12:30, got %v", reservation.ReservationTime)
	}
	if reservation.GuestCount != 2 {
		t.Errorf("expected GuestCount 2, got %v", reservation.GuestCount)
	}
	if reservation.Status != "confirmed" {
		t.Errorf("expected Status confirmed, got %v", reservation.Status)
	}

	// Limpiar el documento de prueba
	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		t.Fatalf("Failed to delete test document: %v", err)
	}
}
