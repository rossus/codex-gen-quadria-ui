package types

import "html/template"

// Server stores pre-parsed templates used by handlers.
type Server struct {
	IndexTmpl *template.Template
	GameTmpl  *template.Template
}

// NewServer loads templates and returns a Server instance.
func NewServer() *Server {
	return &Server{
		IndexTmpl: template.Must(template.ParseFiles("frontend/index.html")),
		GameTmpl:  template.Must(template.ParseFiles("frontend/game.html")),
	}
}
