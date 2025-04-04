package api

import (
	"database/sql"
	"net/http"

	"github.com/cyberhawk12121/Saarthi/internal/service"
	types "github.com/cyberhawk12121/Saarthi/internal/shared"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *sql.DB) {
	userService := service.NewUserService(db)

	router.POST("/register", func(c *gin.Context) {
		var req types.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := userService.RegisterUser(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	})

	router.POST("/upload", func(c *gin.Context) {
		userService.UploadAudio(c)
	})

}
