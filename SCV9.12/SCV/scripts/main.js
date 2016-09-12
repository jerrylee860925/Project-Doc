var main = function(){
	$("#search-bar").keypress(function(event){
		var keycode = (event.keyCode ? event.keyCode : event.which);
		if(keycode == 13){
			makeRequest();
		}
	});
};


function makeRequest(){
	var text = $("#search-bar").val();
	var code = parseInt($("input[name='search-param']:checked").val());
	$(".table").children("tbody").empty();

	if(text === ""){
		$("#select-form").hide();
		$("#no-results").show();
	}
	else{	
		//creating the request object and feeding it the information from the search bar
		//encoding it in the JSON format and creating a new Http request
		var requestObject = {ReqType: code, KeyWord: text};
		var requestText = JSON.stringify(requestObject);
		var xhttp = new XMLHttpRequest();

		xhttp.onreadystatechange = function(){
			if(xhttp.readyState == 4 && xhttp.status == 200){
				var obj = JSON.parse(xhttp.responseText);	
				for(i = 0; i < obj.Name.length ; i++){
				$(".table").children("tbody").append("<tr><td><input type='radio' name='suppliers' value='"+ obj.IpAddr[i] +"'>" + obj.Name[i] + "</td><td>" + obj.IpAddr[i] + "</td></tr>");
				}
				$("#select-form").show();
				$("#no-results").hide();
				
			}
		};

		xhttp.open("POST", "http://localhost:8889/", true);
		xhttp.send(requestText);
	}
}

$(document).ready(main);
