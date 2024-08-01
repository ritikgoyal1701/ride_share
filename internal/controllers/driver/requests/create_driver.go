package requests

type CreateDriverRequest struct {
	Name      string `json:"name" validate:"required"`
	License   string `json:"license" validate:"required"`
	ContactNo string `json:"contact_no" validate:"required"`
	Email     string `json:"email" validate:"required"`
	IsActive  bool   `json:"is_active"`
}
