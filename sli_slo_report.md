# SLI / SLO Report: ToDo App

## SLIs (Service Level Indicators)

| Metric              | Description                     | Tool       |
|---------------------|----------------------------------|------------|
| http_requests_total | Кол-во HTTP-запросов по методам | Prometheus |
| uptime              | Аптайм контейнера                | Docker / Node exporter |
| response_latency    | Задержка ответа сервера         | Добавим позже |

## SLOs (Service Level Objectives)

| Name         | Objective      | Target   | Period    |
|--------------|----------------|----------|-----------|
| Availability | API доступность| ≥ 99.9%  | 30 days   |
| Latency      | Response <200ms| ≥ 95%    | 30 days   |

## Current Metrics (Simulated)

- Availability (за 7 дней): **99.93%**
- Avg latency: **145ms**
- Error rate: **0.01%**

## Tools Used
- Prometheus
- Grafana
- Docker
- Go HTTP server
