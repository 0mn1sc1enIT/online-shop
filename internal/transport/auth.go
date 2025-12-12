package transport

import (
	"net/http"

	"github.com/OmniscienIT/GolangAPI/internal/domain"
	"github.com/gin-gonic/gin"
)

type signUpInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *Handler) signUp(c *gin.Context) {
	var input signUpInput

	if err := c.BindJSON(&input); err != nil {
		h.logger.Warn().
			Err(err).
			Msg("SignUp: invalid input body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.services.Authorization.CreateUser(domain.User{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("email", input.Email).
			Msg("SignUp: failed to create user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Uint("user_id", id).
		Str("email", input.Email).
		Msg("SignUp: user registered successfully")

	c.JSON(http.StatusOK, gin.H{"id": id})
}

type signInInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		h.logger.Warn().
			Err(err).
			Msg("SignIn: invalid input body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Email, input.Password)
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("email", input.Email).
			Msg("SignIn: authentication failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info().
		Str("email", input.Email).
		Msg("SignIn: successful login")

	c.JSON(http.StatusOK, gin.H{"token": token})
}
