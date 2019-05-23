package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
)

var testData = &Data{
	Tempreture: "20.0",
	Humidity:   "50.0",
	Longitude:  "31.1531",
	Latitude:   "133.3214",
}

var testPayload = Payload{
	Namespace: "123",
	Action:    "get",
	Data:      testData,
}

const marshaledData = `eyJUZW1wcmV0dXJlIjoiMTIuMTY0NzExIiwiSHVtaWRpdHkiOiI1MS4xMDAwIiwiTG9naXR1ZGUiOiIzOS4zNjE5IiwiTGF0aXR1ZGUiOiIxMTYuMTUxOSJ9`
const marshaledData2 = `eyJOYW1lc3BhY2UiOiI1YjczNDkiLCJBY3Rpb24iOiJzZXQiLCJEYXRhIjp7IlRlbXByZXR1cmUiOiIyMC4wIiwiSHVtaWRpdHkiOiI1MC4wIiwiTG9naXR1ZGUiOiIzMS4xNTMxIiwiTGF0aXR1ZGUiOiIxMzMuMzIxNCJ9fQ==`
const marshaledData3 = `ZXlKT1lXMWxjM0JoWTJVaU9pSXhNbVl4TWpZMlpUSTRNbU0wTVdKbE5XVTBNalUwWkRnNE1qQTNOekpqTlRVeE9HRXlZelZoT0dNd1l6ZG1OMlZrWVRFNU5UazBZVGRsWWpVek9UUTFNMlV4WldRM0lpd2lRV04wYVc5dUlqb2ljMlYwSWl3aVJHRjBZU0k2ZXlKVVpXMXdjbVYwZFhKbElqb2lNVEl1TUNJc0lraDFiV2xrYVhSNUlqb2lOVEF1TUNJc0lreHZaMmwwZFdSbElqb2lNekV1TVRVek1TSXNJa3hoZEdsMGRXUmxJam9pTVRNekxqTXlNVFFpZlgwPQ==`

func TestMarshall(t *testing.T) {
	dataS, _ := json.Marshal(testPayload)
	var test Payload
	json.Unmarshal(dataS, &test)
	fmt.Println(test.Namespace)
	fmt.Printf("%s\n", dataS)
}

func TestUnmarshal(t *testing.T) {
	var payload Payload
	fmt.Println([]byte(marshaledData))
	json.Unmarshal([]byte(marshaledData), &payload)
	fmt.Println("", payload)
}

func TestDecodeB64(t *testing.T) {
	// decodeB, _ := base64.StdEncoding.DecodeString(marshaledData)
	decodeB2, _ := base64.StdEncoding.DecodeString(marshaledData3)
	// fmt.Printf("%s\n", decodeB)
	fmt.Printf("%s\n", decodeB2)
}
