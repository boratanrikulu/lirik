package meta

import "log"

type Finder interface {
	GetMeta(string, string) (bool, Meta)
}

type finder struct {
	sources []source
}

type source struct {
	name   string
	maxTry int
	f      Finder
}

func NewFinder() Finder {
	return &finder{
		sources: []source{
			{
				name:   "LOCAL",
				maxTry: 1,
				f:      newLocalSource(),
			},
			{
				name:   "FIRST",
				maxTry: 1,
				f:      newFirstSource(),
			},
		},
	}
}

type Meta struct {
	Genre string
	Style string
}

// GetMeta returns meta infos for the Album.
func (f *finder) GetMeta(artistName string, albumName string) (found bool, meta Meta) {
	for _, source := range f.sources {
		for i := 0; i < source.maxTry; i++ {
			if f, meta := source.f.GetMeta(artistName, albumName); f {
				if source.name != "LOCAL" {
					go saveToFile(artistName, albumName, meta)
				}
				log.Printf("[%s] [FOUND META] [x%d] %s", source.name, i+1, meta)
				return true, meta
			}
		}
	}

	log.Printf("[ALL] [NOT FOUND META] \"%s by %s\"", albumName, artistName)
	return false, meta
}
