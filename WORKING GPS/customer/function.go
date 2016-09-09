/*
Package customer contains all the information pertaining
the customer stakeholder
*/
package customer
import (
    "fmt"
    "net"
    "encoding/json"
    "log"
    "gopkg.in/mgo.v2"
    "time"
      "gopkg.in/mgo.v2/bson"
      "bufio"
      "strings"
      "strconv"
)
/*
database error handler
*/
func errhandler(err error) {
  if err != nil {
        log.Fatal(err)
        return
    }
}
/*
update certain database
*/
func update(databaseAddr string,old Order,new Order,database string,collection string ){
     session, err := mgo.Dial(databaseAddr)
     errhandler(err)
     defer session.Close()
     session.SetMode(mgo.Monotonic, true)
     d := session.DB(database).C(collection)
     err = d.Update(old,new)
}

/*find old order by its id*/
func findOldOrder(databaseAddr string,database string,collection string,hexID string) Order{
  var result Order
  session, err := mgo.Dial(databaseAddr)
  if err != nil {
    panic(err)
  }
  defer session.Close()

  session.SetMode(mgo.Monotonic, true)
  c := session.DB(database).C(collection)
  err = c.FindId(bson.ObjectIdHex(hexID)).One(&result) //Works
  if err != nil {
      log.Fatal(err)
  }
  return result
}
/*get all unfinished order*/

func GetUnfinishedOrder(databaseAddr string,database string,collection string) []Order{
  var result []Order
  session, err := mgo.Dial(databaseAddr)
  if err != nil {
    panic(err)
  }
  defer session.Close()

  session.SetMode(mgo.Monotonic, true)
  c := session.DB(database).C(collection)
//  err= c.Find(nil).All(&result)
  err = c.Find(bson.M{"ordersts.status":bson.M{"$ne":"finished"}}).All(&result)
  if err != nil {
      log.Fatal(err)
  }
  return result
}



/*get specific order*/
func Get(hexID string,databaseAddr string,database string,collection string) Order{
     session, err := mgo.Dial(databaseAddr)
     errhandler(err)
     defer session.Close()
     session.SetMode(mgo.Monotonic, true)

     var result Order
     c := session.DB(database).C(collection)
     err = c.FindId(bson.ObjectIdHex(hexID)).One(&result)
     errhandler(err)
     return result
}

/*type indexServer struct{
  ReqType int
  Name string
  ID bson.ObjectId `bson:"_id,omitempty"`
}*/
func GettingIPAddr( IpAddr string, port string,ID string,Name string)string{
  var newReq indexServer
  newReq.ReqType = 0
  newReq.Name = Name
  newReq.ID=bson.ObjectIdHex(ID)
  c, err := net.Dial("tcp", IpAddr+":"+port)
  if err != nil {
        fmt.Println(err)
        return "0.0.0.0"
    }

  b,e := json.Marshal(newReq)
  if e != nil {
    fmt.Println(e)
    c.Close()
    return "0.0.0.0"
  }
  e1 := json.NewEncoder(c).Encode(b)
  if e1 != nil {
    fmt.Println(e1)
    c.Close()
    return "0.0.0.0"
  }
  //
  var ans indexServer
  var msg []byte
  time.Sleep(time.Second*3)
  err = json.NewDecoder(c).Decode(&msg)
  e = json.Unmarshal(msg,&ans)
  if e != nil {
    fmt.Println(e)
    c.Close()
    return "0.0.0.0"
  }
  if err != nil {
    fmt.Println(err)
    c.Close()
    return "0.0.0.0"
  }
  c.Close()
  return ans.IpAddr
}
func GetOrderStatusNum(input Order )int32 {
  var count int32
  count = 0;
  for input.OrderSts.Status != status[count]{
    count++
  }
  return count
}

func sendToOther(ipaddr string,mOrder Order,Port string ){

  c, err := net.Dial("tcp", ipaddr+":"+Port)
  if err != nil {
        fmt.Println(err)
        return
  }
  b,e := json.Marshal(mOrder)
  if e != nil {
      fmt.Println(e)
      c.Close()
      return
    }
    e2 := json.NewEncoder(c).Encode(b)
    if e2 != nil {
      fmt.Println(e2)
      c.Close()
      return
    }
    c.Close()
}

func GetRoll(mOrder Order) int32{
  var result indexServer
  session, err := mgo.Dial("localhost")
  if err != nil {
    panic(err)
  }
  defer session.Close()
  session.SetMode(mgo.Monotonic, true)
  c := session.DB("myDB").C("myInfo")
  err = c.Find(nil).All(&result) //Works
    if err != nil {
      log.Fatal(err)
  }
   switch {
    case mOrder.CustomerCode == result.ID.Hex():
      return 0
    case mOrder.CarrierCode == result.ID.Hex():
      return 2
    case mOrder.SupplierCode == result.ID.Hex():
      return 1
  }
  return -1
}

