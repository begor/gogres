package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/begor/gogres/db"
	"github.com/begor/gogres/web"
	"github.com/jackc/pgx"
)

// Database - represents single connection to PostgreSQL database/schema
type Database struct {
	Config    pgx.ConnConfig
	Schema    string
	PoolSize  int
	Pool      *pgx.ConnPool
	Relations []db.Relation
}

// App - global app configuration
type App struct {
	Databases map[string]*Database // name => *connection
}

func main() {
	confFile := "./conf.json"
	application, err := GetApp(confFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse conf file %v, error: %v\n", confFile, err)
		os.Exit(1)
	}

	printConnections(application)

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

	return app, nil
}

// OpenConnPools - opens connection pools for given settings
func OpenConnPools(app App) error {
	for _, conn := range app.Databases {
		pool, err := db.OpenPool(conn.Config, conn.PoolSize)

		if err != nil {
			return err
		}

		conn.Pool = pool
	}

	return nil
}

// StartAPI - starts REST API endpoints for a given app instance
func StartAPI(app App) {
	relations, err := getRelations(app)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relations: %v\n", err)
		os.Exit(1)
	}

	web.StartWeb(relations, ":5050")
}

func getRelations(app App) (map[string][]db.Relation, error) {
	relations := make(map[string][]db.Relation)

	for name, database := range app.Databases {
		dbRelations, err := db.GetRelations(database.Schema, database.Pool)

		printRelations(name, dbRelations)

		if err != nil {
			return relations, err
		}

		relations[name] = dbRelations
	}

	return relations, nil
}

func printRelations(name string, relations []db.Relation) {
	fmt.Printf("Reflected %v relations:\n", name)
	for _, rel := range relations {
		fmt.Printf("\t- %v\n", rel.Name)
	}
	fmt.Println()
}

func printConnections(app App) {
	fmt.Println("Opened connections: ")
	fmt.Println("(name: database-schema-pool size): ")
	for name, database := range app.Databases {
		fmt.Printf("%v: %v-%v-%v\n", name, database.Config.Database, database.Schema, database.PoolSize)
	}
	fmt.Println()
}
