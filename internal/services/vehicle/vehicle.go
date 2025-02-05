package vehicle

import (
	premitive "go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"vehicles/infra/rabbitmq"
	"vehicles/internal/domain"
	"vehicles/types"
)

type VehicleService struct {
	repo     domain.VehicleRepository
	rabbitMQ *rabbitmq.RabbitMQ
}

func NewVehicleService(vehicleRepo domain.VehicleRepository, rabbitmq *rabbitmq.RabbitMQ) *VehicleService {
	return &VehicleService{
		repo:     vehicleRepo,
		rabbitMQ: rabbitmq,
	}
}

func (service *VehicleService) CreateVehicle(vehicle types.VehicleReq) error {

	req := domain.Vehicle{
		ID:        vehicle.ID,
		DriverID:  vehicle.DriverID,
		Type:      vehicle.Type,
		Color:     vehicle.Color,
		Status:    "in_process",
		CreatedAt: time.Now().UTC(),
	}

	if vehicle.LicenseDetails.LicenseNumber != "" {
		req.LicenseDetails = domain.LicenseDetails{
			LicenseNumber:  vehicle.LicenseDetails.LicenseNumber,
			ExpirationDate: vehicle.LicenseDetails.ExpirationDate,
		}
	}

	if vehicle.InsuranceDetails.PolicyNumber != "" {
		req.InsuranceDetails = domain.InsuranceDetails{
			PolicyNumber:   vehicle.InsuranceDetails.PolicyNumber,
			ExpirationDate: vehicle.InsuranceDetails.ExpirationDate,
		}
	}

	if vehicle.ModelDetails.ID != premitive.NilObjectID {
		req.ModelDetails = domain.VehicleModel{
			ID:       vehicle.ModelDetails.ID,
			BrandID:  vehicle.ModelDetails.BrandID,
			Name:     vehicle.ModelDetails.Name,
			Year:     vehicle.ModelDetails.Year,
			BodyType: vehicle.ModelDetails.BodyType,
		}
	}

	return service.repo.CreateVehicle(req)
}

func (service *VehicleService) CreateBrand(vehicle types.VehicleBrand) error {

	req := domain.VehicleBrand{
		Name: vehicle.Name,
	}

	return service.repo.CreateBrand(req)
}

func (service *VehicleService) ReadBrands() ([]types.VehicleBrand, error) {
	brands, err := service.repo.ReadBrands()
	if err != nil {
		return nil, err
	}

	var res []types.VehicleBrand
	for _, brand := range brands {
		res = append(res, types.VehicleBrand{
			ID:   brand.ID,
			Name: brand.Name,
		})
	}

	return res, nil
}

func (service *VehicleService) CreateModel(vehicle types.VehicleModel) error {
	req := domain.VehicleModel{
		BrandID:  vehicle.BrandID,
		Name:     vehicle.Name,
		Year:     vehicle.Year,
		BodyType: vehicle.BodyType,
	}

	return service.repo.CreateModel(req)
}

func (service *VehicleService) ReadModels() ([]types.VehicleModel, error) {
	models, err := service.repo.ReadModels()
	if err != nil {
		return nil, err
	}

	var res []types.VehicleModel
	for _, model := range models {
		res = append(res, types.VehicleModel{
			ID:       model.ID,
			BrandID:  model.BrandID,
			Name:     model.Name,
			Year:     model.Year,
			BodyType: model.BodyType,
		})
	}

	return res, nil
}

func (service *VehicleService) UpdateVehicleStatus(req types.UpdateStatusReq) error {
	err := service.repo.UpdateVehicleStatus(req.ID, req.Status)
	if err != nil {
		return err
	}

	if req.Status == "approved" {
		filterVehicle := domain.FilterVehicles{
			ID: req.ID.Hex(),
		}

		vehicle, err := service.repo.ReadVehicles(filterVehicle)
		if err != nil {
			return err
		}
		vehicleApprovedMessage := types.VehicleApprovalMessage{
			VehicleID: req.ID.Hex(),
			DriverID:  vehicle[0].DriverID,
		}

		go func() {
			_ = service.rabbitMQ.Publish("vehicle.approved", vehicleApprovedMessage)
		}()
	}
	return nil
}
