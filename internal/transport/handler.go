package transport

import (
	"github.com/OmniscienIT/GolangAPI/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Handler struct {
	services *service.Service
	logger   *zerolog.Logger
}

func NewHandler(services *service.Service, logger *zerolog.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Группа Auth (Публичная)
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	// Группа API (Приватная)
	api := router.Group("/api", h.userIdentity)
	{
		// Товары
		products := api.Group("/products")
		{
			products.POST("/", h.createProduct)
			products.GET("/", h.getAllProducts)
			products.GET("/:id", h.getProductById)
			products.PATCH("/:id", h.updateProduct)
			products.DELETE("/:id", h.deleteProduct)
		}

		// Заказы
		orders := api.Group("/orders")
		{
			orders.POST("/", h.createOrder)
			orders.GET("/", h.getAllOrders)
		}

		// Категории
		categories := api.Group("/categories")
		{
			categories.POST("/", h.createCategory)
			categories.GET("/", h.getAllCategories)
			categories.GET("/:id", h.getCategoryById)
			categories.PUT("/:id", h.updateCategory)
			categories.DELETE("/:id", h.deleteCategory)
		}
	}

	return router
}
