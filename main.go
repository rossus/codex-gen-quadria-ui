package main

import (
	"log"
	"net/http"

	"github.com/rossus/codex-gen-quadria-ui/constants"
	"github.com/rossus/codex-gen-quadria-ui/handlers"
	"github.com/rossus/codex-gen-quadria-ui/router"
	"github.com/rossus/codex-gen-quadria-ui/types"
)

func main() {
	srv := types.NewServer()
	h := handlers.New(srv)
	r := router.NewRouter(h)

	log.Fatal(http.ListenAndServe(constants.ServerPort, r))
}
