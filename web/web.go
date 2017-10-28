package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/begor/gogres/db"
	"github.com/labstack/echo"
)

// SelectParams - limit and offset, parsed from QueryParams
type SelectParams struct {
	Limit  int
	Offset int
}

// StartWeb - starts HTTP server for given collection of relations
func StartWeb(databases map[string]*db.Database, port string) error {
	e := echo.New()

	setupRoutes(databases, e)

	return e.Start(port)
}

func setupRoutes(databases map[string]*db.Database, e *echo.Echo) {
	fmt.Printf("Generating endpoints...\n")

	setupInfoRoute(e)

	fmt.Print("GET /api/\n")

	for name, database := range databases {
		setupRoutesForDatabase(name, database, e)
	}
}

func setupRoutesForDatabase(name string, database *db.Database, e *echo.Echo) {
	for schemaName, relations := range database.Relations {
		for _, relation := range relations {
			path := makeGetPath(name, schemaName, relation.Name)
			handler := makeGetEndpoint(database, schemaName, relation)
			fmt.Printf("GET %v\n", path)
			e.GET(path, handler)
		}
	}
}

func setupInfoRoute(e *echo.Echo) {
	e.GET("/api/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, e.Routes())
	})
}

func makeGetPath(prefix string, schemaName string, relationName string) string {
	return fmt.Sprint("/api/", prefix, "/", schemaName, "/", relationName, "/")
}

func makeGetEndpoint(database *db.Database, schema string, relation db.Relation) func(echo.Context) error {
	handler := func(c echo.Context) error {
		params := parseGetQueryParams(c)
		tuples, err := database.Select(schema, relation, params.Limit, params.Offset)

		if err != nil {
			return c.JSON(http.StatusBadRequest, generateErrorResp(err))
		}

		return c.JSON(http.StatusOK, generateOkResp(params, tuples))
	}

	return handler
}

func generateErrorResp(err error) map[string]error {
	return map[string]error{"error": err}
}

func generateOkResp(params SelectParams, tuples []db.Keyvalue) db.Keyvalue {
	response := make(db.Keyvalue)
	meta := make(db.Keyvalue)

	meta["limit"] = params.Limit
	meta["offset"] = params.Offset

	response["meta"] = meta
	response["tuples"] = tuples

	return response
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
