package main

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/rossus/quadria/board"
	"github.com/rossus/quadria/gameplay"
	"github.com/rossus/quadria/players"
)

func main() {
	indexTmpl := template.Must(template.ParseFiles("frontend/index.html"))
	gameTmpl := template.Must(template.ParseFiles("frontend/game.html"))

	// Register players once when server starts.
	players.InitPlayer("Blue", "blue")
	players.InitPlayer("Red", "red")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := indexTmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
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
		http.Redirect(w, r, "/game", http.StatusSeeOther)
	})

	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		b := board.GetBoard()
		if err := gameTmpl.Execute(w, b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.ListenAndServe(":8080", nil)
}
