package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	postgres "github.com/Metehan1994/HWs/HW4/common/db"
	"github.com/Metehan1994/HWs/HW4/domain/entities"
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
	a.Router = mux.NewRouter()
	s := a.Router.PathPrefix("/authors").Subrouter()
	s.HandleFunc("/", a.GetAuthors).Methods(http.MethodGet)
	s.HandleFunc("/id/{id}", a.GetByAuthorID).Methods(http.MethodGet)
	s.HandleFunc("/name/{name}", a.GetByAuthorName).Methods(http.MethodGet)
	s.HandleFunc("/", a.CreateAuthor).Methods(http.MethodPost)
	//s.HandleFunc("/id/{id}", a.UpdateAuthor).Methods(http.MethodPut)
	log.Fatal(http.ListenAndServe(":8080", s))
}

func (a *App) GetAuthors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authors, _ := a.authorRepo.GetAuthorsWithBookInformation()
	authorApp := Author{}
	for _, author := range authors {
		authorApp = CreateProperFormattedAuthor(author, authorApp)
		json.NewEncoder(w).Encode(authorApp)
	}
}

func (a *App) GetByAuthorID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	author, _ := a.authorRepo.GetByID(id)
	authorApp := Author{}
	authorApp = CreateProperFormattedAuthor(*author, authorApp)
	json.NewEncoder(w).Encode(authorApp)
}

func (a *App) GetByAuthorName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	authors := a.authorRepo.FindByWord(vars["name"])
	authorApp := Author{}
	for _, author := range authors {
		authorApp = CreateProperFormattedAuthor(author, authorApp)
		json.NewEncoder(w).Encode(authorApp)
	}
}

func (a *App) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var author entities.Author
	json.NewDecoder(r.Body).Decode(&author)
	a.authorRepo.Create(author)
	json.NewEncoder(w).Encode(author)
}

func (a *App) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	author, _ := a.authorRepo.GetByID(id)
	json.NewDecoder(r.Body).Decode(&author)
	a.authorRepo.Update(*author)
	json.NewEncoder(w).Encode(author)
}

func CreateProperFormattedAuthor(a entities.Author, authorApp Author) Author {
	authorApp.Name = a.Name
	authorApp.ID = a.ID
	authorApp.CreatedAt = a.CreatedAt.String()
	authorApp.UpdatedAt = a.UpdatedAt.String()
	authorApp.DeletedAt = a.DeletedAt.Time.String()
	authorApp.BooksName = make([]string, 0)
	for _, book := range a.Book {
		authorApp.BooksName = append(authorApp.BooksName, book.Name)
	}
	return authorApp
}
