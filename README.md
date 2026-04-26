# 🤖 AI Agents Skills CLI

<div align="center">

[![Build Status](https://github.com/cloud-ru/evo-ai-agents-skills-cli/workflows/CI/badge.svg)](https://github.com/cloud-ru/evo-ai-agents-skills-cli/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25.2+-blue.svg)](https://golang.org/)
[![Release](https://img.shields.io/github/v/release/cloud-ru/evo-ai-agents-skills-cli)](https://github.com/cloud-ru/evo-ai-agents-skills-cli/releases)
[![GitHub stars](https://img.shields.io/github/stars/cloud-ru/evo-ai-agents-skills-cli.svg?style=flat-square&label=Stars)](https://github.com/cloud-ru/evo-ai-agents-skills-cli/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/cloud-ru/evo-ai-agents-skills-cli.svg?style=flat-square&label=Forks)](https://github.com/cloud-ru/evo-ai-agents-skills-cli/network/members)
[![Issues](https://img.shields.io/github/issues/cloud-ru/evo-ai-agents-skills-cli.svg?style=flat-square)](https://github.com/cloud-ru/evo-ai-agents-skills-cli/issues)

**CLI инструмент для загрузки, авторизации и локальной установки AI Skills из Cloud.ru BFF**

[Установка](#-установка) • [Запуск и авторизация](#запуск-и-авторизация) • [Быстрый старт](#-быстрый-старт) • [Конфигурация](#-конфигурация) • [Примеры](#-примеры-использования)

</div>

---

## ✨ Ключевые возможности

- 📦 **Загрузка skill archive**: скачивание скиллов из BFF endpoint в формате zip
- 📂 **Автоматическая распаковка**: безопасное извлечение архивов в локальную директорию
- 🔐 **IAM аутентификация**: авторизация по `IAM_KEY_ID` / `IAM_SECRET` через IAM API
- 🪪 **Поддержка Bearer token**: переопределение через `--auth-header` и `AUTH_HEADER`
- ♻️ **Кэш токена**: повторное использование IAM access token до истечения
- 🧩 **Совместимость с ai-agents-cli credentials**: автоподстановка из `~/.ai-agents-cli/credentials.json`
- 🧱 **Модульная архитектура**: `cmd` + `internal/api` + `internal/auth` + `internal/di`

---

## 🚀 Установка

### macOS/Linux (go install)

```bash
go install github.com/cloud-ru/evo-ai-agents-skills-cli@latest
```

### Сборка из исходников

```bash
git clone https://github.com/cloud-ru/evo-ai-agents-skills-cli.git
cd evo-ai-agents-skills-cli
go mod tidy
go build -o bin/skills-cli .
```

### Запуск без установки

```bash
go run . --help
```

---

## Запуск и авторизация

Ниже — рабочий порядок, чтобы **запросы к BFF шли с валидным `Authorization`**.

### Откуда CLI берёт доступ

При вызове `skill add` (и любых запросов через `internal/api`) заголовок выставляется так:

1. Флаг **`--auth-header`** (если задан).
2. Иначе переменная **`AUTH_HEADER`**.
3. Иначе **IAM**: запрос `POST {IAM_ENDPOINT}/api/v1/auth/token` с телом `keyId` / `secret` из **`IAM_KEY_ID`** и **`IAM_SECRET`** (токен кэшируется в памяти процесса).
4. Перед разбором env вызывается **`auth.InitCredentials()`**: если в окружении нет `IAM_KEY_ID` / `IAM_SECRET`, они подставляются из файла **`~/.ai-agents-cli/credentials.json`** (его создаёт `auth login`).

Итого: либо задаёте ключи в `.env` / окружении, либо один раз логинитесь и храните ключи в файле, либо передаёте готовый Bearer.

### Рекомендуемый порядок команд (с нуля)

1. **Клонировать / перейти в каталог** и подтянуть зависимости:

   ```bash
   cd evo-ai-agents-skills-cli
   go mod tidy
   ```

2. **Создать `.env`** из шаблона и заполнить хотя бы **`PUBLIC_BFF_URL`** (базовый URL public BFF без завершающего `/`).

   ```bash
   cp .env.example .env
   ```

3. **Выбрать способ авторизации** (достаточно одного):

   | Способ | Что сделать |
   |--------|-------------|
   | **Только `.env`** | В `.env` указать `IAM_KEY_ID`, `IAM_SECRET` (и при необходимости `IAM_ENDPOINT`, `PROJECT_ID`). Команда **`auth login` не обязательна**. |
   | **Файл credentials** | В `.env` те же переменные (удобно копировать из секрет-хранилища) и выполнить **`skills-cli auth login`** — значения подтянутся из конфига/env, лишний раз ключи не вводятся; в `~/.ai-agents-cli/credentials.json` сохранится профиль. Дальше можно не дублировать секреты в `.env`: при пустых `IAM_*` в env подставится файл. |
   | **Готовый токен** | В `.env` задать `AUTH_HEADER="Bearer <access_token>"` или передавать **`--auth-header`** в `skill add`. IAM не используется. |

4. **Проверить**, что окружение и файл (если есть) согласованы:

   ```bash
   go run . auth config    # путь к credentials + маскированные env
   go run . auth status    # есть ли сохранённый логин
   ```

5. **Скачать скилл**:

   ```bash
   go run . skill add --skill-id "<uuid-скилла>" --output ./skills
   ```

   Если `PROJECT_ID` не в `.env`, передайте **`--project-id`** (для маршрута `.../projects/{project_id}/skills/...`).

6. **Сменить аккаунт / убрать сохранённые ключи**:

   ```bash
   go run . auth logout
   ```

### Команда `auth login` и `.env`

`auth login` в начале вызывает **`config.Load()`**: подхватывается **`.env`** (через `godotenv/autoload`) и переменные окружения. Для каждого поля порядок такой: **флаг → значение из конфига (env) → запрос в консоль**.

Практически: если в `.env` уже есть **`IAM_KEY_ID`** и **`IAM_SECRET`**, достаточно выполнить **`skills-cli auth login`** без флагов — интерактивный ввод этих полей не понадобится (опционально спросит только `Project ID` / `Customer ID`, если их нет ни в флагах, ни в env).

Флаг **`--iam-endpoint`** по умолчанию пустой: endpoint берётся из **`IAM_ENDPOINT`** в env или подставляется `https://iam.api.cloud.ru`.

---

## 🎯 Быстрый старт

### 1️⃣ Настройка переменных окружения

```bash
cp .env.example .env
```

Минимум для первого успешного `skill add` через IAM:

```bash
PUBLIC_BFF_URL=https://<public-bff-host>
IAM_KEY_ID=<iam-key-id>
IAM_SECRET=<iam-secret>
```

Часто добавляют:

```bash
PROJECT_ID=<project-uuid>
IAM_ENDPOINT=https://iam.api.cloud.ru
CUSTOMER_ID=<optional>
```

### 2️⃣ (Опционально) Сохранить те же учётные данные в файл

```bash
skills-cli auth login
```

После этого при следующих запусках можно опираться на `~/.ai-agents-cli/credentials.json`, не дублируя секреты в `.env` (см. раздел [Запуск и авторизация](#запуск-и-авторизация)).

### 3️⃣ Загрузка скилла

```bash
skills-cli skill add \
  --skill-id "<skill-id>" \
  --output "./skills"
```

### 4️⃣ Проверка результата

```bash
ls ./skills
```

---

## 📋 Доступные команды

### 🧰 Корневая команда (`skills-cli`)

| Команда | Описание |
|---------|----------|
| `skills-cli --help` | Показать справку |
| `skills-cli --verbose` | Включить подробные логи |

### 🔐 Аутентификация (`auth`)

| Команда | Описание |
|---------|----------|
| `auth login` | Сохранить IAM-профиль в `~/.ai-agents-cli/credentials.json`; значения берутся из флагов, затем из `.env`/окружения (`config.Load`), затем из stdin |
| `auth logout` | Удалить файл credentials и сбросить связанные переменные окружения в процессе |
| `auth status` | Показать, настроен ли сохранённый профиль (маскированный Key ID) |
| `auth config` | Путь к файлу credentials, его содержимое (частично маскировано) и активные env |

### 📦 Управление скиллами (`skill`)

| Команда | Описание |
|---------|----------|
| `skill add` | Скачать и распаковать skill archive |

#### Флаги `skill add`

| Флаг | Описание | Обязательный |
|------|----------|--------------|
| `--skill-id` | ID скилла | ✅ |
| `--project-id` | ID проекта (если нужно scoped обращение) | ❌ |
| `--output` | Директория распаковки | ❌ |
| `--auth-header` | Явный `Authorization` заголовок | ❌ |

---

## 💡 Примеры использования

### Базовый сценарий

```bash
skills-cli skill add \
  --skill-id "6d2f66f6-7f31-44b1-92c7-7b4db978abcd" \
  --output "./skills"
```

### С указанием project id

```bash
skills-cli skill add \
  --project-id "a15e5e77-c8dc-4c4a-bf4d-4e29e7f8abcd" \
  --skill-id "6d2f66f6-7f31-44b1-92c7-7b4db978abcd" \
  --output "./skills"
```

### С явным Bearer token

```bash
skills-cli skill add \
  --project-id "a15e5e77-c8dc-4c4a-bf4d-4e29e7f8abcd" \
  --skill-id "6d2f66f6-7f31-44b1-92c7-7b4db978abcd" \
  --auth-header "Bearer <token>" \
  --output "./skills"
```

### Через environment (IAM flow)

```bash
export PUBLIC_BFF_URL="https://bff.example.com"
export IAM_KEY_ID="<iam-key-id>"
export IAM_SECRET="<iam-secret>"
export IAM_ENDPOINT="https://iam.api.cloud.ru"

skills-cli skill add --skill-id "6d2f66f6-7f31-44b1-92c7-7b4db978abcd"
```

---

## ⚙️ Конфигурация

### Переменные окружения

| Переменная | Описание | Обязательная | По умолчанию |
|------------|----------|--------------|--------------|
| `PUBLIC_BFF_URL` | Базовый URL public BFF | ✅ | - |
| `PROJECT_ID` | Проект по умолчанию для scoped запросов | ❌ | - |
| `HTTP_TIMEOUT_SEC` | Таймаут HTTP-запроса (сек) | ❌ | `60` |
| `AUTH_HEADER` | Готовый `Authorization` header | ❌ | - |
| `IAM_KEY_ID` | IAM Key ID для получения токена | ❌* | - |
| `IAM_SECRET` | IAM Secret для получения токена | ❌* | - |
| `IAM_ENDPOINT` | IAM API endpoint | ❌ | `https://iam.api.cloud.ru` |
| `CUSTOMER_ID` | Опционально для профиля и `auth login` | ❌ | - |

\* Обязательны, если не задан `AUTH_HEADER` и не используется внешний токен через `--auth-header`.

Шаблон всех переменных лежит в `.env.example`.

### Приоритет источников авторизации (для HTTP к BFF)

1. `--auth-header` у команды (`skill add` и т.д.)
2. `AUTH_HEADER` в окружении
3. IAM: `IAM_KEY_ID` + `IAM_SECRET` + опционально `IAM_ENDPOINT`
4. Если `IAM_KEY_ID` / `IAM_SECRET` в env пустые — чтение **`~/.ai-agents-cli/credentials.json`** при старте `config.Load()` → `auth.InitCredentials()`

---

## 🔌 Используемые endpoints

Если передан `project_id`:

```text
GET /api/v1/projects/{project_id}/skills/{skill_id}/.well-known/skills/archive.zip
```

Если `project_id` не передан:

```text
GET /api/v1/skills/{skill_id}/.well-known/skills/archive.zip
```

---

## 📁 Структура проекта

```text
evo-ai-agents-skills-cli/
├── cmd/
│   ├── root.go
│   ├── auth/
│   │   ├── root.go
│   │   ├── login.go
│   │   ├── logout.go
│   │   ├── status.go
│   │   └── config.go
│   └── skill/
│       ├── root.go
│       └── add.go
├── internal/
│   ├── api/
│   │   ├── api.go
│   │   ├── client.go
│   │   └── skill.go
│   ├── auth/
│   │   ├── credentials.go
│   │   ├── iam.go
│   │   └── init.go
│   ├── config/
│   │   └── config.go
│   └── di/
│       ├── container.go
│       └── errors.go
├── main.go
├── go.mod
└── README.md
```

---

## 🧪 Тестирование

```bash
go test ./...
```

---

## 🛠️ Разработка

### Требования

- Go 1.25.2+

### Локальная проверка без `go run` на каждый вызов

`go run .` каждый раз компилирует бинарник — для частых прогонов команд удобнее **один раз собрать** и вызывать готовый файл из каталога репозитория (тогда подхватывается тот же `.env` рядом с бинарником только если вы запускаете из этого каталога; `godotenv` ищет `.env` в **рабочей директории процесса**):

```bash
go build -o skills-cli .
./skills-cli auth config
./skills-cli skill add --project-id "..." --skill-id "..." --output ./skills
```

После правок в коде снова выполните `go build -o skills-cli .`. Для автопересборки можно подключить [air](https://github.com/air-verse/air) или `reflex`, но для CLI обычно достаточно ручного `go build`.

### Полезные команды

```bash
go mod tidy
go test ./...
go build ./...
go run . skill add --help
```

---

## 🆘 Troubleshooting

- `PUBLIC_BFF_URL is required`
  - укажите `PUBLIC_BFF_URL` в env/.env
- `get IAM token` / IAM `400` про пустой `ClientId`
  - раньше: неверные имена полей в JSON тела запроса к IAM; сейчас CLI шлёт **`clientId` / `clientSecret`**. Если ошибка сохраняется — проверьте непустые **`IAM_KEY_ID`** и **`IAM_SECRET`** и корректный **`IAM_ENDPOINT`** для вашего стенда
- `get IAM token` (прочие ошибки)
  - проверьте `IAM_ENDPOINT` и доступность IAM из сети
- `archive request failed with status 401/403`
  - у токена нет прав на проект/скилл, либо неверный `project_id`
- `archive request failed with status 404`
  - проверьте `skill_id` и маршрут

---

## 📄 Лицензия

Этот проект лицензирован под MIT License - см. файл `LICENSE`.

---

<div align="center">

**[⬆ Вернуться к началу](#-ai-agents-skills-cli)**

Made with ❤️ by [Cloud.ru](https://cloud.ru)

</div>