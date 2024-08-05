package rideService

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
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
	"rideShare/pkg/lock"
	"rideShare/pkg/utils"
	"sync"
	"time"
)

const (
	rideLock = time.Second * 10
)

type Service struct {
	rideRepository   interfaces.RideRepository
	driverService    interfaces.DriverService
	driverRepository interfaces.DriverRepository
	ridersRepository interfaces.RiderRepository
	locker           lock.Locker
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
	locker lock.Locker,
) *Service {
	svcOnce.Do(func() {
		svc = &Service{
			rideRepository:   rideRepo,
			driverService:    driverService,
			driverRepository: driverRepo,
			ridersRepository: ridersRepo,
			locker:           locker,
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
		constants.MongoID: {
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

func (s *Service) AcceptRide(
	ctx context.Context,
	rideID string,
	userDetails models.UserDetails,
) (cusErr error2.CustomError) {
	driver, cusErr := s.driverRepository.GetDriver(ctx, map[string]mongo2.QueryFilter{
		constants.MongoID: {
			mongo2.IDQuery,
			userDetails.ID,
		},
	})
	if cusErr.Exists() {
		return
	}

	lockAcquired, lockErr := s.locker.Lock(ctx, getRideLockKey(rideID), rideLock)
	if lockErr != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("unable to acquire lock | err :: %v", lockErr.Error()))
		return
	}

	if !lockAcquired {
		cusErr = error2.NewCustomError(http.StatusBadRequest, fmt.Sprintf("unable to acquire lock"))
		return
	}

	defer func(locker lock.Locker, ctx context.Context, lockKey string) {
		_, releaseErr := locker.Release(ctx, lockKey)
		if releaseErr != nil {
			log.Printf("unable to release lock for %v", lockKey)
		}
	}(s.locker, ctx, getRideLockKey(rideID))

	lockAcquired, lockErr = s.locker.Lock(ctx, getDriverLockKey(userDetails.ID), rideLock)
	if lockErr != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("unable to acquire lock | err :: %v", lockErr.Error()))
		return
	}

	if !lockAcquired {
		cusErr = error2.NewCustomError(http.StatusBadRequest, fmt.Sprintf("unable to acquire lock"))
		return
	}

	defer func(locker lock.Locker, ctx context.Context, lockKey string) {
		_, releaseErr := locker.Release(ctx, lockKey)
		if releaseErr != nil {
			log.Printf("unable to release lock for %v", lockKey)
		}
	}(s.locker, ctx, getDriverLockKey(userDetails.ID))

	cusErr = s.rideRepository.UpdateRide(ctx, map[string]mongo2.QueryFilter{
		constants.MongoID: {
			mongo2.IDQuery,
			rideID,
		},
		constants.Status: {
			mongo2.ExactQuery,
			models.RideStatusPending,
		},
	}, map[string]interface{}{
		constants.Driver: models.User{
			ID:    userDetails.ID,
			Name:  driver.Name,
			Email: userDetails.Email,
		},
		constants.Status:    models.RideStatusAccepted,
		constants.UpdatedAt: time.Now().UTC(),
	})
	if cusErr.Exists() {
		return
	}

	cusErr = s.driverRepository.UpdateDriver(ctx, map[string]mongo2.QueryFilter{
		constants.MongoID: {
			mongo2.IDQuery,
			userDetails.ID,
		},
		"is_on_ride": {
			mongo2.ExactQuery,
			false,
		},
	}, map[string]interface{}{
		"is_on_ride":        true,
		constants.UpdatedAt: time.Now().UTC(),
	})

	return
}

func (s *Service) VerifyRide(
	ctx context.Context,
	rideID string,
	userDetails models.UserDetails,
	req requests2.VerificationRequest,
) (cusErr error2.CustomError) {
	ride, cusErr := s.rideRepository.GetRide(ctx, map[string]mongo2.QueryFilter{
		constants.MongoID: {
			mongo2.IDQuery,
			rideID,
		},
		constants.Status: {
			mongo2.ExactQuery,
			models.RideStatusAccepted,
		},
	})
	if cusErr.Exists() {
		return
	}

	if ride.Verification != req.OTP {
		cusErr = error2.NewCustomError(http.StatusBadRequest, "otp verification error")
		return
	}

	lockAcquired, lockErr := s.locker.Lock(ctx, getRideLockKey(rideID), rideLock)
	if lockErr != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("unable to acquire lock | err :: %v", lockErr.Error()))
		return
	}

	if !lockAcquired {
		cusErr = error2.NewCustomError(http.StatusBadRequest, fmt.Sprintf("unable to acquire lock"))
		return
	}

	defer func(locker lock.Locker, ctx context.Context, lockKey string) {
		_, releaseErr := locker.Release(ctx, lockKey)
		if releaseErr != nil {
			log.Printf("unable to release lock for %v", lockKey)
		}
	}(s.locker, ctx, getRideLockKey(rideID))

	lockAcquired, lockErr = s.locker.Lock(ctx, getDriverLockKey(userDetails.ID), rideLock)
	if lockErr != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("unable to acquire lock | err :: %v", lockErr.Error()))
		return
	}

	if !lockAcquired {
		cusErr = error2.NewCustomError(http.StatusBadRequest, fmt.Sprintf("unable to acquire lock"))
		return
	}

	defer func(locker lock.Locker, ctx context.Context, lockKey string) {
		_, releaseErr := locker.Release(ctx, lockKey)
		if releaseErr != nil {
			log.Printf("unable to release lock for %v", lockKey)
		}
	}(s.locker, ctx, getDriverLockKey(userDetails.ID))

	cusErr = s.rideRepository.UpdateRide(ctx, map[string]mongo2.QueryFilter{
		constants.MongoID: {
			mongo2.IDQuery,
			rideID,
		},
		constants.Status: {
			mongo2.ExactQuery,
			models.RideStatusAccepted,
		},
	}, map[string]interface{}{
		constants.Status:    models.RideStatusInProgress,
		constants.UpdatedAt: time.Now().UTC(),
	})
	if cusErr.Exists() {
		return
	}

	return
}

