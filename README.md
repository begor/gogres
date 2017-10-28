# gogres

**gogres** - REST API generator for PostgreSQL databases.

## Overview

**gogres** generates REST API for existing PostgreSQL databases.

It uses [pgx](https://github.com/jackc/pgx) for establishing connections and [echo](https://github.com/labstack/echo) for API generation.

Features:
- Multiple database support
- Specified schemas (default - `public` schema)
- Configurable size of connection pools
- `LIMIT` and `OFFSET` as query params for every generated endpoint

## Installation

`go get github.com/begor/gogres`

## Configuration

Sample configuration file looks like this:
```
{
    "databases": {
        "first": {
            "database": "first",
            "user": "begor", 
            "host": "localhost",
            "schemas": ["public", "private"], 
            "poolSize": 5
        },
        "second": {
            "database": "second", 
            "user": "begor", 
            "host": "some.host",
            "poolSize": 10
        }
    },
    "port": ":5050"
}
```

It states that `gogres` should generate API for two databases: `first` on `localhost` and `second` on `some.host`. 

For `first` it'll reflect two schemas: `public` and `private` and put them into `/api/first/public/` and `/api/first/private` API namespaces.

For `second` it'll use default schema (namely, `public`).

`gogres` also supports different pool settings for each database in a config. Here, `first` will be server via pool of 5 connections, and `second` of 10 connections.

As a result, API will be available at `localhost:5050`. You can visit `localhost:5050/api/` to see a list of generated endpoints.