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

func GetRawApartment(router *gin.RouterGroup) {
	router.GET("/raw-apartment/:uid", func(c *gin.Context) {
		if a := entity.FindRawApartment(c.Param("uid")); a == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("RawApartment with UUID %s not found", c.Param("uid")),
			})
		} else {
			c.JSON(http.StatusOK, a)
		}
	})
}

func SearchRawApartments(router *gin.RouterGroup) {
	router.GET("/raw-apartments", func(c *gin.Context) {
		var frm form.SearchRawApartment

		err := c.MustBindWith(&frm, binding.Form)

		if err != nil {
			log.Info("Failed to bind form data", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid form data",
			})
			return
		}

		result, err := search.RawApartments(frm)

		if err != nil {
			log.Error("Failed to search raw apartments", "error", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Failed to search raw apartments",
			})
			return
		}

		AddExpectedCountHeader(c)
		AddLoadedCountHeader(c, entity.RawApartmentsCount())
		AddFilteredCountHeader(c, len(result))


		c.JSON(http.StatusOK, result)
	})
}
