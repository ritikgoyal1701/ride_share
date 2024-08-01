package driver

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rideShare/internal/controllers/driver/requests"
	"rideShare/internal/domain/interfaces"
	error2 "rideShare/pkg/error"
	"rideShare/pkg/responses"
	"rideShare/pkg/validate"
	"sync"
)

type Controller struct {
	driverService interfaces.DriverService
}

var (
	ctrl     *Controller
	ctrlOnce sync.Once
)

func NewController(
	driverService interfaces.DriverService,
) *Controller {
	ctrlOnce.Do(func() {
		ctrl = &Controller{
			driverService: driverService,
		}
	})
	return ctrl
}

func (ctrl *Controller) CreateDriver(ctx *gin.Context) {
	var req requests.CreateDriverRequest
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

	cusErr := ctrl.driverService.CreateDriver(ctx, req)
	if cusErr.Exists() {
		error2.NewErrorResponse(ctx, cusErr)
		return
	}

	responses.NewSuccessResponse(ctx, responses.NewSuccessMessage("Driver created successfully"))
}
