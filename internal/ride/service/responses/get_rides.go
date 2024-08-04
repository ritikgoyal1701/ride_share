package responses

type GetRides struct {
	ID            string  `json:"id"`
	RiderDistance float64 `json:"rider_distance"`
	RideDistance  float64 `json:"ride_distance"`
	Price         float64 `json:"price"`
}

type AcceptRide struct {
	Rider User `json:"rider"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
