package main

import (
	"fmt"
	"os"

	"github.com/begor/pgxapi/db"
	"github.com/begor/pgxapi/web"
)

func main() {
	pool, err := db.Connect()
	defer pool.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to establish connection: %v\n", err)
		os.Exit(1)
	}

	relations, err := db.GetRelations("public", pool)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relations: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Relations: %v\n", relations)

	err = web.StartWeb(relations, ":5050")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start web: %v\n", err)
		os.Exit(1)
	}
}
