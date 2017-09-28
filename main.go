package main

import "fmt"
import "github.com/jackc/pgx"
import "os"

type Column struct {
	Name, Type string
	Nullable bool
}

func main() {
	conn := Connect()
	
	defer conn.Close()
	
	tableNames, err := GetTableNames("public", conn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch table names: %v\n", err)
		os.Exit(1)
	}

	tableInfo, err := GetTableColumns(tableNames[0], conn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch %v infos: %v\n", tableNames[0], err)
		os.Exit(1)
	}
	
	fmt.Fprintf(os.Stdout, "Tables: %v\n", tableNames)
	fmt.Fprintf(os.Stdout, "Tables info: %v\n", tableInfo)
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
	rows, err := conn.Query("SELECT table_name AS name FROM information_schema.tables WHERE table_schema = $1;", schema)


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

func Connect() (conn *pgx.Conn) {
	connConfig := pgx.ConnConfig{
		User: "begor",
		Database: "begor",
	}

	conn, err := pgx.Connect(connConfig)
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to establish connection: %v\n", err)
		os.Exit(1)
	}
	
	return conn
}