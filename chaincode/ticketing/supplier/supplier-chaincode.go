package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	//"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	//"github.com/hyperledger/fabric-protos-go"
)

type SmartContract struct {
}
type Event struct {
	ID           int       `json:"id"`
	Issuer       string    `json:"issuer"`
	Price        string    `json:"price"`
	EventName    string    `json:"eventName"`
	Total        int       `json:"total"`
	Sold         int       `json:"sold"`
	CreatedTime  time.Time `json:"createdTime"`
	RedeemedTime time.Time `json:"redeemedTime"`
}
type Key struct{
	KeyEncrypted string `json:"keyEncrypted"`
	Status	bool `json:"status"`
}
type Info struct {
	Number int `json:"number"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

var logger = shim.NewLogger("ticketing")

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	logger.Infof("start ")
	function, args := APIstub.GetFunctionAndParameters()

	info := Info{Number: 5}
	infoAsBytes, _ := json.Marshal(info)
	APIstub.PutState("NUM", infoAsBytes)

	if function == "initEvent" {
		logger.Infof("through init Event")
		return s.initEvent(APIstub)
	} else if function == "queryAllEvent" {
		return s.queryAllEvent(APIstub)
	} else if function == "createEvent" {
		return s.createEvent(APIstub, args)
	} else if function == "redeemTicket" {
		return s.redeemTicket(APIstub, args)
	} else if function == "checkoutTicket" {
		return s.checkoutTicket(APIstub, args)
	}
	return shim.Error("Wrong function name.")
}
func (s *SmartContract) initEvent(APIstub shim.ChaincodeStubInterface) sc.Response {
	logger.Infof("done1 ")

	events := []Event{
		Event{ID: 0, Issuer: "VFF", Price: "220.000", EventName: "Suzuki cup", Total: 20, Sold: 0, time.Now(), time.Now()},
		Event{ID: 1, Issuer: "BFF", Price: "220.000", EventName: "B cup", Total: 20, Sold: 0, time.Now(), time.Now()},
		Event{ID: 2, Issuer: "CFF", Price: "220.000", EventName: "C cup", Total: 20, Sold: 0, time.Now(), time.Now()},
		Event{ID: 3, Issuer: "DFF", Price: "220.000", EventName: "D cup", Total: 20, Sold: 0, time.Now(), time.Now()},
		Event{ID: 4, Issuer: "EFF", Price: "220.000", EventName: "F cup", Total: 20, Sold: 0, time.Now(), time.Now()}}

	return shim.Success(nil)
}

func (s *SmartContract) createEvent(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	numberAsBytes, _ := APIstub.GetState("NUM")
	str := string(numberAsBytes)
	logger.Infof("number as byte:" + str)
	fmt.Printf("number as byte:" + str)
	if numberAsBytes == nil {
		return shim.Error("Could not locate number of events")
	}
	info := Info{}

	json.Unmarshal(numberAsBytes, &info)

	fmt.Printf(strconv.Itoa(info.Number))
	total, _ := strconv.Atoi(args[3])
	createdTime, _ := time.Parse(time.RFC3339, args[4])
	redeemedTime, _ := time.Parse(time.RFC3339, args[5])

	var event = Event{ID: info.Number, Issuer: args[0], Price: args[1], EventName: args[2], Total: total, Sold: 0, createdTime, redeemedTime}
	eventAsBytes, _ := json.Marshal(event)
	APIstub.PutState("EVENT"+strconv.Itoa(event.ID), eventAsBytes)

	info.Number++
	logger.Infof(strconv.Itoa(info.Number))
	fmt.Printf(strconv.Itoa(info.Number))
	numberAsBytes, _ = json.Marshal(info)
	err := APIstub.PutState("NUM", numberAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to add num: "))
	}
	return shim.Success(nil)
}
func (s *SmartContract) triggeredBuyTicket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//0: eventId, 1: num
	fmt.Printf(args[0])

	eventAsBytes, _ := APIstub.GetState(args[0])
	var event = Event{}
	json.Unmarshal(eventAsBytes, &event)
	event.Sold++	
	eventAsBytes, _ = json.Marshal(event)
	APIstub.PutState("EVENT"+strconv.Itoa(event.ID), eventAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) redeemTicket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Printf(args[0])
	thisTicketAsBytes, _ := APIstub.GetState(args[0])
	var thisTicket = Ticket{}
	json.Unmarshal(thisTicketAsBytes, &thisTicket)
	if thisTicket.OnSell == true {
		return shim.Error("This ticket has already been sold!")
	}
	thisTicket.OnSell = false
	thisTicket.IsRedeemed = true
	thisTicketAsBytes, _ = json.Marshal(thisTicket)
	APIstub.PutState("TICKET"+thisTicket.TicketId, thisTicketAsBytes)
	return shim.Success(nil)
}
func (s *SmartContract) queryEvent(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//0: event ID
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	var id = args[0]

	fmt.Printf("ID: " + args[0])
	eventAsBytes, _ := APIstub.GetState(id)
	if eventAsBytes == nil {
		return shim.Error("Could not locate event")
	}
	return shim.Success(eventAsBytes)
}
func (s *SmartContract) queryAllEvent(APIstub shim.ChaincodeStubInterface) sc.Response {
	//var id = args[0]
	var queryString = "{\r\n\"selector\":{\r\n\"total\":{\r\n \"$gt\":0\r\n}\r\n}\r\n}"
	resultsIterator, err := APIstub.GetQueryResult(queryString)
	defer resultsIterator.Close()
	if err != nil {
		return shim.Error(err.Error())
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse,
			err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
