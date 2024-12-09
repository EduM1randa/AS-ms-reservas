package controllers

import (
	"context"
	"log"
	"time"

	m "ms-reservas/models"
	pb "ms-reservas/protos_pb/proto"

	"go.mongodb.org/mongo-driver/bson"
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
	collection := mongoClient.Database("reservations-db").Collection("tables")
	_, err := collection.UpdateOne(context.TODO(), bson.M{"id": id}, bson.M{"$set": update})
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
	collection := mongoClient.Database("reservations-db").Collection("tables")
	cursor, err := collection.Find(context.TODO(), bson.M{"is_reserved": false, "reservation_date": date})
	if err != nil {
		log.Printf("failed to find available tables: %v", err)
		return nil, err
	}
	var tables []m.Table
	if err = cursor.All(context.TODO(), &tables); err != nil {
		log.Printf("failed to decode available tables: %v", err)
		return nil, err
	}
	return tables, nil
}
