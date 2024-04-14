package v1

import (
	"fourth-exam/post-service-clean-arch/internal/usecase"
	"fourth-exam/post-service-clean-arch/pkg/logger"
	"net/http"

	_ "fourth-exam/post-service-clean-arch/docs"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter -.
// Swagger spec:
// @title       Go Clean Post API
// @description Post service
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1

func NewRouter(handler *gin.Engine, l logger.Interface, t usecase.Post) {
	// Options 
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger 
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe 
	handler.GET("/healthz", func(c *gin.Context) {c.Status(http.StatusOK)})

	// Prometheus metrics 
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers 
	h := handler.Group("/v1") 
	{
		newPostRoutes(h, t, l)
	}
}