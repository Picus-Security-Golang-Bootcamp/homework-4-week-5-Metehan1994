package entities

import (
	"fmt"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name              string
	NumOfPages        int
	NumOfBooksInStock int
	Price             int
	StockCode         string
	ISBN              string
	AuthorID          uint
	Author            Author `gorm:"foreignKey:AuthorID;references:ID"`
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
