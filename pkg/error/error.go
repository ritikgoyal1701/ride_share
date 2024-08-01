package error

import (
	"github.com/gin-gonic/gin"
)

type CustomError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	IsExist    bool   `json:"exist"`
}

type ErrorResponse struct {
	IsSuccess  bool  `json:"is_success"`
	StatusCode int   `json:"status_code"`
	Error      Error `json:"error"`
}

type Error struct {
	Message string `json:"message"`
}

func NewErrorResponse(ctx *gin.Context, cusErr CustomError) {
	res := ErrorResponse{
		IsSuccess:  false,
		StatusCode: cusErr.StatusCode,
		Error: Error{
			Message: cusErr.Message,
		},
	}

	ctx.JSON(cusErr.StatusCode, res)
}

func NewCustomError(errorCode int, errorMessage string) (cusErr CustomError) {
	cusErr = CustomError{
		StatusCode: errorCode,
		Message:    errorMessage,
		IsExist:    true,
	}
	return
}

func (cusErr CustomError) Exists() bool {
	return cusErr.IsExist
}
