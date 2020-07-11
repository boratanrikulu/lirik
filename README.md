# S-Lyrics

[![Go Report Card](https://goreportcard.com/badge/github.com/boratanrikulu/S-Lyrics)](https://goreportcard.com/report/github.com/boratanrikulu/S-Lyrics)

Shows lyrics for currently playing song in your spotify account.

## Features

- Auto detects the currently playing song in the spotify account
- 3 resources to take lyrics;
	- Local storage  
	> That's what we call it. Basically, it is the main database that is created by us.
	- If the lyrics is not exist on the database, then checks lyrictranslate.com and genius.com for taking lyrics.
- Shows lyrics and it's translations.
- Simple UI.
- No Ads. Never.

## Technologies

- Go
- Colly
- Mux, net/http
- Bulma

<p align="center">
	<img src="example.png" alt="site-example">
</p>
