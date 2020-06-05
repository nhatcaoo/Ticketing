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
	ID        int    `json:"id"`
	Issuer    string `json:"issuer"`
	Price     string `json:"price"`
	EventName string `json:"eventName"`
	Total     int    `json:"total"`
	Sold      int    `json:"sold"`
}
type Ticket struct {
	EventId      int       `json:"eventId"`
	TicketId     string    `json:"ticketId"`
	Cost         string    `json:"cost"`
	CurrentOwner string    `json:"currentOwner"`
	OnSell       bool      `json:"onSell"`
	TimeStamp    time.Time `json:"timeStamp"`
	IsRedeemed   bool      `json:"isRedeemed"`
}
type Info struct {
	number int `json:"number"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

var logger = shim.NewLogger("ticketing")

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	logger.Infof("start ")
	function, args := APIstub.GetFunctionAndParameters()

	logger.Infof("Invoke is running " + function)
	if function == "queryTicket" {
		return s.queryTicket(APIstub, args)
	} else if function == "initEvent" {
		logger.Infof("through init Event")
		return s.initEvent(APIstub)
	} else if function == "buyTicketFromSupplier" {
		return s.buyTicketFromSupplier(APIstub, args)
	} else if function == "buyTicketFromFromSecondaryMarket" {
		return s.buyTicketFromFromSecondaryMarket(APIstub, args)
	} else if function == "queryAllTicket" {
		return s.queryAllTicket(APIstub, args)
	} else if function == "queryAllEvent" {
		return s.queryAllEvent(APIstub)
	} else if function == "createEvent" {
		return s.createEvent(APIstub, args)
	} else if function == "upTicketToSecondaryMarket" {
		return s.upTicketToSecondaryMarket(APIstub, args)
	} else if function == "removeTicketFromSecondaryMarket" {
		return s.removeTicketFromSecondaryMarket(APIstub, args)
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
		Event{ID: 0, Issuer: "VFF", Price: "220.000", EventName: "Suzuki cup", Total: 20, Sold: 0},
		Event{ID: 1, Issuer: "BFF", Price: "220.000", EventName: "B cup", Total: 20, Sold: 0},
		Event{ID: 2, Issuer: "CFF", Price: "220.000", EventName: "C cup", Total: 20, Sold: 0},
		Event{ID: 3, Issuer: "DFF", Price: "220.000", EventName: "D cup", Total: 20, Sold: 0},
		Event{ID: 4, Issuer: "EFF", Price: "220.000", EventName: "F cup", Total: 20, Sold: 0}}
	logger.Infof("done2 ")

	for j := 0; j < 5; j++ {
		eventAsBytes, _ := json.Marshal(events[j])
		APIstub.PutState("EVENT"+strconv.Itoa(events[j].ID), eventAsBytes)
		for i := 0; i < events[j].Total; i++ {

			var ticket = Ticket{EventId: events[j].ID, TicketId: strconv.Itoa(events[j].ID) + "-" + strconv.Itoa(i), Cost: events[j].Price, CurrentOwner: "N/A", OnSell: true, TimeStamp: time.Now(), IsRedeemed: false}
			ticketAsBytes, _ := json.Marshal(ticket)
			APIstub.PutState("TICKET"+ticket.TicketId, ticketAsBytes)
			logger.Infof(ticket.TicketId)
		}
		logger.Infof("-\n ")
	}
	logger.Infof("done1 ")
	var info = Info{}
	info.number = 5
	infoAsBytes, _ := json.Marshal(info)
	logger.Infof("done1 ")
	APIstub.PutState("NUMBER_EVENTS", infoAsBytes)
	logger.Infof("done1 ")
	return shim.Success(nil)
}
func (s *SmartContract) buyTicketFromSupplier(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//arg: event ID : "EVENT"+strconv.Itoa(i)
	thisEventAsBytes, _ := APIstub.GetState(args[0])
	var thisEvent = Event{}
	json.Unmarshal(thisEventAsBytes, &thisEvent)
	var left = thisEvent.Total - thisEvent.Sold
	num, _ := strconv.Atoi(args[1])
	fmt.Printf(strconv.Itoa(num))
	fmt.Printf((args[0]))
	if num > left {
		return shim.Error("Incorrect number of tickets. Expecting")
	} else {
		//ticketSet := []Ticket{}
		for i := 0; i < num; i++ {
			fmt.Printf("buy ticket - \n")
			eventAsBytes, _ := APIstub.GetState(args[0])
			var event = Event{}
			json.Unmarshal(eventAsBytes, &event)
			thisTicketAsBytes, _ := APIstub.GetState("TICKET" + strconv.Itoa(event.ID) + "-" + strconv.Itoa(event.Sold))
			var thisTicket = Ticket{}
			json.Unmarshal(thisTicketAsBytes, &thisTicket)
			thisTicket.CurrentOwner = args[2]
			thisTicket.OnSell = false
			thisTicket.TimeStamp = time.Now() //timestamp
			event.Sold++
			eventAsBytes, _ = json.Marshal(event)
			APIstub.PutState(args[0], eventAsBytes)
			thisTicketAsBytes, _ = json.Marshal(thisTicket)
			APIstub.PutState("TICKET"+strconv.Itoa(event.ID)+"-"+strconv.Itoa(event.Sold), thisTicketAsBytes)
		}
	}
	return shim.Success(nil)
}
func (s *SmartContract) createEvent(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var info = Info{}
	numberAsBytes, _ := APIstub.GetState("NUMBER_EVENTS")
	fmt.Printf(args[0])
	json.Unmarshal(numberAsBytes, &info)
	var number = info.number
	fmt.Printf(number)
	total, _ := strconv.Atoi(args[3])
	var event = Event{ID: number, Issuer: args[0], Price: args[1], EventName: args[2], Total: total, Sold: 0}
	eventAsBytes, _ := json.Marshal(event)
	APIstub.PutState("EVENT"+strconv.Itoa(events.ID), eventAsBytes)
	for i := 0; i < event.Total; i++ {
		var ticket = Ticket{EventId: event.ID, TicketId: strconv.Itoa(number) + "-" + strconv.Itoa(i), Cost: event.Price, CurrentOwner: "N/A", OnSell: true, TimeStamp: time.Now(), IsRedeemed: false}
		ticketAsBytes, _ := json.Marshal(ticket)
		APIstub.PutState("TICKET"+ticket.TicketId, ticketAsBytes)
		fmt.Printf("-\n")
	}
	number++
	info.number = number
	numberAsBytes, _ = json.Marshal(info)
	APIstub.PutState("NUMBER_EVENTS", numberAsBytes)
	return shim.Success(nil)
}
func (s *SmartContract) upTicketToSecondaryMarket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Printf(args[0])

	thisTicketAsBytes, _ := APIstub.GetState(args[0])
	var thisTicket = Ticket{}
	json.Unmarshal(thisTicketAsBytes, &thisTicket)
	thisTicket.OnSell = true
	thisTicketAsBytes, _ = json.Marshal(thisTicket)
	APIstub.PutState("TICKET"+thisTicket.TicketId, thisTicketAsBytes)
	return shim.Success(nil)
}
func (s *SmartContract) removeTicketFromSecondaryMarket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Printf(args[0])
	thisTicketAsBytes, _ := APIstub.GetState(args[0])
	var thisTicket = Ticket{}
	json.Unmarshal(thisTicketAsBytes, &thisTicket)
	if thisTicket.OnSell == true {
		return shim.Error("This ticket has already been sold!")
	}
	thisTicket.OnSell = false
	thisTicketAsBytes, _ = json.Marshal(thisTicket)
	APIstub.PutState("TICKET"+thisTicket.TicketId, thisTicketAsBytes)
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
func (s *SmartContract) buyTicketFromFromSecondaryMarket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Printf(args[0])
	thisTicketAsBytes, _ := APIstub.GetState(args[0])
	var thisTicket = Ticket{}
	json.Unmarshal(thisTicketAsBytes, &thisTicket)
	if thisTicket.OnSell == false {
		return shim.Error("Ticket is not on selling")
	}
	thisTicket.CurrentOwner = args[1]
	thisTicketAsBytes, _ = json.Marshal(thisTicket)
	APIstub.PutState("TICKET"+thisTicket.TicketId, thisTicketAsBytes)
	return shim.Success(nil)
}
func (s *SmartContract) checkoutTicket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Printf(args[0])
	thisTicketAsBytes, _ := APIstub.GetState(args[0])
	var thisTicket = Ticket{}
	json.Unmarshal(thisTicketAsBytes, &thisTicket)
	if thisTicket.IsRedeemed == true {
		return shim.Error("This ticket has already been redeemed!")
	} else if args[1] != strconv.Itoa(thisTicket.EventId) || args[2] != thisTicket.TicketId || args[3] != thisTicket.CurrentOwner {
		return shim.Error("Ticket fault")
	} else {
		fmt.Printf("Valid ticket")
		return shim.Success(nil)
	}

	return shim.Success(nil)
}
func (s *SmartContract) queryTicket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	var id = args[0]
	if id == "" {
		id = "TICKET1-1"
	}
	fmt.Printf("ID: " + args[0])
	ticketAsBytes, _ := APIstub.GetState(id)
	if ticketAsBytes == nil {
		return shim.Error("Could not locate ticket")
	}
	return shim.Success(ticketAsBytes)
}
func (s *SmartContract) queryAllTicket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//var id = args[0]
	fmt.Printf(args[0])
	var queryString = args[0]
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
