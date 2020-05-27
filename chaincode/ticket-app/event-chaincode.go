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
type Event struct {
	ID int `json:"id"`
	Issuer string `json:"issuer"`
	Price string `json:"price"`
	EventName string `json:"eventName"`
	Total int `json:"total"`
	Sold int `json:"sold"`
}
 type Ticket struct {
	 EventId int `json:"eventId"`
	 TicketId string `json:"ticketId"`	
	 Cost string `json:"cost"`
	 CurrentOwner string `json:"currentOwner"`
	 OnSell bool 	`json:"onSell"`
	 TimeStamp int  `json:"timeStamp"`
	 IsRedeemed bool  `json:"isRedeemed"`
 }
 func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	if function == "queryTicket"{
		return s.queryTicket(APIstub,args)
	}else if function == "initEvent"{
		return s.initEvent(APIstub,args)
	}else if function == "buyTicketFromSupplier"{
		return s.buyTicketFromSupplier(APIstub,args)
	}else if function == "buyTicketFromFromSecondaryMarket"{
		return s.buyTicketFromFromSecondaryMarket(APIstub,args)
	}else if function == "queryAllTicket"{
		return s.queryAllTicket(APIstub,args)
	}else if function == "createEvent"{
		return s.createEvent(APIstub,args)
	}else if function == "upTicketToSecondaryMarket"{
		return s.upTicketToSecondaryMarket(APIstub,args)
	}else if function == "removeTicketFromSecondaryMarket"{
		return s.removeTicketFromSecondaryMarket(APIstub,args)
	}else if function == "redeemTicket"{
		return s.redeemTicket(APIstub,args)
	}else if function == "checkoutTicket"{
		return s.checkoutTicket(APIstub,args)
	}
	return shim.Error("Wrong function name.")	
}
func (s *SmartContract) initEvent(APIstub shim.ChaincodeStubInterface) sc.Response{
	
	events := [] Event{
		Event{ID: 0, Issuer: "VFF", Price: "220.000", EventName: "Suzuki cup", Total: 20, Sold: 0}
		Event{ID: 1, Issuer: "BFF", Price: "220.000", EventName: "B cup", Total: 20, Sold: 0}
		Event{ID: 2, Issuer: "CFF", Price: "220.000", EventName: "C cup", Total: 20, Sold: 0}
		Event{ID: 3, Issuer: "DFF", Price: "220.000", EventName: "D cup", Total: 20, Sold: 0}
		Event{ID: 4, Issuer: "EFF", Price: "220.000", EventName: "F cup", Total: 20, Sold: 0}	}

	i := 0
	for i< len(events)
	{
		eventAsBytes, _ := json.Marshal(events[i])
		APIstub.PutState("EVENT"+strconv.Itoa((events[i].ID), eventAsBytes)		
		for i := 0; i < events[i].Total; i++ {
			var ticket = Ticket{EventId: events[i].ID, TicketId: events[i].ID +"-" + i, Cost: events[i].Price, CurrentOwner: "N/A", OnSell: true,time.Now(),false   }
			ticketAsBytes,_ =json.Marshal(ticket)
			APIstub.PutState("TICKET"+ticket.TicketId,ticketAsBytes)
		}
		i = i+1
	}
	var number int
	number = 5
	APIstub.PutState("NUMBER_EVENTS", number)		
	return shim.Success(nil)
}
func (s *SmartContract) buyTicketFromSupplier(APIstub shim.ChaincodeStubInterface, args []string) sc.Response{
	//arg: event ID : "EVENT"+strconv.Itoa(i)
	thisEventAsBytes, _ = APIstub.GetState(args[0])
	var thisEvent = Event{}
	json.Unmarshal(thisEventAsBytes, &thisEvent)
	var left = thisEvent.Total-thisEvent.Sold
	if args[1]>left {
		return shim.Error("Incorrect number of tickets. Expecting" + left)
	}
	else ticketSet := [] tickets{}
		for i := 0; i < args[1]; i++ {
			eventAsBytes, _ = APIstub.GetState(args[0])
			var event = Event{}
			json.Unmarshal(eventAsBytes, &event)
			thisTicketAsBytes, _ = APIstub.GetState("TICKET"+event.ID+"-"+event.Sold)
			var thisTicket = Ticket{}
			json.Unmarshal(thisTicketAsBytes, &thisTicket)
			thisTicket.CurrentOwner = args[2]
			thisTicket.OnSell = false
			thisTicket.TimeStamp =  time.Now()
			event.Sold++;
			eventAsBytes,_ =json.Marshal(event)
			APIstub.PutState(args[0],eventAsBytes)
			thisTicketAsBytes,_ = json.Marshal(thisTicket)
			APIstub.PutState("TICKET"+event.ID+"-"+event.Sold,thisTicketAsBytes)
		}
		return shim.Success(nil)
}
func (s *SmartContract) createEvent(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var number int;
	numberAsBytes, _ = APIstub.GetState("NUMBER_EVENTS")
	json.Unmarshal(numberAsBytes, &number)
	var event = Event{ID: number, Issuer: args[0], Price: args[1], EventName: args[2], Total: args[3], Sold: 0}
	for i := 0; i < event.Total; i++ {
		var ticket = Ticket{EventId: event.ID, TicketId: event.ID +"-" + i, Cost: event.Price, CurrentOwner: "N/A", OnSell: true,time.Now(),false   }
		ticketAsBytes,_ =json.Marshal(ticket)
		APIstub.PutState("TICKET"+ticket.TicketId,ticketAsBytes)
	}
	return shim.Success(nil)
}
func (s *SmartContract) upTicketToSecondaryMarket(APIstub shim.ChaincodeStubInterface, args [] string) sc.Response {	
	thisTicketAsBytes, _ =  APIstub.GetState(args[0])
	var thisTicket = Ticket{}
	json.Unmarshal(thisTicketAsBytes,&thisTicket)	
	thisTicket.OnSell = true
	thisTicketAsBytes,_ = json.Marshal(thisEvent)
	APIstub.PutState(thisEvent.ID, thisTicketAsBytes)
	return shim.Success(nil)
}
func (s *SmartContract) removeTicketFromSecondaryMarket(APIstub shim.ChaincodeStubInterface, args [] string) sc.Response {	
	thisTicketAsBytes, _ =  APIstub.GetState(args[0])
	var thisTicket = Ticket{}
	json.Unmarshal(thisTicketAsBytes,&thisTicket)	
	if thisTicket.OnSell == true{
		return shim.Error("This ticket has already been sold!")		
	}
	thisTicket.OnSell = false
	thisTicketAsBytes,_ = json.Marshal(thisEvent)
	APIstub.PutState(thisEvent.ID, thisTicketAsBytes)
	return shim.Success(nil)
}
func (s *SmartContract) redeemTicket(APIstub shim.ChaincodeStubInterface, args [] string) sc.Response {	
	thisTicketAsBytes, _ =  APIstub.GetState(args[0])
	var thisTicket = Ticket{}
	json.Unmarshal(thisTicketAsBytes,&thisTicket)	
	if thisTicket.OnSell == true{
		return shim.Error("This ticket has already been sold!")		
	}
	thisTicket.OnSell = false
	thisTicket.IsRedeemed = true
	thisTicketAsBytes,_ = json.Marshal(thisEvent)
	APIstub.PutState(thisEvent.ID, thisTicketAsBytes)
	return shim.Success(nil)
}
func (s *SmartContract) buyTicketFromFromSecondaryMarket(APIstub shim.ChaincodeStubInterface, args [] string) sc.Response {
	thisTicketAsBytes, _ =  APIstub.GetState(args[0])
	var thisTicket = Ticket{}
	json.Unmarshal(thisTicketAsBytes,&thisTicket)	
	if thisTicket.OnSell == false{
		return shim.Error("Ticket is not on selling")		
	}
	thisTicket.CurrentOwner = args[1]
	thisTicketAsBytes,_ = json.Marshal(thisEvent)
	APIstub.PutState(thisEvent.ID, thisTicketAsBytes)
	return shim.Success(nil)
}
func (s *SmartContract) checkoutTicket(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {	
	thisTicketAsBytes, _ =  APIstub.GetState(args[0])
	var thisTicket = Ticket{}
	json.Unmarshal(thisTicketAsBytes,&thisTicket)	
	if thisTicket.IsRedeemed == true{
		return shim.Error("This ticket has already been redeemed!")		
	}
	else if arg[1] != thisTicket.EventId || args[2] != thisTicket.TicketId || args[3] != thisTicket.CurrentOwner
	{
		return shim.Error("Ticket fault")
	}
	else
		return shim.Success("Valid ticket!") 

	
	return shim.Success(nil)
}
func (s *SmartContract) queryTicket(APIstub shim.ChaincodeStubInterface, args [] string) sc.Response{
	
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ticketAsBytes, _ := APIstub.GetState(args[0])
	if ticketAsBytes == nil {
		return shim.Error("Could not locate ticket")
	}
	return shim.Success(ticketAsBytes)
	
}
func (s *SmartContract) queryAllTicket(APIstub shim.ChaincodeStubInterface, args [] string) sc.Response{
	
	var queryString = "{\"selector\":{\"Ticket.EventId\":\"+"args[0]"+\"}"
	resultsIterator, err: = stub.GetQueryResult(queryString)
    defer resultsIterator.Close()
    if err != nil {
        return nil, err
	}
	var buffer bytes.Buffer
    buffer.WriteString("[")
    bArrayMemberAlreadyWritten: = false
    for resultsIterator.HasNext() {
        queryResponse,
        err: = resultsIterator.Next()
        if err != nil {
            return nil, err
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
	return buffer.Bytes(), nil
	
}

func main()
{
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}