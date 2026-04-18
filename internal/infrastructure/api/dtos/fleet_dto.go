package dtos

import "time"

// --- Vehicles ---
type VehicleResponse struct {
	ID             string     `json:"id"`
	Plate          string     `json:"plate"`
	Brand          string     `json:"brand"`
	Model          string     `json:"model"`
	Year           int32      `json:"year"`
	Capacity       int32      `json:"capacity"`
	Status         string     `json:"status"`
	Color          string     `json:"color,omitempty"`
	VehicleType    string     `json:"vehicle_type,omitempty"`
	ChassisNum     string     `json:"chassis_num,omitempty"`
	InsuranceExp   *time.Time `json:"insurance_exp,omitempty"`
	TechReviewExp  *time.Time `json:"tech_review_exp,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type CreateVehicleRequest struct {
	Plate         string `json:"plate" validate:"required"`
	Brand         string `json:"brand" validate:"required"`
	Model         string `json:"model" validate:"required"`
	Year          int32  `json:"year"`
	Capacity      int32  `json:"capacity"`
	Color         string `json:"color"`
	VehicleType   string `json:"vehicle_type"`
	ChassisNum    string `json:"chassis_num"`
	InsuranceExp  string `json:"insurance_exp"`   // RFC3339
	TechReviewExp string `json:"tech_review_exp"` // RFC3339
}

type ListVehiclesResponse struct {
	Vehicles []*VehicleResponse `json:"vehicles"`
	Total    int32              `json:"total"`
}

// --- Schools ---
type SchoolResponse struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Address   string       `json:"address"`
	Location  *LocationDTO `json:"location,omitempty"`
	Phone     string       `json:"phone,omitempty"`
	Email     string       `json:"email,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
}

type CreateSchoolRequest struct {
	Name     string       `json:"name" validate:"required"`
	Address  string       `json:"address" validate:"required"`
	Location *LocationDTO `json:"location"`
	Phone    string       `json:"phone"`
	Email    string       `json:"email"`
}

type UpdateSchoolRequest struct {
	Name     string       `json:"name"`
	Address  string       `json:"address"`
	Location *LocationDTO `json:"location"`
	Phone    string       `json:"phone"`
	Email    string       `json:"email"`
}

type ListSchoolsResponse struct {
	Schools []*SchoolResponse `json:"schools"`
	Total   int32             `json:"total"`
}

// --- Drivers ---
type DriverResponse struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id"`
	LicenseNumber string     `json:"license_number"`
	LicenseType   string     `json:"license_type"`
	LicenseExpiry *time.Time `json:"license_expiry,omitempty"`
	Status        string     `json:"status"`
}

type RegisterDriverRequest struct {
	UserID         string `json:"user_id" validate:"required,uuid"`
	LicenseNumber  string `json:"license_number" validate:"required"`
	LicenseType    string `json:"license_type" validate:"required"`
	LicenseExpiry  string `json:"license_expiry" validate:"required"`
	CedulaID       string `json:"cedula_id" validate:"required"`
	EmergencyPhone string `json:"emergency_phone"`
}

type UpdateDriverRequest struct {
	LicenseNumber  string `json:"license_number"`
	LicenseType    string `json:"license_type"`
	LicenseExpiry  string `json:"license_expiry"`
	EmergencyPhone string `json:"emergency_phone"`
	Status         string `json:"status"`
}
type ListDriversResponse struct {
	Drivers []*DriverResponse `json:"drivers"`
	Total   int32             `json:"total"`
}

// --- Students ---
type StudentResponse struct {
	ID             string       `json:"id"`
	FirstName      string       `json:"first_name"`
	LastName       string       `json:"last_name"`
	Grade          string       `json:"grade"`
	SchoolID       string       `json:"school_id"`
	PickupLocation *LocationDTO `json:"pickup_location,omitempty"`
	PickupAddress  string       `json:"pickup_address,omitempty"`
	IsActive       bool         `json:"is_active"`
	CedulaID       string       `json:"cedula_id"`
}

type RegisterStudentRequest struct {
	FirstName      string       `json:"first_name" validate:"required"`
	LastName       string       `json:"last_name" validate:"required"`
	Grade          string       `json:"grade"`
	SchoolID       string       `json:"school_id" validate:"required,uuid"`
	PickupLocation *LocationDTO `json:"pickup_location"`
	PickupAddress  string       `json:"pickup_address"`
	CedulaID       string       `json:"cedula_id" validate:"required"`
}

type UpdateStudentRequest struct {
	FirstName      string       `json:"first_name"`
	LastName       string       `json:"last_name"`
	Grade          string       `json:"grade"`
	PickupLocation *LocationDTO `json:"pickup_location"`
	PickupAddress  string       `json:"pickup_address"`
	CedulaID       string       `json:"cedula_id"`
}

type ListStudentsResponse struct {
	Students []*StudentResponse `json:"students"`
	Total    int32              `json:"total"`
}

// --- Guardians ---
type GuardianResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	StudentID string `json:"student_id"`
	Relation  string `json:"relation"`
	IsPrimary bool   `json:"is_primary"`
}

type LinkGuardianRequest struct {
	UserID    string `json:"user_id" validate:"required,uuid"`
	StudentID string `json:"student_id" validate:"required,uuid"`
	Relation  string `json:"relation" validate:"required"`
	IsPrimary bool   `json:"is_primary"`
}

// --- Routes ---
type RouteResponse struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Direction    string     `json:"direction"`
	ScheduleTime string     `json:"schedule_time"`
	VehicleID    string     `json:"vehicle_id,omitempty"`
	DriverID     string     `json:"driver_id,omitempty"`
	IsActive     bool       `json:"is_active"`
	Stops        []*StopDTO `json:"stops,omitempty"`
}

type StopDTO struct {
	ID        string       `json:"id"`
	StudentID string       `json:"student_id"`
	Order     int32        `json:"order"`
	Location  *LocationDTO `json:"location"`
	Address   string       `json:"address"`
	EstTime   string       `json:"est_time"`
}

type CreateRouteRequest struct {
	Name         string `json:"name" validate:"required"`
	Direction    string `json:"direction" validate:"required"`
	ScheduleTime string `json:"schedule_time" validate:"required"`
	SchoolID     string `json:"school_id" validate:"required,uuid"`
	VehicleID    string `json:"vehicle_id,omitempty"`
	DriverID     string `json:"driver_id,omitempty"`
}

type UpdateRouteRequest struct {
	Name         string `json:"name"`
	Direction    string `json:"direction"`
	ScheduleTime string `json:"schedule_time"`
	VehicleID    string `json:"vehicle_id"`
	DriverID     string `json:"driver_id"`
}

type AddStopRequest struct {
	StudentID string       `json:"student_id" validate:"required,uuid"`
	Order     int32        `json:"order"`
	Location  *LocationDTO `json:"location"`
	Address   string       `json:"address"`
	EstTime   string       `json:"est_time"`
}

type ListRoutesResponse struct {
	Routes []*RouteResponse `json:"routes"`
	Total  int32            `json:"total"`
}
