package main

import (
	"fmt"
	"net/http"
	"strconv"

	"web-service-gin/db"

	"github.com/gin-gonic/gin"
)

// Define access checkers
func AlbumReadOnly(login db.Login) bool {
	return login.Permission.AccessAlbum == db.ReadOnly || login.Permission.AccessAlbum == db.Write
}

func AlbumWrite(login db.Login) bool {
	return login.Permission.AccessAlbum == db.Write
}

func LoginReadOnly(login db.Login) bool {
	return login.Permission.AccessAdmin == db.ReadOnly || login.Permission.AccessAdmin == db.Write
}

func LoginWrite(login db.Login) bool {
	return login.Permission.AccessAdmin == db.Write
}

// Define auth middleware
func AuthMiddleware(checker func(db.Login) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// send auth header
		c.Header("WWW-Authenticate", "Basic realm=\"Restricted\"")

		username, password, ok := c.Request.BasicAuth()
		fmt.Println(username, password)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Lookup user in your database
		user, found := db.FetchLogin(username)
		if found != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Compare the password
		// err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		err := user.Password != password
		if err {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Check user's permission
		if !checker(user) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}

		// Set the user in the context
		c.Set("user", user)

		c.Next()
	}
}

func main() {
	err := db.Open()

	if err != nil {
		println("Error:", err)
	}

	// adding admin, ignore error
	db.AddOrUpdateLogin(&db.Login{
		User:     "admin",
		Password: "1",
		Permission: db.Permissions{
			AccessAdmin: db.Write,
			AccessAlbum: db.Write,
		},
	})

	router := gin.Default()

	router.GET("/albums", AuthMiddleware(AlbumReadOnly), getAlbums)
	router.GET("/album/:id", AuthMiddleware(AlbumReadOnly), getAlbumById)
	router.POST("/album", AuthMiddleware(AlbumWrite), postAlbum)
	router.DELETE("/album/:id", AuthMiddleware(AlbumWrite), deleteAlbum)

	router.POST("/login", AuthMiddleware(LoginWrite), postLogin)

	router.Run("localhost:8080")
}

func getAlbums(c *gin.Context) {
	albums, err := db.FetchAlbums()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, albums)
}

func getAlbumById(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	album, err := db.FetchAlbumById(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, album)
}

func postAlbum(c *gin.Context) {
	var newAlbum db.Album

	if err := c.BindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, newAlbum)
		return
	}

	if err := db.AddOrUpdateAlbum(&newAlbum); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newAlbum)
}

func deleteAlbum(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	album, err := db.DeleteAlbum(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, album)
}

func postLogin(c *gin.Context) {
	var newLogin db.Login

	if err := c.BindJSON(&newLogin); err != nil {
		c.JSON(http.StatusBadRequest, newLogin)
		return
	}

	if err := db.AddOrUpdateLogin(&newLogin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newLogin)
}
