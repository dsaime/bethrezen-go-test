# NEWSAPI - Выполнение тестового задания

http/json сервер с тремя ручками:
- `POST /create` - создать новость
- `POST /edit/:Id` - изменение новости по Id
- `GET /list` - список новостей

Что использовалось:
- web: fiber
- БД: mysql9.4, sqlx
- Логирование: logrus
- Тестирование: testify, testcontainers, mockery
- CLI: github.com/urfave/cli/v3 v3.4.1
- Запуск/упаковка: Docker, Docker-compose

Реализованные желания из постановки задачи:
- docker
- авторизацию через Authorization заголовок
- грамотная структуризация кода и ручек по группам/папкам
- валидацию полей при редактировании и создании
- грамотное логгирование с использованием любого популярного логгера(напр. logrus)
- грамотная обработка ошибок

### Makefile

`make test` Тестирует все пакеты проекта

`make vet` Запускает Статический анализ кода

`make lint` Запускает golangci-lint

`make check` Запускает vet и lint проверки последовательно

`make run` Запускает основной пакет приложения

`make build` Создает бинарник bin/newsapi

`make genmock` Запускает mockery/v3 для генерации mock-объектов

`make compose-up` Запускает сервис

`make compose-build` Собирает образ

### Запуск

Запуск в докере
```sh
make compose-build && make compose-up
```
- Команда соберет и запустит контейнеры `dsaime.test.mysql` и `dsaime.test.newsapi`
- Сервер будет доступен по адресу `http://localhost:8080`
- Для подключения к БД использовать DSN `root:root@tcp(127.0.0.1:3306)/test_db`


### Автотесты

```sh
make test
```
- Для тестирования mysql используется [testcontainers](https://golang.testcontainers.org/modules/mysql/)
- Для тестирования внешней БД, надо установить переменную окружения `TEST_MYSQL_DSN`. Бд должна иметь суффикс `test_`

> [!CAUTION]
> Таблицы в тестируемой БД очищаются после каждого теста

### Http api

Получить список новостей
```sh
curl --location 'http://localhost:8080/list' \
--header 'Authorization: Bearer foo'
```

Редактировать новость
```sh
curl --location 'http://localhost:8080/edit/1' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer foo' \
--data '{
  "Title": "Новое имя2",
  "Content": "новый контент",
  "Categories": [1,2,5,3]
}'
```

Создать новость
```sh
curl --location 'http://localhost:8080/create' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer foo' \
--data '{
  "Title": "Lorem ipsum",
  "Content": "Dolor sit amet <b>foo</b>",
  "Categories": [88,77,66]
}'
```