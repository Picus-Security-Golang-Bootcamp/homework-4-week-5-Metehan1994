package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Metehan1994/HWs/HW4/domain/entities"
	"github.com/gorilla/mux"
)

//GetBooks list all books with their authors and info together
func (a *App) GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	books, err := a.bookRepo.GetBooksWithAuthorInformation()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(books)
	}
}

//GetByBookID gets a book from the database with its ID if it is available
func (a *App) GetByBookID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	book, err := a.bookRepo.GetByID(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(book)
	}
}

//GetByBookByWord lists the books for the given word included in the their names
func (a *App) GetBookByWord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	books := a.bookRepo.FindByWord(vars["name"])
	if len(books) != 0 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(books)
	} else {
		w.WriteHeader(http.StatusNotFound)
		s := "No books found.\n"
		fmt.Print(s)
		w.Write([]byte(s))
	}

}

//UpdateAuthor updates book info with patch method
func (a *App) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	book, err := a.bookRepo.GetByID(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		err := json.NewDecoder(r.Body).Decode(&book)
		a.bookRepo.Update(*book)
		if err != nil {
			fmt.Print(err)
			w.Write([]byte(err.Error()))
		} else {
			s := "Book is updated.\n"
			w.Write([]byte(s))
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(book)
		}
	}

}

//CreateBook creates a book
func (a *App) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book entities.Book
	json.NewDecoder(r.Body).Decode(&book)
	_, err := a.bookRepo.Create(book)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		s := "Book is created. Book Name: " + book.Name + "\n"
		w.Write([]byte(s))
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(book)
	}
}

//DeleteBookByID implements soft delete to a book for a given ID
func (a *App) DeleteBookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	err := a.bookRepo.DeleteById(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusAccepted)
		s := "Book is successfully deleted.\n"
		w.Write([]byte(s))
	}
}

//DeleteAuthorByName implements soft delete to a book for a given complete book name
func (a *App) DeleteBookByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	book, err := a.bookRepo.DeleteByName(vars["name"])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		s := "Valid book name, deleted: " + book.Name
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(s))
	}
}

//MostExpensiveBook searches for most expensive book and declare it
func (a *App) MostExpensiveBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	books, err := a.bookRepo.MaxPrice()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
	} else {
		for _, b := range books {
			s := fmt.Sprintf("Most expensive book is %s with %d TL.\n", b.Name, b.Price)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(s))
			json.NewEncoder(w).Encode(b)
		}
	}
}

//PriceInRangeInIncreasingOrder list the books found in a range of prices from lowest to greatest
func (a *App) PriceInRangeInIncreasingOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	lower, _ := strconv.Atoi(vars["lower"])
	upper, _ := strconv.Atoi(vars["upper"])
	books, err := a.bookRepo.PriceBetweenFromLowerToUpper(lower, upper)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		for i, book := range books {
			s := fmt.Sprintf("Book %d: %s with %d TL\n", i+1, book.Name, book.Price)
			fmt.Print(s)
			w.Write([]byte(s))
		}
	}
}

//BuyBook creates a purchase for a book with its ID in a given quantity
func (a *App) BuyBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	quantity, _ := strconv.Atoi(vars["quantity"])
	bookBefore, err := a.bookRepo.Buy(quantity, id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		bookAfter, _ := a.bookRepo.GetByID(id)
		w.WriteHeader(http.StatusOK)
		s := "It is successfully bought.\n"
		w.Write([]byte(s))
		s = fmt.Sprintf("Before buying, number of books in stock: %d\n", bookBefore.NumOfBooksInStock)
		w.Write([]byte(s))
		json.NewEncoder(w).Encode(bookAfter)
	}
}
