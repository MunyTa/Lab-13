# Лабораторная работа №13

## Данные студента

- ФИО: Кузьмищев Родион Ильич
- Группа: 221331
- Вариант: 8
- Номер лабораторной: 13

## Тема работы

Вариант 8 в методичке относится к предметной области `Автоматизация HR`. В этой работе он адаптирован под `повышенную сложность`: вместо блока средней сложности реализован расширенный multi-agent сценарий с несколькими агентами, pipeline-оркестрацией, stateful-логикой, аукционным выбором исполнителя, LLM-агентом, мониторингом и инфраструктурой для трассировки.

## Реализованные компоненты

- `resume_parser` на Go: парсинг и нормализация резюме.
- `candidate_matcher` на Go: вычисление соответствия вакансии и кандидата.
- `interview_scheduler` на Go: планирование интервью.
- `feedback_agent` на Go: подготовка сообщения HR.
- `llm_feedback_agent` на Python: интеллектуальная итоговая обратная связь.
- `HROrchestrator` на Python: логика pipeline, аукциона и маршрутизации.
- `FastAPI` API: запуск задач, мониторинг, план масштабирования.
- `docker-compose.yml`: локальный запуск `NATS`, `Redis`, `Jaeger` и API.

## Выполненные задания повышенной сложности

1. Реализована полная система из 5 агентов.
2. Реализована цепочка обработки задач `pipeline`.
3. Подготовлена инфраструктура для распределённой трассировки через `Jaeger`.
4. Реализован stateful-агент с восстанавливаемым состоянием.
5. Реализован сервис планирования масштабирования по глубине очереди.
6. Реализован аукционный механизм выбора агента.
7. Добавлен Python LLM-агент для интеллектуальной обратной связи.
8. Добавлен веб-слой для мониторинга и ручного запуска задач.

## Структура проекта

- [agents/common/models.go](/C:/Users/Kuzmi/OneDrive/Desktop/Новая%20папка%20(3)/Lab-13/agents/common/models.go)
- [agents/resume-parser/main.go](/C:/Users/Kuzmi/OneDrive/Desktop/Новая%20папка%20(3)/Lab-13/agents/resume-parser/main.go)
- [agents/candidate-matcher/main.go](/C:/Users/Kuzmi/OneDrive/Desktop/Новая%20папка%20(3)/Lab-13/agents/candidate-matcher/main.go)
- [agents/interview-scheduler/main.go](/C:/Users/Kuzmi/OneDrive/Desktop/Новая%20папка%20(3)/Lab-13/agents/interview-scheduler/main.go)
- [agents/feedback-agent/main.go](/C:/Users/Kuzmi/OneDrive/Desktop/Новая%20папка%20(3)/Lab-13/agents/feedback-agent/main.go)
- [src/orchestrator/api.py](/C:/Users/Kuzmi/OneDrive/Desktop/Новая%20папка%20(3)/Lab-13/src/orchestrator/api.py)
- [src/orchestrator/services/orchestrator.py](/C:/Users/Kuzmi/OneDrive/Desktop/Новая%20папка%20(3)/Lab-13/src/orchestrator/services/orchestrator.py)
- [docs/architecture.md](/C:/Users/Kuzmi/OneDrive/Desktop/Новая%20папка%20(3)/Lab-13/docs/architecture.md)

## Запуск

Установка Python-зависимостей:

```powershell
python -m pip install -e .[dev]
```

Запуск Python-тестов:

```powershell
python -m pytest -p no:cacheprovider
```

Запуск Go-тестов:

```powershell
$env:GOTELEMETRY='off'
$env:GOCACHE=(Join-Path (Get-Location) '.gocache')
$env:GOMODCACHE=(Join-Path (Get-Location) '.gomodcache')
go test ./agents/...
```

Запуск инфраструктуры:

```powershell
docker compose up --build
```

## Качество кода

Код разбит на небольшие компоненты с понятными именами, общая логика вынесена в отдельные сервисы и модели. Это помогает держать решение ближе к принципам `SOLID`, избегать монолитных обработчиков и упрощает тестирование.
