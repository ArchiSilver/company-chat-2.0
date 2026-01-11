package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yourname/company-superapp/internal/domain"
)

type SearchRepository struct {
	db *sqlx.DB
}

func NewSearchRepository(db *sqlx.DB) *SearchRepository {
	return &SearchRepository{db: db}
}

type userSearchResult struct {
	ID       string  `db:"id"`
	FullName *string `db:"full_name"`
	Email    string  `db:"email"`
	Rank     float64 `db:"rank"`
}

func (r *SearchRepository) SearchUsers(ctx context.Context, query string, limit int) ([]domain.SearchResult, error) {
	sql := `
		SELECT 
			id::text,
			full_name,
			email,
			ts_rank(search_vector, plainto_tsquery('simple', $1)) as rank
		FROM system.users
		WHERE search_vector @@ plainto_tsquery('simple', $1)
		ORDER BY rank DESC
		LIMIT $2
	`

	var results []userSearchResult
	if err := r.db.SelectContext(ctx, &results, sql, query, limit); err != nil {
		return nil, err
	}

	searchResults := make([]domain.SearchResult, len(results))
	for i, u := range results {
		title := u.Email
		if u.FullName != nil && *u.FullName != "" {
			title = *u.FullName
		}
		searchResults[i] = domain.SearchResult{
			Type:     domain.SearchTypeUser,
			ID:       u.ID,
			Title:    title,
			Subtitle: u.Email,
			Rank:     u.Rank,
		}
	}

	return searchResults, nil
}

type messageSearchResult struct {
	ID        int64   `db:"id"`
	ChatID    string  `db:"chat_id"`
	SenderID  string  `db:"sender_id"`
	Content   string  `db:"content"`
	Rank      float64 `db:"rank"`
}

func (r *SearchRepository) SearchMessages(ctx context.Context, query string, limit int) ([]domain.SearchResult, error) {
	sql := `
		SELECT 
			m.id,
			m.chat_id::text,
			m.sender_id::text,
			m.content,
			ts_rank(m.search_vector, plainto_tsquery('simple', $1)) as rank
		FROM messenger.messages m
		WHERE m.search_vector @@ plainto_tsquery('simple', $1)
		ORDER BY rank DESC
		LIMIT $2
	`

	var results []messageSearchResult
	if err := r.db.SelectContext(ctx, &results, sql, query, limit); err != nil {
		return nil, err
	}

	searchResults := make([]domain.SearchResult, len(results))
	for i, m := range results {
		// Обрезаем контент для превью
		content := m.Content
		if len(content) > 100 {
			content = content[:100] + "..."
		}
		
		searchResults[i] = domain.SearchResult{
			Type:     domain.SearchTypeMessage,
			ID:       fmt.Sprintf("%d", m.ID),
			Title:    content,
			Subtitle: "Сообщение в чате",
			Rank:     m.Rank,
		}
	}

	return searchResults, nil
}

type taskSearchResult struct {
	ID          string  `db:"id"`
	Title       string  `db:"title"`
	Description *string `db:"description"`
	Status      string  `db:"status"`
	Rank        float64 `db:"rank"`
}

func (r *SearchRepository) SearchTasks(ctx context.Context, query string, limit int) ([]domain.SearchResult, error) {
	sql := `
		SELECT 
			id::text,
			title,
			description,
			status,
			ts_rank(search_vector, plainto_tsquery('simple', $1)) as rank
		FROM tasks.tasks
		WHERE search_vector @@ plainto_tsquery('simple', $1)
		ORDER BY rank DESC
		LIMIT $2
	`

	var results []taskSearchResult
	if err := r.db.SelectContext(ctx, &results, sql, query, limit); err != nil {
		return nil, err
	}

	searchResults := make([]domain.SearchResult, len(results))
	for i, t := range results {
		subtitle := t.Status
		if t.Description != nil && *t.Description != "" {
			desc := *t.Description
			if len(desc) > 50 {
				desc = desc[:50] + "..."
			}
			subtitle = desc
		}
		
		searchResults[i] = domain.SearchResult{
			Type:     domain.SearchTypeTask,
			ID:       t.ID,
			Title:    t.Title,
			Subtitle: subtitle,
			Rank:     t.Rank,
		}
	}

	return searchResults, nil
}
