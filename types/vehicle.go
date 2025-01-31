package types

import (
	premitive "go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type InsuranceDetails struct {
	PolicyNumber   string    `json:"policy_number" bson:"policy_number"`
	ExpirationDate time.Time `json:"expiration_date" bson:"expiration_date"`
}

type LicenseDetails struct {
	LicenseNumber  string    `json:"license_number" bson:"license_number"`
	ExpirationDate time.Time `json:"expiration_date" bson:"expiration_date"`
}

type VehicleApprovalMessage struct {
	VehicleID string `json:"vehicle_id" bson:"vehicle_id"`
	DriverID  string `json:"driver_id" bson:"driver_id"`
}

type VehicleReq struct {
	ID               premitive.ObjectID `json:"-"`
	DriverID         string             `json:"driver_id" bson:"driver_id"`
	Type             string             `json:"type" bson:"type"` // e.g., "Car", "Truck"
	ModelDetails     VehicleModel       `json:"model_details" bson:"model_details"`
	Color            string             `json:"color" bson:"color"`
	LicenseDetails   LicenseDetails     `json:"license_details" bson:"license_details"`
	InsuranceDetails InsuranceDetails   `json:"insurance_details" bson:"insurance_details"`
	Status           string             `json:"status" bson:"status"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
}

type VehicleBrand struct {
	ID   premitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
}

type VehicleModel struct {
	ID       premitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	BrandID  string             `json:"brand_id" bson:"brand_id"`   // e.g., "Toyota"
	Name     string             `json:"name" bson:"name"`           // e.g., "Corolla"
	Year     int                `json:"year" bson:"year"`           // e.g., 2022
	BodyType string             `json:"body_type" bson:"body_type"` // e.g., "Sedan"
}

type UpdateStatusReq struct {
	ID     premitive.ObjectID `json:"id" bson:"_id"`
	Status string             `json:"status" bson:"status"`
}
