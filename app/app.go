package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/begor/pgxapi/db"
	"github.com/begor/pgxapi/web"
	"github.com/jackc/pgx"
)

type Connection struct {
	Config    pgx.ConnConfig
	Schema    string
	PoolSize  int
	Pool      *pgx.ConnPool
	Relations []db.Relation
}

// App - global app configuration
type App struct {
	Connections map[string]*Connection // name => *connection
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

func OpenPools(app App) error {
	for _, conn := range app.Connections {
		pool, err := db.OpenPool(conn.Config, conn.PoolSize)

		if err != nil {
			return err
		}

		conn.Pool = pool
	}

	return nil
}

func getRelations(app App) (map[string][]db.Relation, error) {
	relations := make(map[string][]db.Relation)

	for _, v := range app.Connections {
		fmt.Println(v.Pool)

		rels, err := db.GetRelations(v.Schema, v.Pool)

		if err != nil {
			return relations, err
		}

		relations[v.Config.Database] = rels
	}

	return relations, nil
}

func StartWeb(app App) {
	relations, err := getRelations(app)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relations: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Relations: %v\n", relations)

	web.StartWeb(relations, ":5050")
}
