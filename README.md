# gogres

**gogres** - read-only REST API generator for PostgreSQL databases.

## Overview

**gogres** generated read-only REST API for existing PostgreSQL instances.

It uses [pgx](https://github.com/jackc/pgx) for establishing connections and [echo](https://github.com/labstack/echo) for REST API generation.

Some of the features:
- Multiple database support
- Specified schemas (default - all schemas in a database) *TODO*
- Specified tables (default - all tables in a schema) *TODO*
- Configurable size of connection pools for each database in a config

## Installation

`go get github.com/begor/gogres`

## Usage

*TODO*