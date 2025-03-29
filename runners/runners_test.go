package runners

import (
	"log"
	"testing"
)

func TestGetFolder(t *testing.T) {
	const folder = "/home/meplos/myFolder"
	const file = "/home/meplos/myFolder/myFile.json"

	res := GetDir(file)

	if res != folder {
		t.Errorf("Expected: %v , Received: %v\n", folder, res)
	}
	log.Printf("%v", res)
}
