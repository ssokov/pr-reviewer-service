# PR Reviewer Service

Сервис для автоматического назначения ревьюверов Pull Request.

## Технологии

- Go `1.25`
- PostgreSQL `17`
- Go-pg `github.com/go-pg/pg/v10`
- Echo `v4`
- Swagger (`swaggo`)
- Genna (`model-named`, `search`, `validation`)

## Запуск (Docker)

```bash
docker compose -f deployments/docker/docker-compose.yml up -d db
docker compose -f deployments/docker/docker-compose.yml run --rm migrate \
  -path /migrations \
  -database "postgres://db:db@db:5432/pr_system?sslmode=disable" \
  up
docker compose -f deployments/docker/docker-compose.yml up -d api
```

- API: `http://localhost:8080`
- DB: `localhost:5433`
- Swagger: `http://localhost:8080/swagger/index.html`

## Локальная разработка

```bash
go test ./...
go run ./cmd/pr-reviewer-service
```

## Make команды

### Docker

```bash
make docker-up        # Запустить все сервисы
make docker-down      # Остановить все сервисы
make docker-build     # Собрать Docker образ
```

### Разработка

```bash
make build            # Собрать бинарник
make run              # Запустить локально
make fmt              # Форматировать код
make lint             # Запустить линтер
```

### Тесты

```bash
make test             # Unit тесты
make test-unit        # Только pr layer тесты
make test-integration # Интеграционные тесты
make test-all         # Все тесты
make test-coverage    # Тесты с отчетом покрытия (coverage.html)
```

### Миграции

```bash
make migrate-up       # Применить миграции
make migrate-down     # Откатить миграции
```

### Git Hooks

```bash
make install-hooks    # Установить pre-commit хуки
```

### Swagger

Доступен по адресу: `http://localhost:8080/swagger/index.html`

### Git Hooks & CI/CD

`pre-commit` хуки:

- Форматирование кода (`go fmt`, `gci`)
- Проверка `go.mod` и `go.sum`
- Unit тесты
- Линтер (`golangci-lint`)

GitHub Actions CI:

- Запускается при `push` и `pull_request` в ветку `main`

## Генерация db-кода (genna)

Модели/поиски/валидации генерируются из схемы БД:

```bash
genna model-named \
  -c "postgres://db:db@localhost:5433/pr_system?sslmode=disable" \
  -o internal/db/model_genna.go \
  -t "pr_system.*" \
  -f

genna search \
  -c "postgres://db:db@localhost:5433/pr_system?sslmode=disable" \
  -o internal/db/search_genna.go \
  -t "pr_system.*" \
  -f

genna validation \
  -c "postgres://db:db@localhost:5433/pr_system?sslmode=disable" \
  -o internal/db/validation_genna.go \
  -t "pr_system.*" \
  -f
```
