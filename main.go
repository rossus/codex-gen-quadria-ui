package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rossus/quadria/board"
	"github.com/rossus/quadria/gameplay"
	"github.com/rossus/quadria/players"
)

func main() {
	// Register players once when server starts.
	players.InitPlayer("Blue", "blue")
	players.InitPlayer("Red", "red")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<html><body><h1>Quadria UI</h1>`+
			`<form action="/start" method="POST">`+
			`Board size: <input type="number" name="size" value="3" min="2"/>`+
			`<input type="submit" value="Start"/>`+
			`</form></body></html>`)
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
		fmt.Fprint(w, "<html><body><h1>Quadria UI</h1><table>")
		for _, row := range b.Tiles {
			fmt.Fprint(w, "<tr>")
			for _, tile := range row {
				fmt.Fprintf(w,
					"<td style='width:30px;height:30px;text-align:center;background:%s'>%d</td>",
					tile.Player.Color, tile.Value)
			}
			fmt.Fprint(w, "</tr>")
		}
		fmt.Fprint(w, "</table></body></html>")
	})

	http.ListenAndServe(":8080", nil)
}
