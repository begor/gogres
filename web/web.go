package web

import (
	"fmt"

	"github.com/begor/pgxapi/db"
	"github.com/valyala/fasthttp"
)

// StartWeb - starts HTTP server for given collection of relations
func StartWeb(relations []db.Relation, port string) error {
	handler := genericRelationGetHandler

	return fasthttp.ListenAndServe(port, handler)
}

func genericRelationGetHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Wow such http very web 2.0")
}
