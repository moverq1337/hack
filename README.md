# HR Avatar - Система автоматизации HR-процессов

Современная система для автоматизации процессов найма, включающая анализ резюме, проведение голосовых интервью и оценку кандидатов с использованием искусственного интеллекта.

## 🚀 Возможности

- **Анализ резюме**: Автоматический парсинг и анализ резюме в формате DOCX
- **Голосовые интервью**: Проведение интервью с использованием речевых технологий
- **Оценка кандидатов**: AI-анализ соответствия кандидата вакансии
- **Управление вакансиями**: Создание и управление вакансиями
- **Современный веб-интерфейс**: Интуитивно понятный пользовательский интерфейс

## 🏗️ Архитектура

Система построена на микросервисной архитектуре с использованием следующих технологий:

### Сервисы
- **API Gateway** (Go + Gin): Основной HTTP API и маршрутизация запросов
- **Resume Service** (Go): Обработка и парсинг резюме
- **Interview Service** (Python): Голосовые интервью с использованием Deepgram
- **Scoring Service** (Python): NLP-анализ и оценка кандидатов
- **AI Service** (Go): Интеграция с OpenRouter для генерации вопросов и анализа

### Инфраструктура
- **PostgreSQL**: Основная база данных
- **Redis**: Кэширование и сессии
- **Kafka**: Очереди сообщений для асинхронной обработки
- **Yandex Disk**: Хранение файлов резюме

## 🛠️ Технологический стек

### Backend
- **Go**: Основной язык для API и сервисов
- **Python**: NLP и AI-обработка
- **gRPC**: Межсервисное взаимодействие
- **Gin**: HTTP веб-фреймворк
- **GORM**: ORM для работы с базой данных

### Frontend
- **HTML5/CSS3**: Современный адаптивный интерфейс
- **JavaScript**: Интерактивность и AJAX-запросы

### AI/ML
- **OpenRouter**: Генерация вопросов и анализ текста
- **Deepgram**: Речевые технологии
- **spaCy**: NLP-обработка текста
- **Sentence Transformers**: Семантический анализ

### Инфраструктура
- **Docker**: Контейнеризация
- **Docker Compose**: Оркестрация сервисов
- **Kafka**: Стриминг данных
- **PostgreSQL**: Реляционная база данных
- **Redis**: Кэш и очереди

## 📋 Требования

- Docker и Docker Compose
- Go 1.21+
- Python 3.9+
- Node.js (для разработки)

## 🚀 Быстрый старт

### 1. Клонирование репозитория
```bash
git clone <repository-url>
cd VTBHack
```

### 2. Настройка переменных окружения
Создайте файл `.env` в корне проекта:
```env
DB_URL=postgres://postgres:password@postgres:5432/hrdb?sslmode=disable
GRPC_PORT=:50051
HTTP_PORT=:8080
REDIS_ADDR=redis:6379
KAFKA_BROKERS=kafka:9092
YANDEX_DISK_TOKEN=your_yandex_disk_token
UNIDOC_LICENSE_API_KEY=your_unidoc_license_key
YANDEX_API_KEY=your_yandex_api_key
DEEPGRAM_API_KEY=your_deepgram_api_key
OPENROUTER_API_KEY=your_openrouter_api_key
```

### 3. Запуск системы
```bash
docker-compose up -d
```

### 4. Доступ к приложению
- **Веб-интерфейс**: http://localhost:8080
- **API Gateway**: http://localhost:8080/api
- **Kafka UI**: http://localhost:8888

## 📁 Структура проекта

```
VTBHack/
├── cmd/                    # Основные сервисы
│   ├── ai-service/         # AI-сервис для генерации вопросов
│   ├── api-gateway/        # API Gateway
│   ├── interview-service/  # Сервис голосовых интервью
│   ├── resume-service/     # Сервис обработки резюме
│   └── scoring-service/    # NLP-сервис для анализа
├── internal/               # Внутренние пакеты
│   ├── config/            # Конфигурация
│   ├── db/                # Работа с базой данных
│   ├── handlers/          # HTTP обработчики
│   ├── models/            # Модели данных
│   ├── pb/                # Protobuf файлы
│   └── utils/             # Утилиты
├── frontend/              # Веб-интерфейс
│   ├── index.html         # Главная страница
│   └── interview.html     # Страница интервью
├── proto/                 # Protobuf схемы
├── scripts/               # Скрипты миграции
└── docker-compose.yml     # Конфигурация Docker
```

## 🔧 API Endpoints

