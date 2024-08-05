package responses

type GetRide struct {
	ID             string   `json:"id"`
	Rider          User     `json:"rider"`
	Driver         User     `json:"driver"`
	DriverLocation Location `json:"driver_location"`
	StartLocation  Location `json:"start_location"`
	DropLocation   Location `json:"drop_location"`
	Price          float64  `json:"price"`
	Distance       float64  `json:"distance"`
}

type Location struct {
	XCoordinate float64 `json:"x_coordinate"`
	YCoordinate float64 `json:"y_coordinate"`
}
