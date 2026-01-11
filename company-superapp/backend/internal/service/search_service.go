package service

import (
	"context"
	"sort"
	"sync"

	"github.com/yourname/company-superapp/internal/domain"
	"github.com/yourname/company-superapp/internal/repository/postgres"
)

type GlobalSearchService struct {
	searchRepo *postgres.SearchRepository
}

func NewGlobalSearchService(searchRepo *postgres.SearchRepository) *GlobalSearchService {
	return &GlobalSearchService{searchRepo: searchRepo}
}

func (s *GlobalSearchService) Search(ctx context.Context, query string) (*domain.SearchResults, error) {
	if query == "" {
		return &domain.SearchResults{
			Query:   query,
			Results: []domain.SearchResult{},
			Total:   0,
		}, nil
	}

	const limitPerType = 10

	var (
		wg           sync.WaitGroup
		userResults  []domain.SearchResult
		msgResults   []domain.SearchResult
		taskResults  []domain.SearchResult
		userErr      error
		msgErr       error
		taskErr      error
	)

	// Параллельный поиск пользователей
	wg.Add(1)
	go func() {
		defer wg.Done()
		userResults, userErr = s.searchRepo.SearchUsers(ctx, query, limitPerType)
	}()

	// Параллельный поиск сообщений
	wg.Add(1)
	go func() {
		defer wg.Done()
		msgResults, msgErr = s.searchRepo.SearchMessages(ctx, query, limitPerType)
	}()

	// Параллельный поиск задач
	wg.Add(1)
	go func() {
		defer wg.Done()
		taskResults, taskErr = s.searchRepo.SearchTasks(ctx, query, limitPerType)
	}()

	wg.Wait()

	// Проверка ошибок (логируем, но не прерываем полностью)
	if userErr != nil {
		userResults = []domain.SearchResult{}
	}
	if msgErr != nil {
		msgResults = []domain.SearchResult{}
	}
	if taskErr != nil {
		taskResults = []domain.SearchResult{}
	}

	// Агрегация результатов
	allResults := make([]domain.SearchResult, 0, len(userResults)+len(msgResults)+len(taskResults))
	allResults = append(allResults, userResults...)
	allResults = append(allResults, msgResults...)
	allResults = append(allResults, taskResults...)

	// Sort by rank (descending)
	sort.Slice(allResults, func(i, j int) bool {
		return allResults[i].Rank > allResults[j].Rank
	})

	// Limit total results
	if len(allResults) > 30 {
		allResults = allResults[:30]
	}

	return &domain.SearchResults{
		Query:   query,
		Results: allResults,
		Total:   len(allResults),
	}, nil
}
