package app

import (
	"log"
	"net/http"

	postgres "github.com/Metehan1994/HWs/HW4/common/db"
	"github.com/Metehan1994/HWs/HW4/domain/repos"
	"github.com/Metehan1994/HWs/HW4/reading_csv"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type App struct {
	Router     *mux.Router
	authorRepo *repos.AuthorRepository
	bookRepo   *repos.BookRepository
}

type Author struct {
	ID        uint     `json:"id"`
	Name      string   `json:"name"`
	CreatedAt string   `json:"createdAt"`
	DeletedAt string   `json:"deletedAt"`
	UpdatedAt string   `json:"updatedAt"`
	BooksName []string `json:"booksName"`
}

type Book struct {
	ID                uint   `json:"id"`
	Name              string `json:"name"`
	NumOfPages        int    `json:"numOfPages"`
	NumOfBooksInStock int    `json:"numOfBooksInStock"`
	Price             int    `json:"price"`
	StockCode         string `json:"stockCode"`
	ISBN              string `json:"isbn"`
	CreatedAt         string `json:"createdAt"`
	DeletedAt         string `json:"deletedAt"`
	UpdatedAt         string `json:"updatedAt"`
	AuthorID          uint   `json:"authorID"`
	AuthorName        string `json:"authorName"`
}

func (a *App) InitializeDBAndRepos(bookList reading_csv.BookList) error {
	//Set environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//Creating DB connection with postgres
	db, err := postgres.NewPsqlDB()
	if err != nil {
		log.Fatal("Postgres cannot init:", err)
	}

	//Autor Repo
	authorRepo := repos.NewAuthorRepository(db)
	authorRepo.Migrations()
	authorRepo.InsertSampleData(bookList)
	a.authorRepo = authorRepo

	//Book Repo
	bookRepo := repos.NewBookRepository(db)
	bookRepo.Migrations()
	bookRepo.InsertSampleData(bookList)
	a.bookRepo = bookRepo
	log.Println("Postgres connected & Repos created")
	return nil
}

func (a *App) InitializeRouter() {
	//Creating router with prefixes
	a.Router = mux.NewRouter()
	s := a.Router.PathPrefix("/authors").Subrouter()
	p := a.Router.PathPrefix("/books").Subrouter()
	//Providing endpoints to Author methods
	s.HandleFunc("/", a.GetAuthors).Methods(http.MethodGet)
	s.HandleFunc("/id/{id}", a.GetByAuthorID).Methods(http.MethodGet)
	s.HandleFunc("/name/{name}", a.GetAuthorByWord).Methods(http.MethodGet)
	s.HandleFunc("/", a.CreateAuthor).Methods(http.MethodPost)
	s.HandleFunc("/id/{id}", a.UpdateAuthor).Methods(http.MethodPatch)
	s.HandleFunc("/id/{id}", a.DeleteAuthorByID).Methods(http.MethodDelete)
	s.HandleFunc("/name/{name}", a.DeleteAuthorByName).Methods(http.MethodDelete)
	//Providing endpoints to Book methods
	p.HandleFunc("/", a.GetBooks).Methods(http.MethodGet)
	p.HandleFunc("/id/{id}", a.GetByBookID).Methods(http.MethodGet)
	p.HandleFunc("/name/{name}", a.GetBookByWord).Methods(http.MethodGet)
	p.HandleFunc("/", a.CreateBook).Methods(http.MethodPost)
	p.HandleFunc("/id/{id}", a.UpdateBook).Methods(http.MethodPatch)
	p.HandleFunc("/id/{id}", a.DeleteBookByID).Methods(http.MethodDelete)
	p.HandleFunc("/name/{name}", a.DeleteBookByName).Methods(http.MethodDelete)
	//p.HandleFunc("/maxPrice", a.MostExpensiveBook).Methods(http.MethodGet)
	//Creating a server URL
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}
