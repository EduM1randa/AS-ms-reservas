package controllers

import (
	"context"
	"log"
	m "ms-reservas/models"

	"go.mongodb.org/mongo-driver/mongo"
)

var mongoClient *mongo.Client

func SetMongoClient(client *mongo.Client) {
	mongoClient = client
}

func create_res(reservation m.Reservation) {
	collection := mongoClient.Database("reservations").Collection("reservations")
	_, err := collection.InsertOne(context.TODO(), reservation)

	if err != nil {
		log.Fatalf("Failed to insert reservation: %v", err)
	}
}
