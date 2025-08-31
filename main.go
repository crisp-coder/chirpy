package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("vim-go")
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app/", http.StatusFound)
	})
	handler.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("public"))))
	handler.Handle("/app/assets/", http.StripPrefix("/app", http.FileServer(http.Dir("public"))))
	handler.HandleFunc("/healthz", readinessHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		fmt.Println(err)
	}
}
