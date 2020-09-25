package routes

import (
	"bytes"
	"image/jpeg"
	"manga-app/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

// PagesCreate is the route for handling the creation of a page
func PagesCreate(c *gin.Context) {
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	file, header, err := c.Request.FormFile("page")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if header == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	decodedFile, err := jpeg.Decode(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, decodedFile, nil)
	imageByteArray := buf.Bytes()

	page := models.Page{}
	c.ShouldBindJSON(&page)

	page.Page = imageByteArray
	page.Title = c.Request.FormValue("title")
	orderNumber := c.Request.FormValue("order_number")
	page.OrderNumber, err = strconv.ParseInt(orderNumber, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	page.ChapterID, err = uuid.FromString(c.Request.FormValue("chapter_id"))
	pageCreateErr := page.Create(&conn)
	if pageCreateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, page)
}
