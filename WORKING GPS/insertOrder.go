package main

import (
 "log"
 "gopkg.in/mgo.v2"
 "gopkg.in/mgo.v2/bson"
 "strconv"
 "fmt"
 "time"
 "math"
 "math/rand"
)
type indexServer struct{
  ReqType int
  Name string
  IpAddr string
  ID bson.ObjectId `bson:"_id,omitempty"`
}
type product struct {
  ProductName,ProductCode,ProductState,UnitMeasure string 
  UnitPrice,Quantity float64
}
type ProductList struct{
  ListOfProduct []product
}
type Order struct {
  ID bson.ObjectId `bson:"_id,omitempty"`
  OrderList ProductList
  CustomerName,CustomerCode,SupplierName,SupplierCode,CarrierName,CarrierCode string

  Origin, Dest string
  OrderSts OrderStatus
  PickUptime string
  Preparation string
  ETA int
  OrderDate string
}
type OrderStatus struct {
  Status string
  Trucks []string
  TimeStamp []string
  GPSCordsX []float64
  GPSCordsY []float64
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



func CalcDay(date string)int {

  timeFormat := "2006-01-02 15:04 MST"
  then, err := time.Parse(timeFormat, date)
  if err != nil {
    fmt.Println(err)
    return (-1)
    }
    duration := time.Since(then)
    day := math.Ceil(duration.Hours()/24)
    return int(day)
}

func errhandler(err error,something string) {
  if err != nil {
        log.Fatal(err)
        fmt.Println(something)
    }
}
type col struct{
     Id string
     Password string
     Uid int
     Other string
}
var m bson.M


func indexPut(in indexServer){


    session, err := mgo.Dial("localhost")
    errhandler(err,"connection")
    defer session.Close()
    session.SetMode(mgo.Monotonic, true)
    d := session.DB("IndexServer").C("Clients")
    err = d.Insert(&in)
    errhandler(err,"db")  
}



func clientPut(in Order, flag int, clientCode string){

    var collection string
    var database string

    if flag == 1{
      database = "Client" + clientCode
      collection = "Supplier"
    }else if flag == 2{
      database = "Client" + clientCode
      collection = "Carrier"
    }else if flag == 3{
      database = "Client" + clientCode
      collection = "Customer"
    }

     session, err := mgo.Dial("localhost")
     errhandler(err,"connection")
     defer session.Close()
     session.SetMode(mgo.Monotonic, true)
     d := session.DB(database).C(collection)
     err = d.Insert(&in)
     errhandler(err,"db")     
}


func dropCol(collectionName string,dbName string ){
     session, err := mgo.Dial("localhost")
     errhandler(err,"connection")
     defer session.Close()
     session.SetMode(mgo.Monotonic, true)
    
     d := session.DB(dbName).C(collectionName)
     err = d.DropCollection()
     errhandler(err,"db")
}

func indexSI(){

     dropCol("Clients","IndexServer")

     var records [6]indexServer

     records[0].Name = "pepsi"
     records[0].IpAddr = "127.0.0.1" 
     records[1].Name = "cokecola"
     records[1].IpAddr = "127.0.0.1"

     records[2].Name = "walmart"
     records[2].IpAddr = "127.0.0.1"      
     records[3].Name = "costco"
     records[3].IpAddr = "127.0.0.1"  

     records[4].Name = "fedex"
     records[4].IpAddr = "127.0.0.1" 
     records[5].Name = "ups"
     records[5].IpAddr = "127.0.0.1"

    for i := 0; i < 6; i++ {
        indexPut(records[i])
    }
}

func findCode(name string)string{
    var result indexServer
    session, err := mgo.Dial("localhost")
    if err != nil {
      panic(err)
    }  
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)
    c := session.DB("IndexServer").C("Clients")

    err = c.Find(bson.M{"name":name}).One(&result)
    if err != nil {
    log.Fatal(err)        
    }
    return result.ID.Hex()

}


func INSERT (coll string, role1 string, role2 string, role3 string,flag int){

      

      year := "2016"
      var tmp [3]Order
      var tmp4 ProductList
      var tmp5 ProductList
      tmp2 := make([]product,10,10) 
      tmp3 := make([]product,10,10)
      for i:=0;i<10;i++{
        tmp2[i].ProductCode = strconv.Itoa((i+1)*rand.Intn(1700))
        tmp2[i].ProductName = strconv.Itoa((i+1)*rand.Intn(200))
        tmp2[i].ProductState = "liquid"
        tmp2[i].UnitMeasure = "liter"
        tmp2[i].UnitPrice = float64(i+1)*(1.7)
        tmp2[i].Quantity  = float64(i+1)*(2)
      }
       for i:=0;i<10;i++{
        tmp3[i].ProductCode = strconv.Itoa((i+1)*rand.Intn(3000))
        tmp3[i].ProductName = strconv.Itoa((i+1)*rand.Intn(400))
        tmp3[i].ProductState = "gas"
        tmp3[i].UnitMeasure = "gallon"
        tmp3[i].UnitPrice = float64(i+1)*(1.7)
        tmp3[i].Quantity  = float64(i+1)*(2)
      }
      tmp4.ListOfProduct =tmp2 
      tmp5.ListOfProduct =tmp3

      tmp[0].OrderList =  tmp4
      tmp[1].OrderList =  tmp5
      tmp[2].OrderList =  tmp5
     
      for i:=0;i<3;i++{


        tmp[i].CustomerName = role1
        tmp[i].CustomerCode = findCode(tmp[i].CustomerName)
        tmp[i].SupplierName = role2
        tmp[i].SupplierCode = findCode(tmp[i].SupplierName)
        tmp[i].CarrierName = role3
        tmp[i].CarrierCode = findCode(tmp[i].CarrierName)

      if i%2==0{
        tmp[i].OrderSts.Status = status[9]

      }else{
        tmp[i].OrderSts.Status  = status[rand.Intn(9)]
      }

        month := strconv.Itoa(i%5+1)
        day := strconv.Itoa((i*13)%30)
        total := year+"-"+month+"-"+day+" 15:04 MST"
        tmp[i].OrderDate=total


    }
    if flag==1{
          tmp[0].Origin = "IL"
          tmp[1].Origin = "IN"
          tmp[2].Origin = "CA"
          tmp[0].Dest = "ID"
          tmp[1].Dest = "IA"
          tmp[2].Dest = "KS"
    }else if flag==2{
          tmp[0].Origin = "AL"
          tmp[1].Origin = "AZ"
          tmp[2].Origin = "CO"
          tmp[0].Dest = "KS"
          tmp[1].Dest = "KY"
          tmp[2].Dest = "LA"
    }else if flag==3{
          tmp[0].Origin = "DE"
          tmp[1].Origin = "FL"
          tmp[2].Origin = "GA"
          tmp[0].Dest = "MI"
          tmp[1].Dest = "MN"
          tmp[2].Dest = "NY"
    }

    
    for i := 0; i < 3; i++ {

         clientPut(tmp[i],flag,"1")
    }
}


func main(){

  indexSI()
  dropCol("Customer","Client1" )
  dropCol("Supplier","Client1" )
  dropCol("Carrier","Client1" )
  INSERT("Customer","pepsi","walmart","fedex",1)
  INSERT("Supplier","walmart","fedex","pepsi",2)
  INSERT("Carrier","fedex","pepsi","walmart",3)



}