func (s *Service) CancelRide(
	ctx context.Context,
	rideID string,
	userDetails models.UserDetails,
) (cusErr error2.CustomError) {
	ride, cusErr := s.rideRepository.GetRide(ctx, map[string]mongo2.QueryFilter{
		constants.MongoID: {
			mongo2.IDQuery,
			rideID,
		},
		constants.Status: {
			mongo2.ExactQuery,
			models.RideStatusAccepted,
		},
	})
	if cusErr.Exists() {
		return
	}

	lockAcquired, lockErr := s.locker.Lock(ctx, getRideLockKey(rideID), rideLock)
	if lockErr != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("unable to acquire lock | err :: %v", lockErr.Error()))
		return
	}

	if !lockAcquired {
		cusErr = error2.NewCustomError(http.StatusBadRequest, fmt.Sprintf("unable to acquire lock"))
		return
	}

	defer func(locker lock.Locker, ctx context.Context, lockKey string) {
		_, releaseErr := locker.Release(ctx, lockKey)
		if releaseErr != nil {
			log.Printf("unable to release lock for %v", lockKey)
		}
	}(s.locker, ctx, getRideLockKey(rideID))

	lockAcquired, lockErr = s.locker.Lock(ctx, getDriverLockKey(ride.Driver.ID), rideLock)
	if lockErr != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("unable to acquire lock | err :: %v", lockErr.Error()))
		return
	}

	if !lockAcquired {
		cusErr = error2.NewCustomError(http.StatusBadRequest, fmt.Sprintf("unable to acquire lock"))
		return
	}

	defer func(locker lock.Locker, ctx context.Context, lockKey string) {
		_, releaseErr := locker.Release(ctx, lockKey)
		if releaseErr != nil {
			log.Printf("unable to release lock for %v", lockKey)
		}
	}(s.locker, ctx, getDriverLockKey(ride.Driver.ID))

	cusErr = s.rideRepository.UpdateRide(ctx, map[string]mongo2.QueryFilter{
		constants.MongoID: {
			mongo2.IDQuery,
			rideID,
		},
		constants.Status: {
			mongo2.ExactQuery,
			models.RideStatusAccepted,
		},
	}, map[string]interface{}{
		constants.Status:    models.RideStatusCancelled,
		constants.UpdatedAt: time.Now().UTC(),
	})
	if cusErr.Exists() {
		return
	}

	cusErr = s.driverRepository.UpdateDriver(ctx, map[string]mongo2.QueryFilter{
		constants.MongoID: {
			mongo2.IDQuery,
			userDetails.ID,
		},
		"is_on_ride": {
			mongo2.ExactQuery,
			true,
		},
	}, map[string]interface{}{
		"is_on_ride":        false,
		constants.UpdatedAt: time.Now().UTC(),
	})

	return
}

