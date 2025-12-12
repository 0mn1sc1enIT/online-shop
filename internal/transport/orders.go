package transport

import (
	"net/http"

	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"github.com/gin-gonic/gin"
)

type orderItemInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

type createOrderInput struct {
	Items []orderItemInput `json:"items" binding:"required,dive"`
}

func (h *Handler) createOrder(c *gin.Context) {
	userId, exists := c.Get(userCtx)
	if !exists {
		h.logger.Warn().
			Msg("CreateOrder: user id not found in context")

		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}

	var input createOrderInput
	if err := c.BindJSON(&input); err != nil {
		h.logger.Warn().
			Err(err).
			Msg("CreateOrder: invalid input body")

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var items []domain.OrderItem
	for _, i := range input.Items {
		items = append(items, domain.OrderItem{
			ProductID: i.ProductID,
			Quantity:  i.Quantity,
		})
	}

	err := h.services.Orders.Create(userId.(uint), domain.Order{Items: items})
	if err != nil {
		h.logger.Error().
			Err(err).
			Uint("user_id", userId.(uint)).
			Int("items_count", len(items)).
			Msg("CreateOrder: failed to create order")

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Uint("user_id", userId.(uint)).
		Int("items_count", len(items)).
		Msg("CreateOrder: order created successfully")

	c.JSON(http.StatusOK, gin.H{"status": "order created"})
}

func (h *Handler) getAllOrders(c *gin.Context) {
	userId, exists := c.Get(userCtx)
	if !exists {
		h.logger.Warn().
			Msg("GetAllOrders: user id not found in context")

		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}

	orders, err := h.services.Orders.GetAllByUserID(userId.(uint))
	if err != nil {
		h.logger.Error().
			Err(err).
			Uint("user_id", userId.(uint)).
			Msg("GetAllOrders: failed to retrieve orders")

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Uint("user_id", userId.(uint)).
		Int("orders_count", len(orders)).
		Msg("GetAllOrders: orders retrieved successfully")

	c.JSON(http.StatusOK, orders)
}
