package requests

type Location struct {
	XCoordinate float64 `json:"x_coordinate"`
	YCoordinate float64 `json:"y_coordinate"`
}

type PriceRequest struct {
	StartLocation Location `json:"start_location"`
	DropLocation  Location `json:"drop_location"`
}
