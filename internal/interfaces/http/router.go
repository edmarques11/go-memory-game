package http

import (
	"html/template"
	"net/http"
)

func NewRouter(gameHandler *GameHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// API Routes
	mux.HandleFunc("/api/games", gameHandler.CreateGame) // POST to create
	
	// Poor man's routing for path parameters
	mux.HandleFunc("/api/games/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > 11 && path[:11] == "/api/games/" {
			// Check if it's /flip or /check or just get
			if len(path) > 5 && path[len(path)-5:] == "/flip" {
				gameHandler.FlipCard(w, r)
				return
			}
			if len(path) > 6 && path[len(path)-6:] == "/check" {
				gameHandler.CheckMatch(w, r)
				return
			}
			gameHandler.GetGame(w, r)
			return
		}
		http.NotFound(w, r)
	})

	// Static files
	fs := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Frontend page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		tmpl, err := template.ParseFiles("./web/templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	return mux
}
