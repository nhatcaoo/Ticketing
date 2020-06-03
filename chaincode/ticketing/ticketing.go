package main

import (
	"time"

	//"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim"
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

func () Init() {
	return shim.Success(nil)
}
func () Invoke() {

}
