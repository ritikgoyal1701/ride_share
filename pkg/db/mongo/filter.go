package mongo

import (
	"github.com/gin-gonic/gin"
)

type FilterScopes interface {
	ToFilterScopes() map[string]QueryFilter
}

func GetFilterParamScope(ctx *gin.Context, filter FilterScopes) (scopes map[string]QueryFilter) {
	scopes = filter.ToFilterScopes()
	return
}
