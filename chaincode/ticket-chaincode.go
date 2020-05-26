package main
import(
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)
type SmartContract struct{

}
 type Ticket struct {
	 Issuer string `json:"issuer"`
	 Cost string `json:"cost"`
	 Name string `json:"name"`
	 CurrentOwner string `json:"currentOwner"`
 }
 func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	if function == "getATicket"{
		return s.getATicket(APIstub,args)
	}else if function == "initLedger"{
		return s.initLedger(APIstub,args)
	}else if function == "issueTicket"{
		return s.issueTicket(APIstub,args)
	}else if function == "getAllTickets"{
		return s.getAllTickets(APIstub,args)
	}else if function == "sellTicket"{
		return s.sellTicket(APIstub,args)
	} 	
	return shim.Error("Wrong function name.")	
}
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response{
	tickets := [] Ticket{
		Ticket{Issuer: "VFF", Cost: "220.000", Name: "Suzuki cup", CurrentOwner: "NA"},
		Ticket{Issuer: "B", Cost: "220.000", Name: "B cup", CurrentOwner: "NA"},
		Ticket{Issuer: "C", Cost: "220.000", Name: "C cup", CurrentOwner: "NA"},
		Ticket{Issuer: "D", Cost: "220.000", Name: "D cup", CurrentOwner: "NA"},
		Ticket{Issuer: "E", Cost: "220.000", Name: "E cup", CurrentOwner: "NA"},
		Ticket{Issuer: "F", Cost: "220.000", Name: "G cup", CurrentOwner: "NA"},
		Ticket{Issuer: "G", Cost: "220.000", Name: "H cup", CurrentOwner: "NA"},		
	}
	i := 0
	for i< len(tickets)
	{
		ticketAsBytes, _ := json.Marshal(tickets[i])
		APIstub.PutState("TICKET"+strconv.Itoa(i), ticketAsBytes)		
		i = i+1
	}
	return shim.Success(nil)
}
func (s *SmartContract) issueTicket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	var ticket = Ticket{Issuer: args[1], Cost: args[2], Name: args[3], CurrentOwner: args[4]}

	ticketAsBytes, _ := json.Marshal(ticket)
	APIstub.PutState(args[0], ticketAsBytes)

	return shim.Success(nil)
}
func (s *SmartContract) sellTicket(APIstub shim.ChaincodeStubInterface, args [] string) sc.Response {
	
	ticketAsBytes, _ = APIstub.GetState(args[0])
	ticket := Ticket{}
	json.Unmarshal(ticketAsBytes, &ticket)
	ticket.CurrentOwner = args[4]
	ticketAsBytes,_ = json.Marshal(ticket)
	APIstub.PutState(args[0], ticketAsBytes)
	return shim.Success(nil)
	
}
func main()
{
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
func (s *SmartContract) queryAllTuna(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add comma before array members,suppress it for the first array member
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

	fmt.Printf("- queryAllTuna:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}