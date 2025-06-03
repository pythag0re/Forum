package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	router := gin.Default()

	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// post  /like et initialisation de la structure pour le json
	router.POST("/like", func(c *gin.Context) {
		var body struct {
			UserID    int  `json:"user_id"`
			PostID    *int `json:"post_id"`
			CommentID *int `json:"comment_id"`
		}

		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec("INSERT INTO likes (user_id, post_id, comment_id) VALUES (?, ?, ?)",
			body.UserID, body.PostID, body.CommentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Like ajouté"})
	})

	// pour supprimer
	router.DELETE("/like", func(c *gin.Context) {
		var body struct {
			UserID    int  `json:"user_id"`
			PostID    *int `json:"post_id"`
			CommentID *int `json:"comment_id"`
		}

		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec("DELETE FROM likes WHERE user_id = ? AND post_id IS ? AND comment_id IS ?",
			body.UserID, body.PostID, body.CommentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Like retiré"})
	})

	// avoir le nombre de like
	router.GET("/like/count", func(c *gin.Context) {
		postID := c.Query("post_id")
		commentID := c.Query("comment_id")

		row := db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id IS ? AND comment_id IS ?", nullable(postID), nullable(commentID))

		var count int
		if err := row.Scan(&count); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"count": count})
	})

	router.Run(":3001")
}

func nullable(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
