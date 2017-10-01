package main

import (
	"fmt"
	"os"

	"github.com/begor/pgxapi/db"
)

func main() {
	conn, err := db.Connect()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to establish connection: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close()

	relations, err := db.GetRelations("public", conn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relations: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Relations: %v\n", relations)
}
