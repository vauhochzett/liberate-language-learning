package main

import (
	"encoding/json"
	"fmt"

	"os"

	"log"
	"net/http"

	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/joho/godotenv"
)

/* ----- Request Handlers ----- */

type registerCertStruct struct {
	PubKey string
	CertId string
}

/* Register a new education certificate */
func registerCert(w http.ResponseWriter, r *http.Request) {
	var data registerCertStruct
	parseRequestJson(r, &data, "registerCert")

	log.Println("To Implement!")
}

type retrieveCertStruct struct {
	PubKey string
}

/* Retrieve a user's education certificate(s) */
func retrieveCert(w http.ResponseWriter, r *http.Request) {
	var data retrieveCertStruct
	parseRequestJson(r, &data, "retrieveCert")

	log.Println("To Implement!")
}

/* Check validity of a given education certificate */
func checkCert(w http.ResponseWriter, r *http.Request) {
	var data registerCertStruct
	parseRequestJson(r, &data, "checkCert")

	log.Println("To Implement!")
}

/* Create a private key from a seed phrase */
func createKey(w http.ResponseWriter, r *http.Request) {
	log.Println("To Implement!")
}

/* ----- Logic ----- */

func parseRequestJson(r *http.Request, v any, funcName string) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&v)
	if err != nil {
		log.Printf("Could not decode %s as JSON for %s", r.Body, funcName)
	}
	log.Printf("%s got: %s\n", funcName, v)
}

func connect() {
	// Load the .env file and throws an error if it cannot load the variables from that file correctly
	err := godotenv.Load(".env")
	if err != nil {
		panic(fmt.Errorf("Unable to load environment variables from .env file. Error:\n%v\n", err))
	}

	// Grab testnet account ID and private key from the .env file
	myAccountId, err := hedera.AccountIDFromString(os.Getenv("HEDERA_ACCOUNT_ID"))
	if err != nil {
		panic(err)
	}

	myPrivateKey, err := hedera.PrivateKeyFromString(os.Getenv("HEDERA_PRIVATE_KEY"))
	if err != nil {
		panic(err)
	}

	// Print your testnet account ID and private key to the console to make sure there was no error
	fmt.Printf("The account ID is = %v\n", myAccountId)
	fmt.Printf("The private key is = %v\n", myPrivateKey)

	// Create your testnet client
	client := hedera.ClientForTestnet()
	client.SetOperator(myAccountId, myPrivateKey)

	// Set default max transaction fee
	client.SetDefaultMaxTransactionFee(hedera.HbarFrom(100, hedera.HbarUnits.Hbar))

	// Set max query payment
	client.SetDefaultMaxQueryPayment(hedera.HbarFrom(50, hedera.HbarUnits.Hbar))
}

/* ----- Main ----- */

func main() {
	http.HandleFunc("/registerCert", registerCert)
	http.HandleFunc("/retrieveCert", retrieveCert)
	http.HandleFunc("/checkCert", checkCert)
	http.HandleFunc("/createKey", createKey)

	// Serve website
	port := "8080"
	fmt.Printf("Serving on http://localhost:%s\n\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
