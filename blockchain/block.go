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

/* ----- Globals ----- */

var client hedera.Client
var treasuryId hedera.AccountID
var treasuryKey hedera.PrivateKey
var certificateBaseNftTokenId hedera.TokenID
var NFT_EnglishDailyB2_CID = []byte("ipfs://bafybeihxnvasdek52refjxoarbltzghbrooma7abpdetcofqjimbrhfpw4")

/* ----- Request Handlers ----- */

type registerCertStruct struct {
	AccId   string
	PrivKey string
	CertId  string
}

/* Register a new education certificate */
func registerCert(w http.ResponseWriter, r *http.Request) {
	var data registerCertStruct
	parseRequestJson(r, &data, "registerCert")

	// Parse user data
	userId, err := hedera.AccountIDFromString(data.AccId)
	if err != nil {
		http.Error(w, "Unable to parse account ID: "+err.Error(), http.StatusInternalServerError)
		return
	}
	userKey, err := hedera.PrivateKeyFromString(data.PrivKey)
	if err != nil {
		http.Error(w, "Unable to parse private key: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Mint certificate NFT
	serials := createCertNft(certificateBaseNftTokenId, NFT_EnglishDailyB2_CID)
	log.Printf("Created certificate of NFT %s with serial(s): %d\n", certificateBaseNftTokenId, serials)
	if len(serials) != 1 {
		log.Fatal("Multiple serial numbers received after concrete certificate NFT creation!")
	}
	serial := serials[0]

	// Associate and transfer NFT
	associateStatus := associateCertNft(userId, userKey, certificateBaseNftTokenId)
	transferStatus := transferCertNft(certificateBaseNftTokenId, serial, userId)
	log.Println("NFT association with Alice's account:", associateStatus)
	log.Println("NFT transfer from Treasury to User:", transferStatus)

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

type createKeyStruct struct {
	Passphrase string
}

/* Create a private key from a seed phrase */
func createKey(w http.ResponseWriter, r *http.Request) {
	var data registerCertStruct
	parseRequestJson(r, &data, "createKey")

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

func createCertBaseNft() hedera.TokenID {
	// Create the NFT
	nftCreate := hedera.NewTokenCreateTransaction().
		SetTokenName("liberatedLanguageLearner_CERT").
		SetTokenSymbol("LLL CERTIFICATE").
		SetTreasuryAccountID(treasuryId).
		SetAdminKey(treasuryKey).
		SetSupplyKey(treasuryKey).
		SetTokenType(hedera.TokenTypeNonFungibleUnique).
		SetInitialSupply(0)

	// Sign the transaction with the treasury key
	nftCreateTxSign := nftCreate.Sign(treasuryKey)

	// Submit the transaction to a Hedera network
	nftCreateSubmit, err := nftCreateTxSign.Execute(&client)
	if err != nil {
		log.Fatal("Unable to create base certificate NFT. Error:\n%v\n", err)
	}

	// Get the transaction receipt
	nftCreateRx, err := nftCreateSubmit.GetReceipt(&client)
	if err != nil {
		log.Fatal("Unable to get transaction receipt. Error:\n%v\n", err)
	}

	// Return the token ID
	return *nftCreateRx.TokenID
}

func createCertNft(tokenId hedera.TokenID, CID []byte) []int64 {

	// Mint new NFT
	mintTx := hedera.NewTokenMintTransaction().
		SetTokenID(tokenId).
		SetMetadata(CID)

	// Sign the transaction with the supply key
	mintTxSign := mintTx.Sign(treasuryKey)

	// Submit the transaction to a Hedera network
	mintTxSubmit, err := mintTxSign.Execute(&client)
	if err != nil {
		log.Fatal("Unable to create concrete certificate NFT. Error:\n%v\n", err)
	}

	// Get the transaction receipt
	mintRx, err := mintTxSubmit.GetReceipt(&client)
	if err != nil {
		log.Fatal("Unable to get transaction receipt. Error:\n%v\n", err)
	}

	return mintRx.SerialNumbers
}

func associateCertNft(userAccountId hedera.AccountID, userAccountKey hedera.PrivateKey, tokenId hedera.TokenID) hedera.Status {
	// Create the associate transaction
	associateAliceTx := hedera.NewTokenAssociateTransaction().
		SetAccountID(userAccountId).
		SetTokenIDs(tokenId)

	//Sign with Alice's key
	signTx := associateAliceTx.Sign(userAccountKey)

	// Submit the transaction to a Hedera network
	associateUserTxSubmit, err := signTx.Execute(&client)
	if err != nil {
		log.Fatal("Unable to associate certificate to user. Error:\n%v\n", err)
	}

	// Get the transaction receipt
	associateUserRx, err := associateUserTxSubmit.GetReceipt(&client)
	if err != nil {
		log.Fatal("Unable to get transaction receipt. Error:\n%v\n", err)
	}

	// Return transaction status
	return associateUserRx.Status
}

func transferCertNft(tokenId hedera.TokenID, serial int64, userId hedera.AccountID) hedera.Status {
	// Transfer the NFT from treasury to user
	tokenTransferTx := hedera.NewTransferTransaction().
		AddNftTransfer(hedera.NftID{TokenID: tokenId, SerialNumber: serial}, treasuryId, userId)

	// Sign with the treasury key to authorize the transfer
	signTransferTx := tokenTransferTx.Sign(treasuryKey)

	// Submit the transaction
	tokenTransferSubmit, err := signTransferTx.Execute(&client)
	if err != nil {
		log.Fatal("Unable to submit transaction. Error:\n%v\n", err)
	}

	//Get the transaction receipt
	tokenTransferRx, err := tokenTransferSubmit.GetReceipt(&client)
	if err != nil {
		log.Fatal("Unable to get transaction receipt. Error:\n%v\n", err)
	}

	// Log the transaction status
	return tokenTransferRx.Status
}

/* ----- Main ----- */

func main() {
	// Routes
	http.HandleFunc("/registerCert", registerCert)
	http.HandleFunc("/retrieveCert", retrieveCert)
	http.HandleFunc("/checkCert", checkCert)
	http.HandleFunc("/createKey", createKey)

	// Load the .env file and throw an error if it cannot load the variables from that file correctly
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Unable to load environment variables from .env file. Error:\n%v\n", err)
	}

	// Grab testnet account ID and private key from the .env file
	treasuryId, err := hedera.AccountIDFromString(os.Getenv("HEDERA_ACCOUNT_ID"))
	if err != nil {
		log.Fatal(err)
	}

	treasuryKey, err := hedera.PrivateKeyFromString(os.Getenv("HEDERA_PRIVATE_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	// Print your testnet account ID and private key to the console to make sure there was no error
	fmt.Printf("The account ID is = %v\n", treasuryId)
	fmt.Printf("The private key is = %v\n", treasuryKey)

	// Create testnet client and configure
	client = *hedera.ClientForTestnet()
	client.SetOperator(treasuryId, treasuryKey)

	// Create base NFT for all certificate NFTs
	tokenId := createCertBaseNft()
	fmt.Println("Created NFT with token ID ", tokenId)

	// Serve website
	port := "8080"
	fmt.Printf("Serving on http://localhost:%s\n\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
