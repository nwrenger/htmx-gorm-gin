package api

import (
	"fmt"
	"htmx-gorm-gin/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetBooks(c *gin.Context) {
	books, err := db.FetchBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}

func GetBookById(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	book, err := db.FetchBookById(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

func PostBook(c *gin.Context) {
	var newBook db.Book

	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, newBook)
		return
	}

	if _, err := db.AddBook(&newBook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newBook)
}

func DeleteBook(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	book, err := db.DeleteBook(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

func PostLogin(c *gin.Context) {
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
