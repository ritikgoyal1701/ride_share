package driverService

import (
	"context"
	"net/http"
	"rideShare/constants"
	"rideShare/internal/controllers/driver/requests"
	"rideShare/internal/domain/interfaces"
	"rideShare/internal/domain/models"
	"rideShare/pkg/db/mongo"
	error2 "rideShare/pkg/error"
	"sync"
	"time"
)

type Service struct {
	driverRepository interfaces.DriverRepository
}

var (
	svc     *Service
	svcOnce sync.Once
)

func NewService(driverRepo interfaces.DriverRepository) *Service {
	svcOnce.Do(func() {
		svc = &Service{
			driverRepository: driverRepo,
		}
	})

	return svc
}

func (s *Service) CreateDriver(
	ctx context.Context,
	req requests.CreateDriverRequest,
) (cusErr error2.CustomError) {
	count, cusErr := s.driverRepository.GetDriversCount(ctx, map[string]mongo.QueryFilter{
		constants.Email: {
			Query: mongo.ExactQuery,
			Value: req.Email,
		},
	})
	if cusErr.Exists() {
		return
	}

	if count > 0 {
		cusErr = error2.NewCustomError(http.StatusNotAcceptable, "Email already exists")
		return
	}

	driver := &models.Driver{
		Name:      req.Name,
		License:   req.License,
		ContactNo: req.ContactNo,
		Email:     req.Email,
		IsActive:  req.IsActive,
		CreatedAt: time.Now().UTC(),
	}

	cusErr = s.driverRepository.CreateDriver(ctx, driver)
	if cusErr.Exists() {
		return
	}

	return
}
