package presenter

import (
	"Le-BlitzZz/streaming-etl-app/internal/dataserver/api"

	"github.com/gin-gonic/gin"
)

var API *gin.RouterGroup

func registerRoutes() {
	api.GetProcessedApartment(API)
	api.SearchRawApartments(API)

	// api.GetProcessedApartment(API)
	// api.SearchProcessedApartments(API)
}
