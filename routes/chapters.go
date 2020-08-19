package routes

import (
	"fmt"
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

func ChaptersUpdate(c *gin.Context)  {
	userID := c.GetString("user_id")
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	chapterSent := models.Chapter{}
	err := c.ShouldBindJSON(&chapterSent)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form sent"})
		return
	}

	chapterBeingUpdated, err := models.FindChapterById(chapterSent.ID, &conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collectionOfChapterBeingDelete, err := models.FindCollectionById(chapterBeingUpdated.CollectionID, &conn)
	if collectionOfChapterBeingDelete.UserID.String() != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this collection"})
		return
	}

	err = chapterSent.Update(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": chapterSent})
}

func ChaptersDelete(c *gin.Context)  {
	userID := c.GetString("user_id")
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	chapterSent := models.Chapter{}
	err := c.ShouldBindJSON(&chapterSent)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form sent"})
		return
	}

	chapterBeingDeleted, err := models.FindChapterById(chapterSent.ID, &conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	collectionOfChapterBeingDelete, err := models.FindCollectionById(chapterBeingDeleted.CollectionID, &conn)
	if collectionOfChapterBeingDelete.UserID.String() != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this collection"})
		return
	}

	err = chapterSent.Delete(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chapter": "chapter deleted"})
}