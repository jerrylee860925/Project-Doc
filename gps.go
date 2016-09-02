package main
import (
    "fmt"
    "net"
    "bufio"
   "time"
)


func ClientListen(port string) {
    ln, err := net.Listen("tcp", port)
    if err != nil {

        fmt.Println("listening error: ", err)
        return
    }

    go accept(ln)
    
}

func accept(ln net.Listener) {

    i:=0
    nc, err := ln.Accept()
    if err != nil {
        fmt.Println("accepting error: ", err)
        nc.Close()
        //return
    }else{
        fmt.Println("Connected with Client")
    }

       recivemsg(nc, i, ln)
       ln.Close()

}
func recivemsg(nc net.Conn, i int, ln net.Listener){
    //var msg []byte   

    for{
        message, err := bufio.NewReader(nc).ReadString('\n')
            if(err!=nil){
            fmt.Println("receiving error: ", err)
            fmt.Println("Client has disconnted before, try to open again...\n" )
            nc.Close()
            ln.Close()
            ClientListen(":4444")
            break
            }
        i++
    fmt.Println(i,": ", message)
    }
}


func main(){

    go ClientListen(":4444")
    
  time.Sleep(time.Second*1000)
}