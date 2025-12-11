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
	Items []orderItemInput `json:"items" binding:"required,dive"` // dive валидирует элементы массива
}

func (h *Handler) createOrder(c *gin.Context) {
	// Достаем ID пользователя из контекста (положили его туда в middleware)
	userId, exists := c.Get(userCtx)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}

	var input createOrderInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Маппинг DTO -> Domain
	var items []domain.OrderItem
	for _, i := range input.Items {
		items = append(items, domain.OrderItem{
			ProductID: i.ProductID,
			Quantity:  i.Quantity,
		})
	}

	// Вызов сервиса
	err := h.services.Orders.Create(userId.(uint), domain.Order{Items: items})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "order created"})
}

func (h *Handler) getAllOrders(c *gin.Context) {
	userId, exists := c.Get(userCtx)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}

	orders, err := h.services.Orders.GetAllByUserID(userId.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}
