//SPDX-License-Identifier: Apache-2.0

var ticketing = require('./controller.js');

module.exports = function(app){
  app.get('/get_all_event/', function(req, res){
    ticketing.get_all_event(req, res);
  });
  app.get('/get_ticket/:id', function(req, res){
    ticketing.get_ticket(req, res);
  });
  app.get('/create_event/:event', function(req, res){
    ticketing.create_event(req, res);
  });
  app.get('/get_all_ticket/:id', function(req, res){
    ticketing.get_all_ticket(req, res);
  });
  app.get('/buy_ticket_from_supplier/:holder', function(req, res){
    ticketing.buy_ticket_from_supplier(req, res);
  });
}
