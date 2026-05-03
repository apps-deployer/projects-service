# Projects Service

Go gRPC-сервис для управления проектами, окружениями, шаблонами деплоя и переменными приложений.

## Назначение

- Хранит проекты и URL GitHub-репозиториев.
- Проверяет, что репозиторий проекта принадлежит GitHub-пользователю из JWT.
- Хранит отображаемое имя проекта и отдельный технический `slug`.
- Хранит окружения и целевые ветки.
- Хранит шаблоны деплоя и project-level overrides.
- Хранит зашифрованные переменные проекта и окружения.
- Отдает resolved variables для deployment runs.

## gRPC API

Сервис реализует API, сгенерированные из пакета `github.com/apps-deployer/protos`.

Основные сервисы:

- `ProjectService`
- `EnvService`
- `DeployConfigService`
- `FrameworkService`
- `VarService`

Порт gRPC по умолчанию:

```text
50051
```

## Конфигурация

Сервис читает YAML-конфиг из `CONFIG_PATH` или из аргумента `-config`.

Основные переменные окружения, которые используются в конфиге:

- `DB_PASSWORD`
- `JWT_SECRET`
- `ENCRYPTION_KEY`
- `DB_ADMIN_URL` - используется migrator init container.

## Правила для URL репозитория

Проекты принимают только GitHub HTTPS clone URL:

```text
https://github.com/owner/repo.git
```

Отклоняются SSH URL, URL без `.git`, не-GitHub hosts, query string, fragment и trailing slash.

## Имена проектов

`name` - отображаемое имя. Оно может содержать пробелы:

```text
Python Hello App
```

`slug` генерируется сервисом и используется для технических имен в деплое:

```text
python-hello-app
```

`slug` уникален в рамках владельца проекта.

## Локальный запуск

```bash
go mod download
go test ./...
go run ./cmd/projects_service -config ./config/config_docker.yaml
```

Локальный запуск миграций:

```bash
go run ./cmd/migrator -path ./migrations -database "$DB_ADMIN_URL" up
```

## Деплой

Helm chart находится в `charts/projects-service`. Chart запускает migrator init container перед стартом gRPC-сервера.

Сервис внутренний для кластера. Обычно он доступен по адресу:

```text
projects-service:50051
```
