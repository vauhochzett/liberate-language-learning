package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"log"
	"net/http"

	"bytes"
	"net/url"

	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/joho/godotenv"
)

/* ----- Globals ----- */

var client hedera.Client
var treasuryId hedera.AccountID
var treasuryKey hedera.PrivateKey
var certificateBaseNftTokenId hedera.TokenID
var NFT_ONLY_ID = "bafybeihxnvasdek52refjxoarbltzghbrooma7abpdetcofqjimbrhfpw4"
var NFT_EnglishDailyB2_CID = []byte("ipfs://" + NFT_ONLY_ID)
var translateKey string
var progressCounter = make(map[string]int)

/* ----- Request Handlers ----- */

type registerCertStruct struct {
	AccId   string
	PrivKey string
	CertId  string
}

/* Register a new education certificate */
func registerCert(w http.ResponseWriter, r *http.Request) {
	var data registerCertStruct
	err := parseBodyJson(r.Body, &data, "registerCert")
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
	log.Println("NFT association with User's account:", associateStatus)
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
	err := parseBodyJson(r.Body, &data, "retrieveCert")
	if err != nil {
		http.Error(w, "Unable to parse request body as JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	http.Error(w, "Method not implemented!", http.StatusInternalServerError)
}

type checkCertStruct struct {
	AccId  string
	CertId string
	Serial string
}

type responseCheckStruct struct {
	Valid bool
}

/* Check validity of a given education certificate */
func checkCert(w http.ResponseWriter, r *http.Request) {
	var data checkCertStruct
	err := parseBodyJson(r.Body, &data, "checkCert")
	if err != nil {
		http.Error(w, "Unable to parse request body as JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	accId, err := hedera.AccountIDFromString(data.AccId)
	if data.AccId == "" {
		http.Error(w, "Unable to parse \""+data.AccId+"\" as Account ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	nftIdStr := fmt.Sprintf("%s@%s", data.Serial, data.CertId)
	nftId, err := hedera.NftIDFromString(nftIdStr)
	if err != nil {
		http.Error(w, "Unable to parse \""+data.CertId+"\" as NFT ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	nftInfos, err := hedera.NewTokenNftInfoQuery().
		SetNftID(nftId).
		Execute(&client)

	if len(nftInfos) > 1 {
		http.Error(w, fmt.Sprintf("Got more than one NFT info: %s", nftInfos), http.StatusBadRequest)
		return
	}

	valid := false

	fmt.Printf("nftInfos: %s\n", nftInfos)

	if len(nftInfos) == 1 {
		fmt.Printf("accID=%s\nnftInfos[0].AccountID=%s", accId, nftInfos[0].AccountID)
		valid = accId == nftInfos[0].AccountID
	}

	// Prepare response
	w.Header().Set("Content-Type", "application/json")
	response := responseCheckStruct{
		Valid: valid,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Could not encode check result as JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	not := ""
	if !valid {
		not = "not "
	}
	log.Printf("Account ID does %sown the given certificate NFT", not)
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

type verifyWordStruct struct {
	AccId            string
	OriginalString   string
	TranslatedString string
	Language         string
}

type responseWordStruct struct {
	Correct     bool
	CorrectWord string
	Certificate string
}

/* Verify the correctness of the word */
func verifyWord(w http.ResponseWriter, r *http.Request) {
	var data verifyWordStruct
	parseBodyJson(r.Body, &data, "verifyWord")

	if data.OriginalString == "" || data.TranslatedString == "" || data.Language == "" {
		http.Error(w, "Missing required request parameter. Expecting: OriginalString, TranslatedString, Language", http.StatusBadRequest)
		return
	}

	// Parse only to check if a correct account ID was sent
	_, err := hedera.AccountIDFromString(data.AccId)
	if data.AccId == "" {
		http.Error(w, "Unable to parse \""+data.AccId+"\" as Account ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	correct, correctWord, err := verifyWordAzure(data.OriginalString, data.TranslatedString, data.Language)
	if err != nil {
		http.Error(w, "Error on word verification: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if correct {
		progressCounter[data.AccId] += 1
	}

	certificate := ""
	if progressCounter[data.AccId] > 0 && progressCounter[data.AccId]%5 == 0 {
		certificate = NFT_ONLY_ID
	}
	log.Printf("Word correct: %t\nTotal correct words of user: %d\nCertificate:%s", correct, progressCounter[data.AccId], certificate)

	// Prepare response
	w.Header().Set("Content-Type", "application/json")
	response := responseWordStruct{
		Correct:     correct,
		CorrectWord: correctWord,
		Certificate: certificate,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Could not encode check result as JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

/* ----- Logic ----- */

func parseBodyJson(body io.ReadCloser, v any, funcName string) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&v)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not decode %s as JSON for %s", body, funcName))
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
	associateUserTx := hedera.NewTokenAssociateTransaction().
		SetAccountID(userAccountId).
		SetTokenIDs(tokenId)

	//Sign with user's key
	signTx := associateUserTx.Sign(userAccountKey)

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

type translationStruct struct {
	Text string `json:"text"`
	To   string `json:"to"`
}

type translationResultStruct struct {
	Translations []translationStruct `json:"translations"`
}

func verifyWordAzure(originalString string, translatedString string, language string) (bool, string, error) {
	if language != "fr" && language != "de" && language != "es" {
		return false, "", errors.New("Unsupported language string " + language)
	}

	endpoint := "https://api.cognitive.microsofttranslator.com/"
	uri := endpoint + "translate?api-version=3.0"
	location := "westeurope"

	// Build the request URL. See: https://go.dev/pkg/net/url/#example_URL_Parse
	u, _ := url.Parse(uri)
	q := u.Query()
	q.Add("from", "en")
	q.Add("to", language)
	u.RawQuery = q.Encode()

	// Create an anonymous struct for your request body and encode it to JSON
	body := []struct {
		Text string
	}{
		{Text: originalString},
	}
	b, _ := json.Marshal(body)

	// Build the HTTP POST request
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(b))
	if err != nil {
		return false, "", err
	}
	// Add required headers to the request
	req.Header.Add("Ocp-Apim-Subscription-Key", translateKey)
	req.Header.Add("Ocp-Apim-Subscription-Region", location)
	req.Header.Add("Content-Type", "application/json")

	// Call the Translator API
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, "", err
	}

	// Decode the JSON response
	var data []translationResultStruct
	err = parseBodyJson(res.Body, &data, "verifyWordAzure")
	if err != nil {
		return false, "", errors.New(fmt.Sprintf("Unable to parse translation response \"%s\" as JSON: "+err.Error(), res.Body))
	}

	if len(data) != 1 || len(data[0].Translations) != 1 {
		return false, "", errors.New(fmt.Sprintf("Received more or less than one translation result! See: %s", data))
	}

	translation := data[0].Translations[0]
	correct := translation.Text == translatedString
	return correct, translation.Text, nil
}

/* ----- Main ----- */

func main() {
	// Routes
	http.HandleFunc("/registerCert", registerCert)
	http.HandleFunc("/retrieveCert", retrieveCert)
	http.HandleFunc("/checkCert", checkCert)
	http.HandleFunc("/createKey", createKey)
	http.HandleFunc("/verifyWord", verifyWord)

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

	translateKey = os.Getenv("AZURE_TRANSLATE_KEY")

	// Print your testnet account ID and private key to the console to make sure there was no error
	log.Printf("The treasury account ID is = %v\n", treasuryId)
	log.Printf("The treasury private key is = %v\n", treasuryKey)
	log.Printf("The Azure translate key is = %s\n", translateKey)

	// Create testnet client and configure
	client = *hedera.ClientForTestnet()
	client.SetOperator(treasuryId, treasuryKey)

	// Create base NFT for all certificate NFTs
	certificateBaseNftTokenId = createCertBaseNft()
	log.Println("Created NFT with token ID ", certificateBaseNftTokenId)

	// Serve website
	port := "8080"
	log.Printf("Serving on http://localhost:%s\n\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
