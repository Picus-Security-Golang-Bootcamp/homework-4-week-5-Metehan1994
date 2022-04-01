package repos

import (
	"errors"
	"fmt"

	"github.com/Metehan1994/HWs/HW4/domain/entities"
	"github.com/Metehan1994/HWs/HW4/reading_csv"
	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

//NewBookRepository create a database for book
func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

//List gives output for books
func (b *BookRepository) List() []entities.Book {
	var books []entities.Book
	b.db.Find(&books)

	return books
}

//GetByID provides the book info for a given ID
func (b *BookRepository) GetByID(ID int) (*entities.Book, error) {
	var book entities.Book
	result := b.db.First(&book, ID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return &book, nil
}

//FindByWord lists the books with the given word case-insensitively
func (b *BookRepository) FindByWord(name string) {
	var books []entities.Book
	b.db.Where("name ILIKE ? ", "%"+name+"%").Find(&books)

	for _, author := range books {
		fmt.Println(author.ToString())
	}
}

//FindByName provides the book with the input of full name
func (b *BookRepository) FindByName(name string) {
	var book entities.Book
	b.db.Where("name = ? ", name).Find(&book)

	fmt.Println("found:", book.Name)
}

//Create creates a new book
func (b *BookRepository) Create(book entities.Book) error {
	result := b.db.Where("name = ?", book.Name).FirstOrCreate(&book)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//DeleteByName deletes the book with the given full name
func (b *BookRepository) DeleteByName(name string) error {
	var book entities.Book
	result := b.db.Unscoped().Where("name = ?", name).Find(&book)
	if result.Error != nil {
		return result.Error
	} else if book.Name != "" && !book.DeletedAt.Valid {
		fmt.Println("Valid book name, deleted:", name)
	} else if book.Name != "" && book.DeletedAt.Valid {
		fmt.Println("It has been already deleted.")
	} else {
		fmt.Println("Invalid book name, not deleted.")
	}
	result = b.db.Where("name = ?", name).Delete(&entities.Book{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

//DeleteByID applies a soft delete to a book with given ID
func (b *BookRepository) DeleteById(id int) error {
	var book entities.Book
	result := b.db.First(&book, id)
	if result.Error != nil {
		return result.Error
	} else {
		fmt.Println("Valid ID, deletion:", id)
	}
	result = b.db.Delete(&entities.Book{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

//Buy creates a purchase for a book with its ID in a given quantity
func (b *BookRepository) Buy(quantity, id int) error {
	var book entities.Book
	result := b.db.First(&book, id)
	if result.Error != nil {
		return result.Error
	} else if book.NumOfBooksInStock < quantity {
		return fmt.Errorf("not Enough Book. Books: %d", book.NumOfBooksInStock)
	} else {
		fmt.Println("It is successfully bought.")
	}

	result = b.db.Model(&book).Where("id = ? AND num_of_books_in_stock >= ?", id, quantity).
		Update("num_of_books_in_stock", gorm.Expr("num_of_books_in_stock - ?", quantity))
	if result.Error != nil {
		return result.Error
	}

	return nil
}

//MaxPrice finds the most expensive book
func (b *BookRepository) MaxPrice() error {
	var maxPrice int
	var book entities.Book
	err := b.db.Model(&book).Select("max(price)").Row().Scan(&maxPrice)
	if err != nil {
		fmt.Println("It could not be found.")
		return err
	}

	var books []entities.Book
	result := b.db.Where("price = ?", maxPrice).Find(&books)
	if result.Error != nil {
		return result.Error
	}
	for _, b := range books {
		fmt.Printf("Most expensive book is %s with %d TL.\n", b.Name, maxPrice)
	}

	return nil
}

//PriceBetweenFromLowerToUpper lists the books in a price range
func (b *BookRepository) PriceBetweenFromLowerToUpper(lower, upper int) error {
	var books []entities.Book

	result := b.db.Where("price > ? AND price < ?", lower, upper).Order("price").Find(&books)
	if result.Error != nil {
		return result.Error
	}
	for i, book := range books {
		fmt.Printf("Book %d: %s with %d TL\n", i+1, book.Name, book.Price)
	}

	return nil
}

//GetBooksWithBookInformation gives output of books with their author info
func (b *BookRepository) GetBooksWithAuthorInformation() ([]entities.Book, error) {
	var books []entities.Book
	result := b.db.Preload("Author").Find(&books)
	if result.Error != nil {
		return nil, result.Error
	}
	return books, nil
}

//Migrations form a book table in db
func (b *BookRepository) Migrations() {
	b.db.AutoMigrate(&entities.Book{})
}

//InsertSampleData creates a list of books
func (b *BookRepository) InsertSampleData(bookList reading_csv.BookList) {

	books := []entities.Book{}
	for _, book := range bookList {
		newBook := entities.Book{
			Name:              book.BookName,
			NumOfPages:        book.NumOfPages,
			NumOfBooksInStock: book.NumOfBooksinStock,
			Price:             book.Price,
			StockCode:         book.StockCode,
			ISBN:              book.ISBN,
			AuthorID:          uint(book.Author.AuthorID),
		}
		books = append(books, newBook)
	}

	for _, eachBook := range books {
		b.db.Unscoped().Where(entities.Book{Name: eachBook.Name}).FirstOrCreate(&eachBook)
	}

}
