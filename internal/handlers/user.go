package handler

import (
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/models"
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/repository"
	"2025-Hackathon-Group/2025-lambda-pdf-translator-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	userService *service.UserService
	repo        *repository.UserRepository
	DB          *gorm.DB
}

func NewUserHandler(userService *service.UserService, repo *repository.UserRepository, db *gorm.DB) *UserHandler {
	return &UserHandler{
		userService: userService,
		repo:        repo,
		DB:          db,
	}
}

// Register handles user registration
func (h *UserHandler) Register(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

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
	user, err := h.repo.Create(input.Name, input.Email, hashedPassword)
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
		"user":  user,
		"token": token,
	})
}

// Login handles user login
func (h *UserHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email
	user, err := h.repo.GetByEmail(input.Email)
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
		"user":  user,
		"token": token,
	})
}

// GetUser retrieves the current user's profile
func (h *UserHandler) GetUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	// Get fresh user data from database
	currentUser := user.(*models.User)
	freshUser, err := h.repo.GetByID(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": freshUser})
}

// GetUserOrganisations retrieves all organisations the current user belongs to
func (h *UserHandler) GetUserOrganisations(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	currentUser := user.(*models.User)
	orgs, err := h.repo.GetOrganisations(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve organisations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"organisations": orgs})
}

type UpdateUserInput struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
		return
	}

	dbUser, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
		return
	}

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != nil {
		dbUser.Name = *input.Name
	}
	if input.Email != nil {
		dbUser.Email = *input.Email
	}
	if input.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		dbUser.Password = string(hashedPassword)
	}

	if err := h.DB.Save(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Do not return password
	response := gin.H{
		"id":    dbUser.ID,
		"name":  dbUser.Name,
		"email": dbUser.Email,
	}

	c.JSON(http.StatusOK, response)
}
