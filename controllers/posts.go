package controllers

import (
	"forum/db"
	"forum/utils"
	"html/template"
	"log"
	"net/http"
)

type PostData struct {
	ID        int
	Title     string
	Content   string
	Author    string
	UserID    int    
	CreatedAt string
	Likes     int
	Comments  int
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	ok, _ := utils.IsAuthenticated(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rows, err := db.DB.Query(`
		SELECT posts.id, posts.title, posts.content, users.pseudo, users.id, posts.created_at,
			(SELECT COUNT(*) FROM likes WHERE likes.post_id = posts.id) as like_count,
			(SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) as comment_count
		FROM posts
		JOIN users ON posts.user_id = users.id
		ORDER BY posts.created_at DESC`)
	if err != nil {
		log.Println("Erreur lors de la récupération des posts:", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []PostData
	for rows.Next() {
		var p PostData
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.UserID, &p.CreatedAt, &p.Likes, &p.Comments)
		if err != nil {
			log.Println("Erreur de scan post:", err)
			continue
		}
		posts = append(posts, p)
	}

	tmpl := template.Must(template.ParseFiles("_templates_/posts.html"))
	err = tmpl.Execute(w, posts)
	if err != nil {
		log.Println("Erreur template posts:", err)
	}
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	ok, userID := utils.IsAuthenticated(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postID := r.FormValue("post_id")
	if postID == "" {
		http.Error(w, "ID de post manquant", http.StatusBadRequest)
		return
	}
	_, err := db.DB.Exec("DELETE FROM posts WHERE id = ? AND user_id = ?", postID, userID)
	if err != nil {
		log.Println("Erreur lors de la suppression du post :", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/my-posts", http.StatusSeeOther)
}

