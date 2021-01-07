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

async function getCurrentSongID() {
  var accessToken = getCookie("AccessToken");
  var token = "Bearer " + accessToken;

  const response = await fetch(
    "https://api.spotify.com/v1/me/player/currently-playing",
    {
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

async function getCurrentPage() {
  fetch("/spotify")
    .then(function (response) {
      return response.text();
    })
    .then(function (html) {
      document.querySelector("html").innerHTML = html;
    })
    .catch(function (err) {
      console.warn("Something went wrong.", err);
    });
}

function checkChanges() {
  if (document.hidden) return;

  getCurrentSongID()
    .then((data) => {
      var cID = getCookie("CurrentSongID");
      if (cID != data.item.id) {
        getCurrentPage();
      }
    })
    .catch((err) => {
      console.log("No song.");
    });
}

function checkChangesTimer() {
  checkChanges();
  setTimeout(checkChangesTimer, 2500);
}

var accessToken = getCookie("AccessToken");
if (accessToken != "") {
  checkChangesTimer();
}
