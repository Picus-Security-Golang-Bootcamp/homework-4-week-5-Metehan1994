package entities

import (
	"fmt"

	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"unique"`
	Book []Book `json:"book,omitempty"`
}

//TableName() returns the table header of author
func (Author) TableName() string {
	return "authors"
}

//ToString() converts the data to readable info
func (a *Author) ToString() string {
	return fmt.Sprintf("ID : %d, Name : %s, CreatedAt : %s", a.ID, a.Name, a.CreatedAt.Format("2006-01-02 15:04:05"))
}
