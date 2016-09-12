var main = function() {
	$("#destination-drop").click(function(){
			$("#destination-icon-drop").toggleClass("glyphicon-chevron-right");
			$("#destination-icon-drop").toggleClass("glyphicon-chevron-down");
	});

	$("#origin-drop").click(function(){
			$("#origin-icon-drop").toggleClass("glyphicon-chevron-right");
			$("#origin-icon-drop").toggleClass("glyphicon-chevron-down");
	});

	$("#carriers-drop").click(function(){
			$("#carriers-icon-drop").toggleClass("glyphicon-chevron-right");
			$("#carriers-icon-drop").toggleClass("glyphicon-chevron-down");
	});

	$("#suppliers-drop").click(function(){
			$("#suppliers-icon-drop").toggleClass("glyphicon-chevron-right");
			$("#suppliers-icon-drop").toggleClass("glyphicon-chevron-down");
	});

    $("#search-bar-input").keypress(function(event) {
        var keycode = (event.keyCode ? event.keyCode : event.which);
        if (keycode == 13) {
            processRequest();
        }
    });
};


//this will create a request based on the search parameters that are selected
function processRequest() {

    //catching all the information needed for the request
    var text = $("#search-bar-input").val();
    var supplier = [];
    var carrier = [];
	var origin = [];
	var destination = [];
    var startYear = parseInt($("select[name='start-year']").val());
    var startMonth = parseInt($("select[name='start-month']").val());
    var startDay = parseInt($("select[name='start-day']").val());
    var endYear = parseInt($("select[name='end-year']").val());
    var endMonth = parseInt($("select[name='end-month']").val());
    var endDay = parseInt($("select[name='end-day']").val());
    var j = 0;
    var k = 0;
	var l = 0;
	var m = 0;

    $("input[name='supplier']:checked").each(function() {
        var values = $(this).val();
        supplier[j++] = values;
    });

    $("input[name='carrier']:checked").each(function() {
        var values = $(this).val();
        carrier[k++] = values;
    });

    $("input[name='from-state']:checked").each(function() {
        var values = $(this).val();
        origin[l++] = values;
    });

    $("input[name='to-state']:checked").each(function() {
        var values = $(this).val();
        destination[m++] = values;
    });
	
    //creating the request object [a struct that can be understood by the Go program]
    var requestObject = {
        Text: text,
        Suppliers: supplier,
        Carriers: carrier,
		FromState: origin,
		ToState: destination,
        StartDate: {
            Year: startYear,
            Month: startMonth,
            Day: startDay
        },
        EndDate: {
            Year: endYear,
            Month: endMonth,
            Day: endDay
        }
    };

    //converting that request into a JSON string
    var requestText = JSON.stringify(requestObject);

	alert(requestText);

    //creating a new HTTPRequest that will inteface with the local HTTP server
    var xhttp = new XMLHttpRequest();

    xhttp.onreadystatechange = function() {
        if (xhttp.readyState == 4 && xhttp.status == 200) {
			displayFilters(requestObject);
            processResponse(xhttp.responseText);
        } //end if
    }; //end function

    //sending the request to the local HTTP server
    xhttp.open("POST", "http://localhost:8889/orders/", true);
    xhttp.send(requestText);

    /*
    	if(text == ""){
    		$("#filters-used").empty();
    		$("#filters-used").removeClass("alert alert-info");
    		$("#filters-used").addClass("alert alert-danger");
    		$("#filters-used").append('<a href="#" class="close" data-dismiss="alert" aria-label="close">&times;</a>No search parameters selected');	
    	}else{
    		$("#filters-used").empty();
    		$("#filters-used").removeClass("alert alert-danger");
    		$("#filters-used").addClass("alert alert-info");
    		$("#filters-used").append('<a href="#" class="close" data-dismiss="alert" aria-label="close">&times;</a>Searched for <strong>'+ text +'</strong> using the following filters' + carrier + " " + supplier);
    	}*/


} //end procesRequest


function displayFilters(requestObject){
	$("#filters-used").empty();
	$("#filters-used").removeClass("alert alert-danger");
	$("#filters-used").addClass("alert alert-info");

	if(requestObject.Text == "")
	{
		$("#filters-used").append('<a href="#" class="close" data-dismiss="alert" aria-label="close">&times;</a>Searched using the following filters: ');
	}else{
	   	$("#filters-used").append('<a href="#" class="close" data-dismiss="alert" aria-label="close">&times;</a>Searched for <strong>'+ requestObject.Text +'</strong> using the following filters');
	}//end if
}//displayFilters

function processResponse(responseText) {
    var obj = JSON.parse(responseText);
	var i;

    //removing all the list elements previously displayed
    $("#shipment-list").empty();

    //listing the results obtained from the local HTTP server
    for (i = 0; i < obj.length; i++) {

       $("#shipment-list").append('<a href="/orders/order/' + obj[i].Ordernumber + '" class="list-group-item"><h4 class="list-group-item-heading">Order #' + obj[i].Ordernumber + '</h4><p class="list-group-item-text">Order Status: ' + obj[i].OrderSts.Status + ' ETA: ' + obj[i].ETA + '</p></a>');

    } //end for
} //end processResponse

$(document).ready(main);
