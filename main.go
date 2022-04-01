package main

import (
	"log"

	"github.com/Metehan1994/HWs/HW4/app"
	"github.com/Metehan1994/HWs/HW4/reading_csv"
)

var filename string = "book.csv"

func main() {
	//CSV to book struct
	bookList, err := reading_csv.ReadCSV(filename)
	if err != nil {
		log.Fatal(err)
	}

	//App
	app := &app.App{}
	app.InitializeDBAndRepos(bookList)
	app.InitializeRouter()
}
