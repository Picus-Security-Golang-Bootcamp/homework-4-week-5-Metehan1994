package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Metehan1994/HWs/HW4/domain/entities"
	"github.com/gorilla/mux"
)

//GetAuthors list all authors with their books and info together
func (a *App) GetAuthors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authors, err := a.authorRepo.GetAuthorsWithBookInformation()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(authors)
	}
}

//GetByAuthorID gets an author from the database with its ID if it is available
func (a *App) GetByAuthorID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	author, err := a.authorRepo.GetByID(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(author)
	}
}

//GetAuthorByWord lists the authors for the given word included in the their names
func (a *App) GetAuthorByWord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	authors := a.authorRepo.FindByWord(vars["name"])
	if len(authors) != 0 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(authors)
	} else {
		s := "No authors found.\n"
		fmt.Print(s)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(s))
	}
}

//CreateAuthor creates an author
func (a *App) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var author entities.Author
	json.NewDecoder(r.Body).Decode(&author)
	_, err := a.authorRepo.Create(author)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		s := "Author is created. Author Name: " + author.Name + "\n"
		w.Write([]byte(s))
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(author)
	}
}

//UpdateAuthor updates author info with patch method
func (a *App) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	author, err := a.authorRepo.GetByID(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		err := json.NewDecoder(r.Body).Decode(&author)
		a.authorRepo.Update(*author)
		if err != nil {
			fmt.Print(err)
			w.Write([]byte(err.Error()))
		} else {
			s := "Author is updated.\n"
			w.Write([]byte(s))
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(author)
		}
	}
}

//DeleteAuthorByID implements soft delete to an author for a given ID
func (a *App) DeleteAuthorByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	err := a.authorRepo.DeleteById(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusAccepted)
		s := fmt.Sprintf("Valid ID, deleted: %d", id)
		w.Write([]byte(s))
	}
}

//DeleteAuthorByName implements soft delete to an author for a given complete author name
func (a *App) DeleteAuthorByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	author, err := a.authorRepo.DeleteByName(vars["name"])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		s := "Valid author name, deleted: " + author.Name
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(s))
	}
}
