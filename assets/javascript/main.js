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

async function getCurrentSongID(token = "") {
  const response = await fetch(
    "https://api.spotify.com/v1/me/player/currently-playing",
    {
      method: "GET", // *GET, POST, PUT, DELETE, etc.
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
        Authorization: token,
      },
    }
  );

  return response.json();
}

function getCookie(name) {
  var parsedCookie = document.cookie.split(";");
  var currentCookie = parsedCookie.find((cookie) =>
    cookie.includes(name + "=")
  );
  var value = currentCookie.split(name + "=");

  return value[1];
}

function checkChanges() {
  setTimeout(checkChanges, 2500);

  if (!document.hidden) {
    accessToken = getCookie("AccessToken");
    cID = getCookie("CurrentSongID");
    getCurrentSongID("Bearer " + accessToken)
      .then((data) => {
        console.log("Current Song ID: ", data.item.id);

        if (cID != data.item.id) {
          console.log("Song was changed.");

          var xhttp = new XMLHttpRequest();
          xhttp.onreadystatechange = function () {
            if (this.readyState == 4 && this.status == 200) {
              document.getElementsByTagName(
                "html"
              )[0].innerHTML = this.responseText;
            }
          };

          xhttp.open("GET", "/spotify", true);
          xhttp.send();
        }
      })
      .catch((err) => {
        console.log("No song.");
      });
  }
}

accessToken = getCookie("AccessToken");
if (accessToken != "") {
  window.onload = checkChanges();
}
