package main

import (
	"fmt"
	"os"

	"github.com/begor/pgxapi/app"
)

func main() {
	confFile := "./conf.json"
	application, err := app.GetApp(confFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse conf file %v, error: %v\n", confFile, err)
		os.Exit(1)
	}

	err = app.OpenPools(application)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to establish connection: %v\n", err)
		os.Exit(1)
	}

	app.StartWeb(application)
}
