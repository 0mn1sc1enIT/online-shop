package transport

import (
	"net/http"
	"strconv"

	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"github.com/gin-gonic/gin"
)

type createProductInput struct {
	Name        string  `json:"name" binding:"required,min=3"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	CategoryID  uint    `json:"category_id" binding:"required"`
}

func (h *Handler) createProduct(c *gin.Context) {
	var input createProductInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.services.Products.Create(domain.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
		CategoryID:  input.CategoryID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "created"})
}

func (h *Handler) getAllProducts(c *gin.Context) {
	products, err := h.services.Products.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *Handler) getProductById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	product, err := h.services.Products.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *Handler) updateProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input createProductInput // Используем ту же структуру для простоты

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.services.Products.Update(uint(id), domain.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		CategoryID:  input.CategoryID,
		Stock:       input.Stock,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *Handler) deleteProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.services.Products.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
