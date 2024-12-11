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
)

// CREATE
func CreateTableHandler(req *pb.CreateTableRequest) (*pb.Response, error) {
	if (req.Number == 0) || (req.Capacity == 0) {
		return &pb.Response{Message: "number and capacity are required", Success: false}, nil
	}
	if (req.Number < 0) || (req.Capacity < 0) {
		return &pb.Response{Message: "number and capacity must be greater than 0", Success: false}, nil
	}
	table := m.Table{
		Number:     int(req.Number),
		Capacity:   int(req.Capacity),
		IsReserved: req.IsReserved,
	}

	err := CreateTable(table)
	if err != nil {
		return &pb.Response{Message: "failed to create table", Success: false}, err
	}
	return &pb.Response{Message: "table created successfully", Success: true}, nil
}

func CreateTable(table m.Table) error {
	collection := mongoClient.Database("reservations-db").Collection("tables")
	_, err := collection.InsertOne(context.TODO(), table)
	if err != nil {
		log.Printf("failed to insert table: %v", err)
		return err
	}
	return nil
}

// GET ALL
func GetTablesHandler(req *pb.Empty) (*pb.Tables, error) {
	tables, err := GetTables()
	if err != nil {
		return nil, err
	}
	var pbTables []*pb.Table
	for _, table := range tables {
		pbTables = append(pbTables, &pb.Table{
			Id:         table.ID,
			Number:     int32(table.Number),
			Capacity:   int32(table.Capacity),
			IsReserved: table.IsReserved,
		})
	}
	return &pb.Tables{Tables: pbTables}, nil
}

func GetTables() ([]m.Table, error) {
	collection := mongoClient.Database("reservations-db").Collection("tables")
	cursor, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		log.Printf("failed to find tables: %v", err)
		return nil, err
	}
	var tables []m.Table
	if err = cursor.All(context.TODO(), &tables); err != nil {
		log.Printf("failed to decode tables: %v", err)
		return nil, err
	}
	return tables, nil
}

// UPDATE
func UpdateTableHandler(req *pb.UpdateTableRequest) (*pb.Response, error) {
	id := req.Id

	update := bson.M{}
	if req.Capacity != 0 {
		update["capacity"] = req.Capacity
	}
	update["is_reserved"] = req.IsReserved
	update["update_at"] = time.Now()

	err := UpdateTable(id, update)
	if err != nil {
		return &pb.Response{Message: "failed to update table", Success: false}, err
	}
	return &pb.Response{Message: "table updated successfully", Success: true}, nil
}

func UpdateTable(id string, update bson.M) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("invalid id format: %v", err)
		return err
	}

	collection := mongoClient.Database("reservations-db").Collection("tables")
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, bson.M{"$set": update})
	if err != nil {
		log.Printf("failed to update table: %v", err)
		return err
	}
	return nil
}

// GET AVAILABLE TABLES
func GetAvailableTablesHandler(req *pb.GetAvailableTablesRequest) (*pb.Tables, error) {
	date := req.ReservationDate
	tables, err := GetAvailableTables(date)
	if err != nil {
		return nil, err
	}
	var pbTables []*pb.Table
	for _, table := range tables {
		pbTables = append(pbTables, &pb.Table{
			Id:         table.ID,
			Number:     int32(table.Number),
			Capacity:   int32(table.Capacity),
			IsReserved: table.IsReserved,
		})
	}
	return &pb.Tables{Tables: pbTables}, nil
}

func GetAvailableTables(date string) ([]m.Table, error) {
	collectionReservations := mongoClient.Database("reservations-db").Collection("reservations")
	cursorReservations, err := collectionReservations.Find(context.TODO(), bson.M{"reservationdate": date})
	if err != nil {
		log.Printf("failed to find reservations: %v", err)
		return nil, err
	}
	var reservations []m.Reservation
	if err = cursorReservations.All(context.TODO(), &reservations); err != nil {
		log.Printf("failed to decode reservations: %v", err)
		return nil, err
	}

	reservedTables := make(map[string]bool)
	for _, reservation := range reservations {
		reservedTables[reservation.TableId] = true
	}

	collectionTables := mongoClient.Database("reservations-db").Collection("tables")
	cursorTables, err := collectionTables.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Printf("failed to find tables: %v", err)
		return nil, err
	}
	var tables []m.Table
	if err = cursorTables.All(context.TODO(), &tables); err != nil {
		log.Printf("failed to decode tables: %v", err)
		return nil, err
	}

	var availableTables []m.Table
	for _, table := range tables {
		if !reservedTables[table.ID] {
			availableTables = append(availableTables, table)
		}
	}

	return availableTables, nil
}

func UpdateTableIsReserved(tableID string, isReserved bool) error {
	objectId, err := primitive.ObjectIDFromHex(tableID)
	if err != nil {
		log.Printf("failed to convert table id to object id: %v", err)
	}

	collection := mongoClient.Database("reservations-db").Collection("tables")
	update := bson.M{"isreserved": isReserved}
	ok, err := collection.UpdateOne(context.TODO(), bson.M{"_id": objectId}, bson.M{"$set": update})

	fmt.Println("err: ", err)
	fmt.Println(ok)

	if err != nil {
		log.Printf("failed to update table status: %v", err)
		return err
	}
	return nil
}
