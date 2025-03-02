package api

import (
	"database/sql"
	"net/http"

	"blazperic/radionica/config"
	"blazperic/radionica/internal/models"
	"blazperic/radionica/internal/repository"
	"blazperic/radionica/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Server manages dependencies for HTTP handlers
type Server struct {
	authService AuthService
	newsService NewsService
}

// AuthService defines authentication operations
type AuthService interface {
	Register(username, password string) (*models.User, error)
	Login(username, password string) (*service.TokenPair, error)
	RefreshToken(refreshToken string) (*service.TokenPair, error)
}

// NewsService defines news-related operations
type NewsService interface {
	GetAllNews() ([]*models.News, error)
	CreateNews(title, content string, userID uuid.UUID) (*models.News, error)
}

// NewServer initializes a Server with injected dependencies
func NewServer(db *sql.DB, cfg *config.Config) *Server {
	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.TokenDuration, cfg.RefreshTokenDuration)
	newsRepo := repository.NewNewsRepository(db)
	newsSvc := service.NewNewsService(newsRepo)
	return &Server{
		authService: authSvc,
		newsService: newsSvc,
	}
}

// RegisterHandler handles user registration
// @Summary Register a new user
// @Description Creates a new user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration details"
// @Success 201 {object} RegisterResponse "User created"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /auth/register [post]
func (s *Server) RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	user, err := s.authService.Register(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to register user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, RegisterResponse{UserID: user.ID.String()})
}

// LoginHandler handles user login
// @Summary Login a user
// @Description Authenticates a user and returns access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} service.TokenPair "Login successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /auth/login [post]
func (s *Server) LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	tokens, err := s.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// RefreshTokenHandler refreshes an access token
// @Summary Refresh access token
// @Description Generates a new access and refresh token pair using a valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh body RefreshRequest true "Refresh token"
// @Success 200 {object} service.TokenPair "Tokens refreshed"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /auth/refresh [post]
func (s *Server) RefreshTokenHandler(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	tokens, err := s.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// GetNewsHandler retrieves all news items
// @Summary Get all news
// @Description Fetches a list of all news items
// @Tags news
// @Produce json
// @Success 200 {array} models.News "News list"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /news [get]
func (s *Server) GetNewsHandler(c *gin.Context) {
	news, err := s.newsService.GetAllNews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch news: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, news)
}

// CreateNewsHandler creates a new news item
// @Summary Create a news item
// @Description Adds a new news item (requires authentication)
// @Tags news
// @Accept json
// @Produce json
// @Param news body CreateNewsRequest true "News details"
// @Security BearerAuth
// @Success 201 {object} models.News "News created"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /news [post]
func (s *Server) CreateNewsHandler(c *gin.Context) {
	var req CreateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Authentication required"})
		return
	}

	news, err := s.newsService.CreateNews(req.Title, req.Content, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create news: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, news)
}

// SetupRouter configures the Gin router with grouped endpoints
func SetupRouter(server *Server, jwtSecret string) *gin.Engine {
	r := gin.Default()

	// API version 1 group
	apiV1 := r.Group("/api/v1")
	{
		// Authentication routes
		auth := apiV1.Group("/auth")
		{
			auth.POST("/register", server.RegisterHandler)
			auth.POST("/login", server.LoginHandler)
			auth.POST("/refresh", server.RefreshTokenHandler)
		}

		// News routes
		news := apiV1.Group("/news")
		{
			news.GET("", server.GetNewsHandler)
			news.POST("", JWTAuth(jwtSecret), server.CreateNewsHandler)
		}
	}

	return r
}

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterResponse represents the response for a successful registration
type RegisterResponse struct {
	UserID string `json:"user_id"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RefreshRequest represents the request body for token refresh
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ErrorResponse represents a generic error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// CreateNewsRequest represents the request body for creating news
type CreateNewsRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}
