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
	result := b.db.Preload("Author").First(&book, ID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return &book, nil
}

//FindByWord lists the books with the given word case-insensitively
func (b *BookRepository) FindByWord(name string) []entities.Book {
	var books []entities.Book
	b.db.Where("name ILIKE ? ", "%"+name+"%").Find(&books)

	return books
}

//FindByName provides the book with the input of full name
func (b *BookRepository) FindByName(name string) {
	var book entities.Book
	b.db.Where("name = ? ", name).Find(&book)

	fmt.Println("found:", book.Name)
}

//Create creates a new book
func (b *BookRepository) Create(book entities.Book) (*entities.Book, error) {
	result := b.db.Where("name = ?", book.Name).FirstOrCreate(&book)
	if result.Error != nil {
		return nil, result.Error
	}
	return &book, nil
}

func (b *BookRepository) Update(book entities.Book) error {
	result := b.db.Save(&book)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

//DeleteByName deletes the book with the given full name
func (b *BookRepository) DeleteByName(name string) (*entities.Book, error) {
	var book entities.Book
	result := b.db.Unscoped().Where("name = ?", name).Find(&book)
	if result.Error != nil {
		return nil, result.Error
	} else if book.Name != "" && !book.DeletedAt.Valid {
		result = b.db.Where("name = ?", name).Delete(&entities.Book{})
		if result.Error != nil {
			return nil, result.Error
		} else {
			return &book, nil
		}
	} else if book.Name != "" && book.DeletedAt.Valid {
		return nil, errors.New("it has been already deleted")
	} else {
		return nil, errors.New("invalid book name, no deletion")
	}
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
func (b *BookRepository) Buy(quantity, id int) (*entities.Book, error) {
	var book entities.Book
	result := b.db.First(&book, id)
	if result.Error != nil {
		return nil, result.Error
	} else if book.NumOfBooksInStock < quantity && book.Name != "" {
		return nil, fmt.Errorf("not Enough Book. NumofBooks: %d", book.NumOfBooksInStock)
	} else if book.Name == "" {
		return nil, fmt.Errorf("invalid ID")
	}

	result = b.db.Model(&book).Where("id = ? AND num_of_books_in_stock >= ?", id, quantity).
		Update("num_of_books_in_stock", gorm.Expr("num_of_books_in_stock - ?", quantity))
	if result.Error != nil {
		return nil, result.Error
	}

	return &book, nil
}

//MaxPrice finds the most expensive book
func (b *BookRepository) MaxPrice() ([]entities.Book, error) {
	var maxPrice int
	var book entities.Book
	err := b.db.Model(&book).Select("max(price)").Row().Scan(&maxPrice)
	if err != nil {
		return nil, errors.New("it could not be found")
	}

	var books []entities.Book
	result := b.db.Where("price = ?", maxPrice).Find(&books)
	if result.Error != nil {
		return nil, result.Error
	}

	return books, nil
}

//PriceBetweenFromLowerToUpper lists the books in a price range
func (b *BookRepository) PriceBetweenFromLowerToUpper(lower, upper int) ([]entities.Book, error) {
	var books []entities.Book

	result := b.db.Where("price > ? AND price < ?", lower, upper).Order("price").Find(&books)
	if result.Error != nil {
		return nil, result.Error
	}
	return books, nil
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
