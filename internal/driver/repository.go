package driver

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
			db: db.Collection("drivers"),
		}
	})

	return repo
}

func (r *Repository) CreateDriver(ctx context.Context, driver *models.Driver) (cusErr error2.CustomError) {
	result, err := repo.db.InsertOne(ctx, driver)
	if err != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, err.Error())
		return
	}

	driver.ID = result.InsertedID.(primitive.ObjectID)
	return
}

func (r *Repository) GetDriver(
	ctx context.Context,
	filters map[string]mongo2.QueryFilter,
) (driver *models.Driver, cusErr error2.CustomError) {
	drivers, cusErr := r.GetDrivers(ctx, filters, map[string]interface{}{})
	if cusErr.Exists() {
		return
	}

	if len(drivers) == 0 {
		cusErr = error2.NewCustomError(http.StatusBadRequest, "driver not found")
		return
	}

	driver = &drivers[0]
	return
}

func (r *Repository) GetDrivers(
	ctx context.Context,
	filters map[string]mongo2.QueryFilter,
	fields map[string]interface{},
) (drivers []models.Driver, cusErr error2.CustomError) {
	queryFilter := mongo2.BuildMongoQuery(ctx, filters)
	findOptions := mongo2.NewFindOptions().SetProjection(fields)

	cursor, err := repo.db.Find(ctx, queryFilter, findOptions.GetMongoFindOptions())
	if err != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Unable to get drivers | err :: %v", err.Error()))
		return
	}

	if err = cursor.All(ctx, &drivers); err != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Unable to decode drivers | err :: %v", err.Error()))
		return
	}

	return
}

func (r *Repository) GetDriversCount(
	ctx context.Context,
	filters map[string]mongo2.QueryFilter,
) (count int64, cusErr error2.CustomError) {
	queryFilter := mongo2.BuildMongoQuery(ctx, filters)
	count, err := repo.db.CountDocuments(ctx, queryFilter)
	if err != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Error in getting driver count | err :: %v", err.Error()))
		return
	}

	return
}

func (r *Repository) UpdateDriver(
	ctx context.Context,
	filters map[string]mongo2.QueryFilter,
	updates map[string]interface{},
) (cusErr error2.CustomError) {
	queryFilter := mongo2.BuildMongoQuery(ctx, filters)
	queryUpdate := mongo2.BuildMongoSetQuery(updates)
	res, err := r.db.UpdateOne(ctx, queryFilter, queryUpdate)
	if err != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Error in updating driver | err :: %v", err.Error()))
		return
	}

	if res.ModifiedCount == 0 {
		cusErr = error2.NewCustomError(http.StatusBadRequest, fmt.Sprintf("Driver Not updated"))
		return
	}

	return
}