func GetCarrOrd(input []Order)[]Order{

  var CarrOrdr []Order
  //reusult:=GetUnfinishedOrder("mycc.cit.iupui.edu",Mdatabase,"Carrier")
  var statusNum int32
  for i := 0; i < len(input); i++ {
      statusNum =  GetOrderStatusNum(input[i])
      if(statusNum>=7 && statusNum<9){
          CarrOrdr= append (CarrOrdr, input[i])
      }
 	}
  return CarrOrdr
}
func GetSuppOrd(input []Order)[]Order{
  var suppOrdr   = make([]Order,len(input),len(input))
  var statusNum int32
  count :=0
  for i := 0; i < len(input); i++ {
      statusNum =  GetOrderStatusNum(input[i])
      if(statusNum<7){
          suppOrdr[count] = input[i]
          count++
      }

  }
  var res = make([]Order,count,count)
  for i := 0; i < count; i++ {
      res[i] = suppOrdr[i]
  }
  return res

}
func sendOrderStatus(){
  for {
        time.Sleep(time.Second*10)
        //func getUnfinishedOrder(databaseAddr string,database string,collection string) []Order{
        CarrierResult := GetUnfinishedOrder("mycc.cit.iupui.edu",Mdatabase,"Carrier")
        SupplierResult := GetUnfinishedOrder("mycc.cit.iupui.edu",Mdatabase,"Supplier")
       // CustomerResult := getUnfinishedOrder("mycc.cit.iupui.edu",Mdatabase,"Customer")
        CarrierResult = GetCarrOrd(CarrierResult)
        SupplierResult = GetSuppOrd(SupplierResult)
        //func gettingIPAddr( IpAddr string, port string,ID string,Name string)string{

        for i:=0; i< len(CarrierResult); i++{
          addCu := GettingIPAddr("bmcj.cit.iupui.edu","9999",CarrierResult[i].CustomerCode,CarrierResult[i].CustomerName)
          addSupp := GettingIPAddr("bmcj.cit.iupui.edu","9999",CarrierResult[i].SupplierCode,CarrierResult[i].SupplierName)
          sendToOther(addCu, CarrierResult[i], "9998")
          sendToOther(addSupp, CarrierResult[i], "9998")
        }
        for i:=0; i< len(SupplierResult); i++{
          addCu := GettingIPAddr("bmcj.cit.iupui.edu","9999",SupplierResult[i].CustomerCode,CarrierResult[i].CustomerName)
          addSupp := GettingIPAddr("bmcj.cit.iupui.edu","9999",SupplierResult[i].SupplierCode,SupplierResult[i].SupplierName)
          sendToOther(addCu, SupplierResult[i], "9998")
          sendToOther(addSupp, SupplierResult[i], "9998")
      }
    }
  }
func ClientListen(port string) {
    ln, err := net.Listen("tcp", port)
    if err != nil {
        fmt.Println("error\n")
        fmt.Println(err)
      return
    }
    for {
        nc, err := ln.Accept()
        if err != nil {
            fmt.Println(err)
            continue
        }
      go recivemsg(nc,Mdatabase)
    }
}


