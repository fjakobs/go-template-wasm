package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir := "./html" // Change to your directory
	fileServer := http.FileServer(http.Dir(dir))
	
	log.Printf("Serving %s on http://localhost:8080\n", dir)
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Always good practice to set 'Vary: Accept-Encoding'
		w.Header().Add("Vary", "Accept-Encoding")

		// Optionally disable caching (for testing/development):
		w.Header().Set("Cache-Control", "no-cache")

		// Check if the request is for a .wasm file
		if strings.HasSuffix(r.URL.Path, ".wasm") {
			w.Header().Set("Content-Type", "application/wasm")

			// Check for Brotli support
			if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
				brFilePath := filepath.Join(dir, r.URL.Path+".br")
				if fileExists(brFilePath) {
					// Serve the pre-compressed .wasm.br file
					w.Header().Set("Content-Encoding", "br")
					http.ServeFile(w, r, brFilePath)
					return
				}
			}
			// Fallback: serve uncompressed .wasm
		}

		// For all other files (or if .wasm.br is missing), use default file server
		fileServer.ServeHTTP(w, r)
	}))
}

// Helper function to check if a file exists
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
