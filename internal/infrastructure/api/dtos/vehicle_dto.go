package dtos

// CreateVehicleRequest represents the payload to create a new vehicle.
type CreateVehicleRequest struct {
	Plate         string `json:"plate" example:"ABC-1234"`
	Brand         string `json:"brand" example:"Toyota"`
	Model         string `json:"model" example:"Hiace"`
	Year          int32  `json:"year" example:"2024"`
	Capacity      int32  `json:"capacity" example:"15"`
	Color         string `json:"color,omitempty" example:"White"`
	VehicleType   string `json:"vehicle_type,omitempty" example:"van"`
	ChassisNum    string `json:"chassis_num,omitempty" example:"WDB9066351L123456"`
	InsuranceExp  string `json:"insurance_exp,omitempty" example:"2026-12-31T00:00:00Z"`
	TechReviewExp string `json:"tech_review_exp,omitempty" example:"2026-06-30T00:00:00Z"`
}

// VehicleResponse represents a vehicle returned by the API.
type VehicleResponse struct {
	ID            string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Plate         string `json:"plate" example:"ABC-1234"`
	Brand         string `json:"brand" example:"Toyota"`
	Model         string `json:"model" example:"Hiace"`
	Year          int32  `json:"year" example:"2024"`
	Capacity      int32  `json:"capacity" example:"15"`
	Status        string `json:"status" example:"active"`
	Color         string `json:"color,omitempty" example:"White"`
	VehicleType   string `json:"vehicle_type,omitempty" example:"van"`
	ChassisNum    string `json:"chassis_num,omitempty" example:"WDB9066351L123456"`
	InsuranceExp  string `json:"insurance_exp,omitempty" example:"2026-12-31T00:00:00Z"`
	TechReviewExp string `json:"tech_review_exp,omitempty" example:"2026-06-30T00:00:00Z"`
	CreatedAt     string `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt     string `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// ListVehiclesResponse represents a paginated list of vehicles.
type ListVehiclesResponse struct {
	Vehicles   []*VehicleResponse `json:"vehicles"`
	TotalCount int32              `json:"total_count" example:"42"`
}

// ErrorResponse represents an API error.
type ErrorResponse struct {
	Error string `json:"error" example:"invalid request body"`
}
