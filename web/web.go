package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/begor/pgxapi/db"
	"github.com/labstack/echo"
)

type Endpoint struct {
	Path   string
	Method string
}

type SelectParams struct {
	Limit  int
	Offset int
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
		params := parseGetQueryParams(c)
		tuples := db.Select(relation, params.Limit, params.Offset)

		return c.JSON(http.StatusOK, tuples)
	}

	return handler
}

func parseGetQueryParams(c echo.Context) SelectParams {
	// Sane defaults
	// TODO: move to config
	limit := parseStrParamToInt(c, "limit", 10)
	offset := parseStrParamToInt(c, "offset", 0)

	return SelectParams{limit, offset}
}

func parseStrParamToInt(c echo.Context, param string, dflt int) int {
	result := dflt
	strParam := c.QueryParam(param)

	if strParam != "" {
		maybeResult, err := strconv.Atoi(strParam)

		if err == nil {
			result = maybeResult
		}
	}

	return result
}
