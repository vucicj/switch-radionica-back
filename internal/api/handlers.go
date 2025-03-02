package api

import (
	"blazperic/radionica/config"
	"blazperic/radionica/internal/repository"
	"blazperic/radionica/internal/service"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	authService *service.AuthService
}

func NewServer(db *sql.DB, cfg *config.Config) *Server {
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	return &Server{authService: authService}
}

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

func SetupRouter(server *Server) *gin.Engine {
	r := gin.Default()
	r.POST("/register", server.RegisterHandler)
	r.POST("/login", server.LoginHandler)
	return r
}
