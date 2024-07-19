package api

import (
	"net/http"

	"binance-integrate/external"
	"binance-integrate/pkg/config"

	"github.com/gin-gonic/gin"
)

type Server struct {
	*http.Server
}

type api struct {
	router         *gin.Engine
	config         config.ServerConfig
	binanceService *external.BinanceService

	binanceHandler *BinanceHandler
}

func NewServer(
	config config.ServerConfig,
	binanceService *external.BinanceService,
) *Server {
	handler := newAPI(config, binanceService)

	srv := &http.Server{
		Addr:              config.ServerPort,
		Handler:           handler,
		ReadHeaderTimeout: config.ReadTimeout.Duration,
	}

	return &Server{
		Server: srv,
	}
}

func newAPI(
	config config.ServerConfig,
	binanceService *external.BinanceService,
) *api {
	api := &api{
		config:         config,
		binanceService: binanceService,
	}

	api.router = configureRouter(api)

	return api
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding,"+
			"X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)

			return
		}

		c.Next()
	}
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func (a *api) Binance() *BinanceHandler {
	if a.binanceHandler == nil {
		a.binanceHandler = NewBinanceHandler(a)
	}

	return a.binanceHandler
}
