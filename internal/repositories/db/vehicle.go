package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	premitive "go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"vehicles/internal/domain"
)

func (repo *Repository) CreateVehicle(vehicle domain.Vehicle) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.client.Database("test").Collection("vehicles")

	_, err := collection.InsertOne(ctx, vehicle)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) CreateBrand(brand domain.VehicleBrand) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.client.Database("test").Collection("brands")

	_, err := collection.InsertOne(ctx, brand)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) ReadBrands() ([]domain.VehicleBrand, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.client.Database("test").Collection("brands")

	var brands []domain.VehicleBrand
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &brands); err != nil {
		return nil, err
	}

	return brands, nil
}

func (repo *Repository) CreateModel(model domain.VehicleModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.client.Database("test").Collection("models")

	_, err := collection.InsertOne(ctx, model)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) ReadModels() ([]domain.VehicleModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.client.Database("test").Collection("models")

	var models []domain.VehicleModel
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	return models, nil
}

func (repo *Repository) UpdateVehicleStatus(id premitive.ObjectID, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.client.Database("test").Collection("vehicles")

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) ReadVehicles(filter domain.FilterVehicles) ([]domain.Vehicle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.client.Database("test").Collection("vehicles")

	// Build the filter dynamically based on the provided filter values
	query := bson.M{}

	if filter.ID != "" {
		objectID, err := premitive.ObjectIDFromHex(filter.ID)
		if err != nil {
			return nil, err
		}
		query["_id"] = objectID
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	if filter.DriverID != "" {
		query["driver_id"] = filter.DriverID
	}

	var vehicles []domain.Vehicle
	cursor, err := collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &vehicles); err != nil {
		return nil, err
	}

	return vehicles, nil
}
