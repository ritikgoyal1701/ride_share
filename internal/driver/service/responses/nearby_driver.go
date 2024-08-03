package responses

type PriceResponse struct {
	Drivers  []Driver `json:"drivers"`
	Price    float64  `json:"price"`
	Distance float64  `json:"distance"`
	IsSurge  bool     `json:"is_surge"`
}

type Driver struct {
	User     User     `json:"user"`
	Location Location `json:"location"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Location struct {
	XCoordinate float64 `json:"x_coordinate"`
	YCoordinate float64 `json:"y_coordinate"`
}
