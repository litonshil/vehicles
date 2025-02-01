package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
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

func (repo *Repository) UpdateVehicleStatus(id string, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.client.Database("test").Collection("vehicles")

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		return err
	}
	return nil
}

//func (repo *Repository) ReadVehicles() error {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	collection := repo.client.Database("test").Collection("models")
//
//	var models []domain.VehicleModel
//	cursor, err := collection.Find(ctx, bson.M{})
//	if err != nil {
//		return nil, err
//	}
//
//	if err := cursor.All(ctx, &models); err != nil {
//		return nil, err
//	}
//
//	return models, nil
//}