func recivemsg(nc net.Conn,database string){
    var msg []byte
    var nOrder Order
    err := json.NewDecoder(nc).Decode(&msg)
    errhandler(err)
    e := json.Unmarshal(msg,&nOrder)
    errhandler(e)
    nc.Close()
    //databaseAddr string,old Order,new Order,database string,collection string
    var collection string
    switch GetRoll(nOrder){
      case  0:
        collection = "Customer"
        break;
      case 1:
        collection = "Supplier"
        break;
      case 2 :
        collection = "Carrier"
        break;
      case -1:
        fmt.Println("wrong")
        return
    }
    oldOrder :=findOldOrder("mycc.cit.iupui.edu",database,collection,nOrder.ID.Hex())
   //oldOrderNum := getOrderStatusNum(oldOrder)
    //newOrderNum  := getOrderStatusNum(nOrder)

    update("mycc.cit.iupui.edu",oldOrder,nOrder,"client1",collection)



   /* if(newOrderNum != oldOrderNum){
      if newOrderNum< oldOrderNum {
        //send update message
      }else{
        if(newOrderNum == 8 ){
          if collection == "Customer" {
          //finalize
          }else{
          //keep updating
          }
        }
       if(newOrderNum ==9){
        //update

       }
    }

  }
*/
}
//index server
//db clients
//
func GetSupplierList(collection []string) []string{
  var mOrder []Order
  for i := 0; i < len(collection); i++ {
    //fmt.Println(GetUnfinishedOrder("mycc.cit.iupui.edu","client1",collection[i]))
    mOrder = append(mOrder, GetUnfinishedOrder("mycc.cit.iupui.edu","Client1",collection[i])...)

  }
  //:= GetUnfinishedOrder(databaseAddr string,database string,collection string)
  //func getUnfinishedOrder(databaseAddr string,database string,collection string) []Order{

  var SuppList []string
  index := -1
  for i := 0; i < len(mOrder); i++ {
    for j := 0; j < len(SuppList); j++ {
      if(mOrder[i].SupplierName==SuppList[j]){
        index = j;

        break;
      }
    }
    if(index == -1){
      SuppList = append(SuppList,mOrder[i].SupplierName)
    }
    index = -1
  }
  return SuppList
}
func GetCarrierList(collection []string) []string{
  var mOrder []Order
  for i := 0; i < len(collection); i++ {
    mOrder = append(mOrder, GetUnfinishedOrder("mycc.cit.iupui.edu","Client1",collection[i])...)

  }
  var CarrierLst []string
  index := -1
  for i := 0; i < len(mOrder); i++ {

    for j := 0; j < len(CarrierLst); j++ {
      if(mOrder[i].CarrierName==CarrierLst[j]){
        index = j;
        break;
      }
    }
    if(index == -1){
      CarrierLst = append(CarrierLst,mOrder[i].CarrierName)
    }
    index = -1
  }
  return CarrierLst
}


func GetDest(collection []string)[]string{
  var allOrder []Order//:= GetUnfinishedOrder()
   for i := 0; i < len(collection); i++ {
    allOrder = append(allOrder, GetUnfinishedOrder("mycc.cit.iupui.edu","Client1",collection[i])...)

  }
  var Dest []string
  var ifExist bool
  ifExist = false
  for i := 0; i < len(allOrder); i++ {
    for j := 0; j < len(Dest); j++ {
      if(Dest[j] == allOrder[i].Dest ){
        ifExist = true
        break;
      }
    }
    if(ifExist == false){
      Dest = append(Dest,allOrder[i].Dest)
    }
    ifExist = false

  }
  return Dest
}

func GetOrigine(collection []string)[]string{
 var allOrder []Order//:= GetUnfinishedOrder()
   for i := 0; i < len(collection); i++ {
    allOrder = append(allOrder, GetUnfinishedOrder("mycc.cit.iupui.edu","Client1",collection[i])...)

  }
  fmt.Println(len(allOrder))
    fmt.Println(allOrder)
  var Origin []string
  var ifExist bool
  ifExist = false
  for i := 0; i < len(allOrder); i++ {
    for j := 0; j < len(Origin); j++ {
      if(Origin[j] == allOrder[i].Origin ){
        ifExist = true
        break;
      }
    }
    if(ifExist == false){
      Origin = append(Origin,allOrder[i].Origin)
    }
    ifExist = false
  }
  fmt.Println(Origin)

  return Origin
}
/*


}

func GetConditionalOrder(supplier []string, carrier []string,origine []string, dest []string, startYear int, startMonth int,startDay int, endYear int, endMonth int, endDay int,databaseAddr string,database string,collection string ) ([]Order){
         //func getUnfinishedOrder(databaseAddr string,database string,collection string) []Order{

  result := getUnfinishedOrder(databaseAddr ,database ,collection)
  var res1 []Order
  var res2 []Order
  var res3 []Order
  var res4 []Order
  var res5 []Order
  if(supplier[0] == "any"){
    for i := 0; i < len(result); i++ {
      res1  = append(res1,result[i])
    }
  }else{
          for i:=0;i<len(result);i++{
            for j:=0;j<len(supplier);j++{
              if(result[i].SupplierName == supplier[j]){
                res1 = append(res1,result[i])
                break;
                }
            }
          }
  }
  if(carrier[0] == "any"){
    for i := 0; i < len(res1); i++ {
      res2  = append(res2,result[i])
    }
  }else{
  for i:=0;i<len(res1);i++{
      for j:=0;j<len(carrier);j++{
          if(res1[i].CarrierName == carrier[j]){
            res2 = append(res2,res1[i])
            break;
        }
      }
    }
  }
  if(origine[0] =="any"){
      for i := 0; i < len(res2); i++ {
        res3  = append(res3,res2[i])
      }
    }else{
      for i := 0; i< len(res2); i++ {
        for j:=0; j< len(origine); j++ {
            if(res2[i].Origin == origine[j]){
              res3 = append(res3,res2[i])
              break;
          }
        }
      }
    }
  if(dest[0] =="any"){
      for i := 0; i < len(res3); i++ {
        res4  = append(res4,res3[i])
      }
    }else{
      for i := 0; i< len(res3); i++ {
        for j:=0; j< len(dest); j++ {
            if(res3[i].Dest == dest[j]){
              res4 = append(res4,res3[i])
              break;
          }
        }
      }
    }

  startdate := parseTimeStart(startYear, startMonth,startDay)
  enddate := parseTimeStart(endYear, endMonth, endDay)
  for i := 0; i < len(res4); i++ {
    checkDate := ParseTime(res4[i].OrderDate)
    if Compare(startdate,enddate,checkDate)==true {
      res5= append(res5,res4[i])
    }
  }
  return res5
}

func parseTimeStart(year int,month int,day int)time.Time{
  locationTime :=time.Now()
  dateParse := time.Date(year , time.Month(month), day, 0, 0, 0, 0, locationTime.Location())
  fmt.Println(dateParse)
  return dateParse
}

func parseTimeEnd(year int,month int,day int) time.Time{
  locationTime :=time.Now()
  dateParse := time.Date(year , time.Month(month),day,23,59,59,999999999, locationTime.Location())
  fmt.Println(dateParse)
  return dateParse
}
func ParseTime(date string)time.Time {
  timeFormat := "2006-01-02 15:04 MST"
  then, _ := time.Parse(timeFormat, date)

  return then
}
func Compare(start time.Time,end time.Time,Mdate time.Time)bool{
  if(Mdate.Before(end)==true && Mdate.After(start)==true){
    return true
  }
  return false
}



*/

