package main

import (
	"fmt"

	"log"
	"net/http"
)

/* Register a new education certificate */
func registerCert(w http.ResponseWriter, r *http.Request) {
	log.Println("To Implement!")
}

/* Retrieve a user's education certificate(s) */
func retrieveCert(w http.ResponseWriter, r *http.Request) {
	log.Println("To Implement!")
}

/* Check validity of a given education certificate */
func checkCert(w http.ResponseWriter, r *http.Request) {
	log.Println("To Implement!")
}

/* Create a private key from a seed phrase */
func createKey(w http.ResponseWriter, r *http.Request) {
	log.Println("To Implement!")
}

func main() {
	http.HandleFunc("/registerCert", registerCert)
	http.HandleFunc("/retrieveCert", retrieveCert)
	http.HandleFunc("/checkCert", checkCert)
	http.HandleFunc("/createKey", createKey)

	// Serve website
	port := "8080"
	fmt.Printf("Serving on localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
