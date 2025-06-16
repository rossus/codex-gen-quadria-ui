package handlers

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/rossus/codex-gen-quadria-ui/constants"
	"github.com/rossus/codex-gen-quadria-ui/types"
	"github.com/rossus/quadria/board"
	"github.com/rossus/quadria/gameplay"
)

// Index returns a handler for the index page.
func Index(s *types.Server) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if err := s.IndexTmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Start handles starting a new game and redirects to the game page.
func Start(s *types.Server) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		size := 3
		if val := r.FormValue("size"); val != "" {
			if i, err := strconv.Atoi(val); err == nil && i > 0 {
				size = i
			}
		}
		board.InitNewBoard(size)
		gameplay.StartNewGame()
		http.Redirect(w, r, constants.RouteGame, http.StatusSeeOther)
	}
}

// Game renders the current game board.
func Game(s *types.Server) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		b := board.GetBoard()
		if err := s.GameTmpl.Execute(w, b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
