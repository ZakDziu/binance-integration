package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BinanceHandler struct {
	api *api
}

func NewBinanceHandler(a *api) *BinanceHandler {
	return &BinanceHandler{
		api: a,
	}
}

//nolint:nestif,gocognit,cyclop
func (h *BinanceHandler) GetPrices(c *gin.Context) {
	c.JSON(http.StatusOK, h.api.binanceService.PriceService.LastPrice.GetAll())
}
