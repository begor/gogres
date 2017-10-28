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
}

func main() {
	// TODO: CLI
	confFile := "./conf.json"
	application, err := GetApp(confFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse conf file %v, error: %v\n", confFile, err)
		os.Exit(1)
	}

	err = OpenConnPools(application)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to establish connection: %v\n", err)
		os.Exit(1)
	}

	StartAPI(application)
}

// GetApp - returns App configuration
func GetApp(filepath string) (App, error) {
	dat, err := ioutil.ReadFile(filepath)

	if err != nil {
		return App{}, err
	}

	app := App{}

	json.Unmarshal(dat, &app)

	printDatabases(app)

	return app, nil
}

// OpenConnPools - opens connection pools for given settings
func OpenConnPools(app App) error {
	for _, database := range app.Databases {
		err := db.OpenPool(database)

		if err != nil {
			return err
		}
	}

	return nil
}

// StartAPI - starts REST API endpoints for a given app instance
func StartAPI(app App) {
	err := getRelations(app)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relations: %v\n", err)
		os.Exit(1)
	}

	web.StartWeb(app.Databases, ":5050")
}

func getRelations(app App) error {
	for _, database := range app.Databases {
		err := db.FetchRelations(database)

		if err != nil {
			return err
		}
	}

	return nil
}

func printDatabases(app App) {
	fmt.Println("Parsed config:")
	for name, databases := range app.Databases {
		fmt.Println(name)
		fmt.Println(databases)
	}
}
