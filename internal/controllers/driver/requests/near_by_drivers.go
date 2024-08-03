package requests

type NearByDriversRequest struct {
	XCoordinate float64 `json:"x_coordinate"`
	YCoordinate float64 `json:"y_coordinate"`
}
