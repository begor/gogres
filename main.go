package main

import "fmt"
import "github.com/jackc/pgx"
import "os"

func main() {
	conn := Connect("connect test")
	defer conn.Close()
	tableNames := GetTableNames("public", conn)
	
	fmt.Fprintf(os.Stdout, "Tables: %v\n", tableNames)
}


func GetTableNames(schema string, conn *pgx.Conn) ([]string) {
	var tableNames []string;
	rows, err := conn.Query("SELECT table_name AS name FROM information_schema.tables WHERE table_schema = $1;", schema)


	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
	}

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		tableNames = append(tableNames, name)
	}

	return tableNames
}

func Connect(applicationName string) (conn *pgx.Conn) {
	var runtimeParams map[string]string;
	runtimeParams = make(map[string]string)
	runtimeParams["application_name"] = applicationName

	connConfig := pgx.ConnConfig{
		User: "begor",
		Database: "begor",
		RuntimeParams: runtimeParams,
	}

	conn, err := pgx.Connect(connConfig)
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to establish connection: %v\n", err)
		os.Exit(1)
	}
	
	return conn
}