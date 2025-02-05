package domain

import (
	premitive "go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"vehicles/types"
)

type InsuranceDetails struct {
	PolicyNumber   string    `json:"policy_number" bson:"policy_number"`
	ExpirationDate time.Time `json:"expiration_date" bson:"expiration_date"`
}

type LicenseDetails struct {
	LicenseNumber  string    `json:"license_number" bson:"license_number"`
	ExpirationDate time.Time `json:"expiration_date" bson:"expiration_date"`
}

type Vehicle struct {
	ID               premitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DriverID         string             `json:"driver_id" bson:"driver_id"`
	Type             string             `json:"type" bson:"type"`
	Color            string             `json:"color" bson:"color"`
	ModelDetails     VehicleModel       `json:"model_details" bson:"model_details"`
	LicenseDetails   LicenseDetails     `json:"license_details" bson:"license_details"`
	InsuranceDetails InsuranceDetails   `json:"insurance_details" bson:"insurance_details"`
	Status           string             `json:"status" bson:"status"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
}

type VehicleBrand struct {
	ID   premitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
}

type VehicleModel struct {
	ID       premitive.ObjectID `json:"id" bson:"_id,omitempty"`
	BrandID  string             `json:"brand_id" bson:"brand_id"`
	Name     string             `json:"name" bson:"name"`
	Year     int                `json:"year" bson:"year"`
	BodyType string             `json:"body_type" bson:"body_type"`
}

type FilterVehicles struct {
	ID       string `json:"id" bson:"id" query:"id"`
	DriverID string `json:"driver_id" bson:"driver_id" query:"driver_id"`
	Status   string `json:"status" bson:"status" query:"status"`
}

type VehicleRepository interface {
	CreateVehicle(vehicle Vehicle) error
	ReadVehicles(vehicles FilterVehicles) ([]Vehicle, error)
	CreateBrand(brand VehicleBrand) error
	ReadBrands() ([]VehicleBrand, error)
	CreateModel(model VehicleModel) error
	ReadModels() ([]VehicleModel, error)
	UpdateVehicleStatus(id premitive.ObjectID, status string) error
}

type VehicleUseCase interface {
	CreateVehicle(vehicle types.VehicleReq) error
	CreateBrand(brand types.VehicleBrand) error
	ReadBrands() ([]types.VehicleBrand, error)
	CreateModel(model types.VehicleModel) error
	ReadModels() ([]types.VehicleModel, error)
	UpdateVehicleStatus(req types.UpdateStatusReq) error
}
