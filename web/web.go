package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/begor/gogres/db"
	"github.com/labstack/echo"
)

// Endpoint - Represents API Endpoint
type Endpoint struct {
	Path   string
	Method string
}

// SelectParams - limit and offset, parsed from QueryParams
type SelectParams struct {
	Limit  int
	Offset int
}

// StartWeb - starts HTTP server for given collection of relations
func StartWeb(relations map[string][]db.Relation, port string) {
	e := echo.New()

	endpoints := setupRoutes(relations, e)

	printEndpoints(endpoints)

	e.Start(port)
}

func setupRoutes(relations map[string][]db.Relation, e *echo.Echo) []Endpoint {
	var endpoints []Endpoint

	for databaseName, databaseRelations := range relations {
		for _, relation := range databaseRelations {
			path := makeGetPath(databaseName, relation.Name)
			handler := makeGetEndpoint(relation)
			endpoints = append(endpoints, Endpoint{path, "GET"})
			e.GET(path, handler)
		}
	}

	return endpoints
}

func printEndpoints(endpoints []Endpoint) {
	fmt.Println("Generated endpoints:")

	for _, endpoint := range endpoints {
		fmt.Println(endpoint.Method, endpoint.Path)
	}
}

func makeGetPath(databaseName string, relationName string) string {
	return fmt.Sprint("/", databaseName, "/", relationName, "/")
}

func makeGetEndpoint(relation db.Relation) func(echo.Context) error {
	handler := func(c echo.Context) error {
		params := parseGetQueryParams(c)
		tuples, err := db.Select(relation, params.Limit, params.Offset)

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]error{"error": err})
		}

		response := addMeta(tuples, params.Limit, params.Offset)

		return c.JSON(http.StatusOK, response)
	}

	return handler
}

func parseGetQueryParams(c echo.Context) SelectParams {
	limit := parseStrParamToInt(c, "limit", 10)
	offset := parseStrParamToInt(c, "offset", 0)

	return SelectParams{limit, offset}
}

func parseStrParamToInt(c echo.Context, param string, deflt int) int {
	result := deflt
	strParam := c.QueryParam(param)

	if strParam != "" {
		if maybeResult, err := strconv.Atoi(strParam); err == nil {
			result = maybeResult
		}
	}

	return result
}

func addMeta(tuples []db.Keyvalue, limit int, offset int) db.Keyvalue {
	response := make(db.Keyvalue)
	meta := make(db.Keyvalue)

	meta["limit"] = limit
	meta["offset"] = offset

	response["meta"] = meta
	response["tuples"] = tuples

	return response
}
