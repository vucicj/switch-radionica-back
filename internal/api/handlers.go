package api

import (
	"database/sql"
	"net/http"

	"blazperic/radionica/config"
	"blazperic/radionica/internal/repository"
	"blazperic/radionica/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Server holds the dependencies for HTTP handlers
type Server struct {
	authService *service.AuthService
	newsService *service.NewsService
}

// NewServer creates a new Server instance
func NewServer(db *sql.DB, cfg *config.Config) *Server {
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	newsRepo := repository.NewNewsRepository(db)
	newsService := service.NewNewsService(newsRepo)
	return &Server{authService: authService, newsService: newsService}
}

// RegisterHandler godoc
// @Summary Register a new user
// @Description Creates a new user with the provided username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration details"
// @Success 201 {object} RegisterResponse "User created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /register [post]
func (s *Server) RegisterHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := s.authService.Register(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user_id": user.ID})
}

// LoginHandler godoc
// @Summary Login a user
// @Description Authenticates a user and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User login credentials"
// @Success 200 {object} LoginResponse "Login successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /login [post]
func (s *Server) LoginHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := s.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetNewsHandler godoc
// @Summary Get all news
// @Description Retrieves a list of all news items
// @Tags news
// @Produce json
// @Success 200 {array} models.News "List of news"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /news [get]
func (s *Server) GetNewsHandler(c *gin.Context) {
	news, err := s.newsService.GetAllNews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, news)
}

// CreateNewsHandler godoc
// @Summary Create a new news item
// @Description Creates a new news item (requires authentication)
// @Tags news
// @Accept json
// @Produce json
// @Param news body CreateNewsRequest true "News details"
// @Security BearerAuth
// @Success 201 {object} models.News "News created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /news [post]
func (s *Server) CreateNewsHandler(c *gin.Context) {
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	news, err := s.newsService.CreateNews(req.Title, req.Content, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, news)
}

// SetupRouter configures the Gin router with all endpoints
func SetupRouter(server *Server, jwtSecret string) *gin.Engine {
	r := gin.Default()
	r.POST("/register", func(c *gin.Context) { server.RegisterHandler(c) })
	r.POST("/login", func(c *gin.Context) { server.LoginHandler(c) })
	r.GET("/news", func(c *gin.Context) { server.GetNewsHandler(c) })
	r.POST("/news", JWTAuth(jwtSecret), func(c *gin.Context) { server.CreateNewsHandler(c) })
	return r
}

// Define additional structs for Swagger
type RegisterResponse struct {
	UserID string `json:"user_id"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateNewsRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}
