package main

import (
       "fmt"
       "net/http"

       "github.com/rossus/quadria/board"
       "github.com/rossus/quadria/gameplay"
       "github.com/rossus/quadria/players"
)

func boardHTML() string {
       b := board.GetBoard()
       html := "<table style='border-collapse:collapse'>"
       for _, row := range b.Tiles {
               html += "<tr>"
               for _, t := range row {
                       color := t.Player.Color
                       if color == "" {
                               color = "white"
                       }
                       html += fmt.Sprintf("<td style='border:1px solid #999;width:30px;height:30px;text-align:center;background:%s'>%d</td>", color, t.Value)
               }
               html += "</tr>"
       }
       html += "</table>"
       return html
}

func main() {
       players.InitPlayer("Blue", "blue")
       players.InitPlayer("Red", "red")
       board.InitNewBoard(5)
       gameplay.StartNewGame()

       http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
               fmt.Fprintf(w, "<html><body><h1>Quadria UI</h1>%s</body></html>", boardHTML())
       })
       http.ListenAndServe(":8080", nil)
}
