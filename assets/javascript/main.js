function openTab(evt, tabName) {
	var i, x, tablinks;
	x = document.getElementsByClassName("card-content");
	for (i = 0; i < x.length; i++) {
		x[i].style.display = "none";
	}
	tablinks = document.getElementsByClassName("tab");
	for (i = 0; i < x.length; i++) {
		tablinks[i].className = tablinks[i].className.replace(" is-active", "");
	}
	document.getElementById(tabName).style.display = "block";
 	evt.currentTarget.className += " is-active";
}

DarkReader.enable({
    brightness: 100,
    contrast: 90,
    sepia: 10
});