### Основные эндпоинты
- `GET /api/vacancies` - Получение списка вакансий
- `POST /api/upload-resume` - Загрузка резюме
- `POST /api/analyze-resume` - Анализ резюме
- `POST /api/upload/vacancy` - Создание вакансии
- `POST /api/analyze-resume-vacancy` - Анализ резюме для конкретной вакансии

### Примеры использования

#### Создание вакансии
```bash
curl -X POST http://localhost:8080/api/upload/vacancy \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Поедатель чипсов",
    "requirements": "Глубокие знания Windows, Linux, MacOS. Опыт администрирования Active Directory, DNS, DHCP, GPO. Навыки автоматизации с помощью PowerShell и VBScript. Знание ITIL, опыт работы с Service Desk системами (1C Itilium, HP Service Manager, Jira). Администрирование антивирусных решений (Kaspersky Security Center). Уверенный английский язык (B2).",
    "responsibilities": "Техническая поддержка пользователей (on-site и remote). Подготовка и сопровождение рабочих станций на Windows, Linux, MacOS. Администрирование локальной сети и серверов. Организация и сопровождение рабочих мест сотрудников (включая кроссировку и настройку сетевого оборудования). Управление печатной инфраструктурой и серверами печати. Разработка и сопровождение корпоративных образов ОС. Тестирование и внедрение систем Service Desk. Обучение и наставничество новых сотрудников.",
    "region": "Санкт-Петербург",
    "city": "Санкт-Петербург",
    "employment_type": "Полная занятость",
    "work_schedule": "Комбинированный (гибридный) формат работы",
    "experience": "15+ лет",
    "education": "Высшее + курсы Microsoft, ITIL, Linux, PowerShell, Service Desk",
    "salary_min": 120000,
    "salary_max": 160000,
    "languages": "Русский, Английский",
    "skills": "Active Directory, Windows Server, Linux, MacOS, DHCP, DNS, GPO, ITIL, Service Desk (1C Itilium, HP Service Manager, Jira), PowerShell, VBScript, Kaspersky Security Center, администрирование рабочих мест, поддержка пользователей"
  }'
```

#### Получение списка вакансий
```bash
curl -X GET http://localhost:8080/api/vacancies
```

### AI-сервис
- `POST /generate-questions` - Генерация вопросов для интервью
- `POST /analyze-answer` - Анализ ответов кандидата
- `POST /analyze-resume` - AI-анализ резюме

## 🗄️ Модели данных

### Vacancy (Вакансия)
- ID, заголовок, требования, обязанности
- Регион, город, тип занятости
- Опыт, образование, зарплата
- Языки, навыки

### Resume (Резюме)
- ID кандидата, текст резюме
- Парсированные данные (JSON)
- URL файла, дата создания

### AnalysisResult (Результат анализа)
- Связь резюме и вакансии
- Оценка соответствия
- Детали анализа

## 🔄 Процесс работы

1. **Создание вакансии**: HR-специалист создает вакансию через веб-интерфейс
2. **Загрузка резюме**: Кандидат загружает резюме в формате DOCX
3. **Парсинг резюме**: Система извлекает структурированные данные
4. **Анализ соответствия**: AI оценивает соответствие кандидата вакансии
5. **Голосовое интервью**: Проведение интервью с генерацией вопросов
6. **Финальная оценка**: Комплексная оценка кандидата

## 🧪 Разработка

### Локальная разработка
```bash
# Запуск только инфраструктуры
docker-compose up -d postgres redis kafka1 kafka2 kafka3

# Запуск сервисов локально
go run cmd/api-gateway/main.go
go run cmd/resume-service/main.go
python cmd/scoring-service/main.py
```


## 📊 Мониторинг

- **Kafka UI**: http://localhost:8888 - Мониторинг очередей
- **Логи сервисов**: `docker-compose logs <service-name>`
- **Метрики Redis**: `redis-cli info`

## 🔐 Безопасность

- Все API ключи хранятся в переменных окружения
- Файлы резюме загружаются на защищенный Yandex Disk
- Валидация входных данных на всех уровнях

## 🤝 Вклад в проект

1. Форкните репозиторий
2. Создайте ветку для новой функции
3. Внесите изменения
4. Создайте Pull Request

## 📄 Лицензия

Проект разработан для VTB Hackathon 2025

## 👥 Команда

- Backend разработка (Go)
- AI/ML интеграция (Python)
- Frontend разработка (HTML/CSS/JS)
- DevOps и инфраструктура

---

**HR Avatar** - революционизируем процесс найма с помощью искусственного интеллекта! 🚀