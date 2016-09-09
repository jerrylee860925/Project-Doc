package main

import (
"time"
"./customer"
//"fmt"
)

func main(){

	customer.ListenGPS(":4444")
    
  time.Sleep(time.Second*1000)
}
