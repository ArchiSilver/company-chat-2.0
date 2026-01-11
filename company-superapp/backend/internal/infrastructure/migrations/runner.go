package migrations

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

// Runner выполняет миграции
type Runner struct {
	db             *sql.DB
	migrationsPath string
}

// NewRunner создаёт новый Runner миграций
func NewRunner(db *sql.DB, migrationsPath string) *Runner {
	return &Runner{
		db:             db,
		migrationsPath: migrationsPath,
	}
}

// Run выполняет все миграции
func (r *Runner) Run() error {
	// Создаём таблицу миграций если не существует
	if err := r.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Получаем выполненные миграции
	applied, err := r.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Получаем все файлы миграций
	files, err := r.getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Выполняем новые миграции
	for _, file := range files {
		if _, ok := applied[file]; ok {
			continue // Уже выполнена
		}

		if err := r.runMigration(file); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", file, err)
		}

		fmt.Printf("✓ Applied migration: %s\n", file)
	}

	return nil
}

func (r *Runner) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *Runner) getAppliedMigrations() (map[string]bool, error) {
	applied := make(map[string]bool)

	rows, err := r.db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

func (r *Runner) getMigrationFiles() ([]string, error) {
	pattern := filepath.Join(r.migrationsPath, "*.up.sql")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	// Извлекаем только имена файлов
	var names []string
	for _, f := range files {
		names = append(names, filepath.Base(f))
	}

	// Сортируем по имени (номеру миграции)
	sort.Strings(names)

	return names, nil
}

func (r *Runner) runMigration(filename string) error {
	// Читаем файл миграции
	path := filepath.Join(r.migrationsPath, filename)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// Начинаем транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Выполняем миграцию
	statements := strings.Split(string(content), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := tx.Exec(stmt); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	// Записываем в таблицу миграций
	if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", filename); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
