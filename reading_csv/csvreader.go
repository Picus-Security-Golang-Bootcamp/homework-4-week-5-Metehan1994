package reading_csv

import (
	"encoding/csv"
	"os"
	"strconv"
)

//ReadCSV reads csv and returns a book list
func ReadCSV(filename string) (BookList, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var bookList BookList
	for _, line := range records[1:] {
		book := BookInfo{}
		book.BookName = line[0]
		book.NumOfPages, _ = strconv.Atoi(line[1])
		book.NumOfBooksinStock, _ = strconv.Atoi(line[2])
		book.Price, _ = strconv.Atoi(line[3])
		book.StockCode = line[4]
		book.ISBN = line[5]
		book.Author.AuthorID, _ = strconv.Atoi(line[6])
		book.Author.AuthorName = line[7]
		bookList = append(bookList, book)
	}

	return bookList, nil
}
