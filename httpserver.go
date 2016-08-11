/*
Package httpserver runs the local HTTP server for a client
this is the presentation layer in the three layered structure
of the Supply Chain Visibility Project.
*/
package main

import(
	"fmt"
	"net/http"
	"html/template"
	"strings"
	customer "github.com/alvarosness/customer"
	"io/ioutil"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

//Date is composed of Year, Month and Day
//not much I can say about that.
type Date struct{
	Year int
	Month int
	Day int
}

//Request contains the information necessary
//to query a shipment from the database.
type Request struct{
	//ID to be searched.
	ID string
	//Suppliers filter.
	Suppliers []string
	//Carriers filter.
	Carriers []string
	//Shipment's origin filter.
	FromState []string
	//Shipment's destination filter.
	ToState []string
	//Lower bound for the shipment date placement filter.
	StartDate Date
	//Upper bound for the shipment date placement filter.
	EndDate Date
}//end Request

//loadPageInfo contains information necessary to load
//the page for thr first time.
type loadPageInfo struct{
	//Carriers will contain a list of
	//all the carriers present in any of the unfinished
	//shipments stored in the user's database.
	Carriers []string
	//Suppliers will contain a list of
	//all the suppliers present in any of the unfinished
	//shipments stored in the user's database.
	Suppliers []string
	//List of states of origin of all the unfinished shipments
	//stored in the user's database.
	Origins []string
	//List of destination states of all the unfinished shipments
	//stored in the user's database.
	Destinations []string
	//List of all the unfinished shipments in which the user
	//plays the role of Customer.
	Shipments []customer.Order
	//List of all the unfinished shipments in which the user
	//plays the role of Carrier.
	CarrierShipments []customer.Order
	//List of all the unfinished shipments in which the user
	//plays the role of Supplier.
	SupplierShipments []customer.Order
}//end loadPage

//ordersHandler will load and display the page that will list all of
//the shipments stored in the user's database and handle query requests
//issued by the browser.
func ordersHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println(r.Method)
	if r.Method == "GET" {

		var orders loadPageInfo

		//loading all the data to be displayed on the page.
		orders.Shipments = customer.GetUnfinishedOrder("mycc.cit.iupui.edu", "Client1", "Customer")
		orders.CarrierShipments = customer.GetUnfinishedOrder("mycc.cit.iupui.edu", "Client1", "Carrier")
		orders.SupplierShipments = customer.GetUnfinishedOrder("mycc.cit.iupui.edu", "Client1", "Supplier")
		orders.Suppliers = customer.GetSupplierList([]string{"Customer","Supplier","Carrier"})
		orders.Carriers = customer.GetCarrierList([]string{"Customer","Supplier","Carrier"})
		orders.Origins = customer.GetOrigine([]string{"Customer","Supplier","Carrier"})
		orders.Destinations = customer.GetDest([]string{"Customer","Supplier","Carrier"})

		//parsing the HTML template from the local resources.
		t := template.New("page")
		t = t.Funcs(template.FuncMap{"makeHex":func (v bson.ObjectId) string {return v.Hex()}})
		t.ParseFiles("templates/page.gohtml")
		t.ExecuteTemplate(w, "page", orders)
		//Writing the HTML to the browser?? -- not so sure about this comment

	}else if r.Method == "POST" {

		var msg Request

		//Reading the request Body.
		//The Body should be a JSON string.
		message, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		//Converting the JSON string to an object.
		// The object should be of type Request.
		json.Unmarshal(message, &msg)

		//response := customer.GetConditionalOrder(msg.Suppliers, msg.Carriers, msg.FromState, msg.ToState, msg.StartDate.Year, msg.StartDate.Month, msg.StartDate.Day, msg.EndDate.Year, msg.EndDate.Month, msg.EndDate.Day, "mycc.cit.iupui.edu", "client1",  "customer")
		//b, err := json.Marshal(response)
		//if err != nil {
		//	panic(err)
		//}
		//fmt.Fprintf(w, "%s", b)
	}//end elseif
}//end ordersHandler

//orderHandler will load and display a page with detailed information about
//a specific unfinished shipment stored in the user's local database.
func orderHandler(w http.ResponseWriter, r *http.Request){
	//text contains the order number and the client's role in the shipment.
	//This information is extracted from the URL path.
	text := r.URL.Path[len("/orders/order/"):]
	//textarr will contain the split text string.
	textarr := strings.Split(text, "-")
	orderNum := textarr[0]
	collection := textarr[1]

	//querying the database for the selected shipment
	order := customer.Get(orderNum, "mycc.cit.iupui.edu", "Client1", collection)

	t := template.New("order")

	t = t.Funcs(template.FuncMap{"makeHex":func (v bson.ObjectId) string {return v.Hex()}})
	t.ParseFiles("templates/order.gohtml")
	t.ExecuteTemplate(w, "order.gohtml", order)
	if(r.Method == "POST"){

	}
}//end orderHandler


func main(){
	//listening on the port for requests
	go customer.ClientListen(":9999")

	//creating static file servers
	http.Handle("/stylesheets/", http.StripPrefix("/stylesheets/", http.FileServer(http.Dir("stylesheets"))))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("scripts"))))

	//defining the handlers for the different routes
	http.HandleFunc("/orders/", ordersHandler)
	http.HandleFunc("/orders/order/", orderHandler)

	//Running the HTTP server
	http.ListenAndServe(":8889", nil)

}
