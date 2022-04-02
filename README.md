# Homework-4 | Week 5 | Booklist App

## Overview

A program reading CSV file, connecting database on PostGreSQL and creating two tables on it which are "Author" and "Book" is introduced. And then, its connection with a server is established.

## How to Use the App ?

The program works with some server queries which can be reached from app.go.

* List names of book and author names
* Search books and authors with a word and IDs
* Create new books and authors with their info needed
* Delete books and authors from their tables with their names and IDs
* Buy books through book repo and update its quantity in stock
* Update them with new entries for the info and so on

### Some Notes for Usage

1. Program produces error messages when it is executed without considering its usage. Http status codes inform about the bad requests and info which is not found.

2. After deleting a book or an author, its status is changed and kept in the database (soft deleting).

### Postman Examples for Server Queries

* List books & authors:

```Postman
http://localhost:8080/books/ + Get Method

http://localhost:8080/authors/ + Get Method
```

* Search books & authors with an ID and a word:

```Postman
Books:
http://localhost:8080/books/id/3 + Get Method
http://localhost:8080/books/name/and + Get Method

Authors:
http://localhost:8080/authors/id/5 + Get Method
http://localhost:8080/authors/name/el + Get Method
```

* Create a book & an author:

```Postman
Book:
http://localhost:8080/books/ + Post Method
Body:
{
    "name":"book1",
    "authorid":3
}

Author:
http://localhost:8080/authors/ + Post Method

Body:
{
    "name":"author1",
    "ID":10
}
```

* Update book & author info

```Postman
Book:
http://localhost:8080/books/id/6 + Patch Method
Body:
{
    "name":"book2"
}

Author:
http://localhost:8080/authors/id/7 + Patch Method

Body:
{
    "name":"author2"
}
```

* Delete book & author with ID or name

```Postman
Book:
http://localhost:8080/books/id/6 + Delete Method
http://localhost:8080/books/name/War and Peace + Delete Method

Author:
http://localhost:8080/authors/id/4 + Delete Method
http://localhost:8080/authors/name/Fyodor Dostoyevski + Delete Method
```

* Buy book and list books in a price range

```Postman
Buy Book:
http://localhost:8080/books/buy/?id=5&quantity=10 + Patch Method

List books in a price range
http://localhost:8080/books/price/?lower=15&upper=25 + Get Method
```

## Packages Used

* The program is created with **GO main package & GORM & Godotenv & Gorilla Mux**.
