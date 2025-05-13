package server

import (
	"fmt"
	"forum/db"
	"html"
	"net/http"
	"text/template"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/home" {
		http.NotFound(w, r)
		fmt.Printf("Error: handler for %s not found\n", html.EscapeString(r.URL.Path))
		return
	}

	tmpl := template.Must(template.ParseFiles("_templates_/home.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Printf("erreur de template %s:", err)
	}
}

func Start() {
	db.InitDB()
	defer db.CloseDB()

	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("_templates_/css"))))
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("_templates_/"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	})

	http.HandleFunc("/home", homeHandler)

	fmt.Println("Serveur démarré sur le port 8080 ")
	http.ListenAndServe(":8080", nil)
}
