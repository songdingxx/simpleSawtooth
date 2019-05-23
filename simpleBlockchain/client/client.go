package client

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"../tp/handler"
	proto2 "github.com/golang/protobuf/proto"
	bat "github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	txn "github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
)

var testData = &handler.Data{
	Tempreture: "20.0",
	Humidity:   "50.0",
	Longitude:  "31.1531",
	Latitude:   "133.3214",
}

var testPayload = handler.Payload{
	Namespace: "1cf1266e282c41be5e4254d8820772c5518a2c5a8c0c7f7eda19594a7eb539453e1ed7",
	Action:    "get",
	Data:      testData,
}

const ValidatorUrl = `http://localhost:8008/batches`

func createSHA512(byte []byte) string {
	var result string
	h := sha512.New()
	h.Write(byte)
	result = hex.EncodeToString(h.Sum(nil))
	return result
}

func Submit(data *handler.Data) {
	payload, err := createPayload(data)
	if err != nil {
		fmt.Println(err)
	}
	batchesB := createBatches(payload)
	send(batchesB)
}

func createPayload(data *handler.Data) (handler.Payload, error) {
	var payload handler.Payload
	//Marshal the data
	dataB, err := json.Marshal(data)
	if err != nil {
		fmt.Println("(Create Payload)Marshal Error: ", err)
		return payload, errors.New("Marshal error")
	}
	dataHash := createSHA512(dataB)[:70]
	payload.Namespace = dataHash
	payload.Action = "set"
	payload.Data = data

	return payload, nil
}

func createBatches(payload handler.Payload) []byte {
	var batchesB []byte

	//Get the namespace
	ns := payload.Namespace
	// Marshal the payload
	payloadB, _ := json.Marshal(payload)
	// Set the secp256k1 context
	context := signing.NewSecp256k1Context()
	radPrivateKey := context.NewRandomPrivateKey()
	radPublicKey := context.GetPublicKey(radPrivateKey)
	// Set the crypto factory
	factory := signing.NewCryptoFactory(context)
	signer := factory.NewSigner(radPrivateKey)

	// Construct the transaction
	transactionHeaderB, err := proto2.Marshal(&txn.TransactionHeader{
		FamilyName:       "supplychain",
		FamilyVersion:    "1.0",
		SignerPublicKey:  radPublicKey.AsHex(),
		BatcherPublicKey: radPublicKey.AsHex(),
		Dependencies:     []string{},
		Inputs:           []string{ns},
		Outputs:          []string{ns},
		Nonce:            time.Now().String(),
		PayloadSha512:    createSHA512(payloadB),
	})
	if err != nil {
		fmt.Println(err)
		return batchesB
	}

	// Sign the header
	txnHeadSignature := hex.EncodeToString(signer.Sign(transactionHeaderB))

	// Construct the transaction
	transaction := txn.Transaction{
		Header:          transactionHeaderB,
		HeaderSignature: txnHeadSignature,
		Payload:         payloadB,
	}

	transactions := []*txn.Transaction{&transaction}

	// Create the batch header
	batchHeaderB, err := proto2.Marshal(&bat.BatchHeader{
		SignerPublicKey: signer.GetPublicKey().AsHex(),
		TransactionIds:  []string{transaction.HeaderSignature},
	})

	// Sign the batch header
	batchHeadSignature := hex.EncodeToString(signer.Sign(batchHeaderB))

	// Construct the batch
	batch := bat.Batch{
		Header:          batchHeaderB,
		HeaderSignature: batchHeadSignature,
		Transactions:    transactions,
	}

	batches := []*bat.Batch{&batch}
	batchList := bat.BatchList{
		Batches: batches,
	}
	batchesB, err = proto2.Marshal(&batchList)
	if err != nil {
		fmt.Println("Unable to marshal batchlist: ", err)
	}

	return batchesB
}

func send(batchlist []byte) error {
	response, err := http.Post(
		ValidatorUrl,
		"application/octet-stream",
		bytes.NewBuffer(batchlist),
	)
	if err != nil {
		fmt.Println("Response failed!\n", err)
		return errors.New("Unable to post the operation")
	}
	defer response.Body.Close()

	resultString, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Io error!\n", err)
	}
	fmt.Printf("Results: %s\n", resultString)
	return nil
}
