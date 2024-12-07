package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("starting the server")

	server()
}

func server() {
	http.HandleFunc("/health/", health)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Printf("failed to start server: %v", err)

		return
	}
}

func health(w http.ResponseWriter, req *http.Request) {
	fmt.Println("processing health request")

	response, err := json.Marshal(map[string]string{"status": "OK"})
	if err != nil {
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)

		return
	}

	_, err = w.Write(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
