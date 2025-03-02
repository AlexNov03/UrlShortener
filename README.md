# Репозиторий проекта URLShortener 
Приложение должно принимать оригинальный URL и создавать на основе него сокращенный. 

В качестве базы данных необходимо привести две реализации:
- PostgreSql
- самостоятельно реализованный пакет для хранения ссылок в памяти приложения
## Техническое задание 
[Ссылка на задание](https://docs.google.com/document/d/1gPAgIpscDjXrczlDdzLfS-XJqpu59HjcgRgO0eRsTvM/edit?tab=t.0)
## Режимы запуска
Для изменения режима запуска необходимо изменить docker-compose.yaml файл

Запуск приложения с базой данных Postgres 
```yaml
mainservice:
    build: .
    container_name: url_shortener
    ports:
        - "8080:8080"
    command: ["./output", "-in-memory=false"]
    depends_on:
        - db1
```
Запуск приложения с in-memory базой данных 
```yaml
mainservice:
    build: .
    container_name: url_shortener
    ports:
        - "8080:8080"
    command: ["./output", "-in-memory=true"]
    depends_on:
        - db1
```
## Примеры входных и выходных данных
Входные данные 
```json
{
  "original_url":"http://ya.ru"
}
```
Выходные данные 
```json
{
  "shortened_url":"http://<IP_ADDRESS>:8080/Ab_Cgf_edB"
}
```
## Работа с приложением
Запуск приложения
```shell
make start
```
Остановка всех контейнеров приложения
```shell
make stop
```
## Работа с миграциями 
Применить миграцию
```shell
make migrate-up
```
Откатить миграцию
```shell
make migrate-down
```
## Покрытие тестами по пакетам
delivery ![Coverage](https://img.shields.io/badge/Coverage-92.6%25-90EE90)


usecase  ![Coverage](https://img.shields.io/badge/Coverage-95.0%25-90EE90)


repository/pg ![Coverage](https://img.shields.io/badge/Coverage-90.5%25-c5e384)
