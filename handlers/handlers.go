package handlers

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/rossus/codex-gen-quadria-ui/constants"
	"github.com/rossus/codex-gen-quadria-ui/types"
	"github.com/rossus/quadria/board"
	"github.com/rossus/quadria/gameplay"
	"github.com/rossus/quadria/players"
	"github.com/rossus/quadria/session"
)

// Handler wraps the application server so route handlers can be defined as
// methods without repeatedly passing the server around.
type Handler struct {
	Server *types.Server
}

// New returns a Handler bound to the provided Server.
func New(s *types.Server) *Handler {
	return &Handler{Server: s}
}

// Index renders the index page.
func (h *Handler) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := h.Server.IndexTmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Start handles starting a new game and redirects to the game page.
func (h *Handler) Start(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	name1 := r.FormValue("name1")
	if name1 == "" {
		name1 = "Player One"
	}
	name2 := r.FormValue("name2")
	if name2 == "" {
		name2 = "Player Two"
	}
	color1 := r.FormValue("color1")
	if color1 == "" {
		color1 = "blue"
	}
	color2 := r.FormValue("color2")
	if color2 == "" {
		color2 = "red"
	}

	if name1 == name2 || color1 == color2 {
		http.Error(w, "players must have distinct names and colors", http.StatusBadRequest)
		return
	}

	pls := players.InitPlayers()
	pls.AddPlayer(name1, color1)
	pls.AddPlayer(name2, color2)

	b := board.InitNewBoard(size, pls)
	g := gameplay.InitializeNewGame(pls)
	h.Server.Session = session.InitializeNewSession(pls, b, g)
	http.Redirect(w, r, constants.RouteGame, http.StatusSeeOther)
}

// Game renders the current game board.
func (h *Handler) Game(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if h.Server.Session == nil {
		http.Error(w, "no active game", http.StatusBadRequest)
		return
	}
	tiles := h.Server.Session.Board.GetTiles()
	if err := h.Server.GameTmpl.Execute(w, tiles); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
