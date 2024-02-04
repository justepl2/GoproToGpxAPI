package tools

import (
	"log"
	"os"
)

func DeleteTempFiles(filename string) {
	rawFilesToRemove := os.Getenv("RAW_VIDEO_DEST_DIR")
	gpxFilesToRemove := os.Getenv("GPX_FILES_DEST_DIR")

	e := os.Remove(rawFilesToRemove + filename + ".bin")
	if e != nil {
		log.Fatal(e)
	}
	e = os.Remove(gpxFilesToRemove + filename + ".gpx")
	if e != nil {
		log.Fatal(e)
	}
}
