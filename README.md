# Сервис баннеров

## Описание

Сервис, который позволяет показывать пользователям баннеры, в зависимости от требуемой фичи и тега пользователя, а также управлять баннерами и связанными с ними тегами и фичами.

## Общие вводные

Баннер — это документ, описывающий какой-либо элемент пользовательского интерфейса. Технически баннер представляет собой  JSON-документ неопределенной структуры.  Тег — это сущность для обозначения группы пользователей; представляет собой число (ID тега).  Фича — это домен или функциональность; представляет собой число (ID фичи).  

1. Один баннер может быть связан только с одной фичей и несколькими тегами
2. При этом один тег, как и одна фича, могут принадлежать разным баннерам одновременно
3. Фича и тег однозначно определяют баннер

Так как баннеры являются для пользователя вспомогательным функционалом, допускается, если пользователь в течение короткого срока будет получать устаревшую информацию. При этом существует часть пользователей (порядка 10%), которым обязательно получать самую актуальную информацию. Для таких пользователей нужно предусмотреть механизм получения информации напрямую из БД.

## Стэк технологий

- Язык программирования: Go 1.22
- СУБД: PostgreSQL для хранения баннеров, Redis для кэширования
- Деплой: Docker и Docker Compose
- Для отладки использовалась утилита Postman

Используемые библиотеки: 
1. godotenv для работы с файлами окружения
2. pq для работы с базой данных PostgreSQL
3. go-redis для работы с базой данных Redis

Работа с сетью выполнялась с помощью net/http

## Запуск
Запуск:
```
make run
```
Запуск тестов:
```
make test
```
Запросы к сервису:
- GET /user_banner?tag_id={int}&feature_id={int}&use_last_revision={bool}
```
docker-compose exec app curl --location 'http://localhost:8000/user_banner?tag_id={int}&feature_id={int}&use_last_revision={bool}' \
--header 'token: {string}'
```
- GET /banner?tag_id={int}&feature_id={int}&limit={int}&offset={int}
```
docker-compose exec app curl --location 'http://localhost:8000/banner?tag_id={int}&feature_id={int}&limit={int}&offset={int}' \
--header 'token: {string}'
```
- POST /banner
```
docker-compose exec app curl --location 'http://localhost:8080/banner' \
--header 'token: {string}' \
--header 'Content-Type: application/json' \
--data '{
    "tag_ids": [
        int, int, ...
    ],
    "feature_id": int,
    "content": {
        JSON
    },
    "is_active": bool
}'
```
- PATCH /banner/{id}
```
docker-compose exec app curl --location --request PATCH 'http://localhost:8000/banner/{id}' \
--header 'token: {string}' \
--header 'Content-Type: application/json' \
--data '{
    "tag_ids": [
        int, int, ...
    ],
    "feature_id": int,
    "content": {
        JSON
    },
    "is_active": bool
}'
```
- DELETE /banner/{id}
```
docker-compose exec app curl --location --request DELETE 'http://localhost:8000/banner/{id}' \
--header 'token: {string}'
```
(curl запросы просто скопировал из Postman)
## Вопросы/проблемы:

- Если сервис при запросе "GET /user_banner" просто бегает в базу данных, то поле "use_last_revision" бесполезно

Решение: Ответы от бд на 5 минут кэшируются в более быстрой для доступа бд Redis. Благодаря этому, можно сказать, система также адаптируется для значительного увеличения количества тегов и фичей, потому что часто используемые баннеры лишь раз в 5 минут обновляются в базе данных, остальное время получаются из Redis кэша.

- Непонятно зачем в запросе "GET /user_banner" ответ 403: Пользователь не имеет доступа.

Решение: просто не давал такой ответ. По логике, разница у админа и обычного пользователя в том, что админ видит отключенные баннеры, значит на отключенный баннер обычный пользователь должен получать 404.

- Сложные запросы и парсинг JSON:

У меня был крутой запрос на  "GET /banner", который прямо в базе данных генерировал массив из тэгов, благодаря чему одним запросом получалось то, что нужно, но библиотека наотрез отказывалась его воспринимать нормально, выдавая пустоту (тот же запрос в psql работал отлично), пришлось его разделить на 2 запроса.

Также была проблема с тем, что в content может находиться любой json, и получать его сразу в нужный тип оказалось нереально. Как итог получается строка и конвыертируется в тип с помощью библиотеки json.
