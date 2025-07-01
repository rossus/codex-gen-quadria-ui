package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/rossus/codex-gen-quadria-ui/constants"
	"github.com/rossus/codex-gen-quadria-ui/handlers"
)

// NewRouter returns a configured httprouter.Router using the provided handler.
func NewRouter(h *handlers.Handler) *httprouter.Router {
	r := httprouter.New()

	r.GET(constants.RouteIndex, h.Index)
	r.POST(constants.RouteStart, h.Start)
	r.GET(constants.RouteGame, h.Game)
	r.ServeFiles("/static/*filepath", http.Dir("static"))
	return r
}
