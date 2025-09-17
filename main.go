package main

import (
	"net/http"
	
)




func main() {
	handler := *http.NewServeMux()

	server := &http.Server{
		Addr: ":8080",
		Handler: &handler,
	}
	handler.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	

	handler.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader((http.StatusOK))
		
		w.Write([]byte("OK"))
	})

	server.ListenAndServe()
}