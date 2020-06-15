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

type Ticket struct {
	EventId      int       `json:"eventId"`
	TicketId     string    `json:"ticketId"`
	Cost         string    `json:"cost"`
	CurrentOwner string    `json:"currentOwner"`
	Status		 string	   `json:"status"`
	TimeStamp    time.Time `json:"timeStamp"`
	Key	string	`json:"key"`
}
var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}                                                                                                                             
                                                                                                                                                                                                            
func encodeBase64(b []byte) string {                                                                                                                                                                        
    return base64.StdEncoding.EncodeToString(b)                                                                                                                                                             
}  
func Encrypt(key, text string) string {                                                                                                                                                                     
    fmt.Println(text)                                                                                                                                                                                       
    block, err := aes.NewCipher([]byte(key))                                                                                                                                                                
    if err != nil { panic(err) }                                                                                                                                                                            
    plaintext := []byte(text)                                                                                                                                                                               
    cfb := cipher.NewCFBEncrypter(block, iv)                                                                                                                                                                
    ciphertext := make([]byte, len(plaintext))                                                                                                                                                              
    cfb.XORKeyStream(ciphertext, plaintext)                                                                                                                                                                
    return encodeBase64(ciphertext)                                                                                                                                                                         
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

	logger.Infof("Invoke is running " + function)
	if function == "queryTicket" {
		return s.queryTicket(APIstub, args)
	} else if function == "buyTicketFromSupplier" {
		return s.buyTicketFromSupplier(APIstub, args)
	} else if function == "buyTicketFromFromSecondaryMarket" {
		return s.buyTicketFromFromSecondaryMarket(APIstub, args)
	} else if function == "queryAllTicket" {
		return s.queryAllTicket(APIstub, args)
	} else if function == "upTicketToSecondaryMarket" {
		return s.upTicketToSecondaryMarket(APIstub, args)
	} else if function == "removeTicketFromSecondaryMarket" {
		return s.removeTicketFromSecondaryMarket(APIstub, args)
	} else if function == "redeemTicket" {
		return s.redeemTicket(APIstub, args)
	}
	return shim.Error("Wrong function name.")
}

func (s *SmartContract) buyTicketFromSupplier(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//arg: event ID : "EVENT"+strconv.Itoa(i)
	//0: eventID, 1: number, 2: current owner, 3: redeemed time, 4: total, 5: sold, 6 cost
	redeemedTime, _ := time.Parse(time.RFC3339, args[3])
	if redeemedTime.Before(time.Now())
	{
		return shim.Error("Redeemed dated exceeded")
	}
	total := strconv.Atoi(args[4])
	sold := strconv.Atoi(args[5])
	var left = thisEvent.Total - thisEvent.Sold
	num, _ := strconv.Atoi(args[1])
	count := 0
	if num > left {
		return shim.Error("Incorrect number of tickets.")
	} else {
		for i := 0; i < num; i++ {
			fmt.Printf("buy ticket - \n")
			var thisTicket = Ticket{EventId:args[0],TicketId: args[0] + "-" + sold, Cost: args[6], CurrentOwner: args[2], Status: "locked", TimeStamp: time.Now(), Key: ""}		
			thisTicketAsBytes, _ = json.Marshal(thisTicket)
			APIstub.PutState("TICKET"+thisTicket.TicketId, thisTicketAsBytes)
			count++;
		}
	}
	return shim.Success(count)
}
func (s *SmartContract) buyTicketFromFromSecondaryMarket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Printf(args[0])
	thisTicketAsBytes, _ := APIstub.GetState(args[0])
	var thisTicket = Ticket{}
	json.Unmarshal(thisTicketAsBytes, &thisTicket)
	if thisTicket.Status == "locked" {
		return shim.Error("Ticket is not on selling")
	}
	if thisTicket.Status == "redeemed" {
		return shim.Error("Ticket is not available to buy")
	}
	thisTicket.CurrentOwner = args[1]
	thisTicketAsBytes, _ = json.Marshal(thisTicket)
	APIstub.PutState("TICKET"+thisTicket.TicketId, thisTicketAsBytes)
	return shim.Success(nil)
}
func (s *SmartContract) upTicketToSecondaryMarket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Printf(args[0])

	thisTicketAsBytes, _ := APIstub.GetState(args[0])
	var thisTicket = Ticket{}
	json.Unmarshal(thisTicketAsBytes, &thisTicket)
	if thisTicket.Status == "trading" {
		return shim.Error("This ticket has already putted on secondary market !")
	}
	if thisTicket.Status == "redeemed" {
		return shim.Error("Ticket is not available to sell !")
	}
	thisTicket.Status = "trading"
	thisTicketAsBytes, _ = json.Marshal(thisTicket)
	APIstub.PutState("TICKET"+thisTicket.TicketId, thisTicketAsBytes)
	return shim.Success(nil)
}
func (s *SmartContract) removeTicketFromSecondaryMarket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Printf(args[0])
	thisTicketAsBytes, _ := APIstub.GetState(args[0])
	var thisTicket = Ticket{}
	json.Unmarshal(thisTicketAsBytes, &thisTicket)

	if thisTicket.Status == "trading" {
		return shim.Error("This ticket has already putted on secondary market !")
	}
	if thisTicket.Status == "redeemed" {
		return shim.Error("Ticket is not available to sell !")
	}

	thisTicket.Status = "locked"
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


func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
