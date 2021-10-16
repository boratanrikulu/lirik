navigator.serviceWorker.getRegistrations().then( function(registrations) { for(let registration of registrations) { registration.unregister(); } }); 

function openTab(event, tabName) {
  var cardContentElems = document.getElementsByClassName("card-content");
  var tabLinkElems = document.getElementsByClassName("tab");
  var currentContentEl = document.getElementById(tabName);

  for (cardContentEl of cardContentElems) {
    cardContentEl.style.display = "none";
  }

  for (tabLinkEl of tabLinkElems) {
    tabLinkEl.classList.remove("is-active");
  }

  currentContentEl.style.display = "block";
  event.currentTarget.classList.add("is-active");
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
  if (currentCookie == null) {
    return "";
  }

  var value = currentCookie.split(name + "=");
  return value[1];
}

async function getCurrentPage() {
  const loadingEl = document.querySelector(".lirik-loading");

  loadingEl.style.visibility = "visible";

  fetch("/spotify")
    .then(function (response) {
      return response.text();
    })
    .then(function (html) {
      loadingEl.style.display = "hidden";
      document.querySelector("html").innerHTML = html;
    })
    .catch(function (err) {
      loadingEl.style.display = "hidden";
      console.warn("Something went wrong.", err);
    });
}

function checkChanges() {
  if (document.hidden) return;

  var accessToken = getCookie("AccessToken");
  var refreshToken = getCookie("RefreshToken");

  if (accessToken == "" && refreshToken != "") {
    location.reload();
  }

  getCurrentSongID()
    .then((data) => {
      if (currentSongID != data.item.id) {
        currentSongID = data.item.id;
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

var currentSongID = "";
if (getCookie("AccessToken") != "" && window.location.pathname == "/spotify") {
  getCurrentSongID().then((data) => {
    currentSongID = data.item.id;
  });

  checkChangesTimer();
}
