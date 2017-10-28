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
func StartWeb(databases map[string]*db.Database, port string) {
	e := echo.New()

	endpoints := setupRoutes(databases, e)

	printEndpoints(endpoints)

	e.Start(port)
}

func setupRoutes(databases map[string]*db.Database, e *echo.Echo) []Endpoint {
	var endpoints []Endpoint

	for name, database := range databases {
		databaseEndpoints := setupRoutesForDatabase(name, database, e)
		endpoints = append(endpoints, databaseEndpoints...)
	}

	return endpoints
}

func setupRoutesForDatabase(name string, database *db.Database, e *echo.Echo) []Endpoint {
	var endpoints []Endpoint

	for schemaName, relations := range database.Relations {
		for _, relation := range relations {
			path := makeGetPath(name, schemaName, relation.Name)
			handler := makeGetEndpoint(database, schemaName, relation)
			endpoints = append(endpoints, Endpoint{path, "GET"})
			fmt.Printf("Generated: %v\n", path)
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

func makeGetPath(prefix string, schemaName string, relationName string) string {
	return fmt.Sprint("/", prefix, "/", schemaName, "/", relationName, "/")
}

func makeGetEndpoint(database *db.Database, schema string, relation db.Relation) func(echo.Context) error {
	handler := func(c echo.Context) error {
		params := parseGetQueryParams(c)
		tuples, err := db.Select(database, schema, relation, params.Limit, params.Offset)

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
