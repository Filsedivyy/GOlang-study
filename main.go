package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"errors"
	"fmt"
	"strings"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "Honzíkova cesta", Author: "Pepek Námořník", Quantity: 1},
	{ID: "2", Title: "Jebingos", Author: "Jebač", Quantity: 6},
	{ID: "3", Title: "Hmm geysex", Author: "Mrdac", Quantity: 2},
	{ID: "4", Title: "Fortnite", Author: "Willem Dafoe", Quantity: 3},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)

}

func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}
	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("book not found")

}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {

		errorMessage := fmt.Sprintf("book with id %v does not exist", id)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": errorMessage})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

// funguje pouze na 1 slovo, je potřeba doopravit

func getBookByTitle(title string) (*book, error) {
	for t, b := range books {
		if strings.Contains(strings.ToLower(b.Title), strings.ToLower(title)) {
			return &books[t], nil
		}
	}
	return nil, errors.New("book not found")
}

func bookByTitle(c *gin.Context) {
	title := c.Param("title")
	book, err := getBookByTitle(title)

	if err != nil {
		errorMessage := fmt.Sprintf("title %v does not exist", title)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": errorMessage})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "missing id query parameter"})
		return
	}

	book, err := getBookById(id)
	if err != nil {
		errorMessage := fmt.Sprintf("book with id %v does not exist", id)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": errorMessage})
		return
	}
	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not available"})
		return
	}
	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)

}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/books", createBook)
	router.GET("/books/id/:id", bookById)
	router.GET("/books/title/:title", bookByTitle)
	router.PATCH("/checkout", checkoutBook)
	router.Run("localhost:8080")

}
