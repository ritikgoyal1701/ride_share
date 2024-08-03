package rideService

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"rideShare/constants"
	"rideShare/internal/controllers/driver/requests"
	requests2 "rideShare/internal/controllers/ride/requests"
	"rideShare/internal/domain/interfaces"
	"rideShare/internal/domain/models"
	"rideShare/internal/driver/service/responses"
	"rideShare/internal/ride/service/adapter"
	responses2 "rideShare/internal/ride/service/responses"
	mongo2 "rideShare/pkg/db/mongo"
	error2 "rideShare/pkg/error"
	"rideShare/pkg/utils"
	"sync"
	"time"
)

type Service struct {
	rideRepository   interfaces.RideRepository
	driverService    interfaces.DriverService
	driverRepository interfaces.DriverRepository
	ridersRepository interfaces.RiderRepository
}

var (
	svc     *Service
	svcOnce sync.Once
)

func NewService(
	rideRepo interfaces.RideRepository,
	driverService interfaces.DriverService,
	driverRepo interfaces.DriverRepository,
	ridersRepo interfaces.RiderRepository,
) *Service {
	svcOnce.Do(func() {
		svc = &Service{
			rideRepository:   rideRepo,
			driverService:    driverService,
			driverRepository: driverRepo,
			ridersRepository: ridersRepo,
		}
	})

	return svc
}

func (s *Service) GetRidePrice(
	ctx context.Context,
	userDetails models.UserDetails,
	req requests2.PriceRequest,
) (resp responses.PriceResponse, cusErr error2.CustomError) {
	rideDistance := utils.CalculateDistance(
		req.StartLocation.XCoordinate,
		req.StartLocation.YCoordinate,
		req.DropLocation.XCoordinate,
		req.DropLocation.YCoordinate,
	)

	isSurge := false
	if utils.IsTimeBetween(
		time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(0, 1, 1, 11, 0, 0, 0, time.UTC),
		time.Now().UTC(),
	) || utils.IsTimeBetween(
		time.Date(0, 1, 1, 19, 0, 0, 0, time.UTC),
		time.Date(0, 1, 1, 22, 0, 0, 0, time.UTC),
		time.Now().UTC(),
	) {
		isSurge = true
	}
	price := utils.CalculatePrice(rideDistance, isSurge)

	nearbyDrivers, cusErr := s.driverService.GetNearbyDrivers(ctx, userDetails, requests.NearByDriversRequest{
		XCoordinate: req.StartLocation.XCoordinate,
		YCoordinate: req.StartLocation.YCoordinate,
	})
	if cusErr.Exists() {
		return
	}

	resp = adapter.GetRidePrice(rideDistance, price, isSurge, nearbyDrivers)
	return
}

func (s *Service) CreateRide(
	ctx context.Context,
	userDetails models.UserDetails,
	req requests2.CreateRideRequest,
) (cusErr error2.CustomError) {
	rider, cusErr := s.ridersRepository.GetRider(ctx, map[string]mongo2.QueryFilter{
		constants.ID: {
			mongo2.IDQuery,
			userDetails.ID,
		},
	})
	if cusErr.Exists() {
		return
	}

	otp, err := utils.GenerateOTP(6)
	if err != nil {
		cusErr = error2.NewCustomError(http.StatusNotImplemented, fmt.Sprintf("Error while generating OTP | err :: %v", err.Error()))
	}

	rideDistance := utils.CalculateDistance(
		req.StartLocation.XCoordinate,
		req.StartLocation.YCoordinate,
		req.DropLocation.XCoordinate,
		req.DropLocation.YCoordinate,
	)

	isSurge := false
	if utils.IsTimeBetween(
		time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(0, 1, 1, 11, 0, 0, 0, time.UTC),
		time.Now().UTC(),
	) || utils.IsTimeBetween(
		time.Date(0, 1, 1, 19, 0, 0, 0, time.UTC),
		time.Date(0, 1, 1, 22, 0, 0, 0, time.UTC),
		time.Now().UTC(),
	) {
		isSurge = true
	}

	price := utils.CalculatePrice(rideDistance, isSurge)

	ride := &models.Ride{
		Rider: models.User{
			ID:    rider.ID.Hex(),
			Name:  rider.Name,
			Email: rider.Email,
		},
		StartLocation: models.Location{
			Coordinates: []float64{req.StartLocation.XCoordinate, req.StartLocation.YCoordinate},
		},
		DropLocation: models.Location{
			Coordinates: []float64{req.DropLocation.XCoordinate, req.DropLocation.YCoordinate},
		},
		Status:       models.RideStatusPending,
		Verification: otp,
		Price:        price,
		Distance:     rideDistance,
		CreatedAt:    time.Now().UTC(),
	}

	cusErr = s.rideRepository.CreateRide(ctx, ride)
	if cusErr.Exists() {
		return
	}

	return
}

func (s *Service) GetRides(
	ctx context.Context,
	userDetails models.UserDetails,
	req requests2.Location,
) (resp []responses2.GetRides, cusErr error2.CustomError) {
	rides, cusErr := s.rideRepository.GetRides(ctx, map[string]mongo2.QueryFilter{
		constants.MongoExpression: {
			Query: mongo2.CustomQuery,
			Value: bson.M{
				"$lt": bson.A{
					bson.M{"$add": bson.A{
						bson.M{"$abs": bson.M{"$subtract": bson.A{bson.M{"$arrayElemAt": bson.A{"$start_location.coordinates", 0}}, req.XCoordinate}}},
						bson.M{"$abs": bson.M{"$subtract": bson.A{bson.M{"$arrayElemAt": bson.A{"$start_location.coordinates", 1}}, req.YCoordinate}}},
					}},
					10,
				},
			}},

		constants.Status: {
			mongo2.ExactQuery,
			models.RideStatusPending,
		},
	}, map[string]interface{}{})
	if cusErr.Exists() {
		return
	}

	resp = adapter.GetRides(rides, req)
	return
}
