package api

import (
	"github.com/gin-gonic/gin"
)

func configureRouter(api *api) *gin.Engine {
	router := gin.Default()

	router.Use(CORSMiddleware())

	public := router.Group("api/v1")

	public.GET("/get-prices", api.Binance().GetPrices)

	return router
}
