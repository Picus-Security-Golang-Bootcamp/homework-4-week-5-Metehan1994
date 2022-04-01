package reading_csv

type BookInfo struct {
	BookName          string
	NumOfPages        int
	NumOfBooksinStock int
	Price             int
	StockCode         string
	ISBN              string
	Author            struct {
		AuthorID   int
		AuthorName string
	}
}

type BookList []BookInfo
