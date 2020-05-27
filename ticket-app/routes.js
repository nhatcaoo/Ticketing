//SPDX-License-Identifier: Apache-2.0

var ticketing = require('./controller.js');

module.exports = function(app){

  app.get('/get_ticket/:ticketId', function(req, res){
    ticketing.get_ticket(req, res);
  });
  app.get('/create_event/:info', function(req, res){
    ticketing.create_event(req, res);
  });
  app.get('/get_all_ticket/:eventId', function(req, res){
    ticketing.get_all_ticket(req, res);
  });
  app.get('/buyTicketFromSupplier/:info', function(req, res){
    ticketing.change_holder(req, res);
  });
}
