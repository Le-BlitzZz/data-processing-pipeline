package api

import (
	"Le-BlitzZz/streaming-etl-app/internal/dataserver/get"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddExpectedCountHeader(c *gin.Context) {
	c.Header("X-Expected-Count", strconv.FormatInt(get.Config().DataSize(), 10))
}

func AddLoadedCountHeader(c *gin.Context, count int64) {
	c.Header("X-Loaded-Count", strconv.FormatInt(count, 10))
}

func AddFilteredCountHeader(c *gin.Context, count int) {
	c.Header("X-Filtered-Count", strconv.Itoa(count))
}

func AddLimitHeader(c *gin.Context, limit int) {
	c.Header("X-Limit", strconv.Itoa(limit))
}
