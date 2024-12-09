package models

type Table struct {
	ID         string `json:"id,omitempty" bson:"_id,omitempty"`
	Number     int    `json:"number"`
	Capacity   int    `json:"capacity"`
	IsReserved bool   `json:"is_reserved"`
}

type Tables []Table
