package handler

import (
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/helper"
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/models"
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/repository"
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	userService    *service.UserService
	userRepository *repository.UserRepository
}

func NewUserHandler(userService *service.UserService, repo *repository.UserRepository, db *gorm.DB) *UserHandler {
	return &UserHandler{
		userService:    userService,
		userRepository: repo,
	}
}

// @Summary Register a new user
// @Description Register a new user with the given name, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body handler.Register.RegisterInput true "User registration data"
// @Success 201 {object} object{user=response.UserBasicResponse,token=string} "User registered successfully"
// @Failure 400 {object} object{error=string} "Invalid input"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	type RegisterInput struct {
		Name     string `json:"name" binding:"required" example:"John Doe"`
		Email    string `json:"email" binding:"required,email" example:"john.doe@example.com"`
		Password string `json:"password" binding:"required,min=6" example:"password"`
	}

	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := h.userService.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Create the user
	user, err := h.userRepository.Create(input.Name, input.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := h.userService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":  helper.ToUserResponse(*user),
		"token": token,
	})
}

// @Summary Login a user
// @Description Login a user with the given email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body handler.Login.LoginInput true "User login data"
// @Success 200 {object} object{user=response.UserBasicResponse,token=string} "User logged in successfully"
// @Failure 400 {object} object{error=string} "Invalid input"
// @Failure 401 {object} object{error=string} "Invalid credentials"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	type LoginInput struct {
		Email    string `json:"email" binding:"required,email" example:"john.doe@example.com"`
		Password string `json:"password" binding:"required" example:"password"`
	}

	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email
	user, err := h.userRepository.GetByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password
	if !h.userService.VerifyPassword(user.Password, input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := h.userService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  helper.ToUserResponse(*user),
		"token": token,
	})
}

// @Summary Get the current user's profile
// @Description Get the current user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} object{user=models.User} "User profile retrieved successfully"
// @Failure 401 {object} object{error=string} "User not found in context"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /me [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": helper.ToUserResponse(*user.(*models.User))})
}

// @Summary Get the current user's organisations
// @Description Get the current user's organisations
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} object{organisations=[]response.OrganisationBasicResponse} "User organisations retrieved successfully"
// @Failure 401 {object} object{error=string} "User not found in context"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /me/organisations [get]
func (h *UserHandler) GetUserOrganisations(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	currentUser := user.(*models.User)
	orgs, err := h.userRepository.GetOrganisations(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve organisations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"organisations": helper.ToOrganisationsResponse(orgs)})
}

// @Summary Update the current user's profile
// @Description Update the current user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param input body handler.UpdateUser.UpdateUserInput true "User update data"
// @Success 200 {object} object{user=response.UserBasicResponse} "User profile updated successfully"
// @Failure 400 {object} object{error=string} "Invalid input"
// @Failure 401 {object} object{error=string} "User not found in context"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /me [patch]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	type UpdateUserInput struct {
		Name     string `json:"name" example:"John Doe"`
		Email    string `json:"email" example:"john.doe@example.com"`
		Password string `json:"password" binding:"required,min=6" example:"password"`
	}

	var input UpdateUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbUser := c.MustGet("user").(*models.User)

	if input.Name != "" {
		dbUser.Name = input.Name
	}

	if input.Email != "" {
		dbUser.Email = input.Email
	}

	if input.Password != "" {
		hashedPassword, err := h.userService.HashPassword(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		dbUser.Password = hashedPassword
	}

	if err := h.userRepository.Save(dbUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": helper.ToUserResponse(*dbUser),
	})
}
