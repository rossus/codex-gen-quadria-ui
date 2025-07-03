package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/rossus/codex-gen-quadria-ui/constants"
	"github.com/rossus/codex-gen-quadria-ui/types"
	"github.com/rossus/quadria/board"
	qtypes "github.com/rossus/quadria/common/types"
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

	h.Server.Players = []qtypes.Player{
		{Name: name1, Color: color1},
		{Name: name2, Color: color2},
	}

	b := board.InitNewBoard(size, pls)
	g := gameplay.InitializeNewGame(pls)
	h.Server.Session = session.InitializeNewSession(pls, b, g)
	h.Server.Winner = nil
	http.Redirect(w, r, constants.RouteGame, http.StatusSeeOther)
}

// Game renders the current game board.
func (h *Handler) Game(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if h.Server.Session == nil {
		http.Error(w, "no active game", http.StatusBadRequest)
		return
	}
	data := types.GamePageData{
		Tiles:   h.Server.Session.Board.GetTiles(),
		Turn:    h.Server.Session.Game.GetTurnNum(),
		Players: h.Server.Players,
		Active:  *h.Server.Session.Players.GetActivePlayer(),
		Winner:  h.Server.Winner,
	}
	if err := h.Server.GameTmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Move executes a player's turn on the selected tile then redirects to the game board.
func (h *Handler) Move(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if h.Server.Session == nil {
		http.Error(w, "no active game", http.StatusBadRequest)
		return
	}
	if h.Server.Winner != nil {
		http.Redirect(w, r, constants.RouteGame, http.StatusSeeOther)
		return
	}

	x, err := strconv.Atoi(ps.ByName("x"))
	if err != nil {
		http.Error(w, "invalid coordinates", http.StatusBadRequest)
		return
	}
	y, err := strconv.Atoi(ps.ByName("y"))
	if err != nil {
		http.Error(w, "invalid coordinates", http.StatusBadRequest)
		return
	}

	if h.Server.Session.Go(x, y) {
		h.Server.Winner = h.Server.Session.Players.GetActivePlayer()
	}
	http.Redirect(w, r, constants.RouteGame, http.StatusSeeOther)
}

// Dice serves a dice image for the provided value. If a custom image exists in
// constants.CustomDiceDir it is served. Otherwise an SVG is generated with the
// dice background colored via the "color" query parameter.
func (h *Handler) Dice(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	val, err := strconv.Atoi(ps.ByName("value"))
	if err != nil || val < 1 || val > 6 {
		http.NotFound(w, r)
		return
	}

	custom := filepath.Join(constants.CustomDiceDir, fmt.Sprintf("dice-%d.svg", val))
	if _, err := os.Stat(custom); err == nil {
		http.ServeFile(w, r, custom)
		return
	}

	color := r.URL.Query().Get("color")
	if color == "" {
		color = "white"
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write([]byte(renderDiceSVG(val, color)))
}

// renderDiceSVG returns a minimal dice SVG with the given value and fill color.
func renderDiceSVG(val int, color string) string {
	dots := map[int]string{
		1: "<circle cx='15' cy='15' r='3' fill='black'/>",
		2: "<circle cx='9' cy='9' r='3' fill='black'/><circle cx='21' cy='21' r='3' fill='black'/>",
		3: "<circle cx='9' cy='9' r='3' fill='black'/><circle cx='15' cy='15' r='3' fill='black'/><circle cx='21' cy='21' r='3' fill='black'/>",
		4: "<circle cx='9' cy='9' r='3' fill='black'/><circle cx='21' cy='9' r='3' fill='black'/><circle cx='9' cy='21' r='3' fill='black'/><circle cx='21' cy='21' r='3' fill='black'/>",
		5: "<circle cx='9' cy='9' r='3' fill='black'/><circle cx='21' cy='9' r='3' fill='black'/><circle cx='15' cy='15' r='3' fill='black'/><circle cx='9' cy='21' r='3' fill='black'/><circle cx='21' cy='21' r='3' fill='black'/>",
		6: "<circle cx='9' cy='7' r='3' fill='black'/><circle cx='9' cy='15' r='3' fill='black'/><circle cx='9' cy='23' r='3' fill='black'/><circle cx='21' cy='7' r='3' fill='black'/><circle cx='21' cy='15' r='3' fill='black'/><circle cx='21' cy='23' r='3' fill='black'/>",
	}

	return fmt.Sprintf("<svg xmlns='http://www.w3.org/2000/svg' width='30' height='30'><rect width='30' height='30' rx='5' ry='5' fill='%s' stroke='black'/>%s</svg>", color, dots[val])
}
