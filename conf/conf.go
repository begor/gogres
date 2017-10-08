package conf

import (
	"encoding/json"
	"io/ioutil"

	"github.com/jackc/pgx"
)

// App - global app configuration
type App struct {
	Connection pgx.ConnConfig // TODO: multi connections
	Pool       int
}

// GetApp - factory function that returns App configuration
func GetApp(filepath string) (App, error) {
	dat, err := ioutil.ReadFile(filepath)

	if err != nil {
		return App{}, err
	}

	app := App{}

	json.Unmarshal(dat, &app)

	return app, nil
}
