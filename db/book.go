package db

type Book struct {
	ID     uint    `gorm:"primaryKey" json:"id" form:"id"`
	Title  string  `json:"title" form:"title"`
	Author string  `json:"author" form:"author"`
	Price  float64 `json:"price" form:"price"`
}

func FetchBooks() ([]Book, error) {
	var books []Book
	rs := db.Find(&books)
	return books, rs.Error
}

func FetchBookById(id uint) (Book, error) {
	var book Book
	rs := db.Where("id = ?", id).Find(&book)
	return book, rs.Error
}

func SearchBook(query string) ([]Book, error) {
	var books []Book
	rs := db.Where("title LIKE ? OR author LIKE ? OR id LIKE ? OR price LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%").Order("title ASC").Find(&books)
	return books, rs.Error
}

func AddBook(newBook *Book) (Book, error) {
	rs := db.Create(newBook)
	return *newBook, rs.Error
}

func UpdateBook(newBook *Book, id int) (Book, error) {
	rs := db.Where("id = ?", id).Updates(&newBook)
	return *newBook, rs.Error
}

func DeleteBook(id uint) (Book, error) {
	var book Book
	rs := db.Where("id = ?", id).Delete(&book)
	return book, rs.Error
}