/*
func ListenGPS(port string) {
    ln, err := net.Listen("tcp", port)
    if err != nil {

        fmt.Println("listening error: ", err)
        return
    }

    for{
    	   nc,err := ln.Accept()
    	   if err != nil {
        	       fmt.Println("accepting error: ", err)

    	   }else{
                go receiveStr(nc, ln)  
                }
    }

    //nc.Close()
    
}

func receiveStr(nc net.Conn, ln net.Listener){
    //var msg []byte   
    i:=0

    for{
        message, err := bufio.NewReader(nc).ReadString('\n')
            if(err!=nil){
            fmt.Println("receiving error: ", err)
            fmt.Println("Client has disconnted before, try to open again...\n" )
            nc.Close()
            ln.Close()
            //ClientListen(":4444")
            break
            }
        i++
    fmt.Println("Content: ", message)
    updateLoc(message,i)
    }
}
*/





func ListenGPS (port string) {
    ln, err := net.Listen("tcp", port)
    if err != nil {

        fmt.Println("listening error: ", err)
        return
    }

    go gpsAccept(ln)
    
}
func gpsAccept(ln net.Listener) {

    i:=0
    nc, err := ln.Accept()
    if err != nil {
        fmt.Println("accepting error: ", err)
        nc.Close()
        //return
    }else{
        fmt.Println("============================================================")
        fmt.Println("=                 Connected with Client                    =")
        fmt.Println("============================================================")
        reciveGps(nc, i, ln)
    }

}
func reciveGps(nc net.Conn, i int, ln net.Listener){
    //var msg []byte  
    for{ 
      message, err := bufio.NewReader(nc).ReadString('\n')
      if(err!=nil){
            //fmt.Println("receiving error: ", err)
            fmt.Println("============================================================")
            fmt.Println("=  Client has disconnected. Waiting for new connection...  =")
            fmt.Println("============================================================\n\n")
            nc.Close()
            ln.Close()
            go ListenGPS(":4444")

            break
      }else{

        updateLoc(message,i)
        i++
      }
    }
}


func updateLoc(message string, i int){



  var upOrder Order
  if(message != ""){
    fmt.Println(message)
    GPS:=strings.SplitN(message,",",5)
    x,_ := strconv.ParseFloat(GPS[3],64)
    y,_ := strconv.ParseFloat(strings.Replace(GPS[4],"\n","",1),64)

    old:=findOldOrder("127.0.0.1","Client1","Carrier","57d2df6b636bd32f205559e9")
    upOrder=old
    upOrder.OrderSts.Trucks= append(old.OrderSts.Trucks,GPS[0])
    upOrder.OrderSts.TimeStamp= append(old.OrderSts.TimeStamp,GPS[2])

    upOrder.OrderSts.GPSCordsX= append(old.OrderSts.GPSCordsX,x)
    upOrder.OrderSts.GPSCordsY= append(old.OrderSts.GPSCordsY,y)
    update("127.0.0.1",old,upOrder,"Client1","Carrier")

  }
  

}


func main(){
 

}
