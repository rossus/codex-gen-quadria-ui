package router

import (
	"github.com/julienschmidt/httprouter"

	"github.com/rossus/codex-gen-quadria-ui/constants"
	"github.com/rossus/codex-gen-quadria-ui/handlers"
	"github.com/rossus/codex-gen-quadria-ui/types"
)

// NewRouter returns a configured httprouter.Router.
func NewRouter(s *types.Server) *httprouter.Router {
	r := httprouter.New()
	r.GET(constants.RouteIndex, handlers.Index(s))
	r.POST(constants.RouteStart, handlers.Start(s))
	r.GET(constants.RouteGame, handlers.Game(s))
	return r
}
