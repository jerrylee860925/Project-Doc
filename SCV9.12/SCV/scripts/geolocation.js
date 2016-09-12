function initialize(){
	var mapProp = {
		center: new google.maps.LatLng(0,0),
		zoom: 5,
		mapTypeId: google.maps.MapTypeId.ROADMAP
	};
	
	var map = new google.maps.Map(document.getElementById("map"),mapProp);
}

google.maps.event.addDomListener(window, 'load', initialize);