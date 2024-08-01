package requests

type CreateRiderRequest struct {
	Name      string `json:"name" validate:"required"`
	ContactNo string `json:"contact_no" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
}
