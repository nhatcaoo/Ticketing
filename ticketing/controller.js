var express = require("express");
var app = express();
var bodyParser = require("body-parser");
var http = require("http");
var fs = require("fs");
var Fabric_Client = require("fabric-client");
var path = require("path");
var util = require("util");
var os = require("os");

module.exports = (function () {
  return {
    get_all_event: function (req, res) {
    console.log("getting all event: ");
   
    var fabric_client = new Fabric_Client();
    var channel = fabric_client.newChannel("mychannel");
    var peer = fabric_client.newPeer("grpc://localhost:7051");
    channel.addPeer(peer);

    member_user = null;
    var store_path = path.join(os.homedir(), ".hfc-key-store");
    console.log("Store path:" + store_path);
    var tx_id = null;

    Fabric_Client.newDefaultKeyValueStore({ path: store_path })
      .then((state_store) => {
        fabric_client.setStateStore(state_store);
        var crypto_suite = Fabric_Client.newCryptoSuite();
        var crypto_store = Fabric_Client.newCryptoKeyStore({
          path: store_path,
        });
        crypto_suite.setCryptoKeyStore(crypto_store);
        fabric_client.setCryptoSuite(crypto_suite);

        return fabric_client.getUserContext("user1", true);
      })
      .then((user_from_store) => {
        if (user_from_store && user_from_store.isEnrolled()) {
          console.log("Successfully loaded user1 from persistence");
          member_user = user_from_store;
        } else {
          throw new Error("Failed to get user1.... run registerUser.js");
        }

        const request = {
          chaincodeId: "ticketing",
          txId: tx_id,
          fcn: "queryAllEvent",
          args: [''],
        };

        return channel.queryByChaincode(request);
      })
      .then((query_responses) => {
        console.log("Query has completed, checking results");

        if (query_responses && query_responses.length == 1) {
          if (query_responses[0] instanceof Error) {
            console.error("error from query = ", query_responses[0]);
          } else {
            console.log("Response is ", query_responses[0].toString());
            res.json(JSON.parse(query_responses[0].toString()));
          }
        } else {
          console.log("No payloads were returned from query");
        }
      })
      .catch((err) => {
        console.error("Failed to query :: " + err);
      });
  },
    get_all_ticket: function (req, res) {
      console.log("getting all ticket of an event: ");
      var fabric_client = new Fabric_Client();
      var key = req.params.id;
      var query = "{\r\n\"selector\":{\r\n\"total\":{\r\n \"$gt\":0\r\n}\r\n}\r\n}"
      console.log("KEY:",key);   
      var channel = fabric_client.newChannel("mychannel");
      var peer = fabric_client.newPeer("grpc://localhost:7051");
      channel.addPeer(peer);

      member_user = null;
      var store_path = path.join(os.homedir(), ".hfc-key-store");
      console.log("Store path:" + store_path);
      var tx_id = null;

      Fabric_Client.newDefaultKeyValueStore({ path: store_path })
        .then((state_store) => {
          fabric_client.setStateStore(state_store);
          var crypto_suite = Fabric_Client.newCryptoSuite();
          var crypto_store = Fabric_Client.newCryptoKeyStore({
            path: store_path,
          });
          crypto_suite.setCryptoKeyStore(crypto_store);
          fabric_client.setCryptoSuite(crypto_suite);

          return fabric_client.getUserContext("user1", true);
        })
        .then((user_from_store) => {
          if (user_from_store && user_from_store.isEnrolled()) {
            console.log("Successfully loaded user1 from persistence");
            member_user = user_from_store;
          } else {
            throw new Error("Failed to get user1.... run registerUser.js");
          }

          const request = {
            chaincodeId: "ticketing",
            txId: tx_id,
            fcn: "queryAllTicket",
            args: [query],
          };

          return channel.queryByChaincode(request);
        })
        .then((query_responses) => {
          console.log("Query has completed, checking results");

          if (query_responses && query_responses.length == 1) {
            if (query_responses[0] instanceof Error) {
              console.error("error from query = ", query_responses[0]);
            } else {
              console.log("Response is ", query_responses[0].toString());
              res.json(JSON.parse(query_responses[0].toString()));
            }
          } else {
            console.log("No payloads were returned from query");
          }
        })
        .catch((err) => {
          console.error("Failed to query :: " + err);
        });
    },
    get_ticket: function (req, res) {
      var fabric_client = new Fabric_Client();
      var key = req.params.id;

      var channel = fabric_client.newChannel("mychannel");
      var peer = fabric_client.newPeer("grpc://localhost:7051");
      channel.addPeer(peer);
      var member_user = null;
      var store_path = path.join(os.homedir(), ".hfc-key-store");
      console.log("Store path:" + store_path);
      var tx_id = null;

      Fabric_Client.newDefaultKeyValueStore({ path: store_path })
        .then((state_store) => {
          fabric_client.setStateStore(state_store);
          var crypto_suite = Fabric_Client.newCryptoSuite();
          var crypto_store = Fabric_Client.newCryptoKeyStore({
            path: store_path,
          });
          crypto_suite.setCryptoKeyStore(crypto_store);
          fabric_client.setCryptoSuite(crypto_suite);
          return fabric_client.getUserContext("user1", true);
        })
        .then((user_from_store) => {
          if (user_from_store && user_from_store.isEnrolled()) {
            console.log("Successfully loaded user1 from persistence");
            member_user = user_from_store;
          } else {
            throw new Error("Failed to get user1.... run registerUser.js");
          }

          const request = {
            chaincodeId: "ticketing",
            txId: tx_id,
            fcn: "queryTicket",
            args: [key],
          };

          return channel.queryByChaincode(request);
        })
        .then((query_responses) => {
          console.log("Query has completed, checking results");
          if (query_responses && query_responses.length == 1) {
            if (query_responses[0] instanceof Error) {
              console.error("error from query = ", query_responses[0]);
              res.send("Could not locate tuna");
            } else {
              console.log("Response is ", query_responses[0].toString());
              res.send(query_responses[0].toString());
            }
          } else {
            console.log("No payloads were returned from query");
            res.send("Could not locate ticket");
          }``
        })
        .catch((err) => {
          console.error("Failed to query successfully :: " + err);
          res.send("Could not locate ticket");
        });
    },
    create_event: function (req, res) {
      console.log("CREATE EVENT ");

      var array = req.params.event.split("-");
      console.log(array);

      var issuer = array[0];
      var price = array[1];
      var eventName = array[2];
      var total = array[3];

      var fabric_client = new Fabric_Client();

      var channel = fabric_client.newChannel("mychannel");
      var peer = fabric_client.newPeer("grpc://localhost:7051");
      channel.addPeer(peer);
      var order = fabric_client.newOrderer("grpc://localhost:7050");
      channel.addOrderer(order);

      var member_user = null;
      var store_path = path.join(os.homedir(), ".hfc-key-store");
      console.log("Store path:" + store_path);
      var tx_id = null;

      Fabric_Client.newDefaultKeyValueStore({ path: store_path })
        .then((state_store) => {
          fabric_client.setStateStore(state_store);
          var crypto_suite = Fabric_Client.newCryptoSuite();
          var crypto_store = Fabric_Client.newCryptoKeyStore({
            path: store_path,
          });
          crypto_suite.setCryptoKeyStore(crypto_store);
          fabric_client.setCryptoSuite(crypto_suite);

          return fabric_client.getUserContext("user1", true);
        })
        .then((user_from_store) => {
          if (user_from_store && user_from_store.isEnrolled()) {
            console.log("Successfully loaded user1 from persistence");
            member_user = user_from_store;
          } else {
            throw new Error("Failed to get user1.... run registerUser.js");
          }

          tx_id = fabric_client.newTransactionID();
          console.log("Assigning transaction_id: ", tx_id._transaction_id);
          const request = {
            chaincodeId: "ticketing",
            fcn: "createEvent",
            args: [issuer, price, eventName, total],
            chainId: "mychannel",
            txId: tx_id,
          };
          

          return channel.sendTransactionProposal(request);
        })
        .then((results) => {
          var proposalResponses = results[0];
          var proposal = results[1];
          let isProposalGood = false;
          if (
            proposalResponses &&
            proposalResponses[0].response &&
            proposalResponses[0].response.status === 200
          ) {
            isProposalGood = true;
            console.log("Transaction proposal was good");
          } else {
            console.error("Transaction proposal was bad");
          }
          if (isProposalGood) {
            console.log(
              util.format(
                'Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s"',
                proposalResponses[0].response.status,
                proposalResponses[0].response.message
              )
            );

            var request = {
              proposalResponses: proposalResponses,
              proposal: proposal,
            };

            var transaction_id_string = tx_id.getTransactionID(); //Get the transaction ID string to be used by the event processing
            var promises = [];

            var sendPromise = channel.sendTransaction(request);
            promises.push(sendPromise); //we want the send transaction first, so that we know where to check status

            let event_hub = fabric_client.newEventHub();
            event_hub.setPeerAddr("grpc://localhost:7053");

            let txPromise = new Promise((resolve, reject) => {
              let handle = setTimeout(() => {
                event_hub.disconnect();
                resolve({ event_status: "TIMEOUT" }); //we could use reject(new Error('Trnasaction did not complete within 30 seconds'));
              }, 3000);
              event_hub.connect();
              event_hub.registerTxEvent(
                transaction_id_string,
                (tx, code) => {
                  clearTimeout(handle);
                  event_hub.unregisterTxEvent(transaction_id_string);
                  event_hub.disconnect();

                  var return_status = {
                    event_status: code,
                    tx_id: transaction_id_string,
                  };
                  if (code !== "VALID") {
                    console.error(
                      "The transaction was invalid, code = " + code
                    );
                    resolve(return_status); // we could use reject(new Error('Problem with the tranaction, event status ::'+code));
                  } else {
                    console.log(
                      "The transaction has been committed on peer " +
                        event_hub._ep._endpoint.addr
                    );
                    resolve(return_status);
                  }
                },
                (err) => {
                  reject(
                    new Error("There was a problem with the eventhub ::" + err)
                  );
                }
              );
            });
            promises.push(txPromise);

            return Promise.all(promises);
          } else {
            console.error(
              "Failed to send Proposal or receive valid response. Response null or status is not 200. exiting..."
            );
            throw new Error(
              "Failed to send Proposal or receive valid response. Response null or status is not 200. exiting..."
            );
          }
        })
        .then((results) => {
          console.log(
            "Send transaction promise and event listener promise have completed"
          );
          // check the results in the order the promises were added to the promise all list
          if (results && results[0] && results[0].status === "SUCCESS") {
            console.log("Successfully sent transaction to the orderer.");
            res.send(tx_id.getTransactionID());
          } else {
            console.error(
              "Failed to order the transaction. Error code: " + response.status
            );
          }

          if (results && results[1] && results[1].event_status === "VALID") {
            console.log(
              "Successfully committed the change to the ledger by the peer"
            );
            res.send(tx_id.getTransactionID());
          } else {
            console.log(
              "Transaction failed to be committed to the ledger due to ::" +
                results[1].event_status
            );
          }
        })
        .catch((err) => {
          console.error("Failed to invoke successfully :: " + err);
        });
    },

    buyTicketFromSupplier: function (req, res) {
      console.log("Change current owner of a ticket: ");

      var Fabric_Client = require("fabric-client");
      var path = require("path");
      var util = require("util");
      var os = require("os");

      var array = req.params.info.split("_");
      var key = array[0];
      var number = array[1];
      var owner = array[2];

      var fabric_client = new Fabric_Client();
      var channel = fabric_client.newChannel("mychannel");
      var peer = fabric_client.newPeer("grpc://localhost:7051");
      channel.addPeer(peer);
      var order = fabric_client.newOrderer("grpc://localhost:7050");
      channel.addOrderer(order);

      var member_user = null;
      var store_path = path.join(os.homedir(), ".hfc-key-store");
      console.log("Store path:" + store_path);
      var tx_id = null;

      Fabric_Client.newDefaultKeyValueStore({ path: store_path })
        .then((state_store) => {
          fabric_client.setStateStore(state_store);
          var crypto_suite = Fabric_Client.newCryptoSuite();
          var crypto_store = Fabric_Client.newCryptoKeyStore({
            path: store_path,
          });
          crypto_suite.setCryptoSuite(crypto_store);
          fabric_client.setCryptoSuite(crypto_suite);
          return fabric_client.getUserContext("user1", true);
        })
        .then((user_from_store) => {
          if (user_from_store && user_from_store.isEnrolled()) {
            console.log("Successfully loaded user1 from persistence");
            member_user = user_from_store;
          } else {
            throw new Error("Failed to get user 1.");
          }
          tx_id = fabric_client.newTransactionID();
          console.log("Assigning transaction_id: ", tx_id._transaction_id);

          var request = {
            //targets : --- letting this default to the peers assigned to the channel
            chaincodeId: "ticketing",
            fcn: "buyTicketFromSupplier",
            args: [key, number, holder, owner],
            chainId: "mychannel",
            txId: tx_id,
          };
          return channel.sendTransactionProposal(request);
        })
        .then((results) => {
          var proposalResponses = results[0];
          var proposal = results[1];
          let isProposalGood = false;
          if (
            proposalResponses &&
            proposalResponses[0].response &&
            proposalResponses[0].response.status === 200
          ) {
            isProposalGood = true;
            console.log("Transaction proposal was good");
          } else {
            console.error("Transaction proposal was bad");
          }

          if (isProposalGood) {
            console.log(
              util.format(
                'Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s"',
                proposalResponses[0].response.status,
                proposalResponses[0].response.message
              )
            );
            var request = {
              proposalResponses: proposalResponses,
              proposal: proposal,
            };
            var transaction_id_string = tx_id.getTransactionID();
            var promises = [];
            var sendPromise = channel.sendTransaction(request);
            promises.push(sendPromise);
            let event_hub = fabric_client.newEventHub();
            event_hub.setPeerAddr("grpc://localhost:7053");

            let txPromise = new Promise((resolve, reject) => {
              let handle = setTimeout(() => {
                event_hub.disconnect();
                resolve({ event_status: "TIMEOUT" });
              }, 3000);
              event_hub.connect();
              event_hub.registerTxEvent(
                transaction_id_string,
                (tx, code) => {
                  clearTimeout(handle);
                  event_hub.unregisterTxEvent(transaction_id_string);
                  event_hub.disconnect();

                  var return_status = {
                    event_status: code,
                    tx_id: transaction_id_string,
                  };
                  if (code !== "VALID") {
                    console.error(
                      "The transaction was invalid, code = " + code
                    );
                    resolve(return_status); // we could use reject(new Error('Problem with the tranaction, event status ::'+code));
                  } else {
                    console.log(
                      "The transaction has been committed on peer " +
                        event_hub._ep._endpoint.addr
                    );
                    resolve(return_status);
                  }
                },
                (err) => {
                  reject(new Error("Problem with eventhub"));
                }
              );
            });
            promises.push(txtPromise);
          } else {
            console.error(
              "Failed to send Proposal or receive valid response. Response null or status is not 200. exiting..."
            );
            throw new Error(
              "Failed to send Proposal or receive valid response. Response null or status is not 200. exiting..."
            );
          }
        })
        .then((results) => {
          console.log(
            "Send transaction promise and event listener promise have completed"
          );
          if (results && results[0] && results[0].status === "SUCCESS") {
            console.log("Successfully sent transaction to the orderer.");
            res.json(tx_id.getTransactionID());
          } else {
            console.error(
              "Failed to order the transaction. Error code: " + response.status
            );
          }

          if (results && results[1] && results[1].event_status === "VALID") {
            console.log(
              "Successfully committed the change to the ledger by the peer"
            );
            res.json(tx_id.getTransactionID());
          } else {
            console.log(
              "Transaction failed to be committed to the ledger due to ::" +
                results[1].event_status
            );
          }
        })
        .catch((err) => {
          console.error("Failed to invoke :" + err);
        });
    },
  };
})();
