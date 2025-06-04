package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"forum/utils"
)

func main() {
	router := gin.Default()

	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router.POST("/like", func(c *gin.Context) {
		auth, userID := utils.IsAuthenticated(c.Request)
		if !auth {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var body struct {
			PostID int `json:"post_id"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		// Vérifie si le like existe
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id=? AND post_id=?", userID, body.PostID).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
			return
		}

		if count > 0 {
			// Unlike
			_, err := db.Exec("DELETE FROM likes WHERE user_id=? AND post_id=?", userID, body.PostID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"liked": false})
		} else {
			// Like
			_, err := db.Exec("INSERT INTO likes (user_id, post_id) VALUES (?, ?)", userID, body.PostID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"liked": true})
		}
	})

	router.Run(":8080")
}
