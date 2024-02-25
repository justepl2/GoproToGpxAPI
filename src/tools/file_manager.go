package tools

import (
	"os"
)

func DeleteTempFiles(filename string) error {
	rawFilesToRemove := os.Getenv("RAW_VIDEO_DEST_DIR")
	gpxFilesToRemove := os.Getenv("GPX_FILES_DEST_DIR")

	e := os.Remove(rawFilesToRemove + filename + ".bin")
	if e != nil {
		return e
	}
	e = os.Remove(gpxFilesToRemove + filename + ".gpx")
	if e != nil {
		return e
	}

	return nil
}
