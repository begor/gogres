package main

import (
	"fmt"
	"os"
	"pgxapi/db"
)

func main() {
	conn, err := db.Connect()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to establish connection: %v\n", err)
		os.Exit(1)
	}
	
	defer conn.Close()
	
	tableNames, err := db.GetTableNames("public", conn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch table names: %v\n", err)
		os.Exit(1)
	}

	tableInfo, err := db.GetTableColumns(tableNames[0], conn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch %v infos: %v\n", tableNames[0], err)
		os.Exit(1)
	}
	
	fmt.Fprintf(os.Stdout, "Tables: %v\n", tableNames)
	fmt.Fprintf(os.Stdout, "Tables info: %v\n", tableInfo)
}


