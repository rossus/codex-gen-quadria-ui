package types

import (
	"html/template"

	"github.com/rossus/quadria/session"
)

// Server stores pre-parsed templates used by handlers.
type Server struct {
	IndexTmpl *template.Template
	GameTmpl  *template.Template
	Session   *session.Session
}

// NewServer loads templates and returns a Server instance.
func NewServer() *Server {
	return &Server{
		IndexTmpl: template.Must(template.ParseFiles("frontend/index.html")),
		GameTmpl:  template.Must(template.ParseFiles("frontend/game.html")),
	}
}
