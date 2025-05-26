package api

import (
	"Le-BlitzZz/streaming-etl-app/internal/dataserver/entity"
	"Le-BlitzZz/streaming-etl-app/internal/dataserver/entity/search"
	"Le-BlitzZz/streaming-etl-app/internal/dataserver/form"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
)

func GetProcessedApartment(router *gin.RouterGroup) {
	router.GET("/processed-apartment/:uid", func(c *gin.Context) {
		if a := entity.FindProcessedApartment(c.Param("uid")); a == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("ProcessedApartment with UUID %s not found", c.Param("uid")),
			})
		} else {
			c.JSON(http.StatusOK, a)
		}
	})
}

func SearchProcessedApartments(router *gin.RouterGroup) {
	router.GET("/processed-apartments", func(c *gin.Context) {
		var frm form.SearchProcessedApartment

		err := c.MustBindWith(&frm, binding.Form)

		if err != nil {
			log.Info("Failed to bind form data", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid form data",
			})
			return
		}

		result, err := search.ProcessedApartments(frm)

		if err != nil {
			log.Error("Failed to search processed apartments", "error", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Failed to search processed apartments",
			})
			return
		}

		AddExpectedCountHeader(c)
		AddLoadedCountHeader(c, entity.ProcessedApartmentsCount())
		AddFilteredCountHeader(c, len(result))

		c.JSON(http.StatusOK, result)
	})
}
