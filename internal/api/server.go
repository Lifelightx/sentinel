package api

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
	
	"sentinel/internal/store"
)

func Start(addr string, mem *store.MemoryStore) error {

	tmpl := template.Must(template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/dashboard.html",
		"web/templates/server.html",
	))

	/* JSON APIs */

	http.HandleFunc("/api/servers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mem.GetAll())
	})

	http.HandleFunc("/api/containers", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("ServerID")
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mem.GetContainers(id))
	})

	http.HandleFunc("/api/servers/alerts", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("ServerID")
		w.Header().Set("Content-Type", "application/json")
		server, ok := mem.GetByID(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(server)
	})

	/* Dashboard */

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		tmpl.ExecuteTemplate(w, "layout", map[string]any{
			"Page": "dashboard",
			"Data": nil,
		})
	})

	/* Server Details */

	http.HandleFunc("/server/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/server/")

		server, ok := mem.GetByID(id)
		if !ok {
			http.NotFound(w, r)
			return
		}

		tmpl.ExecuteTemplate(w, "layout", map[string]any{
			"Page": "server",
			"Data": server,
		})
	})



	/* Static CSS */

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	return http.ListenAndServe(addr, nil)
}