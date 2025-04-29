package db

import (
	"database/sql"
	"fmt"
	"log"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "app.db")
	if err != nil {
		log.Fatal("Erreur de connexion à SQLite:", err)
	}
	fmt.Println("Connexion à SQLite réussie!")

	migrate()
}

func migrate() {
	query := ``

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Erreur lors de la migration:", err)
	}
	fmt.Println("Tables créées avec succès!")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		fmt.Println("Connexion SQLite fermée.")
	}
}
