# ==========================================
# Company SuperApp - Makefile
# ==========================================

.PHONY: help up down build logs db-shell redis-shell api-logs migrate test clean

# Default target
help:
	@echo "Company SuperApp - Available commands:"
	@echo ""
	@echo "  make up        - Start all services (Docker Compose)"
	@echo "  make down      - Stop all services"
	@echo "  make build     - Build and start services"
	@echo "  make logs      - Show logs from all services"
	@echo "  make api-logs  - Show API server logs"
	@echo "  make db-shell  - Open PostgreSQL shell"
	@echo "  make redis-cli - Open Redis CLI"
	@echo "  make migrate   - Run database migrations"
	@echo "  make test      - Run backend tests"
	@echo "  make clean     - Remove all containers and volumes"
	@echo ""

# Start all services
up:
	docker-compose up -d
	@echo ""
	@echo "✓ Services started!"
	@echo ""
	@echo "  API:        http://localhost:8080"
	@echo "  PostgreSQL: localhost:5432"
	@echo "  Redis:      localhost:6379"
	@echo "  MinIO:      http://localhost:9001 (admin: minioadmin/minioadminpassword)"
	@echo "  Prometheus: http://localhost:9090"
	@echo "  Grafana:    http://localhost:3000 (admin/admin)"
	@echo "  Jaeger:     http://localhost:16686"
	@echo ""

# Stop all services
down:
	docker-compose down

# Build and start
build:
	docker-compose up -d --build

# Show all logs
logs:
	docker-compose logs -f

# Show API logs
api-logs:
	docker-compose logs -f api

# PostgreSQL shell
db-shell:
	docker exec -it superapp_db psql -U admin -d company_db

# Redis CLI
redis-cli:
	docker exec -it superapp_cache redis-cli

# Run migrations manually (if needed)
migrate:
	cd backend && go run ./cmd/migrate

# Run tests
test:
	cd backend && go test ./...

# Clean up everything
clean:
	docker-compose down -v --remove-orphans
	rm -rf docker-data/postgres/*
	rm -rf docker-data/redis/*
	rm -rf docker-data/prometheus/*
	rm -rf docker-data/grafana/*
	@echo "✓ Cleaned up all data"

# Development - run backend locally
dev-backend:
	cd backend && go run ./cmd/api

# Development - run frontend
dev-frontend:
	cd frontend && npm start

# Database backup
db-backup:
	docker exec superapp_db pg_dump -U admin company_db > backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "✓ Database backed up"

# Database restore (usage: make db-restore FILE=backup.sql)
db-restore:
	cat $(FILE) | docker exec -i superapp_db psql -U admin company_db
	@echo "✓ Database restored from $(FILE)"
