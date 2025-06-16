package main

import (
	"log"
	"net/http"

	"github.com/rossus/codex-gen-quadria-ui/constants"
	"github.com/rossus/codex-gen-quadria-ui/router"
	"github.com/rossus/codex-gen-quadria-ui/types"
	"github.com/rossus/quadria/players"
)

func main() {
	// Register players once when server starts.
	players.InitPlayer("Blue", "blue")
	players.InitPlayer("Red", "red")

	srv := types.NewServer()
	r := router.NewRouter(srv)

	log.Fatal(http.ListenAndServe(constants.ServerPort, r))
}
