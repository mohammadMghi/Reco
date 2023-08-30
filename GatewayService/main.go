package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

 
)


func main(){
	// Define the target URL of the backend server
	targetURL, _ := url.Parse("http://localhost:8081") // Change this to your backend server's URL

	// Create a reverse proxy handler
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Create a custom handler function that adds any necessary headers before forwarding the request
	customHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/test" {
			// You can modify headers here if needed
			// r.Header.Add("X-Custom-Header", "Value")

			// Forward the request to the backend server
			proxy.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	}

	// Create an HTTP server using the custom handler
	http.HandleFunc("/", customHandler)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

 