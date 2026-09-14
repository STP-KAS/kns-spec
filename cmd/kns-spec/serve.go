package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/STP-KAS/kns-spec/internal/kns"
)

func serve() {
	addr := os.Getenv("KNS_SPEC_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8083"
	}
	docs := filepath.Dir(findFile("docs/index.html"))
	if _, err := os.Stat(filepath.Join(docs, "index.html")); err != nil {
		fail("docs/ not found; run from the kns-spec repo")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/resolve", handleResolve)
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"name":"kns-spec","listen":"` + addr + `"}`))
	})
	mux.Handle("/", http.FileServer(http.Dir(docs)))
	log.Printf("kns-spec companion  http://%s/open.html", addr)
	log.Printf("resolve             http://%s/api/resolve?q=kns.kas", addr)
	log.Printf("no seed. not a wallet. indexer uniqueness.")
	s := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 8 * time.Second}
	if err := s.ListenAndServe(); err != nil {
		fail(err.Error())
	}
}

func handleResolve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		http.Error(w, `{"error":"q required"}`, http.StatusBadRequest)
		return
	}
	c := kns.New("")
	s, err := c.Snapshot(q)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error(), "q": q})
		return
	}
	_ = json.NewEncoder(w).Encode(s)
}
