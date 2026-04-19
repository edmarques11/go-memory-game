package http

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/edmarqueslima/memorygame/web"
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
	staticFS, err := fs.Sub(web.FS, "static")
	if err != nil {
		panic(err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Frontend page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		tmpl, err := template.ParseFS(web.FS, "templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	return mux
}
