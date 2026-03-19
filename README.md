# URL Shortener with SSO

Микросервисный проект, состоящий из двух сервисов:
- **REST API URL Shortener** — сервис для сокращения ссылок
- **gRPC SSO** — сервис аутентификации и авторизации

Структура проекта

├── sso/                   # gRPC сервис аутентификации
│   ├── cmd/               # Точки входа
│   │   ├── migrator/      # Утилита миграций
│   │   └── sso/           # Основной сервис
│   ├── internal/         
│   │   ├── app/           # Инициализация приложения
│   │   ├── config/        # Конфигурация
│   │   ├── domain/        # Модели данных
│   │   ├── grpc/          # gRPC хендлеры
│   │   ├── lib/           # Утилиты (jwt, логгер)
│   │   ├── services/      # Бизнес-логика
│   │   └── storage/       # Работа с БД
│   ├── migrations/        # SQL миграции
│   ├── tests/             # Интеграционные тесты
│   └── proto/             # Protobuf файлы
│
├── rest_API/              # REST API URL shortener
│   ├── cmd/               # Точка входа
│   ├── internal/          
│   │   ├── config/        # Конфигурация
│   │   ├── http-server/   # HTTP хендлеры
│   │   ├── lib/           # Утилиты
│   │   ├── storage/       # Работа с БД
│   │   └── clients/       # gRPC клиенты
│   └── mocks/             # Сгенерированные моки
│
└── docker-compose.yaml    
