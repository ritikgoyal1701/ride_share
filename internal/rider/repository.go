package rider

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
			db: db.Collection("riders"),
		}
	})

	return repo
}

func (r *Repository) CreateRider(ctx context.Context, rider *models.Rider) (cusErr error2.CustomError) {
	result, err := repo.db.InsertOne(ctx, rider)
	if err != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Error in creating rider | err :: %v", err.Error()))
		return
	}

	rider.ID = result.InsertedID.(primitive.ObjectID)
	return
}

func (r *Repository) GetRider(
	ctx context.Context,
	filters map[string]mongo2.QueryFilter,
) (rider *models.Rider, cusErr error2.CustomError) {
	riders, cusErr := r.GetRiders(ctx, filters, map[string]interface{}{})
	if cusErr.Exists() {
		return
	}

	if len(riders) == 0 {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Rider not found"))
		return
	}

	rider = &riders[0]
	return
}

func (r *Repository) GetRiders(
	ctx context.Context,
	filters map[string]mongo2.QueryFilter,
	fields map[string]interface{},
) (riders []models.Rider, cusErr error2.CustomError) {
	queryFilter := mongo2.BuildMongoQuery(ctx, filters)
	findOptions := mongo2.NewFindOptions().SetProjection(fields)

	cursor, err := repo.db.Find(ctx, queryFilter, findOptions.GetMongoFindOptions())
	if err != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Fail to find riders | err :: %v", err.Error()))
		return
	}

	if err = cursor.All(ctx, &riders); err != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Unable to decode mongo response for rider | err :: %v", err.Error()))
		return
	}

	return
}

func (r *Repository) GetRidersCount(
	ctx context.Context,
	filters map[string]mongo2.QueryFilter,
) (count int64, cusErr error2.CustomError) {
	queryFilter := mongo2.BuildMongoQuery(ctx, filters)
	count, err := repo.db.CountDocuments(ctx, queryFilter)
	if err != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Error in getting rider count | err :: %v", err.Error()))
		return
	}

	return
}
