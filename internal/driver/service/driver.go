package driverService

import (
	"context"
	"fmt"
	"net/http"
	"rideShare/constants"
	"rideShare/internal/controllers/driver/requests"
	"rideShare/internal/domain/interfaces"
	"rideShare/internal/domain/models"
	"rideShare/internal/driver/service/responses"
	"rideShare/pkg/db/mongo"
	error2 "rideShare/pkg/error"
	"rideShare/pkg/hash"
	"rideShare/pkg/jwt"
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
		Password:  hash.GetHashedPassword(req.Password),
		CreatedAt: time.Now().UTC(),
	}

	cusErr = s.driverRepository.CreateDriver(ctx, driver)
	if cusErr.Exists() {
		return
	}

	return
}

func (s *Service) Login(
	ctx context.Context,
	req requests.LoginRequest,
) (resp responses.LoginResp, cusErr error2.CustomError) {
	driver, cusErr := s.driverRepository.GetDriver(ctx, map[string]mongo.QueryFilter{
		constants.Email: {
			Query: mongo.ExactQuery,
			Value: req.Email,
		},
	})
	if cusErr.Exists() {
		return
	}

	if driver.Password != hash.GetHashedPassword(req.Password) {
		cusErr = error2.NewCustomError(http.StatusBadRequest, "Incorrect password")
		return
	}

	isValid, err := isTokenValid(driver.Jwt)
	if err != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Error validating token: %v", err.Error()))
		return
	}

	if isValid {
		resp = responses.LoginResp{
			Token: driver.Jwt,
		}

		return
	}

	token, err := jwt.GenerateToken(driver.Email, driver.ID.Hex(), models.TitleDriver)
	if err != nil {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("Unable to generate JWT token | err :: %v", err.Error()))
		return
	}

	cusErr = s.driverRepository.UpdateDriver(ctx, map[string]mongo.QueryFilter{
		constants.MongoID: {
			mongo.ExactQuery,
			driver.ID,
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

func (s *Service) Logout(ctx context.Context, userDetails models.UserDetails) (cusErr error2.CustomError) {
	cusErr = s.driverRepository.UpdateDriver(ctx, map[string]mongo.QueryFilter{
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

func isTokenValid(token string) (isValid bool, err error) {
	if len(token) == 0 {
		return
	}

	_, isValid, err = jwt.ValidateToken(token)
	return
}
