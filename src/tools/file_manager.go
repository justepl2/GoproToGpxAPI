package tools

import (
	"log"
	"os"
)

func DeleteTempFiles(filename string) {
	fileToRemove := os.Getenv("RAW_VIDEO_DEST_DIR")
	e := os.Remove(fileToRemove + filename + ".bin")

	if e != nil {
		log.Fatal(e)
	}
}
