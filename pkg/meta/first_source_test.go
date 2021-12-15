package meta

import (
	"fmt"
	"testing"
)

func TestGetMeta_firstSource(t *testing.T) {
	finder := firstSource{}
	found, meta := finder.GetMeta("Britney Spears", "Blackout")
	if !found {
		t.Fatalf("Not found")
	}
	fmt.Println(meta)
}
