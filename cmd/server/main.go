package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ABAlosaimi/ParseInsight/internal/api"
)

func main() {
	handler := api.NewHandler()

	// Setup routes
	http.HandleFunc("/", handler.HandleIndex)
	http.HandleFunc("/api/benchmark", handler.HandleBenchmark)
	http.HandleFunc("/api/libraries", handler.HandleLibraries)

	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start server
	port := "8080"
	fmt.Printf("🚀 ParseInsight server starting on http://localhost:%s\n", port)
	fmt.Println("📊 Ready to benchmark HTTP parsers!")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
