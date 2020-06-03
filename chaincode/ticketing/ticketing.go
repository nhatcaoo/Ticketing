package main

import (
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
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()
	if function == "queryTicket" {
		//	return s.queryTicket(APIstub, args)
	} else if function == "initEvent" {
		//	return s.initEvent(APIstub)
	} else if function == "buyTicketFromSupplier" {
		//	return s.buyTicketFromSupplier(APIstub, args)
	} else if function == "buyTicketFromFromSecondaryMarket" {
		return s.buyTicketFromFromSecondaryMarket(APIstub, args)
	} else if function == "queryAllTicket" {
		return s.queryAllTicket(APIstub, args)
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
