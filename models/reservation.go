package models

import "time"

type Reservation struct {
	ID              string    `json:"id,omitempty" bson:"_id,omitempty"`
	UserId          string    `json:"user_id"`
	TableId         string    `json:"table_id"`
	ReservationDate string    `json:"reservation_date"`
	ReservationTime string    `json:"reservation_time"`
	GuestCount      int       `json:"guest_count"`
	Status          string    `json:"status"`
	CreateAt        time.Time `json:"create_at"`
	UpdateAt        time.Time `json:"update_at,omitempty"`
}

type Reservations []Reservation
