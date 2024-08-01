package riderService

import (
	"context"
	"fmt"
	"net/http"
	"rideShare/constants"
	"rideShare/internal/controllers/driver/requests"
	requests2 "rideShare/internal/controllers/rider/requests"
	"rideShare/internal/domain/interfaces"
	"rideShare/internal/domain/models"
	"rideShare/internal/rider/service/responses"
	"rideShare/pkg/db/mongo"
	error2 "rideShare/pkg/error"
	"rideShare/pkg/hash"
	"rideShare/pkg/jwt"
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
		Password:  hash.GetHashedPassword(req.Password),
		CreatedAt: time.Now().UTC(),
	}

	cusErr = s.riderRepository.CreateRider(ctx, rider)
	if cusErr.Exists() {
		return
	}

	return
}

func (s *Service) Login(
	ctx context.Context,
	req requests.LoginRequest,
) (resp responses.LoginResp, cusErr error2.CustomError) {
	rider, cusErr := s.riderRepository.GetRider(ctx, map[string]mongo.QueryFilter{
		constants.Email: {
			Query: mongo.ExactQuery,
			Value: req.Email,
		},
	})
	if cusErr.Exists() {
		return
	}

	if rider.Password != hash.GetHashedPassword(req.Password) {
		cusErr = error2.NewCustomError(http.StatusBadRequest, "Incorrect password")
		return
	}

	isValid, err := isTokenValid(rider.Jwt)
	if err != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Error validating token: %v", err.Error()))
		return
	}

	if isValid {
		resp = responses.LoginResp{
			Token: rider.Jwt,
		}

		return
	}

	token, err := jwt.GenerateToken(rider.Email, rider.ID.Hex(), models.TitleRider)
	if err != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Unable to generate JWT token | err :: %v", err.Error()))
		return
	}

	cusErr = s.riderRepository.UpdateRider(ctx, map[string]mongo.QueryFilter{
		constants.MongoID: {
			mongo.ExactQuery,
			rider.ID,
		},
	}, map[string]interface{}{
		constants.Jwt: token,
	})
	if cusErr.Exists() {
		return
	}

	resp = responses.LoginResp{
		Token: token,
	}

	return
}

func isTokenValid(token string) (isValid bool, err error) {
	if len(token) == 0 {
		return
	}

	_, isValid, err = jwt.ValidateToken(token)
	return
}

func (s *Service) Logout(ctx context.Context, userDetails models.UserDetails) (cusErr error2.CustomError) {
	cusErr = s.riderRepository.UpdateRider(ctx, map[string]mongo.QueryFilter{
		constants.MongoID: {
			mongo.IDQuery,
			userDetails.ID,
		},
	}, map[string]interface{}{
		constants.Jwt: "",
	})
	if cusErr.Exists() {
		return
	}

	return
}
