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
	books, err := a.bookRepo.GetBooksWithAuthorInformation()
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(err.Error()))
	} else {
		json.NewEncoder(w).Encode(books)
	}
	//bookApp := Book{}
	// for _, book := range books {
	// 	bookApp = CreateProperFormattedBook(book, bookApp)
	// 	json.NewEncoder(w).Encode(bookApp)
	// }
}

func (a *App) GetByBookID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	book, err := a.bookRepo.GetByID(id)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(err.Error()))
	} else {
		json.NewEncoder(w).Encode(book)
	}
	// bookApp := Book{}
	// bookApp = CreateProperFormattedBook(*book, bookApp)

}

func (a *App) GetBookByWord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	books := a.bookRepo.FindByWord(vars["name"])
	// bookApp := Book{}
	// for _, book := range books {
	// 	bookApp = CreateProperFormattedBook(book, bookApp)
	// 	json.NewEncoder(w).Encode(bookApp)
	// }
	if len(books) != 0 {
		json.NewEncoder(w).Encode(books)
	} else {
		s := "No books found.\n"
		fmt.Print(s)
		w.Write([]byte(s))
	}

}

func (a *App) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	book, err := a.bookRepo.GetByID(id)
	if err != nil {
		fmt.Println(err)
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
			json.NewEncoder(w).Encode(book)
		}
	}

}

func (a *App) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book entities.Book
	json.NewDecoder(r.Body).Decode(&book)
	_, err := a.bookRepo.Create(book)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(err.Error()))
	} else {
		s := "Book is created. Book Name: " + book.Name + "\n"
		w.Write([]byte(s))
	}
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
		w.Write([]byte(err.Error()))
	} else {
		s := "Valid book name, deleted: " + book.Name
		w.Write([]byte(s))
	}
}

func (a *App) MostExpensiveBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	books, err := a.bookRepo.MaxPrice()
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(err.Error()))
	} else {
		for _, b := range books {
			s := fmt.Sprintf("Most expensive book is %s with %d TL.\n", b.Name, b.Price)
			w.Write([]byte(s))
			json.NewEncoder(w).Encode(b)
		}
	}
}

func (a *App) PriceInRangeInIncreasingOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	lower, _ := strconv.Atoi(vars["lower"])
	upper, _ := strconv.Atoi(vars["upper"])
	books, err := a.bookRepo.PriceBetweenFromLowerToUpper(lower, upper)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(err.Error()))
	} else {
		for i, book := range books {
			s := fmt.Sprintf("Book %d: %s with %d TL\n", i+1, book.Name, book.Price)
			fmt.Print(s)
			w.Write([]byte(s))
		}
	}
}

func (a *App) buyBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	quantity, _ := strconv.Atoi(vars["quantity"])
	book, err := a.bookRepo.Buy(quantity, id)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(err.Error()))
	} else {
		s := "It is successfully bought.\n"
		w.Write([]byte(s))
		json.NewEncoder(w).Encode(book)
	}
}

// func CreateProperFormattedBook(b entities.Book, bookApp Book) Book {
// 	bookApp.Name = b.Name
// 	bookApp.ID = b.ID
// 	bookApp.NumOfPages = b.NumOfPages
// 	bookApp.NumOfBooksInStock = b.NumOfBooksInStock
// 	bookApp.Price = b.Price
// 	bookApp.StockCode = b.StockCode
// 	bookApp.ISBN = b.ISBN
// 	bookApp.CreatedAt = b.CreatedAt.String()
// 	bookApp.UpdatedAt = b.UpdatedAt.String()
// 	bookApp.DeletedAt = b.DeletedAt.Time.String()
// 	bookApp.AuthorID = b.Author.ID
// 	bookApp.AuthorName = b.Author.Name

// 	return bookApp
// }
