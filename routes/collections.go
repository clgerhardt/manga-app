package routes

import (
	"fmt"
	"manga-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

// CollectionsIndex is the route for handling returning all the collections
func CollectionsIndex(c *gin.Context) {
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	collections, err := models.GetAllCollections(&conn)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{"collections": collections})
}

// CollectionsCreate is the route for handling creation of a collection
func CollectionsCreate(c *gin.Context) {
	userID := c.GetString("user_id")
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	collection := models.Collection{}
	c.ShouldBindJSON(&collection)
	err := collection.Create(&conn, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, collection)
}

// CollectionsUpdate is the route for handling the updating of a collection
func CollectionsUpdate(c *gin.Context) {
	userID := c.GetString("user_id")
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	collectionSent := models.Collection{}
	err := c.ShouldBindJSON(&collectionSent)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form sent"})
		return
	}

	collectionBeingUpdated, err := models.FindCollectionById(collectionSent.ID, &conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if collectionBeingUpdated.UserID.String() != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this collection"})
		return
	}

	collectionSent.UserID = collectionBeingUpdated.UserID
	err = collectionSent.Update(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": collectionSent})
}

// CollectionsDelete is the route for handling the deletion of a collection
func CollectionsDelete(c *gin.Context) {
	userID := c.GetString("user_id")
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	collectionSent := models.Collection{}
	err := c.ShouldBindJSON(&collectionSent)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form sent"})
		return
	}

	collectionBeingDeleted, err := models.FindCollectionById(collectionSent.ID, &conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if collectionBeingDeleted.UserID.String() != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this collection"})
		return
	}

	collectionSent.UserID = collectionBeingDeleted.UserID
	err = collectionSent.Delete(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"collection": "collection deleted"})
}

// CollectionsChapters is the route for handling return of all the chapters under this collection
func CollectionsChapters(c *gin.Context) {
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	collection := models.Collection{}
	c.ShouldBindJSON(&collection)
	chapters, err := collection.GetAllChapters(&conn)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{"chapters": chapters})
}
