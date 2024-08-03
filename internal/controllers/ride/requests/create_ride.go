package requests

type CreateRideRequest struct {
	StartLocation Location `json:"start_location"`
	DropLocation  Location `json:"drop_location"`
}
