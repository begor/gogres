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

// Keyvalue - represents {column: value} mapping, without any type-checking
type Keyvalue map[string]interface{}

// Relation - represents PostgreSQL relation
type Relation struct {
	Name       string
	Attributes []column
}

// SchemaRelations - mapping from schema name to its relations
type SchemaRelations map[string][]Relation

// Database - represents Schema with Schemas
type Database struct {
	pgx.ConnConfig
	Schemas   []string
	Relations SchemaRelations
	PoolSize  int
	Pool      *pgx.ConnPool
}

// FetchRelations - sets existing relations for database
func (d *Database) FetchRelations() error {
	d.checkSchemas()

	err := d.openPool()
	if err != nil {
		return err
	}

	relations := make(SchemaRelations)

	// TODO: rewrite to one query
	for _, schema := range d.Schemas {
		tableNames, err := getTableNames(schema, d.Pool)

		if err != nil {
			return err
		}

		for _, name := range tableNames {
			cols, err := getTableColumns(name, d.Pool)

			if err != nil {
				return err
			}

			relations[schema] = append(relations[schema], Relation{name, cols})
		}
	}

	d.Relations = relations

	return nil
}

// Select - selects given relation with limit and offset
func (d *Database) Select(schema string, relation Relation, limit int, offset int) ([]Keyvalue, error) {
	// TODO: this is kinda awful, revisit
	// TODO: https://github.com/Masterminds/squirrel
	template := "SELECT * FROM %v.%v LIMIT %d OFFSET %d;"
	query := fmt.Sprintf(template, schema, relation.Name, limit, offset)

	rows, err := d.Pool.Query(query)

	if err != nil {
		return make([]Keyvalue, 0), err
	}

	return parseSelectResult(rows)
}

func (d *Database) openPool() error {
	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     d.buildConnConfig(),
		MaxConnections: d.PoolSize,
	}

	pool, err := pgx.NewConnPool(poolConfig)

	if err != nil {
		return err
	}

	d.Pool = pool

	return nil
}

func (d *Database) checkSchemas() {
	var defaultSchema = []string{"public"}

	if d.Schemas == nil {
		d.Schemas = defaultSchema
	}
}

func (d *Database) buildConnConfig() pgx.ConnConfig {
	return pgx.ConnConfig{
		Host:     d.Host,
		Port:     d.Port,
		User:     d.User,
		Database: d.Database,
		Password: d.Password,
	}
}

func parseSelectResult(rows *pgx.Rows) ([]Keyvalue, error) {
	var rawTuples []tuple
	var columnValueMap []Keyvalue

	for rows.Next() {
		vals, _ := rows.Values()
		rawTuples = append(rawTuples, vals)
	}

	fields := rows.FieldDescriptions()

	for _, tuple := range rawTuples {
		kv := make(Keyvalue)

		for index, field := range fields {
			kv[field.Name] = tuple[index]
		}

		columnValueMap = append(columnValueMap, kv)
	}

	return columnValueMap, nil
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
