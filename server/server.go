package server

import (
	"encoding/json"
	"fmt"
	"forum/controllers"
	"forum/db"
	"forum/utils"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"text/template"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.NotFound(w, r)
		fmt.Printf("Error: handler for %s not found\n", html.EscapeString(r.URL.Path))
		return
	}

	switch r.Method {
	case http.MethodGet:
		tmpl := template.Must(template.ParseFiles("_templates_/login.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			fmt.Printf("erreur de template : %s\n", err)
		}
	case http.MethodPost:
		controllers.LoginUser(w, r)
	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		http.NotFound(w, r)
		fmt.Printf("Error: handler for %s not found\n", html.EscapeString(r.URL.Path))
		return
	}

	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("_templates_/register.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			fmt.Printf("erreur de template %s:", err)
		}
	} else if r.Method == http.MethodPost {
		controllers.RegisterUser(w, r)
	}
}

func landingHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/landing" {
		http.NotFound(w, r)
		fmt.Printf("Error: handler for %s not found\n", html.EscapeString(r.URL.Path))
		return
	}

	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("_templates_/landing.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			fmt.Printf("erreur de template %s:", err)
		}
	}
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("_templates_/create_post.html"))
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")

		ok, userID := utils.IsAuthenticated(r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		_, err = db.DB.Exec(`
			INSERT INTO posts (title, content, user_id, created_at) 
			VALUES (?, ?, ?, datetime('now'))`, title, content, userID)

		if err != nil {
			log.Println("Erreur d'insertion du post :", err)
			http.Error(w, "Error creating post", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/landing", http.StatusSeeOther)
	}
}

func likeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// ✅ Vérifie si l'utilisateur est authentifié
	auth, userID := utils.IsAuthenticated(r)
	if !auth {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	// ✅ Ne pas accepter le user_id depuis le client
	var req struct {
		PostID int `json:"post_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	var count int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = ? AND post_id = ?", userID, req.PostID).Scan(&count)
	if err != nil {
		log.Println(err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		_, err = db.DB.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userID, req.PostID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(struct {
			Liked bool `json:"liked"`
		}{Liked: false})
	} else {
		_, err = db.DB.Exec("INSERT INTO likes (user_id, post_id) VALUES (?, ?)", userID, req.PostID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(struct {
			Liked bool `json:"liked"`
		}{Liked: true})
	}
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	ok, userID := utils.IsAuthenticated(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseMultipartForm(10 << 20)

		action := r.FormValue("action")

		switch action {
		case "change_photo":
			file, handler, err := r.FormFile("profilePic")
			if err == nil {
				defer file.Close()
				filename := fmt.Sprintf("user%d_%s", userID, handler.Filename)
				filePath := path.Join("_templates_/uploads", filename)

				dst, err := os.Create(filePath)
				if err == nil {
					defer dst.Close()
					io.Copy(dst, file)

					_, err = db.DB.Exec("UPDATE users SET profile_picture = ? WHERE id = ?", filename, userID)
					if err != nil {
						log.Println("Erreur update photo:", err)
					}
				}
			}

		case "change_username":
			newPseudo := r.FormValue("username")
			if newPseudo != "" {
				_, err := db.DB.Exec("UPDATE users SET pseudo = ? WHERE id = ?", newPseudo, userID)
				if err != nil {
					log.Println("Erreur update pseudo:", err)
				}
			}

		case "change_email":
			newEmail := r.FormValue("email")
			if newEmail != "" {
				_, err := db.DB.Exec("UPDATE users SET email = ? WHERE id = ?", newEmail, userID)
				if err != nil {
					log.Println("Erreur update email:", err)
				}
			}
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}
	var pseudo, email, profilePic, createdAt string
	err := db.DB.QueryRow(`SELECT pseudo, email, profile_picture, created_at FROM users WHERE id = ?`, userID).
		Scan(&pseudo, &email, &profilePic, &createdAt)

	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	if profilePic == "" {
		profilePic = "avatar.png"
	}

	data := struct {
		Username    string
		Email       string
		ProfilePic  string
		MemberSince string
	}{
		Username:    pseudo,
		Email:       email,
		ProfilePic:  profilePic,
		MemberSince: createdAt,
	}

	tmpl := template.Must(template.ParseFiles("_templates_/prrofile.html"))
	tmpl.Execute(w, data)
}

func Start() {
	db.InitDB()
	defer db.CloseDB()

	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("_templates_/css"))))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("_templates_/uploads"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("likeAPI"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/landing", http.StatusSeeOther)
	})

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/landing", landingHandler)
	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/logout", controllers.LogoutHandler)
	http.HandleFunc("/delete-profile", controllers.DeleteProfileHandler)
	http.HandleFunc("/change-password", controllers.ChangePasswordHandler)
	http.HandleFunc("/create-post", createPostHandler)
	http.HandleFunc("/posts", controllers.PostsHandler)
	http.HandleFunc("/post", controllers.PostDetailHandler)
	http.HandleFunc("/user/", controllers.PublicProfileHandler)
	http.HandleFunc("/my-posts", controllers.MyPostsHandler)
	http.HandleFunc("/delete-post", controllers.DeletePostHandler)
	http.HandleFunc("/edit-post", controllers.EditPostHandler)
	http.HandleFunc("/update-post", controllers.UpdatePostHandler)
	http.HandleFunc("/add-comment", controllers.AddCommentHandler)
	http.HandleFunc("/my-comments", controllers.MyCommentsHandler)
	http.HandleFunc("/edit-comment", controllers.EditCommentHandler)
	http.HandleFunc("/delete-comment", controllers.DeleteCommentHandler)
	http.HandleFunc("/update-comment", controllers.UpdateCommentHandler)
	http.HandleFunc("/user-posts", controllers.UserPostsHandler)
	http.HandleFunc("/user-comments", controllers.UserCommentsHandler)
	http.HandleFunc("/user-likes", controllers.UserLikedPostsHandler)
	http.HandleFunc("/like", likeHandler)

	fmt.Println("Serveur démarré sur le port 8080 ")
	http.ListenAndServe(":8080", nil)
}
