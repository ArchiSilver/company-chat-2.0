package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/google/uuid"
	"github.com/yourname/company-superapp/internal/domain"
)

type ReportService struct {
	taskRepo domain.TaskRepository
}

func NewReportService(taskRepo domain.TaskRepository) *ReportService {
	return &ReportService{taskRepo: taskRepo}
}

// GenerateTasksReport генерирует PDF-отчёт по задачам пользователя за указанный период
func (s *ReportService) GenerateTasksReport(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]byte, error) {
	// Получение задач за указанный период
	tasks, err := s.taskRepo.GetByDateRange(ctx, userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tasks: %w", err)
	}

	// Создание PDF-документа
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	// Заголовок
	pdf.SetFont("Arial", "B", 20)
	pdf.SetTextColor(51, 51, 51)
	pdf.CellFormat(0, 15, "Tasks Report", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Информация о периоде
	pdf.SetFont("Arial", "", 11)
	pdf.SetTextColor(100, 100, 100)
	dateRangeText := fmt.Sprintf("Period: %s - %s", from.Format("02.01.2006"), to.Format("02.01.2006"))
	pdf.CellFormat(0, 8, dateRangeText, "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Дата генерации
	pdf.SetFont("Arial", "I", 9)
	generatedText := fmt.Sprintf("Generated: %s", time.Now().Format("02.01.2006 15:04"))
	pdf.CellFormat(0, 6, generatedText, "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Сводка статистики
	todoCount, inProgressCount, doneCount := s.countTasksByStatus(tasks)
	s.renderStatistics(pdf, len(tasks), todoCount, inProgressCount, doneCount)
	pdf.Ln(10)

	// Таблица задач
	if len(tasks) > 0 {
		s.renderTasksTable(pdf, tasks)
	} else {
		pdf.SetFont("Arial", "I", 12)
		pdf.SetTextColor(150, 150, 150)
		pdf.CellFormat(0, 20, "No tasks found for the selected period", "", 1, "C", false, 0, "")
	}

	// Нижний колонтитул
	pdf.SetY(-20)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(150, 150, 150)
	pdf.CellFormat(0, 10, "Company SuperApp - Confidential", "", 0, "C", false, 0, "")

	// Вывод в байты
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

func (s *ReportService) countTasksByStatus(tasks []domain.Task) (todo, inProgress, done int) {
	for _, task := range tasks {
		switch task.Status {
		case domain.TaskStatusTodo:
			todo++
		case domain.TaskStatusInProgress:
			inProgress++
		case domain.TaskStatusDone:
			done++
		}
	}
	return
}

func (s *ReportService) renderStatistics(pdf *fpdf.Fpdf, total, todo, inProgress, done int) {
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(51, 51, 51)
	pdf.CellFormat(0, 8, "Summary", "", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "", 10)

	// Statistics boxes
	boxWidth := 42.0
	startX := 15.0

	// Total Tasks
	pdf.SetFillColor(240, 240, 240)
	pdf.SetXY(startX, pdf.GetY())
	pdf.CellFormat(boxWidth, 20, fmt.Sprintf("Total: %d", total), "1", 0, "C", true, 0, "")

	// Todo
	pdf.SetFillColor(255, 235, 235)
	pdf.SetXY(startX+boxWidth+5, pdf.GetY())
	pdf.CellFormat(boxWidth, 20, fmt.Sprintf("Todo: %d", todo), "1", 0, "C", true, 0, "")

	// In Progress
	pdf.SetFillColor(255, 250, 235)
	pdf.SetXY(startX+2*(boxWidth+5), pdf.GetY())
	pdf.CellFormat(boxWidth, 20, fmt.Sprintf("In Progress: %d", inProgress), "1", 0, "C", true, 0, "")

	// Done
	pdf.SetFillColor(235, 255, 235)
	pdf.SetXY(startX+3*(boxWidth+5), pdf.GetY())
	pdf.CellFormat(boxWidth, 20, fmt.Sprintf("Done: %d", done), "1", 1, "C", true, 0, "")
}

func (s *ReportService) renderTasksTable(pdf *fpdf.Fpdf, tasks []domain.Task) {
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(51, 51, 51)
	pdf.CellFormat(0, 8, "Tasks List", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// Table Header
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(255, 75, 51) // Primary color #FF4B33
	pdf.SetTextColor(255, 255, 255)

	colWidths := []float64{70, 30, 40, 40}
	headers := []string{"Title", "Status", "Created", "Due Date"}

	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Table Rows
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(51, 51, 51)

	for i, task := range tasks {
		// Alternate row colors
		if i%2 == 0 {
			pdf.SetFillColor(250, 250, 250)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		// Title (truncate if too long)
		title := task.Title
		if len(title) > 35 {
			title = title[:32] + "..."
		}
		pdf.CellFormat(colWidths[0], 8, title, "1", 0, "L", true, 0, "")

		// Status with color coding
		statusText := s.formatStatus(task.Status)
		pdf.CellFormat(colWidths[1], 8, statusText, "1", 0, "C", true, 0, "")

		// Created Date
		pdf.CellFormat(colWidths[2], 8, task.CreatedAt.Format("02.01.2006"), "1", 0, "C", true, 0, "")

		// Due Date
		dueDate := "-"
		if task.DueDate != nil {
			dueDate = task.DueDate.Format("02.01.2006")
		}
		pdf.CellFormat(colWidths[3], 8, dueDate, "1", 1, "C", true, 0, "")
	}
}

func (s *ReportService) formatStatus(status domain.TaskStatus) string {
	switch status {
	case domain.TaskStatusTodo:
		return "Todo"
	case domain.TaskStatusInProgress:
		return "In Progress"
	case domain.TaskStatusDone:
		return "Done"
	default:
		return string(status)
	}
}
