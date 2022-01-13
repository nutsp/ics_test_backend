# ics_test_backend

## Full list what has been used

* [fiber](https://github.com/gofiber/fiber) - Web framework
* [sqlx](https://github.com/jmoiron/sqlx) - Extensions to database/sql.
* [validator](https://github.com/go-playground/validator) - Go Struct and Field validation
* [apm](https://go.elastic.co/apm) - APM Agent for Go
* [CompileDaemon](https://github.com/githubnemo/CompileDaemon) - Compile daemon for Go
* [Docker](https://www.docker.com/) - Docker

## Usage
### Env

```
MY_SERVER_HOST= `Host API (http://localhost:3000)`
MYSQL_HOST= `Host Database (db)`
MYSQL_PORT= `PORT Database (3306)`
MYSQL_USER= `USER Database (root)`
MYSQL_PASS= `PASS Database (Icstest1234)`
MYSQL_DB= `Database Name (ics_test_db)`
```

### Build

```
docker build -t ics_test_backend:latest .
```

### Run

#### Docker-compose:

```
docker-compose up -d
```

API: _http://locahost:3000_

Adminer For Manage Database: _http://localhost:8080_

_For more examples, please refer to the [Documentation](https://documenter.getpostman.com/view/3109803/UVXgNdEz)_

