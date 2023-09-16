package main

import (
	"encoding/json"
	"errors"
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
	err := parseRequestJson(r, &data, "registerCert")
	if err != nil {
		http.Error(w, "Unable to parse request body as JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Parse user data
	userId, err := hedera.AccountIDFromString(data.AccId)
	if err != nil {
		http.Error(w, "Unable to parse account ID: "+err.Error(), http.StatusBadRequest)
		return
	}
	userKey, err := hedera.PrivateKeyFromString(data.PrivKey)
	if err != nil {
		http.Error(w, "Unable to parse private key: "+err.Error(), http.StatusBadRequest)
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

	if associateStatus != hedera.StatusSuccess || transferStatus != hedera.StatusSuccess {
		http.Error(w, "Did not manage to associate OR transfer NFT", http.StatusInternalServerError)
		return
	}
}

type retrieveCertStruct struct {
	PubKey string
}

/* Retrieve a user's education certificate(s) */
func retrieveCert(w http.ResponseWriter, r *http.Request) {
	var data retrieveCertStruct
	err := parseRequestJson(r, &data, "retrieveCert")
	if err != nil {
		http.Error(w, "Unable to parse request body as JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("To Implement!")
}

/* Check validity of a given education certificate */
func checkCert(w http.ResponseWriter, r *http.Request) {
	var data registerCertStruct
	err := parseRequestJson(r, &data, "checkCert")
	if err != nil {
		http.Error(w, "Unable to parse request body as JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("To Implement!")
}

type responseAccountStruct struct {
	AccId   string
	PrivKey string
	PubKey  string
}

/* Create a private key from a seed phrase */
func createKey(w http.ResponseWriter, r *http.Request) {
	privateKey, err := hedera.PrivateKeyGenerateEd25519()
	if err != nil {
		http.Error(w, "Private key could not be created: "+err.Error(), http.StatusBadRequest)
		return
	}
	publicKey := privateKey.PublicKey()
	log.Printf("Created account. \nPUBLIC key: %v\nPRIVATE key: %v\n", publicKey, privateKey)

	// Create new account and assign the public key
	newAccount, err := hedera.NewAccountCreateTransaction().
		SetKey(publicKey).
		SetInitialBalance(hedera.HbarFrom(1000, hedera.HbarUnits.Tinybar)).
		Execute(&client)

	// Request the receipt of the transaction
	receipt, err := newAccount.GetReceipt(&client)
	if err != nil {
		http.Error(w, "Transaction receipt could not be requested: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the new account ID from the receipt
	accountId := *receipt.AccountID

	// Log the account ID
	fmt.Printf("The new account ID is %v\n", accountId)

	// Prepare response
	w.Header().Set("Content-Type", "application/json")
	response := responseAccountStruct{
		AccId:   accountId.String(),
		PrivKey: privateKey.String(),
		PubKey:  publicKey.String(),
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Could not encode new account data as JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

/* ----- Logic ----- */

func parseRequestJson(r *http.Request, v any, funcName string) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&v)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not decode %s as JSON for %s", r.Body, funcName))
	}
	log.Printf("%s got: %s\n", funcName, v)
	return nil
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
		log.Fatal("[createCertNft] Unable to create concrete certificate NFT. Error:\n%v\n", err)
	}

	// Get the transaction receipt
	mintRx, err := mintTxSubmit.GetReceipt(&client)
	if err != nil {
		log.Fatal("[createCertNft] Unable to get transaction receipt. Error:\n%v\n", err)
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
		log.Fatal("[associateCertNft] Unable to associate certificate to user. Error:\n%v\n", err)
	}

	// Get the transaction receipt
	associateUserRx, err := associateUserTxSubmit.GetReceipt(&client)
	if err != nil {
		log.Fatal("[associateCertNft] Unable to get transaction receipt. Error:\n%v\n", err)
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
		log.Fatal("[transferCertNft] Unable to submit transaction. Error:\n%v\n", err)
	}

	//Get the transaction receipt
	tokenTransferRx, err := tokenTransferSubmit.GetReceipt(&client)
	if err != nil {
		log.Fatal("[transferCertNft] Unable to get transaction receipt. Error:\n%v\n", err)
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
	treasuryId, err = hedera.AccountIDFromString(os.Getenv("HEDERA_ACCOUNT_ID"))
	if err != nil {
		log.Fatal(err)
	}

	treasuryKey, err = hedera.PrivateKeyFromString(os.Getenv("HEDERA_PRIVATE_KEY"))
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
	certificateBaseNftTokenId = createCertBaseNft()
	fmt.Println("Created NFT with token ID ", certificateBaseNftTokenId)

	// Serve website
	port := "8080"
	fmt.Printf("Serving on http://localhost:%s\n\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
