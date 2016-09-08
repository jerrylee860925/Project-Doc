/*************************************************
struct.go contains all struct that are needed 
in both http server and middleware. 
 *
* @author  zhijie li 
* @author  haoyan wu
* @data 09-02-2016
**************************************************/
package customer

import (
    "gopkg.in/mgo.v2/bson"
)
/**
 * @description the indexServer struct for index server
 * distinguishes different kinds of requests, and responses
 * clients with another client's name and IP address.
*/
type indexServer struct{
  ReqType int                                     //to differentiate request for responding
  Name string                                     //clients' name
  IpAddr string                                   //corresponding client's IP address 
  ID bson.ObjectId `bson:"_id,omitempty"`         //unique ID of each object which contains client's names and IPs
}

/**
 * @description the Order struct for information of order 
 * which contains the info of stakeholders, the products
 * the shipment and order itself.
*/
type Order struct {
  ID        bson.ObjectId `bson:"_id,omitempty"`  //unique ID of each order
  OrderList ProductList                           //substruct containing a list of info of product
  CustomerName,CustomerCode string                //customer's name and it's unique ID
  SupplierName,SupplierCode string                //supplier's name and it's unique ID
  CarrierName,CarrierCode string                  //carrier's name and it's unique ID
  Origin,Dest string                              //origin and destination of the shipment
  OrderSts OrderStatus                            //substruct containing the status of shipment
  PickUptime string                               //the time fo r supplier when the truck will arrive 
  ETA int                                         //estimated time of arrival
  OrderDate string                                //the time the order has been placed
}

/**
 * @description the OrderStatus struct which is
 * substruct in Order struct for status of order 
 * that contains the status of shipment and the 
 * info of geolocation.
*/
type OrderStatus struct{
  Status string                                   //substruct containing the status of shipment
  GPSCordsX []float64                             //An array for GPS only having longitudes
  GPSCordsY []float64                             //An array for GPS only having latitude
}

/**
 * @description the Product struct which is
 * substruct in ProductList struct for recording
 * the basic info of product
*/
type product struct {
  ProductName,ProductCode string                  //product's name and it's unique ID
  ProductState,UnitMeasure string                 
  UnitPrice,Quantity float64
}

/**
 * @description the Product struct which is
 * substruct in Order struct containing a 
 * list of product
*/
type ProductList struct{
  ListOfProduct []product
}

/**
 * @description an array containing each 
 * status of shipment
*/
var status = [...]string {
  "on the ship",
  "off the ship",
  "on the dock",
  "in the storage",
  "getting ready",
  "ready for pickup",
  "picked up",
  "in transit",
  "delivered",
  "finished",
}

var Mdatabase string = "client1"
