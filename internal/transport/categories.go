package transport

import (
	"net/http"
	"strconv"

	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"github.com/gin-gonic/gin"
)

type createCategoryInput struct {
	Name string `json:"name" binding:"required"`
}

func (h *Handler) createCategory(c *gin.Context) {
	var input createCategoryInput

	if err := c.BindJSON(&input); err != nil {
		h.logger.Warn().
			Err(err).
			Msg("CreateCategory: invalid input body")

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.services.Categories.Create(domain.Category{Name: input.Name})
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("name", input.Name).
			Msg("CreateCategory: failed to create category")

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Str("name", input.Name).
		Msg("CreateCategory: category created successfully")

	c.JSON(http.StatusOK, gin.H{"status": "created"})
}

func (h *Handler) getAllCategories(c *gin.Context) {
	cats, err := h.services.Categories.GetAll()
	if err != nil {
		h.logger.Error().
			Err(err).
			Msg("GetAllCategories: failed to get categories list")

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Int("count", len(cats)).
		Msg("GetAllCategories: categories list retrieved successfully")

	c.JSON(http.StatusOK, cats)
}

func (h *Handler) getCategoryById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("id", c.Param("id")).
			Msg("GetCategoryById: invalid id param")

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id param"})
		return
	}

	cat, err := h.services.Categories.GetByID(uint(id))
	if err != nil {
		h.logger.Warn().
			Err(err).
			Int("category_id", id).
			Msg("GetCategoryById: category not found")

		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	h.logger.Info().
		Int("category_id", id).
		Msg("GetCategoryById: category retrieved successfully")

	c.JSON(http.StatusOK, cat)
}

func (h *Handler) updateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("id", c.Param("id")).
			Msg("UpdateCategory: invalid id param")

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id param"})
		return
	}

	var input createCategoryInput
	if err := c.BindJSON(&input); err != nil {
		h.logger.Warn().
			Err(err).
			Msg("UpdateCategory: invalid input body")

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.services.Categories.Update(uint(id), domain.Category{Name: input.Name})
	if err != nil {
		h.logger.Error().
			Err(err).
			Int("category_id", id).
			Str("name", input.Name).
			Msg("UpdateCategory: failed to update category")

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Int("category_id", id).
		Str("name", input.Name).
		Msg("UpdateCategory: category updated successfully")

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *Handler) deleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("id", c.Param("id")).
			Msg("DeleteCategory: invalid id param")

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id param"})
		return
	}

	err = h.services.Categories.Delete(uint(id))
	if err != nil {
		h.logger.Error().
			Err(err).
			Int("category_id", id).
			Msg("DeleteCategory: failed to delete category")

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Int("category_id", id).
		Msg("DeleteCategory: category deleted successfully")

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
