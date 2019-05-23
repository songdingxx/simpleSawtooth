package handler

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
)

type TestHandler struct {
	namespace string
}

type Data struct {
	Tempreture string
	Humidity   string
	Longitude  string
	Latitude   string
}

type Payload struct {
	Namespace string
	Action    string
	Data      *Data
}

const (
	MIN_VALUE       = 0
	MAX_VALUE       = 4294967295
	MAX_NAME_LENGTH = 20
	FAMILY_NAME     = "supplychain"
)

const testAddress = "1cf1266e282c41be5e4254d8820772c5518a2c5a8c0c7f7eda19594a7eb539453e1ed7"

func (self *TestHandler) FamilyName() string {
	return FAMILY_NAME
}

func (self *TestHandler) FamilyVersions() []string {
	return []string{"1.0"}
}

func (self *TestHandler) Namespaces() []string {
	return []string{self.namespace}
}

func (*TestHandler) Apply(request *processor_pb2.TpProcessRequest, context *processor.Context) error {
	// Get the payload []bytes
	PayloadData := request.GetPayload()
	if PayloadData == nil {
		return &processor.InvalidTransactionError{
			Msg: "Must contain payload",
		}
	}

	// Unmarshal the data
	var payload Payload
	err := json.Unmarshal(PayloadData, &payload)
	if err != nil {
		fmt.Println(err)
		return &processor.InvalidTransactionError{
			Msg: "Unable to unmarshal the payload",
		}
	}

	// Check the action
	action := payload.Action
	if !(action == "get" || action == "set") {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Invalid action: %v", action)}
	}
	switch action {
	case "get":
		GetValue(payload.Namespace, context)
		break
	case "set":
		SetValue(payload, context)
		break
	}
	return nil
}

func SetValue(payload Payload, context *processor.Context) error {
	// TODO: define
	dataB, err := json.Marshal(payload.Data)
	if err != nil {
		fmt.Println("(Set value)Marshal Error: ", err)
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Setvalue marshal error: %s", err)}
	}
	addresses, err := context.SetState(map[string][]byte{
		payload.Namespace: dataB,
	})
	if err != nil {
		fmt.Println("(Set value)Setstate Error: ", err)
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("Setvalue setstate error: %s", err)}
	}
	fmt.Println(payload.Data.Latitude)
	fmt.Println(addresses)
	return nil
}

func GetValue(address string, context *processor.Context) (Data, error) {
	var data Data
	results, err := context.GetState([]string{address})

	if err != nil {
		fmt.Println("(Get value error)GetState Error: ", err)
		return data, &processor.InvalidTransactionError{Msg: fmt.Sprintf("Setvalue getstate error: %s", err)}
	}

	if len(results[address]) > 0 {
		err := json.Unmarshal(results[address], &data)
		if err != nil {
			fmt.Println("(Get value error)Marshal Error: ", err)
		}
	}
	fmt.Println(data.Humidity)
	return data, nil
}

func CreateHandler(namespace string) *TestHandler {
	return &TestHandler{
		namespace: namespace,
	}
}
