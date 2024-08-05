package ride

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"rideShare/internal/domain/models"
	mongo2 "rideShare/pkg/db/mongo"
	error2 "rideShare/pkg/error"
	"sync"
)

type Repository struct {
	db *mongo.Collection
}

var (
	repo     *Repository
	repoOnce sync.Once
)

func NewRepository(db *mongo.Database) *Repository {
	repoOnce.Do(func() {
		repo = &Repository{
			db: db.Collection("rides"),
		}
	})

	return repo
}

func (r *Repository) CreateRide(ctx context.Context, ride *models.Ride) (cusErr error2.CustomError) {
	result, err := repo.db.InsertOne(ctx, ride)
	if err != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Error in creating ride | err :: %v", err.Error()))
		return
	}

	ride.ID = result.InsertedID.(primitive.ObjectID)
	return
}

func (r *Repository) GetRide(
	ctx context.Context,
	filters map[string]mongo2.QueryFilter,
) (ride *models.Ride, cusErr error2.CustomError) {
	rides, cusErr := r.GetRides(ctx, filters, map[string]interface{}{})
	if cusErr.Exists() {
		return
	}

	if len(rides) == 0 {
		cusErr = error2.NewCustomError(http.StatusBadRequest, fmt.Sprintf("Ride not found"))
		return
	}

	ride = &rides[0]
	return
}

func (r *Repository) GetRides(
	ctx context.Context,
	filters map[string]mongo2.QueryFilter,
	fields map[string]interface{},
) (riders []models.Ride, cusErr error2.CustomError) {
	queryFilter := mongo2.BuildMongoQuery(ctx, filters)
	findOptions := mongo2.NewFindOptions().SetProjection(fields)

	cursor, err := repo.db.Find(ctx, queryFilter, findOptions.GetMongoFindOptions())
	if err != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Fail to find rides | err :: %v", err.Error()))
		return
	}

	if err = cursor.All(ctx, &riders); err != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Unable to decode mongo response for ride | err :: %v", err.Error()))
		return
	}

	return
}

func (r *Repository) GetRidesCount(
	ctx context.Context,
	filters map[string]mongo2.QueryFilter,
) (count int64, cusErr error2.CustomError) {
	queryFilter := mongo2.BuildMongoQuery(ctx, filters)
	count, err := repo.db.CountDocuments(ctx, queryFilter)
	if err != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Error in getting ride count | err :: %v", err.Error()))
		return
	}

	return
}

func (r *Repository) UpdateRide(
	ctx context.Context,
	filters map[string]mongo2.QueryFilter,
	updates map[string]interface{},
) (cusErr error2.CustomError) {
	queryFilter := mongo2.BuildMongoQuery(ctx, filters)
	queryUpdate := mongo2.BuildMongoSetQuery(updates)
	res, err := r.db.UpdateOne(ctx, queryFilter, queryUpdate)
	if err != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Error in updating ride | err :: %v", err.Error()))
		return
	}

	if res.ModifiedCount == 0 {
		cusErr = error2.NewCustomError(http.StatusBadRequest, fmt.Sprintf("Ride Not updated"))
		return
	}

	return
}
