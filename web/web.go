package web

import (
	"fmt"
	"net/http"

	"github.com/begor/pgxapi/db"
	"github.com/labstack/echo"
)

// StartWeb - starts HTTP server for given collection of relations
func StartWeb(relations []db.Relation, port string) {
	e := echo.New()

	setupRoutes(relations, e)

	e.Logger.Fatal(e.Start(port))
}

func setupRoutes(relations []db.Relation, e *echo.Echo) {
	for _, relation := range relations {
		path := fmt.Sprint("/", relation.Name)

		fmt.Print(path)

		e.GET(path, func(c echo.Context) error {
			xs := db.Select(relation, 10, 0)

			return c.JSON(http.StatusOK, xs)
		})
	}
}
