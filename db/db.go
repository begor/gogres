package db

import "github.com/jackc/pgx"

type Column struct {
	Name, Type string
	Nullable bool
}

func Connect() (*pgx.Conn, error) {
	connConfig := pgx.ConnConfig{
		User: "begor",
		Database: "begor",
	}

	conn, err := pgx.Connect(connConfig)
	
	if err != nil {
		return nil, err
	}
	
	return conn, nil
}

func GetTableColumns(tableName string, conn *pgx.Conn) ([]Column, error) {
	var columns []Column

	rows, err := conn.Query(`SELECT column_name, data_type, is_nullable
							 FROM information_schema.columns 
		                     WHERE table_name = $1;`, tableName)

	if err != nil {
		return columns, nil
	}

	for rows.Next() {
		var name, data_type string
		var nullable bool
		
		err = rows.Scan(&name, &data_type, &nullable)


		columns = append(columns, Column{name, data_type, nullable})
	}

	return columns, nil

}

func GetTableNames(schema string, conn *pgx.Conn) ([]string, error) {
	var tableNames []string;
	rows, err := conn.Query(`SELECT table_name AS name 
							 FROM information_schema.tables 
							 WHERE table_schema = $1;`, schema)


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