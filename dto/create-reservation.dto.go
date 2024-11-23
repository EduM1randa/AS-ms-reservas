package dto

type CreateReservationDTO struct {
	UserId          string `json:"user_id"`
	TableId         string `json:"table_id"`
	ReservationDate string `json:"reservation_date"`
	GuestCount      int    `json:"guest_count"`
	Status          string `json:"status"`
}
