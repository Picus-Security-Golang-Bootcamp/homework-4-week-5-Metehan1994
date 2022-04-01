package entities

import (
	"fmt"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name              string  `json:"name"`
	NumOfPages        int     `json:"numOfPages"`
	NumOfBooksInStock int     `json:"numOfBooksInStock"`
	Price             int     `json:"price"`
	StockCode         string  `json:"stockCode"`
	ISBN              string  `json:"isbn"`
	AuthorID          uint    `json:"authorid"`
	Author            *Author `json:"author,omitempty" gorm:"foreignKey:AuthorID;references:ID"`
}

//TableName() returns the table header of book
func (Book) TableName() string {
	return "books"
}

//ToString() converts the data to readable info
func (b *Book) ToString() string {
	return fmt.Sprintf("ID : %d, Name : %s, Pages: %d, Price: %d, ISBN: %s, AuthorID: %d, CreatedAt : %s",
		b.ID, b.Name, b.NumOfPages, b.Price, b.ISBN, b.AuthorID, b.CreatedAt.Format("2006-01-02 15:04:05"))
}
