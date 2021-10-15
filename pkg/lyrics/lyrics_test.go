package lyrics

import (
	"testing"
)

func Test_GetLyrics(t *testing.T) {
	finder := NewFinder()
	songs := []struct {
		artistName string
		songName   string
	}{
		{
			artistName: "Duman",
			songName:   "Ah",
		},
	}
	for _, s := range songs {
		found, _ := finder.GetLyrics(s.artistName, s.songName)
		if !found {
			t.Error("Lyrics is not found")
		}
	}
}
