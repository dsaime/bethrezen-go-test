
### Makefile

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