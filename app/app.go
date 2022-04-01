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

	//Providing endpoints to Author methods
	s := a.Router.PathPrefix("/authors").Subrouter()
	s.HandleFunc("/", a.GetAuthors).Methods(http.MethodGet)
	s.HandleFunc("/id/{id}", a.GetByAuthorID).Methods(http.MethodGet)
	s.HandleFunc("/name/{name}", a.GetAuthorByWord).Methods(http.MethodGet)
	s.HandleFunc("/", a.CreateAuthor).Methods(http.MethodPost)
	s.HandleFunc("/id/{id}", a.UpdateAuthor).Methods(http.MethodPatch)
	s.HandleFunc("/id/{id}", a.DeleteAuthorByID).Methods(http.MethodDelete)
	s.HandleFunc("/name/{name}", a.DeleteAuthorByName).Methods(http.MethodDelete)

	//Providing endpoints to Book methods
	p := a.Router.PathPrefix("/books").Subrouter()
	p.HandleFunc("/", a.GetBooks).Methods(http.MethodGet)
	p.HandleFunc("/id/{id}", a.GetByBookID).Methods(http.MethodGet)
	p.HandleFunc("/name/{name}", a.GetBookByWord).Methods(http.MethodGet)
	p.HandleFunc("/", a.CreateBook).Methods(http.MethodPost)
	p.HandleFunc("/id/{id}", a.UpdateBook).Methods(http.MethodPatch)
	p.HandleFunc("/id/{id}", a.DeleteBookByID).Methods(http.MethodDelete)
	p.HandleFunc("/name/{name}", a.DeleteBookByName).Methods(http.MethodDelete)
	p.HandleFunc("/maxprice", a.MostExpensiveBook).Methods(http.MethodGet)
	p.HandleFunc("/price/", a.PriceInRangeInIncreasingOrder).Methods(http.MethodGet).Queries("lower", "{lower}", "upper", "{upper}")
	p.HandleFunc("/buy/", a.BuyBook).Methods(http.MethodPatch).Queries("id", "{id}", "quantity", "{quantity}")

	//Creating a server URL
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}
