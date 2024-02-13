package domain

import (
	"encoding/xml"

	"github.com/justepl2/gopro_to_gpx_api/models"
)

type Trkpt struct {
	Lat string  `xml:"lat,attr"`
	Lon string  `xml:"lon,attr"`
	Ele float64 `xml:"ele"`
}

type Trkseg struct {
	Trkpt []Trkpt `xml:"trkpt"`
}

type Trk struct {
	Name   string   `xml:"name"`
	Trkseg []Trkseg `xml:"trkseg"`
}

type GpxStruct struct {
	FileName string
	XMLName  xml.Name `xml:"gpx"`
	Trk      Trk      `xml:"trk"`
}

// ConvertOpenRouteServiceGpxIntoTrkGpx converts OpenRouteService GPX into Trk GPX
func (gpxStruct *GpxStruct) ConvertOpenRouteServiceGpxIntoTrkGpx(openRouteServiceGpxStruct models.OpenRouteServiceResponse) {

	gpxStruct.Trk.Name = "openRouteService"
	gpxStruct.Trk.Trkseg = append(gpxStruct.Trk.Trkseg, Trkseg{})

	for _, rtept := range openRouteServiceGpxStruct.Rte.Rtept {
		trkpt := Trkpt{Lat: rtept.Lat, Lon: rtept.Lon}
		gpxStruct.Trk.Trkseg[0].Trkpt = append(gpxStruct.Trk.Trkseg[0].Trkpt, trkpt)
	}
}
