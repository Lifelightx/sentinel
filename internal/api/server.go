package api

import (
	"encoding/json"
	"net/http"
	"sentinel/internal/store"
)

func Start(addr string, mem *store.MemoryStore) error{


	http.HandleFunc("/api/servers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(mem.GetAll())
	})

	http.HandleFunc("/api/containers", func(w http.ResponseWriter, r *http.Request) {
		serverId := r.URL.Query().Get("serverId")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mem.GetContainers(serverId))
	})

	//serving static files
	http.Handle("/", http.FileServer(http.Dir("./web")))
	
	return  http.ListenAndServe(addr, nil)


}



