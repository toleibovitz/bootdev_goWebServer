package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

func handlerHealthCheck(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader((http.StatusOK))
		w.Write([]byte("OK"))
	}



type apiConfig struct {
		fileserverHits atomic.Int32
	}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			cfg.fileserverHits.Add(1)
			next.ServeHTTP(w, req)

		})
	}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	hits := cfg.fileserverHits.Load()
	fmt.Fprintf(w, "Hits: %d", hits)
	
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
}


func main() {
	
	mux := http.NewServeMux()
	apiCfg := &apiConfig{}
	
	fs := http.FileServer(http.Dir("."))
	handler := apiCfg.middlewareMetricsInc((fs))
	mux.Handle("/app/", http.StripPrefix("/app", handler))
	
	
	mux.HandleFunc("/healthz", handlerHealthCheck)
	mux.HandleFunc("/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("/reset", apiCfg.handlerReset)
	
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}