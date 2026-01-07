package config

import (
	"database/sql"

	//"forum/internal/repositories"
	"os"

	"log"
	"real-time-forum/backend/repositories"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	//connexion à la BdD:
	dbPath := os.Getenv("FORUM_DB_PATH")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("❌ error opening database:", err)
	}
	//vérification de la connexion:
	if err = db.Ping(); err != nil {
		log.Fatal("❌ error connecting to database:", err)
	}
	//exécution des migrations (ie modifications structurelles de la BdD):
	if err := repositories.RunMigrations(db, "backend/migrations"); err != nil {
		log.Fatal("❌ error migrating:", err)
	}
	// if err := repositories.RunMigrations(db, "migrations"); err != nil {
	// 	log.Fatal("❌ error migrating:", err)
	// }
	return db
}
