package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"vehicles/internal/domain"
	"time"
)

//func (repo *Repository) CreateUser(user domain.User) (*types.UserResp, error) {
//	//if err := repo.db.Create(user).Error; err != nil {
//	//	return err
//	//}
//
//	return nil, nil
//}

func (cr *Repository) CreateUser(user *domain.User) error {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	collection := cr.client.Database("test").Collection("users")
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (cr *Repository) GetUsers(filter domain.UserFilter) ([]domain.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := cr.client.Database("test").Collection("users")

	// Build filter
	query := bson.M{}
	if filter.ID != "" {
		objectID, err := primitive.ObjectIDFromHex(filter.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid id format: %w", err)
		}
		query["id"] = objectID
	}

	findOptions := options.Find()
	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, err
	}

	var users []domain.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
