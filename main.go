package main

import (
	"fmt"
	"net/http"
	"strings"

	"manga-app/models"
	"manga-app/routes"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
)

func main() {

	conn, err := connectDB()
	if err != nil {
		return
	}

	router := gin.Default()
	router.Use(CORS())
	router.Use(dbMiddleware(*conn))

	usersGroup := router.Group("users")
	{
		usersGroup.POST("register", routes.UsersRegister)
		usersGroup.POST("login", routes.UsersLogin)
	}
	collectionsGroup := router.Group("collections")
	{
		collectionsGroup.GET("index", routes.CollectionsIndex)
		collectionsGroup.POST("create", authMiddleWare(), routes.CollectionsCreate)
		collectionsGroup.PUT("update", authMiddleWare(), routes.CollectionsUpdate)
		collectionsGroup.DELETE("delete", authMiddleWare(), routes.CollectionsDelete)
		collectionsGroup.GET("get_collection_chapters", routes.CollectionsChapters)
	}
	chaptersGroup := router.Group("chapters")
	{
		chaptersGroup.POST("create", authMiddleWare(), routes.ChaptersCreate)
		chaptersGroup.PUT("update", authMiddleWare(), routes.ChaptersUpdate)
		chaptersGroup.DELETE("delete", authMiddleWare(), routes.ChaptersDelete)
		chaptersGroup.POST("get_chapter_pages", CORS(), routes.ChaptersPages)
	}

	pagesGroup := router.Group("pages")
	{
		pagesGroup.POST("create", authMiddleWare(), routes.PagesCreate)
	}
	router.Run(":8000")
}

func connectDB() (c *pgx.Conn, err error) {
	conn, err := pgx.Connect(context.Background(), "postgresql://postgres:Cooper4fox!@localhost:5432/manga-app")
	if err != nil {
		fmt.Println("Error connnecting to db", err.Error())
	}
	_ = conn.Ping(context.Background())
	return conn, err
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func dbMiddleware(conn pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", conn)
		c.Next()
	}
}

func authMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.Request.Header.Get("Authorization")
		split := strings.Split(bearer, "Bearer ")
		if len(split) < 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated."})
			c.Abort()
			return
		}
		token := split[1]

		isValid, userID := models.IsTokenValid(token)
		if isValid == false {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated."})
			c.Abort()
		} else {
			c.Set("user_id", userID)
			c.Next()
		}
	}
}
