package db

import (
	"fmt"

	"github.com/jackc/pgx"
)

type tuple []interface{}

type column struct {
	Name, Type string
	Nullable   bool
}

// Keyvalue - represents {column: value} mapping
type Keyvalue map[string]interface{}

// Relation - represents PostgreSQL relation
type Relation struct {
	Name       string
	Attributes []column
	Pool       *pgx.ConnPool
}

// Connect - opens connection to PostgreSQL instance
func Connect() (*pgx.ConnPool, error) {
	connConfig := pgx.ConnConfig{
		User:     "begor",
		Database: "begor",
	}

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig: connConfig,
	}

	pool, err := pgx.NewConnPool(poolConfig)

	if err != nil {
		return nil, err
	}

	return pool, nil
}

// GetRelations returns existing relations for schema
func GetRelations(schema string, pool *pgx.ConnPool) ([]Relation, error) {
	var relations []Relation

	// TODO: rewrite to one query
	tableNames, err := getTableNames(schema, pool)

	if err != nil {
		return nil, err
	}

	for _, name := range tableNames {
		cols, err := getTableColumns(name, pool)

		if err != nil {
			return nil, err
		}

		relations = append(relations, Relation{name, cols, pool})
	}

	return relations, nil
}

// Select - selects given relation with limit and offset
// TODO: refactor this mess
func Select(relation Relation, limit int, offset int) []Keyvalue {
	// TODO: this is kinda awful, revisit
	// TODO: https://github.com/Masterminds/squirrel
	template := "SELECT * FROM %v LIMIT %d OFFSET %d;"
	query := fmt.Sprintf(template, relation.Name, limit, offset)

	// TODO: this is also awful, maybe in Go there are ways to represent it better

	var tuples []tuple

	rows, _ := relation.Pool.Query(query)

	for rows.Next() {
		vals, _ := rows.Values()
		tuples = append(tuples, vals)
	}

	fields := rows.FieldDescriptions()

	var keyvals []Keyvalue

	for _, tuple := range tuples {
		kv := make(Keyvalue)

		for i, fd := range fields {

			kv[fd.Name] = tuple[i]
		}

		keyvals = append(keyvals, kv)
	}

	return keyvals

}

func getTableColumns(tableName string, pool *pgx.ConnPool) ([]column, error) {
	var columns []column

	rows, err := pool.Query(`
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns 
		WHERE table_name = $1;
	`, tableName)

	if err != nil {
		return columns, nil
	}

	for rows.Next() {
		var name, dataType string
		var nullable bool

		err = rows.Scan(&name, &dataType, &nullable)

		columns = append(columns, column{name, dataType, nullable})
	}

	return columns, nil

}

func getTableNames(schema string, pool *pgx.ConnPool) ([]string, error) {
	var tableNames []string

	rows, err := pool.Query(`
		SELECT table_name AS name 
		FROM information_schema.tables 
		WHERE table_schema = $1;
	`, schema)

	if err != nil {
		return tableNames, err
	}

	for rows.Next() {
		var name string

		err = rows.Scan(&name)

		if err != nil {
			return tableNames, err
		}

		tableNames = append(tableNames, name)
	}

	return tableNames, nil
}
