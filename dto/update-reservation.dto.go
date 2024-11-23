package dto

type UpdateReservationDTO struct {
	TableId         string `json:"table_id,omitempty"`
	ReservationDate string `json:"reservation_date,omitempty"`
	GuestCount      int    `json:"guest_count,omitempty"`
	Status          string `json:"status,omitempty"`
}
