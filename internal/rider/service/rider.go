package riderService

import (
	"context"
	"net/http"
	"rideShare/constants"
	requests2 "rideShare/internal/controllers/rider/requests"
	"rideShare/internal/domain/interfaces"
	"rideShare/internal/domain/models"
	"rideShare/pkg/db/mongo"
	error2 "rideShare/pkg/error"
	"sync"
	"time"
)

type Service struct {
	riderRepository interfaces.RiderRepository
}

var (
	svc     *Service
	svcOnce sync.Once
)

func NewService(riderRepo interfaces.RiderRepository) *Service {
	svcOnce.Do(func() {
		svc = &Service{
			riderRepository: riderRepo,
		}
	})

	return svc
}

func (s *Service) CreateRider(
	ctx context.Context,
	req requests2.CreateRiderRequest,
) (cusErr error2.CustomError) {
	count, cusErr := s.riderRepository.GetRidersCount(ctx, map[string]mongo.QueryFilter{
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

	rider := &models.Rider{
		Name:      req.Name,
		ContactNo: req.ContactNo,
		Email:     req.Email,
		CreatedAt: time.Now().UTC(),
	}

	cusErr = s.riderRepository.CreateRider(ctx, rider)
	if cusErr.Exists() {
		return
	}

	return
}
