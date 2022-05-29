package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Post struct {
	Id    int
	Title string
	Body  string
}

var db, err = sql.Open("mysql", "root:my-secret-pw@(172.17.0.2:3306)/go_course?charset=utf8")

func main() {

	// stmt, err := db.Prepare("Insert into posts(title, body) values(?,?)")
	// checkError(err)

	// _, err = stmt.Exec("My first Post", "My first content")
	// checkError(err)

	rows, err := db.Query("Select * from posts")
	checkError(err)

	items := []Post{}

	for rows.Next() {
		post := Post{}
		rows.Scan(&post.Id, &post.Title, &post.Body)
		items = append(items, post)
	}

	db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		t := template.Must(template.ParseFiles("templates/index.html"))
		if err := t.ExecuteTemplate(w, "index.html", items); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println(http.ListenAndServe(":8080", nil))
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
