package api

import (
	"encoding/json"
	"html/template"
	
	"net/http"
	"strings"
	"time"

	"sentinel/internal/broker"
	"sentinel/internal/models"
	"sentinel/internal/store"
)

func Start(addr string, mem *store.MemoryStore, client *broker.Client) error {

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

	// action handler
	http.HandleFunc("/api/actions", func(w http.ResponseWriter, r *http.Request) {
		
		if r.Method != http.MethodPost{
			http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
			return 
			
		}
		var req models.CommandRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil{
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return 
		}
		//Basic validation
		if req.Action == "" || req.ContainerID == "" || req.HostName == ""{
			http.Error(w, "missing fields", http.StatusBadRequest)
			return 
		}
		req.ReplyTo = "reply."+ req.HostName
		// log.Println(req)
		err = client.Publish("commands."+req.HostName, req)
		if err != nil{
			http.Error(w, "Failed to send commands", http.StatusInternalServerError)
			return 
		}
		w.WriteHeader(http.StatusAccepted)
	})
	/* action result*/
	http.HandleFunc("/api/result", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet{
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return 
		}
		hostname := r.URL.Query().Get("hostname")
		containerID := r.URL.Query().Get("containerId")
		action := r.URL.Query().Get("action")
		res, ok := mem.GetCommandResult(hostname, containerID, action)
		// log.Println(res)
		if !ok {
			http.Error(w, "result not ready", http.StatusNotFound)
			return 
		}
		if time.Now().Unix() - res.Timestamp > 60{
			http.Error(w, "result expired", http.StatusNotFound)
			return 
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)


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