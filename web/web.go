package web

import (
	"fmt"
	"net/http"

	"github.com/begor/pgxapi/db"
	"github.com/labstack/echo"
)

type Endpoint struct {
	Path   string
	Method string
}

// StartWeb - starts HTTP server for given collection of relations
func StartWeb(relations []db.Relation, port string) {
	e := echo.New()

	endpoints := setupRoutes(relations, e)

	printEndpoints(endpoints)

	e.Start(port)
}

func setupRoutes(relations []db.Relation, e *echo.Echo) []Endpoint {
	var endpoints []Endpoint

	for _, relation := range relations {
		path := fmt.Sprint("/", relation.Name, "/")
		handler := makeRelationGetEndpoint(relation)

		endpoint := Endpoint{path, "GET"}
		endpoints = append(endpoints, endpoint)

		e.GET(path, handler)
	}

	return endpoints
}

func printEndpoints(endpoints []Endpoint) {
	fmt.Println("Generated endpoints:")

	for _, endpoint := range endpoints {
		fmt.Println(endpoint.Method, endpoint.Path)
	}
}

func makeRelationGetEndpoint(relation db.Relation) func(echo.Context) error {
	handler := func(c echo.Context) error {
		xs := db.Select(relation, 10, 0)

		return c.JSON(http.StatusOK, xs)
	}

	return handler
}
