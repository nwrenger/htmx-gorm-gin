package web

import (
	"fmt"
	"htmx-gorm-gin/db"
	"htmx-gorm-gin/web/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	query := c.PostForm("search")

	books, err := db.SearchBook(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "", template.Results(books))
}

func ShowAddBook(c *gin.Context) {
	c.HTML(http.StatusOK, "", template.AddBook(&db.Book{}))
}

func ShowEditBook(c *gin.Context) {
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

	c.HTML(http.StatusOK, "", template.EditBook(&book))
}

func AddBook(c *gin.Context) {
	var newBook db.Book

	if err := c.ShouldBind(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	book, err := db.AddBook(&newBook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Header("HX-Trigger", "update")
	c.HTML(http.StatusOK, "", template.EditBook(&book))
}

func EditBook(c *gin.Context) {
	var newBook db.Book

	if err := c.ShouldBind(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	book, err := db.UpdateBook(&newBook, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Header("HX-Trigger", "update")
	c.HTML(http.StatusOK, "", template.EditBook(&book))
}

func DeleteBook(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	_, err = db.DeleteBook(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Header("HX-Trigger", "update")
	c.HTML(http.StatusOK, "", template.NoBook())
}
