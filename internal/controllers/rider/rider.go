package rider

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rideShare/internal/controllers/driver/requests"
	requests2 "rideShare/internal/controllers/rider/requests"
	"rideShare/internal/domain/interfaces"
	error2 "rideShare/pkg/error"
	"rideShare/pkg/responses"
	"rideShare/pkg/utils"
	"rideShare/pkg/validate"
	"sync"
)

type Controller struct {
	riderService interfaces.RiderService
}

var (
	ctrl     *Controller
	ctrlOnce sync.Once
)

func NewController(
	riderService interfaces.RiderService,
) *Controller {
	ctrlOnce.Do(func() {
		ctrl = &Controller{
			riderService: riderService,
		}
	})
	return ctrl
}

func (ctrl *Controller) CreateRider(ctx *gin.Context) {
	var req requests2.CreateRiderRequest
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

	cusErr := ctrl.riderService.CreateRider(ctx, req)
	if cusErr.Exists() {
		error2.NewErrorResponse(ctx, cusErr)
		return
	}

	responses.NewSuccessResponse(ctx, responses.NewSuccessMessage("Rider created successfully"))
}

func (ctrl *Controller) Login(ctx *gin.Context) {
	var req requests.LoginRequest
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

	resp, cusErr := ctrl.riderService.Login(ctx, req)
	if cusErr.Exists() {
		error2.NewErrorResponse(ctx, cusErr)
		return
	}

	responses.NewSuccessResponse(ctx, resp)
}

func (ctrl *Controller) Logout(ctx *gin.Context) {
	userDetails, cusErr := utils.GetUserDetails(ctx)
	if cusErr.Exists() {
		error2.NewErrorResponse(ctx, cusErr)
		return
	}

	cusErr = ctrl.riderService.Logout(ctx, userDetails)
	if cusErr.Exists() {
		error2.NewErrorResponse(ctx, cusErr)
		return
	}

	responses.NewSuccessResponse(ctx, responses.NewSuccessMessage("Rider logout successfully"))
}
