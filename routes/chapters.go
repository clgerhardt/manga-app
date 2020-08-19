package routes

import (
	"manga-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func ChaptersCreate(c *gin.Context) {
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	chapter := models.Chapter{}
	c.ShouldBindJSON(&chapter)
	err := chapter.Create(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, chapter)
}