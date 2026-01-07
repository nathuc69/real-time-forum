package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func RunMigrations(db *sql.DB, dir string) error {
	// Créer une table de suivi des migrations, pour ne pas toutes les faire à chaque fois
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id TEXT PRIMARY KEY);`)
	if err != nil {
		return fmt.Errorf("❌ unable to create 'schema_migrations' table: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return err
	}

	for _, file := range files {
		/*log.Println("updating database...:", file)*/
		filename := filepath.Base(file)
		id := strings.Split(filename, "_")[0] // "001" dans "001_create_tables.sql"

		// Vérifier si la migration a déjà été exécutée
		var exists string
		err = db.QueryRow(`SELECT id FROM schema_migrations WHERE id = ?`, id).Scan(&exists)
		if err == sql.ErrNoRows {
			log.Println("Applying migration:", filename)

			content, err := os.ReadFile(file)
			if err != nil {
				return err
			}

			// exécuter le SQL:
			_, err = db.Exec(string(content))
			if err != nil {
				return err
			}

			// enregistrer la migration comme appliquée
			_, err = db.Exec(`INSERT INTO schema_migrations (id) VALUES (?)`, id)
			if err != nil {
				return fmt.Errorf("❌ error recording migration %s: %w", filename, err)
			}

			log.Println("✅ migration applied:", filename)
		} else if err != nil {
			return fmt.Errorf("checking migration status for %s failed: %w", filename, err)
		} else {
			log.Println("migration already applied:", filename)
		}

	}

	return nil
}
