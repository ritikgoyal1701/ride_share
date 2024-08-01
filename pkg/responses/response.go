package responses

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type SuccessMessage struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	IsSuccess  bool        `json:"is_success"`
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
	Meta       interface{} `json:"meta"`
}

func NewSuccessResponse(ctx *gin.Context, data interface{}) {
	res := &SuccessResponse{
		IsSuccess:  true,
		StatusCode: http.StatusOK,
		Data:       data,
	}

	ctx.AbortWithStatusJSON(http.StatusOK, res)
	return
}

func NewSuccessResponseWithMeta(ctx *gin.Context, data interface{}, meta interface{}) {
	res := &SuccessResponse{
		IsSuccess:  true,
		StatusCode: http.StatusOK,
		Data:       data,
		Meta:       meta,
	}

	ctx.AbortWithStatusJSON(http.StatusOK, res)
	return
}

func NewSuccessMessage(message string) SuccessMessage {
	res := SuccessMessage{
		Message: message,
	}

	return res
}