func (s *Service) CompleteRide(
	ctx context.Context,
	rideID string,
	userDetails models.UserDetails,
) (cusErr error2.CustomError) {
	ride, cusErr := s.rideRepository.GetRide(ctx, map[string]mongo2.QueryFilter{
		constants.MongoID: {
			mongo2.IDQuery,
			rideID,
		},
		constants.Status: {
			mongo2.ExactQuery,
			models.RideStatusInProgress,
		},
	})
	if cusErr.Exists() {
		return
	}

	lockAcquired, lockErr := s.locker.Lock(ctx, getRideLockKey(rideID), rideLock)
	if lockErr != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("unable to acquire lock | err :: %v", lockErr.Error()))
		return
	}

	if !lockAcquired {
		cusErr = error2.NewCustomError(http.StatusBadRequest, fmt.Sprintf("unable to acquire lock"))
		return
	}

	defer func(locker lock.Locker, ctx context.Context, lockKey string) {
		_, releaseErr := locker.Release(ctx, lockKey)
		if releaseErr != nil {
			log.Printf("unable to release lock for %v", lockKey)
		}
	}(s.locker, ctx, getRideLockKey(rideID))

	lockAcquired, lockErr = s.locker.Lock(ctx, getDriverLockKey(ride.Driver.ID), rideLock)
	if lockErr != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("unable to acquire lock | err :: %v", lockErr.Error()))
		return
	}

	if !lockAcquired {
		cusErr = error2.NewCustomError(http.StatusBadRequest, fmt.Sprintf("unable to acquire lock"))
		return
	}

	defer func(locker lock.Locker, ctx context.Context, lockKey string) {
		_, releaseErr := locker.Release(ctx, lockKey)
		if releaseErr != nil {
			log.Printf("unable to release lock for %v", lockKey)
		}
	}(s.locker, ctx, getDriverLockKey(ride.Driver.ID))

	cusErr = s.rideRepository.UpdateRide(ctx, map[string]mongo2.QueryFilter{
		constants.MongoID: {
			mongo2.IDQuery,
			rideID,
		},
		constants.Status: {
			mongo2.ExactQuery,
			models.RideStatusAccepted,
		},
	}, map[string]interface{}{
		constants.Status:    models.RideStatusCompleted,
		constants.UpdatedAt: time.Now().UTC(),
	})
	if cusErr.Exists() {
		return
	}

	cusErr = s.driverRepository.UpdateDriver(ctx, map[string]mongo2.QueryFilter{
		constants.MongoID: {
			mongo2.IDQuery,
			userDetails.ID,
		},
		"is_on_ride": {
			mongo2.ExactQuery,
			true,
		},
	}, map[string]interface{}{
		"is_on_ride":        false,
		constants.UpdatedAt: time.Now().UTC(),
	})

	return
}

func getDriverLockKey(driverID string) string {
	return "driver_" + driverID
}

func getRideLockKey(rideID string) string {
	return "ride_" + rideID
}
