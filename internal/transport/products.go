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

type updateProductInput struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price" binding:"omitempty,gt=0"`
	Stock       *int     `json:"stock" binding:"omitempty,gte=0"`
	CategoryID  *uint    `json:"category_id"`
}

func (h *Handler) createProduct(c *gin.Context) {
	var input createProductInput

	if err := c.BindJSON(&input); err != nil {
		h.logger.Warn().
			Err(err).
			Msg("CreateProduct: invalid input body")

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
		h.logger.Error().
			Err(err).
			Str("name", input.Name).
			Uint("category_id", input.CategoryID).
			Msg("CreateProduct: failed to create product")

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Str("name", input.Name).
		Uint("category_id", input.CategoryID).
		Msg("CreateProduct: product created successfully")

	c.JSON(http.StatusOK, gin.H{"status": "created"})
}

func (h *Handler) getAllProducts(c *gin.Context) {
	products, err := h.services.Products.GetAll()
	if err != nil {
		h.logger.Error().
			Err(err).
			Msg("GetAllProducts: failed to retrieve products")

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Int("count", len(products)).
		Msg("GetAllProducts: products retrieved successfully")

	c.JSON(http.StatusOK, products)
}

func (h *Handler) getProductById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("id", c.Param("id")).
			Msg("GetProductById: invalid id param")

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id param"})
		return
	}

	product, err := h.services.Products.GetByID(uint(id))
	if err != nil {
		h.logger.Warn().
			Err(err).
			Int("product_id", id).
			Msg("GetProductById: product not found")

		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	h.logger.Info().
		Int("product_id", id).
		Msg("GetProductById: product retrieved successfully")

	c.JSON(http.StatusOK, product)
}

func (h *Handler) updateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("id", c.Param("id")).
			Msg("UpdateProduct: invalid id param")

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id param"})
		return
	}

	var input updateProductInput
	if err := c.BindJSON(&input); err != nil {
		h.logger.Warn().
			Err(err).
			Msg("UpdateProduct: invalid input body")

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingProduct, err := h.services.Products.GetByID(uint(id))
	if err != nil {
		h.logger.Warn().
			Err(err).
			Int("product_id", id).
			Msg("UpdateProduct: product not found")

		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	if input.Name != nil {
		existingProduct.Name = *input.Name
	}
	if input.Description != nil {
		existingProduct.Description = *input.Description
	}
	if input.Price != nil {
		existingProduct.Price = *input.Price
	}
	if input.Stock != nil {
		existingProduct.Stock = *input.Stock
	}
	if input.CategoryID != nil {
		existingProduct.CategoryID = *input.CategoryID
	}

	err = h.services.Products.Update(uint(id), existingProduct)
	if err != nil {
		h.logger.Error().
			Err(err).
			Int("product_id", id).
			Msg("UpdateProduct: failed to update product")

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Int("product_id", id).
		Msg("UpdateProduct: product updated successfully")

	c.JSON(http.StatusOK, gin.H{"status": "updated", "product": existingProduct})
}

func (h *Handler) deleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("id", c.Param("id")).
			Msg("DeleteProduct: invalid id param")

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id param"})
		return
	}

	err = h.services.Products.Delete(uint(id))
	if err != nil {
		h.logger.Error().
			Err(err).
			Int("product_id", id).
			Msg("DeleteProduct: failed to delete product")

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Int("product_id", id).
		Msg("DeleteProduct: product deleted successfully")

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
