package routes

import (
	"fmt"
	"manga-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

// ChaptersCreate is the route for creating a chapter
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

// ChaptersUpdate is the route for updating a chapter
func ChaptersUpdate(c *gin.Context) {
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

	chapterBeingUpdated, err := models.FindChapterByID(chapterSent.ID, &conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collectionOfChapterBeingUpdated, err := models.FindCollectionById(chapterBeingUpdated.CollectionID, &conn)
	if collectionOfChapterBeingUpdated.UserID.String() != userID {
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

// ChaptersDelete is the route for handling deletion of a chapter
func ChaptersDelete(c *gin.Context) {
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

	chapterBeingDeleted, err := models.FindChapterByID(chapterSent.ID, &conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	AuthorizeAction(c, chapterBeingDeleted, userID, conn, "delete")

	err = chapterSent.Delete(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chapter": "chapter deleted"})
}

// AuthorizeAction checks if collection user id against user id provided to see if the action can be performed
func AuthorizeAction(c *gin.Context, chapter models.Chapter, userID string, conn pgx.Conn, action string) {
	chapterCollection, err := models.FindCollectionById(chapter.CollectionID, &conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if chapterCollection.UserID.String() != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to " + action + " this collection"})
		return
	}
	return
}

// ChaptersPages is the route for handling return of all the pages under this chapter
func ChaptersPages(c *gin.Context) {
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	chapter := models.Chapter{}
	fmt.Printf("%s", c.Request.Body)
	fmt.Println(c.GetRawData())
	c.ShouldBindJSON(&chapter)
	convertedUUID, err := uuid.FromString(c.Query("id"))
	chapter.ID = convertedUUID

	pages, err := chapter.GetAllPages(&conn)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{"pages": pages})
}
