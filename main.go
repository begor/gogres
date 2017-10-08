package main

import (
	"fmt"
	"os"

	"github.com/begor/pgxapi/conf"
	"github.com/begor/pgxapi/db"
	"github.com/begor/pgxapi/web"
)

func main() {
	confFile := "./conf.json"
	app, err := conf.GetApp(confFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse conf file %v, error: %v\n", confFile, err)
		os.Exit(1)
	}

	pool, err := db.Connect(app)
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

	web.StartWeb(relations, ":5050")
}
