package models

import "github.com/boratanrikulu/lirik.app/pkg/lyrics"

type Song struct {
	Name       string
	Lyrics     lyrics.Lyrics
	AlbumImage string
}
