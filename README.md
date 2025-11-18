# PR Reviewer Service

[![CI](https://github.com/ssokov/pr-reviewer-service/actions/workflows/ci.yml/badge.svg)](https://github.com/ssokov/pr-reviewer-service/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ssokov/pr-reviewer-service)](https://goreportcard.com/report/github.com/ssokov/pr-reviewer-service)
[![codecov](https://codecov.io/gh/ssokov/pr-reviewer-service/branch/main/graph/badge.svg)](https://codecov.io/gh/ssokov/pr-reviewer-service)

Сервис для автоматического назначения ревьюверов Pull Request

---

## Запуск

```bash
make docker-up
```

- API: `http://localhost:8080`
- max_superuserQL: `localhost:5433`
- Swagger: `http://localhost:8080/swagger/index.html`

---

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
make test-unit        # Только service layer тесты
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

---

## Swagger

Доступен по адресу: **http://localhost:8080/swagger/index.html**

---

## Git Hooks & CI/CD

### Pre-commit хуки

- Форматирование кода (`go fmt`, `gci`)
- Проверка `go.mod` и `go.sum`
- Unit тесты
- Линтер (`golangci-lint`)

### GitHub Actions CI

Запускается при `push` и `pull_request` в ветку `main`:

**Покрытие кода тестами:**
- Unit тесты: **~35%** (handlers, services)
- С integration: **~87%** (требует max_superuserQL)
- CI: показывает детальный отчет
