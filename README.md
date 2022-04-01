# Homework-3 | Week 4 | Booklist App

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

### Postman Usage with Server Queries

## Package Used

* The program is created with **GO main package & GORM & Godotenv**.
