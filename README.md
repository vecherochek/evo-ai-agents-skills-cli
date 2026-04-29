# 🤖 AI Agents Skills CLI

<div align="center">

[![Build Status](https://github.com/vecherochek/evo-ai-agents-skills-cli/workflows/CI/badge.svg)](https://github.com/vecherochek/evo-ai-agents-skills-cli/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25.2+-blue.svg)](https://golang.org/)
[![Release](https://img.shields.io/github/v/release/vecherochek/evo-ai-agents-skills-cli)](https://github.com/vecherochek/evo-ai-agents-skills-cli/releases)
[![GitHub stars](https://img.shields.io/github/stars/vecherochek/evo-ai-agents-skills-cli.svg?style=flat-square&label=Stars)](https://github.com/vecherochek/evo-ai-agents-skills-cli/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/vecherochek/evo-ai-agents-skills-cli.svg?style=flat-square&label=Forks)](https://github.com/vecherochek/evo-ai-agents-skills-cli/network/members)
[![Issues](https://img.shields.io/github/issues/vecherochek/evo-ai-agents-skills-cli.svg?style=flat-square)](https://github.com/vecherochek/evo-ai-agents-skills-cli/issues)

**CLI инструмент для загрузки, авторизации и локальной установки AI Skills из Cloud.ru BFF**

[Установка](#-установка) • [Быстрый старт](#-быстрый-старт) • [Конфигурация](#-конфигурация) • [Примеры](#-примеры-использования)

</div>

---

## ✨ Ключевые возможности

- 📦 **Загрузка skill archive**: скачивание архивов скиллов из public BFF
- 📂 **Автоматическая распаковка**: безопасное извлечение zip в локальную директорию
- 🔐 **IAM аутентификация**: получение access token через `IAM_KEY_ID` / `IAM_SECRET`
- 🪪 **Поддержка Bearer token**: переопределение через `--auth-header` и `AUTH_HEADER`
- ♻️ **Поддержка credentials профиля**: автоподстановка из `~/.ai-agents-cli/credentials.json`
- 🧱 **Модульная архитектура**: `cmd` + `internal/api` + `internal/auth` + `internal/di`

---

## 🚀 Установка

### Windows (winget) (Скоро)

```powershell
winget install CloudRu.AIAgentsSkillsCLI
```

### Windows (Scoop) (Скоро)

```powershell
scoop bucket add cloud-ru https://github.com/cloud-ru/scoop-bucket
scoop install ai-agents-skills-cli
```

### macOS/Linux (Homebrew) (Скоро)

```bash
brew tap cloud-ru/evo-ai-agents-skills-cli
brew install ai-agents-skills-cli
```

### macOS/Linux (go install)

```bash
go install github.com/vecherochek/evo-ai-agents-skills-cli@latest
```

### Ручная установка (через GitHub Releases)

#### Linux

```bash
# Скачайте последнюю версию
wget https://github.com/vecherochek/evo-ai-agents-skills-cli/releases/latest/download/ai-agents-skills-cli-linux

# Установите
sudo mv ai-agents-skills-cli-linux /usr/local/bin/ai-agents-skills-cli
sudo chmod +x /usr/local/bin/ai-agents-skills-cli

# Проверьте установку
ai-agents-skills-cli --version
```

#### macOS

```bash
# Скачайте последнюю версию
curl -L https://github.com/vecherochek/evo-ai-agents-skills-cli/releases/latest/download/ai-agents-skills-cli-darwin -o ai-agents-skills-cli

# Установите
sudo mv ai-agents-skills-cli /usr/local/bin/ai-agents-skills-cli
sudo chmod +x /usr/local/bin/ai-agents-skills-cli

# Проверьте установку
ai-agents-skills-cli --version
```

#### Windows (PowerShell)

```powershell
# Скачайте последнюю версию
Invoke-WebRequest -Uri "https://github.com/vecherochek/evo-ai-agents-skills-cli/releases/latest/download/ai-agents-skills-cli-windows.exe" -OutFile "ai-agents-skills-cli.exe"

# Добавьте директорию в PATH или запускайте из текущей папки
# Проверка установки
.\ai-agents-skills-cli.exe --version
```

### Сборка из исходников

```bash
git clone https://github.com/vecherochek/evo-ai-agents-skills-cli.git
cd evo-ai-agents-skills-cli
go mod tidy
go build -o bin/ai-agents-skills-cli .
```

### Запуск без установки

```bash
go run . --help
```

---

## 🎯 Быстрый старт

### 1️⃣ Настройка переменных окружения

```bash
cp .env.example .env
```

Минимально обязательные параметры для `skill add`:

```bash
PUBLIC_BFF_URL=https://<public-bff-host>
PROJECT_ID=<project-id>
```

Варианты авторизации:

- IAM flow:

```bash
IAM_KEY_ID=<iam-key-id>
IAM_SECRET=<iam-secret>
IAM_ENDPOINT=https://iam.api.cloud.ru
```

- или готовый токен:

```bash
AUTH_HEADER=Bearer <token>
```

### 2️⃣ Проверка авторизации

```bash
ai-agents-skills-cli auth config
ai-agents-skills-cli auth status
```

### 3️⃣ Загрузка скилла

```bash
ai-agents-skills-cli skill add \
  --project-id "<project-id>" \
  --skill-id "<skill-id>" \
  --output "./skills"
```

### 4️⃣ Проверка результата

```bash
ls ./skills
```

---

## 📋 Доступные команды

### 🧰 Корневая команда (`ai-agents-skills-cli`)

| Команда | Описание |
|---------|----------|
| `ai-agents-skills-cli --help` | Показать справку |
| `ai-agents-skills-cli --verbose` | Включить подробные логи |

### 🔐 Аутентификация (`auth`)

| Команда | Описание |
|---------|----------|
| `auth login` | Сохранить IAM credentials в `~/.ai-agents-cli/credentials.json` |
| `auth logout` | Удалить сохраненный профиль |
| `auth status` | Показать статус аутентификации |
| `auth config` | Показать активную auth-конфигурацию |

### 📦 Управление скиллами (`skill`)

| Команда | Описание |
|---------|----------|
| `skill add` | Скачать и распаковать skill archive |
| `skill-marketplace add` | Скачать и распаковать marketplace skill archive |

#### Флаги `skill add`

| Флаг | Описание | Обязательный |
|------|----------|--------------|
| `--skill-id` | ID скилла | ✅ |
| `--project-id` | ID проекта | ✅* |
| `--output` | Директория распаковки | ❌ |
| `--auth-header` | Явный `Authorization` заголовок | ❌ |

\* Обязателен как `--project-id` или через `PROJECT_ID` в окружении.

#### Флаги `skill-marketplace add`

| Флаг | Описание | Обязательный |
|------|----------|--------------|
| `--skill-id` | ID marketplace скилла | ✅ |
| `--output` | Директория распаковки | ❌ |
| `--auth-header` | Явный `Authorization` заголовок | ❌ |

---

## 💡 Примеры использования

### Базовый сценарий

```bash
ai-agents-skills-cli skill add \
  --project-id "a15e5e77-c8dc-4c4a-bf4d-4e29e7f8abcd" \
  --skill-id "6d2f66f6-7f31-44b1-92c7-7b4db978abcd" \
  --output "./skills"
```

### С явным Bearer token

```bash
ai-agents-skills-cli skill add \
  --project-id "a15e5e77-c8dc-4c4a-bf4d-4e29e7f8abcd" \
  --skill-id "6d2f66f6-7f31-44b1-92c7-7b4db978abcd" \
  --auth-header "Bearer <token>" \
  --output "./skills"
```

### Через environment (IAM flow)

```bash
export PUBLIC_BFF_URL="https://bff.example.com"
export PROJECT_ID="a15e5e77-c8dc-4c4a-bf4d-4e29e7f8abcd"
export IAM_KEY_ID="<iam-key-id>"
export IAM_SECRET="<iam-secret>"
export IAM_ENDPOINT="https://iam.api.cloud.ru"

ai-agents-skills-cli skill add --skill-id "6d2f66f6-7f31-44b1-92c7-7b4db978abcd"
```

### Загрузка скилла из marketplace

```bash
ai-agents-skills-cli skill-marketplace add \
  --skill-id "6d2f66f6-7f31-44b1-92c7-7b4db978abcd" \
  --output "./skills"
```

---

## ⚙️ Конфигурация

### Переменные окружения

| Переменная | Описание | Обязательная | По умолчанию |
|------------|----------|--------------|--------------|
| `PUBLIC_BFF_URL` | Базовый URL public BFF | ✅ | - |
| `PROJECT_ID` | Проект по умолчанию для scoped запросов | ✅* | - |
| `HTTP_TIMEOUT_SEC` | Таймаут HTTP-запроса (сек) | ❌ | `60` |
| `AUTH_HEADER` | Готовый `Authorization` header | ❌ | - |
| `IAM_KEY_ID` | IAM Key ID для получения токена | ❌** | - |
| `IAM_SECRET` | IAM Secret для получения токена | ❌** | - |
| `IAM_ENDPOINT` | IAM API endpoint | ❌ | `https://iam.api.cloud.ru` |
| `CUSTOMER_ID` | Поле совместимости с общим credentials профилем | ❌ | - |

\* Обязателен, если не передан `--project-id`.

\** Обязательны, если не задан `AUTH_HEADER` и не используется внешний токен через `--auth-header`.

Шаблон всех переменных находится в `.env.example`.

### Приоритет источников авторизации

1. `--auth-header`
2. `AUTH_HEADER`
3. IAM token через `IAM_KEY_ID` + `IAM_SECRET` + `IAM_ENDPOINT`
4. Автоподстановка IAM данных из `~/.ai-agents-cli/credentials.json`

---

## 🔌 Используемый endpoint

```text
GET /api/v1/projects/{project_id}/skills/{skill_id}/.well-known/skills/archive.zip
GET /api/v1/marketplace/skills/{skill_id}/.well-known/skills/archive.zip
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
├── .env.example
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
- `project ID is required: pass --project-id or set PROJECT_ID`
  - передайте `--project-id` или задайте `PROJECT_ID`
- `get IAM token`
  - проверьте `IAM_KEY_ID`, `IAM_SECRET`, `IAM_ENDPOINT` и сетевую доступность IAM
- `archive request failed with status 401/403`
  - у токена нет прав на проект/скилл или указан неверный `project_id`
- `archive request failed with status 404`
  - проверьте `skill_id`, `project_id` и наличие архива у скилла

---

## 📄 Лицензия

Этот проект лицензирован под MIT License - см. файл [LICENSE](LICENSE).

---

## 🆘 Поддержка

- 🐛 **Баги**: [GitHub Issues](https://github.com/vecherochek/evo-ai-agents-skills-cli/issues)
- 💬 **Обсуждения**: [GitHub Discussions](https://github.com/vecherochek/evo-ai-agents-skills-cli/discussions)

---

<div align="center">

**[⬆ Вернуться к началу](#-ai-agents-skills-cli)**

Made with ❤️ by [Cloud.ru](https://cloud.ru)

</div>
