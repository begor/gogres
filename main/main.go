package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/begor/gogres/db"
	"github.com/begor/gogres/web"
)

// App - global app configuration
type App struct {
	Databases map[string]*db.Database // name => database
	Port      string
}

func main() {
	// TODO: CLI
	confFile := "./conf.json"
	application, err := configure(confFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse conf file %v, error: %v\n", confFile, err)
		os.Exit(1)
	}

	err = openPools(application)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to establish connection: %v\n", err)
		os.Exit(1)
	}

	err = fetchRelations(application)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relations: %v\n", err)
		os.Exit(1)
	}

	err = startAPI(application)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start API: %v\n", err)
		os.Exit(1)
	}
}

func configure(filepath string) (App, error) {
	app := App{}
	data, err := ioutil.ReadFile(filepath)

	if err != nil {
		return app, err
	}

	json.Unmarshal(data, &app)

	return app, nil
}

func openPools(app App) error {
	for _, database := range app.Databases {
		err := db.OpenPool(database)

		if err != nil {
			return err
		}
	}

	return nil
}

func fetchRelations(app App) error {
	for _, database := range app.Databases {
		err := db.FetchRelations(database)

		if err != nil {
			return err
		}
	}

	return nil
}

func startAPI(app App) error {
	return web.StartWeb(app.Databases, ":5050")
}
