package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello{
	return &Hello{
		l: l,
	}
}

func (h *Hello) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	_, err := w.Write([]byte("Hello World!"))
    if err != nil {
        h.l.Warn(err.Error())
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}