package sqlite

import (
	"database/sql"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// MIGRATIONS
func applyMigrations(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			version INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	migrationFiles, err := listMigrationFiles("database/migrations")
	if err != nil {
		return err
	}

	for _, migration := range migrationFiles {
		applied := isMigrationApplied(db, migration)
		if !applied {
			content, err := ioutil.ReadFile(filepath.Join("database/migrations", migration))
			if err != nil {
				return err
			}

			_, err = db.Exec(string(content))
			if err != nil {
				return err
			}

			_, err = db.Exec("INSERT INTO migrations (version) VALUES (?)", extractMigrationVersion(migration))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
func isMigrationApplied(db *sql.DB, migration string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE version = ?", extractMigrationVersion(migration)).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count > 0
}
func extractMigrationVersion(migration string) int {
	version, err := strconv.Atoi(strings.Split(migration, "_")[0])
	if err != nil {
		log.Fatal(err)
	}
	return version
}
func listMigrationFiles(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var migrationFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	sort.Strings(migrationFiles)

	return migrationFiles, nil
}
