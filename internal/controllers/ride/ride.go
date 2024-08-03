package ride

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rideShare/internal/controllers/ride/requests"
	"rideShare/internal/domain/interfaces"
	error2 "rideShare/pkg/error"
	"rideShare/pkg/responses"
	"rideShare/pkg/utils"
	"rideShare/pkg/validate"
	"sync"
)

type Controller struct {
	rideService interfaces.RideService
}

var (
	ctrl     *Controller
	ctrlOnce sync.Once
)

func NewController(
	rideSvc interfaces.RideService,
) *Controller {
	ctrlOnce.Do(func() {
		ctrl = &Controller{
			rideService: rideSvc,
		}
	})
	return ctrl
}

func (ctrl *Controller) GetRidePrice(ctx *gin.Context) {
	var req requests.PriceRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		cusErr := error2.NewCustomError(http.StatusUnprocessableEntity, err.Error())
		error2.NewErrorResponse(ctx, cusErr)
		return
	}

	validateErr := validate.Get().Struct(req)
	if validateErr != nil {
		cusErr := error2.NewCustomError(http.StatusUnprocessableEntity, validateErr.Error())
		error2.NewErrorResponse(ctx, cusErr)
		return
	}

	userDetails, cusErr := utils.GetUserDetails(ctx)
	if cusErr.Exists() {
		error2.NewErrorResponse(ctx, cusErr)
		return
	}

	resp, cusErr := ctrl.rideService.GetRidePrice(ctx, userDetails, req)
	if cusErr.Exists() {
		error2.NewErrorResponse(ctx, cusErr)
		return
	}

	responses.NewSuccessResponse(ctx, resp)
}

func (ctrl *Controller) CreateRide(ctx *gin.Context) {
	var req requests.CreateRideRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		cusErr := error2.NewCustomError(http.StatusUnprocessableEntity, err.Error())
		error2.NewErrorResponse(ctx, cusErr)
		return
	}

	validateErr := validate.Get().Struct(req)
	if validateErr != nil {
		cusErr := error2.NewCustomError(http.StatusUnprocessableEntity, validateErr.Error())
		error2.NewErrorResponse(ctx, cusErr)
		return
	}

	userDetails, cusErr := utils.GetUserDetails(ctx)
	if cusErr.Exists() {
		error2.NewErrorResponse(ctx, cusErr)
		return
	}

	cusErr = ctrl.rideService.CreateRide(ctx, userDetails, req)
	if cusErr.Exists() {
		error2.NewErrorResponse(ctx, cusErr)
		return
	}

	responses.NewSuccessResponse(ctx, responses.NewSuccessMessage("Ride created successfully"))
}

func (ctrl *Controller) GetRides(ctx *gin.Context) {
	var req requests.Location
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		cusErr := error2.NewCustomError(http.StatusUnprocessableEntity, err.Error())
		error2.NewErrorResponse(ctx, cusErr)
		return
	}

	validateErr := validate.Get().Struct(req)
	if validateErr != nil {
		cusErr := error2.NewCustomError(http.StatusUnprocessableEntity, validateErr.Error())
		error2.NewErrorResponse(ctx, cusErr)
		return
	}

	userDetails, cusErr := utils.GetUserDetails(ctx)
	if cusErr.Exists() {
		error2.NewErrorResponse(ctx, cusErr)
		return
	}

	resp, cusErr := ctrl.rideService.GetRides(ctx, userDetails, req)
	if cusErr.Exists() {
		error2.NewErrorResponse(ctx, cusErr)
		return
	}

	responses.NewSuccessResponse(ctx, resp)
}
