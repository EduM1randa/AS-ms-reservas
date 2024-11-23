package main

import (
	"fmt"
	"log"
	"net/http"

	"ms-reservas/controllers"
	"ms-reservas/database"
)

func main() {
	client := database.ConnectMongoDB()

	http.HandleFunc("/", helloHandler)
	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	controllers.SetMongoClient(client)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
