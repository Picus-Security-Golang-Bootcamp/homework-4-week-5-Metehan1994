package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Metehan1994/HWs/HW4/domain/entities"
	"github.com/gorilla/mux"
)

func (a *App) GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	books, _ := a.bookRepo.GetBooksWithAuthorInformation()
	bookApp := Book{}
	for _, book := range books {
		bookApp = CreateProperFormattedBook(book, bookApp)
		json.NewEncoder(w).Encode(bookApp)
	}
}

func (a *App) GetByBookID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	book, _ := a.bookRepo.GetByID(id)
	bookApp := Book{}
	bookApp = CreateProperFormattedBook(*book, bookApp)
	json.NewEncoder(w).Encode(bookApp)
}

func (a *App) GetBookByWord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	books := a.bookRepo.FindByWord(vars["name"])
	bookApp := Book{}
	for _, book := range books {
		bookApp = CreateProperFormattedBook(book, bookApp)
		json.NewEncoder(w).Encode(bookApp)
	}
}

func (a *App) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	book, _ := a.bookRepo.GetByID(id)
	json.NewDecoder(r.Body).Decode(&book)
	a.bookRepo.Update(*book)
	json.NewEncoder(w).Encode(book)
}

func (a *App) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book entities.Book
	json.NewDecoder(r.Body).Decode(&book)
	a.bookRepo.Create(book)
	json.NewEncoder(w).Encode(book)
}

func (a *App) DeleteBookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	a.bookRepo.DeleteById(id)
	var book entities.Book
	json.NewEncoder(w).Encode(book)
}

func (a *App) DeleteBookByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	book, err := a.bookRepo.DeleteByName(vars["name"])
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Valid book name, deleted:", book.Name)
	}
	json.NewEncoder(w).Encode(book)
}

func CreateProperFormattedBook(b entities.Book, bookApp Book) Book {
	bookApp.Name = b.Name
	bookApp.ID = b.ID
	bookApp.NumOfPages = b.NumOfPages
	bookApp.NumOfBooksInStock = b.NumOfBooksInStock
	bookApp.Price = b.Price
	bookApp.StockCode = b.StockCode
	bookApp.ISBN = b.ISBN
	bookApp.CreatedAt = b.CreatedAt.String()
	bookApp.UpdatedAt = b.UpdatedAt.String()
	bookApp.DeletedAt = b.DeletedAt.Time.String()
	bookApp.AuthorID = b.Author.ID
	bookApp.AuthorName = b.Author.Name

	return bookApp
}
