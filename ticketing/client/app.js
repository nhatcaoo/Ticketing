// SPDX-License-Identifier: Apache-2.0

"use strict";

var app = angular.module("application", []);
// Angular Controller
app.controller("appController", function ($scope, appFactory) {
  $("#success_holder").hide();
  $("#success_create").hide();
  $("#error_holder").hide();
  $("#error_query").hide();
  
  $scope.getAllEvent = function () {
    appFactory.getAllEvent(function (data) {
      var array = [];
      console.log("data: ",data)
      for (var i = 0; i < data.length; i++) {
        parseInt(data[i].Key);
        data[i].Record.Key = parseInt(data[i].Key);
        array.push(data[i].Record);
      }
      console.log("a")
      array.sort(function (a, b) {
        return parseFloat(a.Key) - parseFloat(b.Key);
      });
      $scope.all_event = array;
    });
  };
  
  $scope.getAllTicket = function () {
    var id = $scope.event_id;
    appFactory.getAllTicket(id,function (data) {
      var array = [];
      console.log("data: ",data)
      for (var i = 0; i < data.length; i++) {
        parseInt(data[i].Key);
        data[i].Record.Key = parseInt(data[i].Key);
        array.push(data[i].Record);
      }
      console.log("a")
      array.sort(function (a, b) {
        return parseFloat(a.Key) - parseFloat(b.Key);
      });
      $scope.all_ticket = array;
    });
    
  };

 
  $scope.getTicket = function () {
    var id = $scope.ticket_id;

    appFactory.getTicket(id, function (data) {
      $scope.get_ticket = data;

      if ($scope.get_ticket == "Could not locate ticket") {
        console.log();
        $("#error_query").show();
      } else {
        $("#error_query").hide();
      }
    });
  };

  $scope.recordEvent = function () {
    appFactory.recordEvent($scope.event, function (data) {
      $scope.create_event = data;
      console.log(data)
      $("#success_create").show();
    });
  };

  $scope.buyTicketFromSupplier = function () {
    appFactory.buyTicketFromSupplier($scope.holder, function (data) {
      $scope.buyTicketFromSupplier = data;
      console.log(data)
      if ($scope.buyTicketFromSupplier == "Error: no tuna catch found") {
        $("#error_holder").show();
        $("#success_holder").hide();
      } else {
        $("#success_holder").show();
        $("#error_holder").hide();
      }
    });
  };
  
});

// Angular Factory
app.factory("appFactory", function ($http) {
  var factory = {};

  factory.getAllEvent = function (callback) {
    $http.get("/get_all_event/").success(function (output) {
      console.log("output: ",output)
      callback(output);
    });
  };
  factory.getAllTicket = function (id, callback) {
    $http.get("/get_all_ticket/" + id).success(function (output) {
      callback(output);
    });
  };
  factory.getTicket = function (id, callback) {
    $http.get("/get_ticket/" + id).success(function (output) {
      console.log("ok")
      callback(output);
    });
  };

  factory.recordEvent = function (data, callback) {
    var event =
      data.name +
      "_" +
      data.issuer +
      "_" +
      data.price +
      "_" +
      data.total;

    $http.get("/create_event/" + event).success(function (output) {
      callback(output);
    });
  };

  factory.buyTicketFromSupplier = function (data, callback) {
    var holder = data.key + "_" + data.number + "_" + data.owner;

    $http.get("/buyTicketFromSupplier/" + holder).success(function (output) {
      callback(output);
    });
  };

  return factory;
});
