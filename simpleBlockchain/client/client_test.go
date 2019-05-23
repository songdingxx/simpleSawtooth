package client

import (
	"encoding/json"
	"fmt"
	"testing"

	"../tp/handler"
)

var TestData = handler.Data{
	Tempreture: "12.0",
	Humidity:   "50.0",
	Longitude:  "31.1531",
	Latitude:   "133.3214",
}

var TestPayload = handler.Payload{
	Namespace: "12f1266e282c41be5e4254d8820772c5518a2c5a8c0c7f7eda19594a7eb539453e1ed7",
	Action:    "set",
	Data:      &TestData,
}

var marshaled = `eyJUZW1wcmV0dXJlIjoiMTIuMCIsIkh1bWlkaXR5IjoiNTAuMCIsIkxvZ2l0dWRlIjoiMzEuMTUzMSIsIkxhdGl0dWRlIjoiMTMzLjMyMTQifQ`

func TestSubmit(t *testing.T) {
	marshal, _ := json.Marshal(TestData)
	fmt.Println(string(marshal))
	fmt.Println([]byte(marshaled))
	// Submit(&TestData)
}

func TestSendRequest(t *testing.T) {
	batchesB := createBatches(TestPayload)
	send(batchesB)
}
