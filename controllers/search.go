package controllers

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

func getPost(title string) (*sql.Rows, error)  {
	db, err := sql.Open("sqlite3", "../app.db")
	if err != nil {
		fmt.Println("la data base n'a pas fonctionné")
	} else {
		rows, err := db.Query("SELECT id, name FROM post WHERE title LIKE ?", "%"+title+"%")
		if err != nil {
			return nil, err
		}
		return rows, nil
	}
	return nil, err
} 