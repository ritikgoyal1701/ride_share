package requests

type LocationUpdate struct {
	XCoordinate int64 `json:"x_coordinate"`
	YCoordinate int64 `json:"y_coordinate"`
}
