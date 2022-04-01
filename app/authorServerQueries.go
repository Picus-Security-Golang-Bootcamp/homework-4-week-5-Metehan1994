package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Metehan1994/HWs/HW4/domain/entities"
	"github.com/gorilla/mux"
)

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

func (a *App) GetAuthorByWord(w http.ResponseWriter, r *http.Request) {
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

func (a *App) DeleteAuthorByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	a.authorRepo.DeleteById(id)
	var author entities.Author
	json.NewEncoder(w).Encode(author)
}

func (a *App) DeleteAuthorByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	author, err := a.authorRepo.DeleteByName(vars["name"])
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Valid author name, deleted:", author.Name)
	}
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
