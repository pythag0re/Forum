package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "app.db")
	if err != nil {
		log.Fatal("Erreur de connexion à SQLite:", err)
	}
	fmt.Println("Connexion à SQLite réussie !")

	migrate()
}

func migrate() {
	query := `
-- users
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  email TEXT NOT NULL UNIQUE,
  pseudo TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  profile_picture TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- sessions
CREATE TABLE IF NOT EXISTS sessions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  token TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- posts
CREATE TABLE IF NOT EXISTS posts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME,
  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- comments
CREATE TABLE IF NOT EXISTS comments (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  post_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  content TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- likes
CREATE TABLE IF NOT EXISTS likes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  post_id INTEGER,
  comment_id INTEGER,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
  FOREIGN KEY(comment_id) REFERENCES comments(id) ON DELETE CASCADE
);
`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Erreur lors de la migration :", err)
	}
	fmt.Println("✅ Table 'Users' créée ou déjà existante.")
}

func RegisterUser(email, pseudo, password string) error {
	var exists int
	err := DB.QueryRow(`SELECT COUNT(*) FROM Users WHERE email = ? OR pseudo = ?`, email, pseudo).Scan(&exists)
	if err != nil {
		return err
	}
	if exists > 0 {
		return errors.New("Email ou pseudo déjà utilisé")
	}

	stmt, err := DB.Prepare("INSERT INTO Users (email, password, pseudo) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(email, password, pseudo)
	if err != nil {
		return err
	}

	fmt.Println("Utilisateur inscrit avec succès :", pseudo)
	return nil
}

func UpdateUserProfile(userID int, email, pseudo string) error {
	stmt, err := DB.Prepare(`
		UPDATE Users 
		SET email = ?, pseudo = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(email, pseudo, userID)
	return err
}

func DeleteUserByID(userID int) error {
	stmt, err := DB.Prepare("DELETE FROM Users WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID)
	return err
}

func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Println("Erreur lors de la fermeture de la base :", err)
		} else {
			fmt.Println("🔒 Connexion SQLite fermée.")
		}
	}
}
