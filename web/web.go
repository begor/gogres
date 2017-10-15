package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/begor/gogres/db"
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
func StartWeb(relations map[string][]db.Relation, port string) {
	e := echo.New()

	endpoints := setupRoutes(relations, e)

	printEndpoints(endpoints)

	e.Start(port)
}

func setupRoutes(relations map[string][]db.Relation, e *echo.Echo) []Endpoint {
	var endpoints []Endpoint

	for prefix, thisRelations := range relations {
		for _, relation := range thisRelations {
			path := fmt.Sprint("/", prefix, "/", relation.Name, "/")
			getHandler := makeRelationGetEndpoint(relation)

			getEndpoint := Endpoint{path, "GET"}

			endpoints = append(endpoints, getEndpoint)

			e.GET(path, getHandler)
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

func makeRelationGetEndpoint(relation db.Relation) func(echo.Context) error {
	handler := func(c echo.Context) error {
		params := parseGetQueryParams(c)
		tuples := db.Select(relation, params.Limit, params.Offset)
		response := addMeta(tuples, params.Limit, params.Offset)

		return c.JSON(http.StatusOK, response)
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
