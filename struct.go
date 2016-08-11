package customer

import (
    "gopkg.in/mgo.v2/bson"
)

type indexServer struct{
	ReqType int
	Name string
  IpAddr string
  ID bson.ObjectId `bson:"_id,omitempty"`
}

type Order struct {
  ID        bson.ObjectId `bson:"_id,omitempty"`
	OrderList ProductList
	CustomerName,CustomerCode,SupplierName,SupplierCode,CarrierName,CarrierCode string
  Origin,Dest string
  OrderSts OrderStatus
	PickUptime string
	ETA int
	OrderDate string
}

type OrderStatus struct{
  Status string
  GPSCordsX []float64
  GPSCordsY []float64
}

type product struct {
	ProductName,ProductCode,ProductState,UnitMeasure string 
	UnitPrice,Quantity float64
}

type ProductList struct{
	ListOfProduct []product
}

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
