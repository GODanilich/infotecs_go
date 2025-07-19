# RUS Infotecs_go

## Что это?

 Это простое **приложение для обработки транзакций платёжной системы**
в виде **HTTP сервера**, реализующее **REST API** и имеющее три метода:
**POST Send**, **GET GetLast** и **GET GetBalance**.    

>**Данное приложение является тестовым заданием на стажировку**.
## Что было использовано в разработке?

Приложение написано на **Go** версии **1.22.2** с использованием сторонних библиотек, в качестве базы данных используется **PostgreSQL 17**.
 Для миграций таблиц используется [goose](https://github.com/pressly/goose),
  для генерации **sql-кода** используется [sqlc](https://github.com/sqlc-dev/sqlc?tab=readme-ov-file). Генерируемые sqlc запросы строго параметризированны, что практически полностью исключает SQL-инъекции.
 Для создания образов использовался **Docker**  и **Docker Compose** для оркестрации контейнеров.

## Запуск
>При первом запуске приложения создается 10 кошельков со случайными адресами и 100.0 у.е. на счету.

HTTP запросы к приложению можно делать любой удобной утилитой (curl, postman и т.д.)

### Docker-compose
> Запустить приложение вместе с базой данных можно командой:
```bash
docker-compose up --build
 ```
 > В контейнере db-1 можно войти в интерфейс PosgreSQL командой:
```
psql -h localhost -p 5433 -U user -d infotecs
```
> Получить таблицу кошельков для получения адресов кошельков:
```sql
SELECT * FROM wallets;
```
* Далее можно пользоваться приложением, отправля запросы к localhost:8080/***{эндпоинт}***.
### Нативный запуск
> Необходим клиент PostgreSQL

> В файле **.env** необходимо заменить **DB_URL** на собственный. После ссылки обязательно должен быть параметр **?sslmode=disable**.
```
PORT = 8080
DB_URL = postgres://postgres:1@172.24.32.1:5432/infotecs_go?sslmode=disable
MINIMAL_TRANSACTION_AMOUNT = 0.01
```
* По желанию можно изменить порт приложения или минимальную сумму транзакциии. **Важный момент** - все цифры младше второго десятичного разряда в минимальной сумме транзакции игнорируются. По умолчанию установлено минимально возможное значение для этой переменной. 

> Крайне желательно установить goose для удобной миграции таблиц:
```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
```

> Приминить миграции командой:
```sh
goose -dir ./sql/schema postgres "postgres://postgres:1@172.24.32.1:5432/infotecs_go" up
```
* Здесь параметр **?sslmode=disable** не является обязательным.

> Запустить:
```sh
go build && ./infotecs_go
```
* Далее можно пользоваться приложением, отправля запросы к localhost:8080/***{эндпоинт}***.
 ## Методы

* **POST Send** имеет эндпоинт **/api/send**. **Отправляет средства с одного из кошельков на указанный 
кошелек**. Метод принимает в теле запроса **JSON-объект**, содержащий следующие поля:
    * **from** – адрес кошелька, откуда нужно отправить деньги.  
    *Например: cdd750ee-5525-411b-bb7b-be2c23a8b926*
    * **to** – адрес кошелька, куда нужно отправить деньги.  
    *Например: b8169672-d721-465e-a0a0-eebed16c3e42*
    * **amount** – сумма перевода.  
    *Например: 3.50*.

> Возвращает json-объект созданной транзакции:
```json
{
    "id": "d8cf02a6-51d0-4d8d-801b-91a6a573f41c",
    "executed_at": "2025-07-18T05:11:04.459813Z",
    "amount": "10.22",
    "sender_address": "68ea3291-8d2f-4d57-a116-eda6d7466086",
    "recipient_address": "b5e1af31-e333-40e1-a874-0d7c571b110b"
}
```
> Если средств на кошельке недостаточно, то возвращается ошибка в виде json:
```json
{
    "error": "Not enough money: balance is 100, your amount 100.22"
}
```
> Если введенная сумма перевода меньше минимально допустимой, то возврщается ошибка в виде json:
```json
{
    "error": "Amout value is too small: minimum value is 0.01, your amount is 0"
}
```
> Если десятичная часть введенной суммы перевода содержит более 2ух знаков, то возврщается ошибка в виде json:
```json
{
    "error": "Amount can have maximum 2 decimal places, your input has 3 in .234"
}
```
> Поля "from" и "to" не должны совпадать, иначе возврщается ошибка в виде json:
```json
{
    "error": "Sender and recipient addresses can`t be the same"
}
```
***

* **GET GetLast** имеет эндпоинт **/api/transactions?count=N**. **Возвращает информацию о N последних по времени переводах средств**. Метод принимает в query-параметрах число возвращаемых **JSON-объектов** в массиве.
> Пример для N = 2:
```json
[
    {
        "id": "d8cf02a6-51d0-4d8d-801b-91a6a573f41c",
        "executed_at": "2025-07-18T05:11:04.459813Z",
        "amount": "10.22",
        "sender_address": "68ea3291-8d2f-4d57-a116-eda6d7466086",
        "recipient_address": "b5e1af31-e333-40e1-a874-0d7c571b110b"
    },
    {
        "id": "d50035ce-6036-4aed-ade2-fd0c6ff9024c",
        "executed_at": "2025-07-18T03:08:10.620316Z",
        "amount": "1.22",
        "sender_address": "cdd750ee-5525-411b-bb7b-be2c23a8b926",
        "recipient_address": "5ac7fcdb-f373-475c-a402-b7f4bc36a691"
    }
]
```
> Если на момент выполнения запроса не было выполнено ни 1 транзакции, то возвращает ошибку в виде json:
```json
{
    "error": "No transactions found"
}
```

***

* **GetBalance**, имеет эндпоинт **GET /api/wallet/{address}/balance**. **Возвращает информацию о балансе кошелька в JSON-объекте**. Метод принимаент адрес кошелька в пути запроса.


> Пример вывода в ответе:
```json
{
    "balance": "102.23"
}
```











*
*
*
*
*
*
*
*
*













# ENG Infotecs_go

## What is this?

This is a simple **payment transaction processing application** implemented as an **HTTP server** with a **REST API**, providing three methods:
**POST Send**, **GET GetLast**, and **GET GetBalance**.

>**This application was made as a test task for an internship.**

## What was used in development?

The application is written in **Golang 1.22.2**, with external libraries. It uses **PostgreSQL 17** as the database.
[goose](https://github.com/pressly/goose) is used for table migrations,
and [sqlc](https://github.com/sqlc-dev/sqlc?tab=readme-ov-file) is used to generate **SQL code**. The SQL queries generated by sqlc are strictly parameterized, which effectively eliminates the risk of SQL injection.
**Docker** was used to create images and **Docker Compose** for container orchestration.

## Running the App

>When the application starts for the first time, it creates 10 wallets with random addresses and a balance of 100.0 units.

HTTP requests can be made to the application using any convenient tool (curl, Postman, etc.)

### Docker Compose

> Run the application with database using:
```bash
docker-compose up --build
```

> To enter the PostgreSQL interface inside the db-1 container:
```
psql -h localhost -p 5433 -U user -d infotecs
```

> To retrieve the wallets table and get wallet`s addresses:
```sql
SELECT * FROM wallets;
```

* You can now use the application by sending requests to localhost:8080/***{endpoint}***.

### Native Run

> Requires a PostgreSQL client.

> In the **.env** file, replace **DB_URL** with your own connection string. It must include the parameter **?sslmode=disable** at the end.
```
PORT = 8080
DB_URL = postgres://postgres:1@172.24.32.1:5432/infotecs_go?sslmode=disable
MINIMAL_TRANSACTION_AMOUNT = 0.01
```
* You may change the port or minimum transaction amount if needed. **Important note** – any digits beyond the second decimal place in the minimum transaction amount will be ignored. The default value is the lowest allowed.

> It's highly recommended to install goose for easy table migration:
```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
```

> Apply migrations using:
```sh
goose -dir ./sql/schema postgres "postgres://postgres:1@172.24.32.1:5432/infotecs_go" up
```
* The **?sslmode=disable** parameter is optional here.

> Run the application:
```sh
go build && ./infotecs_go
```
* You can now use the application by sending requests to localhost:8080/***{endpoint}***.

## Methods

* **POST Send** has the endpoint **/api/send**. **Transfers funds from one wallet to another**.
It accepts a **JSON object** in the request body with the following fields:
    * **from** – sender's wallet address.  
    *Example: cdd750ee-5525-411b-bb7b-be2c23a8b926*
    * **to** – recipient's wallet address.  
    *Example: b8169672-d721-465e-a0a0-eebed16c3e42*
    * **amount** – amount to transfer.  
    *Example: 3.50*

> Returns a JSON object of the created transaction:
```json
{
    "id": "d8cf02a6-51d0-4d8d-801b-91a6a573f41c",
    "executed_at": "2025-07-18T05:11:04.459813Z",
    "amount": "10.22",
    "sender_address": "68ea3291-8d2f-4d57-a116-eda6d7466086",
    "recipient_address": "b5e1af31-e333-40e1-a874-0d7c571b110b"
}
```

> If the wallet has insufficient funds, returns a JSON error:
```json
{
    "error": "Not enough money: balance is 100, your amount 100.22"
}
```

> If the amount is below the minimum allowed, returns a JSON error:
```json
{
    "error": "Amout value is too small: minimum value is 0.01, your amount is 0"
}
```

> If the decimal part has more than 2 digits, returns a JSON error:
```json
{
    "error": "Amount can have maximum 2 decimal places, your input has 3 in .234"
}
```

> If the "from" and "to" fields are the same, returns a JSON error:
```json
{
    "error": "Sender and recipient addresses can`t be the same"
}
```

***

* **GET GetLast** has the endpoint **/api/transactions?count=N**. **Returns the latest N transactions**.
The method accepts the number of returned **JSON objects** as a query parameter.

> Example for N = 2:
```json
[
    {
        "id": "d8cf02a6-51d0-4d8d-801b-91a6a573f41c",
        "executed_at": "2025-07-18T05:11:04.459813Z",
        "amount": "10.22",
        "sender_address": "68ea3291-8d2f-4d57-a116-eda6d7466086",
        "recipient_address": "b5e1af31-e333-40e1-a874-0d7c571b110b"
    },
    {
        "id": "d50035ce-6036-4aed-ade2-fd0c6ff9024c",
        "executed_at": "2025-07-18T03:08:10.620316Z",
        "amount": "1.22",
        "sender_address": "cdd750ee-5525-411b-bb7b-be2c23a8b926",
        "recipient_address": "5ac7fcdb-f373-475c-a402-b7f4bc36a691"
    }
]
```

> If there have been no transactions, returns a JSON error:
```json
{
    "error": "No transactions found"
}
```

***

* **GET GetBalance** has the endpoint **GET /api/wallet/{address}/balance**. **Returns the balance of the specified wallet in a JSON object**. The method accepts the wallet address in the path.

> Example of a response:
```json
{
    "balance": "102.23"
}
```



